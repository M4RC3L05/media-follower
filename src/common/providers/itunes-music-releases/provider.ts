import type {
  IProvider,
  IProviderFeedGetOutputsFeedProps,
  IProviderRepositoryQueryOutputsProps,
} from "../interfaces.ts";
import type {
  DbInputsTable,
  DbOutputsTable,
} from "#src/common/database/types.ts";
import type { IHttpFetch } from "#src/common/http/mod.ts";
import type { IDatabase } from "#src/common/database/database.ts";
import {
  type Input,
  ITunesLookupEntityType,
  type ItunesMusicReleasesInput,
  itunesMusicReleasesInputSchema,
  itunesMusicReleasesInputWithExtraSchema,
  type ItunesMusicReleasesOutput,
  itunesMusicReleasesOutputAlbumSchema,
  type ItunesMusicReleasesOutputSong,
  itunesMusicReleasesOutputSongSchema,
  type ItunesResponseModel,
  type Output,
} from "./types.ts";
import { EInputProvider } from "#src/common/database/enums/input-provider.ts";
import { inputListItem, outputListItem } from "./components/mod.tsx";
import type { VNode } from "preact";
import z from "@zod/zod";
import { Feed } from "feed";

const resolveArtistImage = async (httpClient: IHttpFetch, url: string) => {
  const textDecoder = new TextDecoder();
  const response = await httpClient.fetchReadable(url);

  let html = "";

  for await (const chunk of response) {
    html += textDecoder.decode(chunk);

    if (
      /<meta property="og:image".*>/im.test(html) || /<\/head>/im.test(html)
    ) {
      break;
    }
  }

  const result = /<meta\s+property="og:image"\s+content="([^"]*)"/im
    .exec(html)
    ?.at(1);

  if (!result || result.includes("apple-music-")) {
    return;
  }

  const imageSplitted = result.split("/");
  const imageFile = imageSplitted.at(-1)?.split(".");

  if (!imageFile) {
    return;
  }

  imageSplitted[imageSplitted.length - 1] = `256x256.${imageFile.at(1)}`;

  return imageSplitted.join("/");
};

const usableOutput = (release: Output) => {
  const isCompilation =
    (release.wrapperType === "track"
      ? (release as ItunesMusicReleasesOutputSong).collectionArtistName
      : release.artistName)
      ?.toLowerCase?.()
      ?.includes?.("Various Artists".toLowerCase());

  const isDjMix = release.collectionName?.toLowerCase?.()
    ?.includes?.("DJ Mix".toLowerCase()) ||
    release.collectionCensoredName?.toLowerCase?.()
      ?.includes?.("DJ Mix".toLowerCase());

  const isStreamable = release.wrapperType === "track"
    ? (release as ItunesMusicReleasesOutputSong).isStreamable
    : true;

  return !isCompilation && !isDjMix && isStreamable;
};

export type ItunesMusicReleasesProviderProps = {
  httpClient: IHttpFetch;
  database: IDatabase;
};

