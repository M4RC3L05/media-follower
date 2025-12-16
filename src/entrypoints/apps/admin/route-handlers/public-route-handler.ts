import { AbstractRouteHandler } from "./route-handler.ts";
import css from "simpledotcss/simple.min.css" with { type: "text" };
import mainCss from "../static/main.css" with { type: "text" };
import { pageToHtmlResponse } from "../pages/page.tsx";
import { notFoundPage } from "../pages/notFound.tsx";

const bundledCss = `${css}${mainCss}`;

export class PublicRouteHandler extends AbstractRouteHandler {
  static override PATH = "/public";

  override GET(request: Request): Response | Promise<Response> {
    const url = new URL(request.url);

    switch (url.searchParams.get("asset")) {
      case "main.css": {
        return new Response(bundledCss, {
          status: 200,
          headers: { "content-type": "text/css" },
        });
      }
      default: {
        return pageToHtmlResponse(notFoundPage(), 404);
      }
    }
  }
}
