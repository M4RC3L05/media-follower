import z from "@zod/zod";
import {
  ReleaseSourceProvider,
} from "../../../common/database/enums/release-source-provider.ts";
import { pageToHtmlResponse } from "./pages/page.tsx";
import type {
  DbReleaseSourcesTable,
  DbReleasesTable,
} from "../../../common/database/types.ts";
import type {
  IBlurayComService,
  IItunesService,
} from "../../../common/services/service.ts";
import { ReleaseType } from "../../../common/database/enums/release-type.ts";
import { indexPage } from "./pages/index.tsx";
import type { IDatabase } from "../../../common/database/database.ts";
import { sourcesCreatePage, sourcesIndexPage } from "./pages/sources/mod.ts";
import { releasesIndexPage } from "./pages/releases/mod.ts";
import * as bluRayComMappers from "#src/common/mappers/blu-ray-com-mappers.ts";
import * as itunesMappers from "#src/common/mappers/itunes-mappers.ts";

export type AppProps = {
  blurayComService: IBlurayComService;
  itunesService: IItunesService;
  database: IDatabase;
};

export class App {
  #props: AppProps;

  constructor(props: AppProps) {
    this.#props = props;
  }

  fetch = async (request: Request) => {
    const url = new URL(request.url);

    if (request.method === "GET" && url.pathname === "/sources") {
      const { provider } = z.object({
        provider: z.enum(ReleaseSourceProvider).optional(),
      }).parse(Object.fromEntries(url.searchParams.entries()));

      const sources = this.#props.database.sql<DbReleaseSourcesTable>`
        select *, json(raw) as raw from release_sources
        where ${provider ? 1 : null} is null or provider = ${provider ?? null}
      `;

      return pageToHtmlResponse(sourcesIndexPage({ sources }));
    }

    if (url.pathname === "/sources/create") {
      if (request.method === "GET") {
        return pageToHtmlResponse(sourcesCreatePage());
      }

      if (request.method === "POST") {
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
            const remote = await this.#props.blurayComService.getCountries();
            const selected = remote.find((item) => item.code === data.country);

            if (!selected) {
              return Response.redirect(new URL("/sources/create", url));
            }

            const mapped = bluRayComMappers.fromReleaseSourceToPersistance(
              selected,
            );

            this.#props.database.sql<DbReleaseSourcesTable>`
              insert or replace into release_sources
                (id, provider, raw)
              values
                (${mapped.id}, ${mapped.provider}, jsonb(${mapped.raw}))
            `;

            return Response.redirect(new URL("/sources", url));
          }
          case ReleaseSourceProvider.ITUNES: {
            const remote = await this.#props.itunesService.lookupArtistById(
              data.artistId,
            );

            if (!remote) {
              return Response.redirect(new URL("/sources/create", url));
            }

            const mapped = itunesMappers.fromReleaseSourceToPersistance(
              remote,
            );

            this.#props.database.sql<DbReleaseSourcesTable>`
              insert or replace into release_sources
                (id, provider, raw)
              values
                (${mapped.id}, ${mapped.provider}, jsonb(${mapped.raw}))
            `;

            return Response.redirect(new URL("/sources", url));
          }
        }
      }
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

      const releases = this.#props.database.sql<DbReleasesTable>`
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

      return pageToHtmlResponse(
        releasesIndexPage({
          releases,
          url,
          paginatio: { page, limit },
        }),
      );
    }

    return pageToHtmlResponse(indexPage());
  };
}
