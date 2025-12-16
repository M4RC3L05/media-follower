import { rootPages } from "../pages/mod.ts";
import { pageToHtmlResponse } from "../pages/page.tsx";
import { AbstractRouteHandler } from "./route-handler.ts";

export class IndexRouteHandler extends AbstractRouteHandler {
  static override PATH = "/";

  override GET(_request: Request): Response | Promise<Response> {
    return pageToHtmlResponse(rootPages.indexPage());
  }
}
