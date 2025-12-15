import { config, initConfig } from "#src/common/config/mod.ts";
import { gracefulShutdown } from "#src/common/process/mod.ts";
import { Server } from "#src/common/server/mod.ts";
import { CustomDatabase } from "#src/common/database/mod.ts";
import { EInputProvider } from "#src/common/database/enums/mod.ts";
import { providerFactory } from "#src/common/providers/provider.ts";
import { HttpFetch } from "#src/common/http/mod.ts";
import { App } from "./app.ts";

initConfig();

const { promise: shutdownPromise, signal: shutdownSignal } = gracefulShutdown();

using database = new CustomDatabase(config().database.path);
const httpClient = new HttpFetch({ signal: shutdownSignal });
await using _server = new Server({
  app: new App({
    database,
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
  }),
  hostname: config().apps.rssFeed.host,
  port: config().apps.rssFeed.port,
});

await shutdownPromise;
