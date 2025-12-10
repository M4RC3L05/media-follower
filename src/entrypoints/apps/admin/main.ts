import z from "@zod/zod";
import { config, initConfig } from "../../../common/config/mod.ts";
import { HttpFetch } from "../../../common/http/mod.ts";
import { ItunesService } from "../../../common/services/itunes-service.ts";
import {
  ReleaseSourceProvider,
} from "../../../common/database/enums/release-source-provider.ts";
import { gracefulShutdown } from "../../../common/process/mod.ts";
import { template } from "./template.ts";
import { BluRayComService } from "../../../common/services/blu-ray-com-service.ts";
import type {
  DbReleaseSourcesTable,
  DbReleasesTable,
} from "../../../common/database/types.ts";
import { makeDatabase } from "../../../common/database/mod.ts";
import {
  type BluRayComCountry,
  bluRayComReleaseWithExtraSchema,
  itunesLookupAlbumModelWithExtraSchema,
  type ItunesLookupArtistModelWithExtra,
  itunesLookupSongModelWithExtraSchema,
} from "../../../common/services/service.ts";
import { ReleaseType } from "../../../common/database/enums/release-type.ts";
import { Server } from "#src/common/server/mod.ts";

initConfig();

const { promise: shutdownPromise, signal: shutdownSignal } = gracefulShutdown();

using db = makeDatabase();
const itunesService = new ItunesService({
  httpClient: new HttpFetch({ signal: shutdownSignal }),
});
const blurayComService = new BluRayComService({
  httpClient: new HttpFetch({ signal: shutdownSignal }),
});

const genSubmits = () => {
  return Object.entries(ReleaseSourceProvider).map(([_, value]) => {
    switch (value) {
      case ReleaseSourceProvider.BLU_RAY_COM: {
        return `
          <h3>Add "${value}" source</h3>
          <form method="post" action="/sources">
            <input type="hidden" value="${value}" id="provider" name="provider" />
            <label for="country">Country: </label>
            <input id="country" name="country" />
            <button type="submit">Add</button>
          </form>
        `;
      }
      case ReleaseSourceProvider.ITUNES: {
        return `
          <h3>Add "${value}" source</h3>
          <form method="post" action="/sources">
            <input type="hidden" value="${value}" id="provider" name="provider" />
            <label for="artistId">Artist ID: </label>
            <input id="artistId" name="artistId" />
            <button type="submit">Add</button>
          </form>
        `;
      }
    }
  });
};

