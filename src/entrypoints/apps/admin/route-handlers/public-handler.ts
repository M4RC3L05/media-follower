import css from "simpledotcss/simple.min.css" with { type: "text" };
import mainCss from "../assets/main.css" with { type: "text" };
import { pageToHtmlResponse } from "../pages/page.tsx";
import { notFoundPage } from "../pages/notFound.tsx";

const bundledCss = `${css}${mainCss}`;

type PublicIndexProps = {
  url: URL;
};

export const publicIndex = ({ url }: PublicIndexProps) => {
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
