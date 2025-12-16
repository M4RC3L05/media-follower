import { config, initConfig } from "#src/common/config/mod.ts";
import { HttpFetch } from "#src/common/http/mod.ts";
import { gracefulShutdown } from "#src/common/process/mod.ts";
import { pageToHtmlResponse } from "#src/entrypoints/apps/admin/pages/page.tsx";
import { CustomDatabase } from "#src/common/database/mod.ts";
import { Server } from "#src/common/server/mod.ts";
import { errorPage } from "#src/entrypoints/apps/admin/pages/error.tsx";
import { makeLogger } from "#src/common/logger/mod.ts";
import { EInputProvider } from "#src/common/database/enums/mod.ts";
import { providerFactory } from "#src/common/providers/provider.ts";
import { App } from "./app.ts";
import { InputsRouteHandler } from "./route-handlers/inputs-route-handler.ts";
import { OutputsRouteHandler } from "./route-handlers/outputs-route-handler.ts";
import { InputsCreateRouteHandler } from "./route-handlers/inputs-create-route-handler.ts";
import { IndexRouteHandler } from "./route-handlers/index-route-handler.ts";
import { HttpError } from "#src/common/errors/mod.ts";
import { PublicRouteHandler } from "./route-handlers/public-route-handler.ts";

initConfig();

const { promise: shutdownPromise, signal: shutdownSignal } = gracefulShutdown();

const log = makeLogger("admin-app");

using database = new CustomDatabase(config().database.path);
const httpClient = new HttpFetch({ signal: shutdownSignal });

const routeHandlerProps = {
  providers: {
    [EInputProvider.BLU_RAY_COM_PHYSICAL_RELEASE]: providerFactory(
      EInputProvider.BLU_RAY_COM_PHYSICAL_RELEASE,
      { database, httpClient },
    ),
    [EInputProvider.ITUNES_MUSIC_RELEASE]: providerFactory(
      EInputProvider.ITUNES_MUSIC_RELEASE,
      { database, httpClient },
    ),
    [EInputProvider.STEAM_GAMES_FREE_PROMOS]: providerFactory(
      EInputProvider.STEAM_GAMES_FREE_PROMOS,
      { database, httpClient },
    ),
  },
  database,
};

await using _server = new Server({
  hostname: config().apps.admin.host,
  port: config().apps.admin.port,
  onListen: (host, port) => {
    log.info(`Serving on http://${host}:${port}`);
  },
  onError: (error) => {
    log.error({ error }, "Something went wrong");

    return pageToHtmlResponse(
      errorPage({
        message: error instanceof HttpError ? error.message : undefined,
      }),
      error instanceof HttpError ? error.status : 500,
    );
  },
  app: new App({
    routeHandles: {
      ...Object.fromEntries(
        [
          PublicRouteHandler,
          IndexRouteHandler,
          InputsRouteHandler,
          InputsCreateRouteHandler,
          OutputsRouteHandler,
        ].map((item) => [item.PATH, new item(routeHandlerProps)]),
      ),
    },
  }),
});

await shutdownPromise;