const genReleaseDisplay = (
  release: DbReleasesTable & { provider: ReleaseSourceProvider },
) => {
  switch (release.provider) {
    case ReleaseSourceProvider.BLU_RAY_COM: {
      const parsed = bluRayComReleaseWithExtraSchema.parse(
        JSON.parse(release.raw),
      );

      return `
        <div>
            <img style="max-width: 100%; height: auto; aspect-ratio: 1/1; max-height: 256px" src="${parsed.extra.artworkUrl}" />
            <h3>${parsed.title} | ${parsed.extra.type}</h3>
            <p>${new Date(parsed.releasedate).toDateString()}${
        new Date(parsed.releasedate) > new Date()
          ? " <em>(To be released)</em>"
          : ""
      }</p>
            <a target="_bank" href="${parsed.extra.link}">View on source</a>
            <br>
            <br>
            <details>
              <summary>Raw:</summary>
              <pre>${JSON.stringify(parsed, null, 2)}</pre>
            </details>
          </div>
      `;
    }
    case ReleaseSourceProvider.ITUNES: {
      switch (release.type) {
        case ReleaseType.SONG: {
          const parsed = itunesLookupSongModelWithExtraSchema.parse(
            JSON.parse(release.raw),
          );
          const image = parsed.artworkUrl100
            .split("/")
            .map((segment, index, array) =>
              index === array.length - 1 ? "512x512bb.jpg" : segment
            ).join("/");

          return `
        <div>
            <img style="max-width: 100%; height: auto; aspect-ratio: 1/1; max-height: 256px" src="${image}" />
            <h3>${parsed.trackName} | ${parsed.kind}</h3>
            <p>${new Date(parsed.releaseDate).toDateString()}${
            new Date(parsed.releaseDate) > new Date()
              ? " <em>(To be released)</em>"
              : ""
          }</p>
            <a target="_bank" href="${parsed.trackViewUrl}">View on source</a>
            <br>
            <br>
            <details>
              <summary>Raw:</summary>
              <pre>${JSON.stringify(parsed, null, 2)}</pre>
            </details>
          </div>
      `;
        }
        case ReleaseType.ALBUM: {
          const parsed = itunesLookupAlbumModelWithExtraSchema.parse(
            JSON.parse(release.raw),
          );
          const image = parsed.artworkUrl100
            .split("/")
            .map((segment, index, array) =>
              index === array.length - 1 ? "512x512bb.jpg" : segment
            ).join("/");

          return `
        <div>
            <img style="max-width: 100%; height: auto; aspect-ratio: 1/1; max-height: 256px" src="${image}" />
            <h3>${parsed.collectionName} | ${parsed.collectionType}</h3>
            <p>${new Date(parsed.releaseDate).toDateString()}${
            new Date(parsed.releaseDate) > new Date()
              ? " <em>(To be released)</em>"
              : ""
          }</p>
            <a target="_bank" href="${parsed.collectionViewUrl}">View on source</a>
            <br>
            <br>
            <details>
              <summary>Raw:</summary>
              <pre>${JSON.stringify(parsed, null, 2)}</pre>
            </details>
          </div>
      `;
        }
        case ReleaseType.DVD:
        case ReleaseType.BLURAY: {
          throw new Error("Not valid");
        }
      }
    }
  }
};

const genRelaseSourceDisplay = (
  source: DbReleaseSourcesTable,
) => {
  switch (source.provider) {
    case ReleaseSourceProvider.BLU_RAY_COM: {
      const parsed = JSON.parse(source.raw) as BluRayComCountry;

      return `
          <div>
            <h3>${parsed.name} | ${source.provider}</h3>
            <br>
            <br>
            <details>
              <summary>Raw:</summary>
              <pre>${JSON.stringify(parsed, null, 2)}</pre>
            </details>
          </div>
        `;
    }
    case ReleaseSourceProvider.ITUNES: {
      const parsed = JSON.parse(source.raw) as ItunesLookupArtistModelWithExtra;

      return `
          <div>
            <img style="max-width: 100%; height: auto; aspect-ratio: 1/1; max-height: 256px" src="${parsed.extra.artistImage}" />
            <h3>${parsed.artistName} | ${source.provider}</h3>
            <br>
            <br>
            <details>
              <summary>Raw:</summary>
              <pre>${JSON.stringify(parsed, null, 2)}</pre>
            </details>
          </div>
        `;
    }
  }
};

