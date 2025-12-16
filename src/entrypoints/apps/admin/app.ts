import { notFoundPage } from "./pages/notFound.tsx";
import type {
  IRouteHandler,
  SupportedHttpMethods,
} from "./route-handlers/route-handler.ts";
import { pageToHtmlResponse } from "./pages/page.tsx";

export type AppProps = {
  routeHandles: Record<string, IRouteHandler>;
};

export class App {
  #props: AppProps;
  #notfoundPage = notFoundPage();

  constructor(props: AppProps) {
    this.#props = props;
  }

  fetch = async (request: Request) => {
    const url = new URL(request.url);

    const handler = this.#props.routeHandles[url.pathname]
      ?.[request.method as SupportedHttpMethods];

    if (!handler) return pageToHtmlResponse(this.#notfoundPage, 404);

    return await handler.call(this.#props.routeHandles[url.pathname], request);
  };
}
