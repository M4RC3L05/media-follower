import z from "@zod/zod";
import { EInputProvider } from "#src/common/database/enums/input-provider.ts";
import { pageToHtmlResponse } from "#src/entrypoints/apps/admin/pages/page.tsx";
import { indexPage } from "#src/entrypoints/apps/admin/pages/index.tsx";
import type { IDatabase } from "#src/common/database/database.ts";
import type { IProvider } from "#src/common/providers/interfaces.ts";
import type {
  DbInputsTable,
  DbOutputsTable,
} from "#src/common/database/types.ts";
import {
  inputPages,
  outputPages,
} from "#src/entrypoints/apps/admin/pages/mod.ts";

export type AppProps = {
  providers: Record<EInputProvider, IProvider>;
  database: IDatabase;
};

const getOutputs = (
  { limit, page, provider, database }: {
    provider: EInputProvider | undefined;
    limit: number;
    page: number;
    database: IDatabase;
  },
) => {
  return database.sql<DbOutputsTable>`
    select id, input_id, provider, json(outputs.raw) as raw
    from outputs
    where (${provider ? 1 : null} is null or provider = ${provider ?? null})
    order by (
      case
        when provider = ${EInputProvider.BLU_RAY_COM_PHYSICAL_RELEASE}
          then outputs.raw->>'releasedate'
        when provider = ${EInputProvider.ITUNES_MUSIC_RELEASE}
          then outputs.raw->>'releaseDate'
        else "rowid"
      end
    ) desc, "rowid" desc
    limit ${limit}
    offset ${page * limit}
  `;
};

export class App {
  #props: AppProps;

  constructor(props: AppProps) {
    this.#props = props;
  }

  fetch = async (request: Request) => {
    const url = new URL(request.url);

    if (request.method === "GET" && url.pathname === "/inputs") {
      const { provider, limit, page } = z.object({
        provider: z.enum(EInputProvider).optional(),
        page: z.string().optional().pipe(z.coerce.number()).pipe(
          z.number().min(0),
        ).default(0),
        limit: z.string().optional().pipe(z.coerce.number()).pipe(
          z.number().min(0),
        ).default(10),
      }).parse(Object.fromEntries(url.searchParams.entries()));

      const sources = this.#props.database.sql<DbInputsTable>`
        select *, json(raw) as raw from inputs
        where ${provider ? 1 : null} is null or provider = ${provider ?? null}
        limit ${limit}
        offset ${page * limit}
      `;

      return pageToHtmlResponse(
        inputPages.indexPage({
          sources,
          providers: this.#props.providers,
          url,
          paginatio: { page, limit },
        }),
      );
    }

    if (url.pathname === "/inputs/create") {
      if (request.method === "GET") {
        return pageToHtmlResponse(inputPages.createPage());
      }

      if (request.method === "POST") {
        const formData = await request.formData();
        const data = z.object({
          provider: z.enum(EInputProvider),
          term: z.string().min(1),
        }).parse(Object.fromEntries(formData.entries()));

        const provider = this.#props.providers[data.provider];
        const input = await provider.lookupInput(data.term);

        if (!input) return Response.redirect(new URL("/inputs", url));

        const dbObj = provider.fromInputToPersistence(input);

        this.#props.database.sql`
          insert or replace into inputs
            (id, provider, raw)
          values
            (${dbObj.id}, ${dbObj.provider}, jsonb(${dbObj.raw}))
        `;

        return Response.redirect(new URL("/inputs", url));
      }
    }

    if (request.method === "GET" && url.pathname === "/outputs") {
      const { provider, page, limit } = z.object({
        provider: z.enum(EInputProvider).optional(),
        page: z.string().optional().pipe(z.coerce.number()).pipe(
          z.number().min(0),
        ).default(0),
        limit: z.string().optional().pipe(z.coerce.number()).pipe(
          z.number().min(0),
        ).default(10),
      }).parse(Object.fromEntries(url.searchParams.entries()));

      const outputs = provider
        ? await this.#props.providers[provider].queryOutputs({
          pagination: { limit, page },
          queries: Object.fromEntries(url.searchParams.entries()),
        })
        : getOutputs({ limit, database: this.#props.database, page, provider });

      return pageToHtmlResponse(
        outputPages.indexPage({
          outputs,
          providers: this.#props.providers,
          url,
          paginatio: { page, limit },
        }),
      );
    }

    return pageToHtmlResponse(indexPage());
  };
}
