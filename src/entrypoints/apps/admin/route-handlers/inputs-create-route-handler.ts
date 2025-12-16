import z from "@zod/zod";
import { inputPages } from "../pages/mod.ts";
import { pageToHtmlResponse } from "../pages/page.tsx";
import { AbstractRouteHandler } from "./route-handler.ts";
import { EInputProvider } from "#src/common/database/enums/mod.ts";

export class InputsCreateRouteHandler extends AbstractRouteHandler {
  static override PATH = "/inputs/create";

  override GET(_request: Request): Response | Promise<Response> {
    return pageToHtmlResponse(inputPages.createPage());
  }

  override async POST(request: Request): Promise<Response> {
    const url = new URL(request.url);
    const formData = await request.formData();
    const data = z.object({
      provider: z.enum(EInputProvider),
      term: z.string().min(1),
    }).parse(Object.fromEntries(formData.entries()));

    const provider = this.props.providers[data.provider];
    const input = await provider.lookupInput(data.term);

    if (!input) return Response.redirect(new URL("/inputs", url));

    const dbObj = provider.fromInputToPersistence(input);

    this.props.database.sql`
      insert or replace into inputs
        (id, provider, raw)
      values
        (${dbObj.id}, ${dbObj.provider}, jsonb(${dbObj.raw}))
    `;

    return Response.redirect(new URL("/inputs", url));
  }
}
