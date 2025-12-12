import z from "@zod/zod";
import { Feed, type Item } from "feed";
import type { IServerApp } from "#src/common/server/mod.ts";
import type { IDatabase } from "#src/common/database/database.ts";
import { ReleaseType } from "#src/common/database/enums/release-type.ts";
import type { DbReleasesTable } from "#src/common/database/types.ts";
import { ReleaseSourceProvider } from "#src/common/database/enums/release-source-provider.ts";
import {
  bluRayComReleaseWithExtraSchema,
  itunesLookupAlbumModelWithExtraSchema,
  itunesLookupSongModelWithExtraSchema,
} from "#src/common/services/service.ts";

type AppProps = {
  database: IDatabase;
};

const releaseToFeedItem = (release: DbReleasesTable): Item => {
  switch (release.provider) {
    case ReleaseSourceProvider.BLU_RAY_COM: {
      switch (release.type) {
        case ReleaseType.DVD:
        case ReleaseType.BLURAY: {
          const parsed = bluRayComReleaseWithExtraSchema.parse(
            JSON.parse(release.raw),
          );
          return {
            date: parsed.releasedate,
            link: parsed.extra.link,
            title: parsed.title,
            id: `${release.provider}@${release.type}@${parsed.id}`,
            image: parsed.extra.artworkUrl,
          };
        }
        case ReleaseType.SONG:
        case ReleaseType.ALBUM: {
          throw new Error("Not valid");
        }
      }

      // @ts-ignore: Fix typechecking
      throw new Error("Not valid");
    }
    case ReleaseSourceProvider.ITUNES: {
      switch (release.type) {
        case ReleaseType.SONG: {
          const parsed = itunesLookupSongModelWithExtraSchema.parse(
            JSON.parse(release.raw),
          );

          return {
            date: parsed.releaseDate,
            link: parsed.trackViewUrl,
            title: parsed.trackName,
            id: `${release.provider}@${release.type}@${parsed.trackId}`,
            image: parsed.artworkUrl100
              .split("/")
              .map((segment, index, array) =>
                index === array.length - 1 ? "512x512bb.jpg" : segment
              )
              .join("/"),
          };
        }
        case ReleaseType.ALBUM: {
          const parsed = itunesLookupAlbumModelWithExtraSchema.parse(
            JSON.parse(release.raw),
          );

          return {
            date: parsed.releaseDate,
            link: parsed.collectionViewUrl,
            title: parsed.collectionName,
            id: `${release.provider}@${release.type}@${parsed.collectionId}`,
            image: parsed.artworkUrl100
              .split("/")
              .map((segment, index, array) =>
                index === array.length - 1 ? "512x512bb.jpg" : segment
              )
              .join("/"),
          };
        }
        case ReleaseType.DVD:
        case ReleaseType.BLURAY: {
          throw new Error("Not valid");
        }
      }
    }
  }
};

export class App implements IServerApp {
  #props: AppProps;

  constructor(props: AppProps) {
    this.#props = props;
  }

  fetch = (request: Request) => {
    const parsedUrl = URL.parse(request.url)!;
    const { success, data, error } = z.object({
      type: z.string().transform((val) => val.split(",")).pipe(
        z.array(z.enum(ReleaseType)),
      ).optional(),
      provider: z.string().transform((val) => val.split(",")).pipe(
        z.array(z.enum(ReleaseSourceProvider)),
      ).optional(),
      country: z.string().optional(),
    })
      .safeParse(
        Object.fromEntries(parsedUrl.searchParams.entries()),
      );

    if (!success) {
      return new Response(
        `Validation error:\n${JSON.stringify(z.treeifyError(error), null, 2)}`
          .trim(),
        {
          status: 400,
          headers: {
            "content-type": "text/plain",
          },
        },
      );
    }

    const releases = this.#props.database.sql<DbReleasesTable>`
      select *, json(raw) as raw
      from releases
      where
          (${(data.type?.length ?? 0) > 0 ? 1 : null} is null or type in (${
      data.type ?? []
    }))
        and (${
      (data.provider?.length ?? 0) > 0 ? 1 : null
    } is null or provider in (${data.provider ?? []}))
        and "releasedAt" <= strftime('%Y-%m-%dT%H:%M:%fZ' , 'now')
        and (${
      data.country ? 1 : null
    } is null or raw->'extra'->>'country' is null or raw->'extra'->>'country' = ${
      data.country ?? null
    })
      order by "releasedAt" desc, "rowid" desc
      limit 200;
    `;

    const prefix = (data.provider || data.type || data.country)
      ? [data.provider, data.type, data.country].filter(Boolean)
      : "";

    const feed = new Feed({
      title: `${prefix ? `[${prefix.join(" | ")}] ` : ""}Media follower`,
      description: `Get the latest${
        prefix ? ` ${prefix.join(" and ")} ` : " "
      }media releases`,
      id: `media_follower${prefix ? `_${prefix.join("_")}` : ""}`,
      copyright: "Media Follower",
      updated: new Date(),
    });

    for (const item of releases) {
      feed.addItem(releaseToFeedItem(item));
    }

    const accepts = request.headers.get("accept");

    if (
      accepts?.includes("application/rss+xml") ??
        accepts?.includes("application/xml")
    ) {
      return new Response(feed.rss2(), {
        status: 200,
        headers: {
          "content-type": accepts?.includes("application/rss+xml")
            ? "application/rss+xml"
            : "application/xml",
        },
      });
    }

    if (accepts?.includes("application/atom+xml")) {
      return new Response(feed.atom1(), {
        status: 200,
        headers: {
          "content-type": "application/atom+xml",
        },
      });
    }

    if (accepts?.includes("application/json")) {
      return new Response(feed.json1(), {
        status: 200,
        headers: { "content-type": "application/json" },
      });
    }

    return new Response(feed.rss2(), {
      status: 200,
      headers: { "content-type": "application/xml" },
    });
  };
}
