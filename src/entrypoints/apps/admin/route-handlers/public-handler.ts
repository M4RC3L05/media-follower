import type { RequestContext } from "@remix-run/fetch-router";
import css from "simpledotcss/simple.min.css" with { type: "text" };
import mainCss from "../static/main.css" with { type: "text" };
import { pageToHtmlResponse } from "../pages/page.tsx";
import { notFoundPage } from "../pages/notFound.tsx";

const bundledCss = `${css}${mainCss}`;

export const publicIndex = ({ url }: RequestContext) => {
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
};
