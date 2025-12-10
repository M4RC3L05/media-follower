import type { IHttpFetch } from "#src/common/http/mod.ts";
import {
  type IItunesService,
  itunesLookupAlbumModelSchema,
  type ItunesLookupAlbumModelWithExtra,
  itunesLookupAlbumModelWithExtraSchema,
  type ItunesLookupArtistModel,
  itunesLookupArtistModelSchema,
  type ItunesLookupArtistModelWithExtra,
  itunesLookupArtistModelWithExtraSchema,
  ITunesLookupEntityType,
  type ItunesLookupSongModel,
  itunesLookupSongModelSchema,
  type ItunesLookupSongModelWithExtra,
  itunesLookupSongModelWithExtraSchema,
  type ItunesLookupType,
  type ItunesLookupTypeWithExtra,
  type ItunesResponseModel,
} from "#src/common/services/service.ts";
import { ReleaseSourceProvider } from "../database/enums/release-source-provider.ts";
import { ReleaseType } from "../database/enums/release-type.ts";
import type {
  DbReleaseSourcesTable,
  DbReleasesTable,
} from "../database/types.ts";

export type ItunesServiceProps = {
  httpClient: IHttpFetch;
};

const usableRelease = <E extends ITunesLookupEntityType>(
  release: ItunesLookupType<E>,
) => {
  const isCompilation =
    (release.wrapperType === "track"
      ? (release as ItunesLookupSongModel).collectionArtistName
      : release.artistName)
      ?.toLowerCase?.()
      ?.includes?.("Various Artists".toLowerCase());

  const isDjMix = release.collectionName?.toLowerCase?.()
    ?.includes?.("DJ Mix".toLowerCase()) ||
    release.collectionCensoredName?.toLowerCase?.()
      ?.includes?.("DJ Mix".toLowerCase());

  return !isCompilation && !isDjMix;
};

export class ItunesService implements IItunesService {
  #lookupUrl = "https://itunes.apple.com/lookup";
  #props: ItunesServiceProps;

  constructor(props: ItunesServiceProps) {
    this.#props = props;
  }

  static toReleaseSourcePersistance(
    item: ItunesLookupArtistModel,
  ): DbReleaseSourcesTable {
    return {
      id: crypto.randomUUID(),
      provider: ReleaseSourceProvider.ITUNES,
      raw: JSON.stringify(item),
    };
  }

  static toReleasePersistance(
    item: ItunesLookupSongModelWithExtra | ItunesLookupAlbumModelWithExtra,
  ): DbReleasesTable {
    return {
      releasedAt: item.releaseDate.toISOString(),
      provider: item.extra.provider,
      type: item.extra.type,
      raw: JSON.stringify(item),
      id: String(
        item.wrapperType === "collection"
          ? item.collectionId
          : (item as ItunesLookupSongModel).trackId,
      ),
    };
  }

  static fromPersistanceToRelease(
    row: DbReleasesTable,
  ): ItunesLookupAlbumModelWithExtra | ItunesLookupSongModelWithExtra {
    if (row.provider !== ReleaseSourceProvider.ITUNES) {
      throw new Error(`Release provider "${row.provider}" not supported`);
    }

    switch (row.type) {
      case ReleaseType.SONG: {
        return itunesLookupSongModelWithExtraSchema.parse(JSON.parse(row.raw));
      }
      case ReleaseType.ALBUM: {
        return itunesLookupAlbumModelWithExtraSchema.parse(JSON.parse(row.raw));
      }
      default: {
        throw new Error(`Relase type "${row.type}" not supported`);
      }
    }
  }

  static fromPersistanceToReleaseSurce(
    row: DbReleaseSourcesTable,
  ): ItunesLookupArtistModel {
    return itunesLookupArtistModelSchema.parse(JSON.parse(row.raw));
  }

  #getArtistImage = async (url: string) => {
    const textDecoder = new TextDecoder();
    const response = await this.#props.httpClient.fetchReadable(url);

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

  async lookupArtistById(
    id: number,
  ): Promise<ItunesLookupArtistModelWithExtra | undefined> {
    const url = new URL(this.#lookupUrl);
    url.searchParams.set("id", String(id));

    const response = await this.#props.httpClient.fetch<
      ItunesResponseModel<ItunesLookupArtistModel>
    >(url);

    const remote = response.results.at(0);

    if (!remote) return;

    const artistImage = await this.#getArtistImage(remote.artistLinkUrl);

    return itunesLookupArtistModelWithExtraSchema.parse({
      ...itunesLookupArtistModelSchema.parse(remote),
      extra: { artistImage: artistImage ?? "https://placehold.co/256" },
    });
  }

  async lookupLatestReleasesByArtist<E extends ITunesLookupEntityType>(
    artistId: string,
    entity: E,
    limit: number,
  ): Promise<Array<ItunesLookupTypeWithExtra<E>>> {
    const path = new URL(this.#lookupUrl);
    path.searchParams.set("id", artistId);
    path.searchParams.set("entity", entity);
    path.searchParams.set("media", "music");
    path.searchParams.set("sort", "recent");
    path.searchParams.set("limit", String(limit));

    const data = await this.#props.httpClient.fetch<
      ItunesResponseModel<ItunesLookupType<E>>
    >(path);

    // Remove artists info from results
    data.results.splice(0, 1);

    return data.results
      .filter((item) => usableRelease(item))
      .map((item) => {
        switch (entity) {
          case ITunesLookupEntityType.ALBUM: {
            return itunesLookupAlbumModelWithExtraSchema.parse({
              ...itunesLookupAlbumModelSchema.parse(item),
              extra: {
                type: ReleaseType.ALBUM,
                provider: ReleaseSourceProvider.ITUNES,
              },
            });
          }
          case ITunesLookupEntityType.SONG: {
            return itunesLookupSongModelWithExtraSchema.parse({
              ...itunesLookupSongModelSchema.parse(item),
              extra: {
                type: ReleaseType.SONG,
                provider: ReleaseSourceProvider.ITUNES,
              },
            });
          }
        }
      }) as Array<ItunesLookupTypeWithExtra<E>>;
  }
}
