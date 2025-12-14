import { config, initConfig } from "#src/common/config/mod.ts";
import { HttpFetch } from "#src/common/http/mod.ts";
import { gracefulShutdown } from "#src/common/process/mod.ts";
import { pageToHtmlResponse } from "#src/entrypoints/apps/admin/pages/page.tsx";
import { CustomDatabase } from "#src/common/database/mod.ts";
import { Server } from "#src/common/server/mod.ts";
import { errorPage } from "#src/entrypoints/apps/admin/pages/error.tsx";
import { App } from "#src/entrypoints/apps/admin/app.ts";
import { makeLogger } from "#src/common/logger/mod.ts";
import { EInputProvider } from "../../../common/database/enums/mod.ts";
import { providerFactory } from "../../../common/providers/provider.ts";

initConfig();

const { promise: shutdownPromise, signal: shutdownSignal } = gracefulShutdown();

const log = makeLogger("admin-app");
log.critical;
using database = new CustomDatabase(config().database.path);
const httpClient = new HttpFetch({ signal: shutdownSignal });

await using _server = new Server({
  hostname: config().apps.admin.host,
  port: config().apps.admin.port,
  onListen: (host, port) => {
    log.info(`Serving on http://${host}:${port}`);
  },
  onError: (error) => {
    log.error("Something went wrong", { error });

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
    },
    database,
  }),
});

await shutdownPromise;