await using _server = new Server({
  hostname: config().apps.admin.host,
  port: config().apps.admin.port,
  app: {
    fetch: async (request) => {
      const url = new URL(request.url);

      if (request.method === "POST" && url.pathname === "/sources") {
        const formData = await request.formData();
        const data = z.discriminatedUnion("provider", [
          z.object({
            provider: z.literal(ReleaseSourceProvider.BLU_RAY_COM),
            country: z.string().min(1),
          }),
          z.object({
            provider: z.literal(ReleaseSourceProvider.ITUNES),
            artistId: z.string().min(1).pipe(z.coerce.number()),
          }),
        ]).parse(Object.fromEntries(formData.entries()));

        switch (data.provider) {
          case ReleaseSourceProvider.BLU_RAY_COM: {
            const remote = await blurayComService.getCountries();
            const selected = remote.find((item) => item.code === data.country);

            if (!selected) return Response.redirect(url);

            const mapped = BluRayComService.toReleaseSourcePersistance(
              selected,
            );

            db.sql<DbReleaseSourcesTable>`
            insert or replace into release_sources
              (id, provider, raw)
            values
              (${mapped.id}, ${mapped.provider}, jsonb(${mapped.raw}))
          `;

            return Response.redirect(url);
          }
          case ReleaseSourceProvider.ITUNES: {
            const remote = await itunesService.lookupArtistById(data.artistId);

            if (!remote) return Response.redirect(url);

            const mapped = ItunesService.toReleaseSourcePersistance(remote);

            db.sql<DbReleaseSourcesTable>`
            insert or replace into release_sources
              (id, provider, raw)
            values
              (${mapped.id}, ${mapped.provider}, jsonb(${mapped.raw}))
          `;

            return Response.redirect(url);
          }
        }
      }

      if (request.method === "GET" && url.pathname === "/sources") {
        const { provider } = z.object({
          provider: z.enum(ReleaseSourceProvider).optional(),
        }).parse(Object.fromEntries(url.searchParams.entries()));

        const sources = db.sql<DbReleaseSourcesTable>`
        select *, json(raw) as raw from release_sources
        where ${provider ? 1 : null} is null or provider = ${provider ?? null}
      `;

        const body = `
        <a href="/sources">All</a> | ${
          Object.values(ReleaseSourceProvider).map((value) =>
            `<a href="?provider=${value}">${value}</a>`
          ).join(" | ")
        }
      <hr>
      ${genSubmits().join("")}
      <hr>
      <h1>Sources</h1>
      ${sources.map((item) => genRelaseSourceDisplay(item)).join("<hr>")}
      `;

        return new Response(template.replace("{{ body }}", body), {
          status: 200,
          headers: { "content-type": "text/html" },
        });
      }

      if (request.method === "GET" && url.pathname === "/releases") {
        const { provider, type, page, limit } = z.object({
          provider: z.enum(ReleaseSourceProvider).optional(),
          type: z.enum(ReleaseType).optional(),
          page: z.string().optional().pipe(z.coerce.number()).pipe(
            z.number().min(0),
          ).default(0),
          limit: z.string().optional().pipe(z.coerce.number()).pipe(
            z.number().min(0),
          ).default(10),
        }).parse(Object.fromEntries(url.searchParams.entries()));

        const releases = db.sql<DbReleasesTable>`
        select *, json(releases.raw) as raw
        from releases
        where
            (${type ? 1 : null} is null or releases.type = ${type ?? null})
        and (${provider ? 1 : null} is null or releases.provider = ${
          provider ?? null
        })
        order by (
          case 
            when provider = ${ReleaseSourceProvider.ITUNES}
              then releases.raw->>'releasedAt'
            when provider = ${ReleaseSourceProvider.BLU_RAY_COM}
              then releases.raw->>'releasedate'
            else null
          end
        ) desc
        limit ${limit}
        offset ${page * limit}
      `;

        const prevPageLink = new URL(url);
        prevPageLink.searchParams.set("page", String(Math.max(page - 1, 0)));

        const nextPageLink = new URL(url);
        nextPageLink.searchParams.set("page", String(page + 1));

        const body = `
        <h4>Filter by provider</h4>
        <a href="/releases">All</a> | ${
          Object.values(ReleaseSourceProvider).map((value) =>
            `<a href="?provider=${value}">${value}</a>`
          ).join(" | ")
        }
        <hr>
        <h4>Filter by type</h4>
        <a href="/releases">All</a> | ${
          Object.values(ReleaseType).map((value) =>
            `<a href="?type=${value}">${value}</a>`
          ).join(" | ")
        }
        <hr>
        <h1>Releases</h1>
        <a href="${prevPageLink.toString()}">Prev</a> | <a href="${nextPageLink.toString()}">Next</a>
        <hr>
        ${releases.map((item) => genReleaseDisplay(item)).join("<hr>")}
      `.trim();

        return new Response(template.replace("{{ body }}", body), {
          status: 200,
          headers: { "content-type": "text/html" },
        });
      }

      return new Response("", {
        status: 200,
        headers: { "content-type": "text/plain" },
      });
    },
  },
});

await shutdownPromise;
