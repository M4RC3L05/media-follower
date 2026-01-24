import { EInputProvider } from "#src/common/database/enums/mod.ts";
import type { IProvider } from "#src/common/providers/interfaces.ts";
import type { IDatabase } from "#src/common/database/database.ts";
import type { DbOutputsTable } from "#src/common/database/types.ts";
import z from "@zod/zod";
import { pageToHtmlResponse } from "../pages/page.tsx";
import { outputPages } from "../pages/mod.ts";

type OutputsIndexProps = {
  url: URL;
  database: IDatabase;
  providers: Record<EInputProvider, IProvider>;
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

export const outputsIndex = async (props: OutputsIndexProps) => {
  const { provider, page, limit } = z.object({
    provider: z.enum(EInputProvider).optional(),
    page: z.string().optional().pipe(z.coerce.number()).pipe(
      z.number().min(0),
    ).default(0),
    limit: z.string().optional().pipe(z.coerce.number()).pipe(
      z.number().min(0),
    ).default(10),
  }).parse(Object.fromEntries(props.url.searchParams.entries()));

  const outputs = provider
    ? await props.providers[provider].queryOutputs({
      pagination: { limit, page },
      queries: Object.fromEntries(props.url.searchParams.entries()),
    })
    : getOutputs({ limit, database: props.database, page, provider });

  return pageToHtmlResponse(
    outputPages.indexPage({
      outputs,
      providers: props.providers,
      url: props.url,
      paginatio: { page, limit },
    }),
  );
};
