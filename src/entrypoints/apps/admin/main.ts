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

initConfig();

const { promise: shutdownPromise, signal: shutdownSignal } = gracefulShutdown();

const log = makeLogger("admin-app");

using database = new CustomDatabase(config().database.path);
const httpClient = new HttpFetch({ signal: shutdownSignal });

await using _server = new Server({
  hostname: config().apps.admin.host,
  port: config().apps.admin.port,
  onListen: (host, port) => {
    log.info(`Serving on http://${host}:${port}`);
  },
  onError: (error) => {
    log.error({ error }, "Something went wrong");

    return pageToHtmlResponse(errorPage(), 500);
  },
  app: new App({
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
  }),
});

await shutdownPromise;