export class ItunesMusicReleasesProvider
  implements IProvider<EInputProvider.ITUNES_MUSIC_RELEASE, Input, Output> {
  #lookupUrl = "https://itunes.apple.com/lookup";
  #props: ItunesMusicReleasesProviderProps;

  constructor(props: ItunesMusicReleasesProviderProps) {
    this.#props = props;
  }

  readonly provider = EInputProvider.ITUNES_MUSIC_RELEASE;

  async lookupInput(term: string): Promise<Input | undefined> {
    const url = new URL(this.#lookupUrl);
    url.searchParams.set("id", term);

    const response = await this.#props.httpClient.fetch<
      ItunesResponseModel<ItunesMusicReleasesInput>
    >(url);

    const remote = response.results.at(0);

    if (!remote) return;

    const artistImage = await resolveArtistImage(
      this.#props.httpClient,
      remote.artistLinkUrl,
    );

    return itunesMusicReleasesInputWithExtraSchema.parse({
      ...itunesMusicReleasesInputSchema.parse(remote),
      extra: { artistImage: artistImage ?? "https://placehold.co/256" },
    });
  }

  async fetchInput(row: DbInputsTable): Promise<Input> {
    const parsed = itunesMusicReleasesInputWithExtraSchema.parse(
      JSON.parse(row.raw),
    );
    const input = await this.lookupInput(String(parsed.artistId));

    if (!input) throw new Error("Unable to fetch by input", { cause: row });

    return input;
  }

  async #lookupLatestReleasesByArtist(
    artistId: string,
    entity: ITunesLookupEntityType,
    limit: number,
  ): Promise<Array<Output>> {
    const path = new URL(this.#lookupUrl);
    path.searchParams.set("id", artistId);
    path.searchParams.set("entity", entity);
    path.searchParams.set("media", "music");
    path.searchParams.set("sort", "recent");
    path.searchParams.set("limit", String(limit));

    const data = await this.#props.httpClient.fetch<
      ItunesResponseModel<ItunesMusicReleasesOutput<typeof entity>>
    >(path);

    // Remove artists info from results
    data.results.splice(0, 1);

    return data.results
      .filter((item) => usableOutput(item))
      .map((item) => {
        switch (entity) {
          case ITunesLookupEntityType.ALBUM: {
            return itunesMusicReleasesOutputAlbumSchema.parse(item);
          }
          case ITunesLookupEntityType.SONG: {
            return itunesMusicReleasesOutputSongSchema.parse(item);
          }
        }
      });
  }

  async fetchOutputs(input: Input): Promise<Output[]> {
    const [songs, albums] = await Promise.all([
      this.#lookupLatestReleasesByArtist(
        String(input.artistId),
        ITunesLookupEntityType.SONG,
        50,
      ),
      this.#lookupLatestReleasesByArtist(
        String(input.artistId),
        ITunesLookupEntityType.ALBUM,
        50,
      ),
    ]);

    return [...songs, ...albums];
  }

  // deno-lint-ignore require-await
  async queryOutputs(
    { pagination }: IProviderRepositoryQueryOutputsProps,
  ): Promise<DbOutputsTable[]> {
    return this.#props.database.sql.all`
      select id, input_id, provider, json(outputs.raw) as raw
      from outputs
      where provider = ${EInputProvider.ITUNES_MUSIC_RELEASE}
      order by outputs.raw->>'releaseDate' desc, "rowid" desc
      limit ${pagination.limit}
      offset ${pagination.page * pagination.limit}
    ` as DbOutputsTable[];
  }

  getOutputsFeed({ queries }: IProviderFeedGetOutputsFeedProps): Feed {
    const { type } = z.object({
      type: z.enum(ITunesLookupEntityType).optional(),
    })
      .parse(queries ?? {});

    const rows = this.#props.database.sql.all`
      select
        outputs.id as id,
        outputs.input_id as input_id,
        outputs.provider as provider,
        json(outputs.raw) as raw,
        related_album.raw as related_album_raw
      from outputs
      left join outputs as related_album on (
            related_album.id != outputs.id
        and related_album.provider = outputs.provider
        and related_album.input_id = outputs.input_id
        and related_album.raw->>'wrapperType' = 'collection'
        and related_album.raw->>'collectionId' = outputs.raw->>'collectionId'
      )
      where outputs.provider = ${EInputProvider.ITUNES_MUSIC_RELEASE}
      and outputs.raw->>'releaseDate' <= strftime('%Y-%m-%dT%H:%M:%fZ' , 'now')
      and (${type ?? null} is null or outputs.raw->>'wrapperType' = ${
      type === ITunesLookupEntityType.ALBUM ? "collection" : "track"
    })
      and (
        case
          when outputs.raw->>'wrapperType' = 'track'
            then
              outputs.raw->>'isStreamable' = 1
              and (
                    related_album_raw is null
                or  related_album_raw->>'releaseDate' > strftime('%Y-%m-%dT%H:%M:%fZ' , 'now')
              )
          else true
        end
      )
      order by outputs.raw->>'releaseDate' desc, outputs."rowid" desc
      limit 200
    ` as DbOutputsTable[];

    const prefix = [EInputProvider.ITUNES_MUSIC_RELEASE, type].filter(Boolean);

    const feed = new Feed({
      title: `${prefix ? `[${prefix.join(" | ")}] ` : ""}Media follower`,
      description: `Get the latest${
        prefix ? ` ${prefix.join(" and ")} ` : " "
      }media releases`,
      id: `media_follower${prefix ? `_${prefix.join("_")}` : ""}`,
      copyright: "Media Follower",
      updated: new Date(),
    });

    const outputs = rows.map((row) => ({
      ...this.fromPersistenceToOutput(row),
      row,
    }));

    for (const output of outputs) {
      feed.addItem({
        date: output.releaseDate,
        link: output.wrapperType === "collection"
          ? output.collectionViewUrl
          : output.trackViewUrl,
        title: `${output.artistName} - ${
          output.wrapperType === "collection"
            ? output.collectionName
            : output.trackName
        }`,
        id: `${EInputProvider.ITUNES_MUSIC_RELEASE}@${
          output.wrapperType === "collection"
            ? ITunesLookupEntityType.ALBUM
            : ITunesLookupEntityType.SONG
        }@${output.row.id}`,
        image: output.artworkUrl100
          .split("/")
          .map((segment, index, array) =>
            index === array.length - 1 ? "512x512bb.jpg" : segment
          )
          .join("/"),
      });
    }

    return feed;
  }

  renderInputListItem(row: DbInputsTable) {
    return inputListItem({
      input: itunesMusicReleasesInputWithExtraSchema.parse(
        JSON.parse(row.raw),
      ),
    });
  }

  renderOutputListItem(row: DbOutputsTable): VNode {
    return outputListItem({
      output: this.fromPersistenceToOutput(row),
    });
  }

  fromInputToPersistence(item: Input): DbInputsTable {
    return {
      id: String(item.artistId),
      provider: EInputProvider.ITUNES_MUSIC_RELEASE,
      raw: JSON.stringify(item),
    };
  }

  fromPersistenceToInput(row: DbInputsTable): Input {
    const parsed = JSON.parse(row.raw);

    return itunesMusicReleasesInputWithExtraSchema.parse(parsed);
  }

  fromOutputToJsonPatchPersistance(
    row: DbInputsTable,
    item: Output,
  ): DbOutputsTable {
    const id = item.wrapperType === "collection"
      ? item.collectionId
      : item.trackId;

    return {
      id: String(id),
      input_id: row.id,
      provider: row.provider,
      raw: JSON.stringify({ ...item, releaseDate: undefined }),
    };
  }

  fromOutputToPersistence(row: DbInputsTable, item: Output): DbOutputsTable {
    const id = item.wrapperType === "collection"
      ? item.collectionId
      : item.trackId;

    return {
      id: String(id),
      input_id: row.id,
      provider: row.provider,
      raw: JSON.stringify(item),
    };
  }

  fromPersistenceToOutput(row: DbOutputsTable): Output {
    const parsed = JSON.parse(row.raw);

    if (parsed.wrapperType === "collection") {
      return itunesMusicReleasesOutputAlbumSchema.parse(parsed);
    }

    return itunesMusicReleasesOutputSongSchema.parse(parsed);
  }
}
