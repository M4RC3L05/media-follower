import z from "@zod/zod";
import type { IServerApp } from "#src/common/server/mod.ts";
import type { IDatabase } from "#src/common/database/database.ts";
import { EInputProvider } from "#src/common/database/enums/input-provider.ts";
import type { IProviderFeed } from "#src/common/providers/interfaces.ts";
import favicon from "#src/common/assets/favicon.ico" with { type: "bytes" };

type AppProps = {
  database: IDatabase;
  providers: Record<EInputProvider, IProviderFeed>;
};

export class App implements IServerApp {
  #props: AppProps;

  constructor(props: AppProps) {
    this.#props = props;
  }

  fetch = (request: Request) => {
    const parsedUrl = URL.parse(request.url)!;
    const queries = Object.fromEntries(parsedUrl.searchParams.entries());

    if (parsedUrl.pathname === "/favicon.ico") {
      return new Response(favicon, {
        headers: { "content-type": "image/x-icon" },
      });
    }

    const { provider } = z.object({ provider: z.enum(EInputProvider) }).parse(
      queries,
    );

    const feed = this.#props.providers[provider].getOutputsFeed({ queries });

    const accepts = request.headers.get("accept");

    if (
      accepts?.includes("application/rss+xml") ??
        accepts?.includes("application/xml")
    ) {
      return new Response(feed.rss2(), {
        status: 200,
        headers: {
          "content-type": accepts?.includes("application/rss+xml")
            ? "application/rss+xml"
            : "application/xml",
        },
      });
    }

    if (accepts?.includes("application/atom+xml")) {
      return new Response(feed.atom1(), {
        status: 200,
        headers: {
          "content-type": "application/atom+xml",
        },
      });
    }

    if (accepts?.includes("application/json")) {
      return new Response(feed.json1(), {
        status: 200,
        headers: { "content-type": "application/json" },
      });
    }

    return new Response(feed.rss2(), {
      status: 200,
      headers: { "content-type": "application/xml" },
    });
  };
}
