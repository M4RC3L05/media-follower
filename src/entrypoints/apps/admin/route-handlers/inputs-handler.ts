import z from "@zod/zod";
import { inputPages } from "../pages/mod.ts";
import { pageToHtmlResponse } from "../pages/page.tsx";
import { EInputProvider } from "#src/common/database/enums/mod.ts";
import type { RequestContext } from "@remix-run/fetch-router";
import type { DbInputsTable } from "#src/common/database/types.ts";
import type { IDatabase } from "#src/common/database/database.ts";
import type { IProvider } from "#src/common/providers/interfaces.ts";
import { createRedirectResponse } from "@remix-run/response/redirect";
import { routes } from "../router.ts";

type InputsIndexProps = {
  database: IDatabase;
  providers: Record<EInputProvider, IProvider>;
};

export const inputsIndex =
  (props: InputsIndexProps) => ({ url }: RequestContext) => {
    const { provider, limit, page } = z.object({
      provider: z.enum(EInputProvider).optional(),
      page: z.string().optional().pipe(z.coerce.number()).pipe(
        z.number().min(0),
      ).default(0),
      limit: z.string().optional().pipe(z.coerce.number()).pipe(
        z.number().min(0),
      ).default(10),
    }).parse(Object.fromEntries(url.searchParams.entries()));

    const sources = props.database.sql<DbInputsTable>`
      select *, json(raw) as raw from inputs
      where ${provider ? 1 : null} is null or provider = ${provider ?? null}
      limit ${limit}
      offset ${page * limit}
    `;

    return pageToHtmlResponse(
      inputPages.indexPage({
        sources,
        providers: props.providers,
        url,
        paginatio: { page, limit },
      }),
    );
  };

export const inputsCreateGet = () =>
  pageToHtmlResponse(inputPages.createPage());

type InputsCreatePostProps = {
  database: IDatabase;
  providers: Record<EInputProvider, IProvider>;
};

export const inputsCreatePost = (props: InputsCreatePostProps) =>
async (
  { formData }: RequestContext,
) => {
  const data = z.object({
    provider: z.enum(EInputProvider),
    term: z.string().min(1),
  }).parse(Object.fromEntries(formData!.entries()));

  const provider = props.providers[data.provider];
  const input = await provider.lookupInput(data.term);

  if (!input) return createRedirectResponse(routes.inputs.create.index.href());

  const dbObj = provider.fromInputToPersistence(input);

  props.database.sql`
    insert or replace into inputs
      (id, provider, raw)
    values
      (${dbObj.id}, ${dbObj.provider}, jsonb(${dbObj.raw}))
  `;

  return createRedirectResponse(routes.inputs.index.href());
};
