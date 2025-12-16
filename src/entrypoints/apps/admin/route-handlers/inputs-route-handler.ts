import z from "@zod/zod";
import { inputPages } from "../pages/mod.ts";
import { pageToHtmlResponse } from "../pages/page.tsx";
import { AbstractRouteHandler } from "./route-handler.ts";
import { EInputProvider } from "#src/common/database/enums/mod.ts";
import type { DbInputsTable } from "#src/common/database/types.ts";

export class InputsRouteHandler extends AbstractRouteHandler {
  static override PATH = "/inputs";

  override GET(request: Request): Response | Promise<Response> {
    const url = new URL(request.url);
    const { provider, limit, page } = z.object({
      provider: z.enum(EInputProvider).optional(),
      page: z.string().optional().pipe(z.coerce.number()).pipe(
        z.number().min(0),
      ).default(0),
      limit: z.string().optional().pipe(z.coerce.number()).pipe(
        z.number().min(0),
      ).default(10),
    }).parse(Object.fromEntries(url.searchParams.entries()));

    const sources = this.props.database.sql<DbInputsTable>`
      select *, json(raw) as raw from inputs
      where ${provider ? 1 : null} is null or provider = ${provider ?? null}
      limit ${limit}
      offset ${page * limit}
    `;

    return pageToHtmlResponse(
      inputPages.indexPage({
        sources,
        providers: this.props.providers,
        url,
        paginatio: { page, limit },
      }),
    );
  }
}
