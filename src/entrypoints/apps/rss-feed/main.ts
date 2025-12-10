import { config, initConfig } from "#src/common/config/mod.ts";
import { gracefulShutdown } from "#src/common/process/mod.ts";
import { App } from "#src/entrypoints/apps/rss-feed/app.ts";
import { Server } from "#src/common/server/mod.ts";
import { CustomDatabase } from "../../../common/database/mod.ts";

initConfig();

const { promise: shutdownPromise } = gracefulShutdown();

using db = new CustomDatabase(config().database.path);
await using _server = new Server({
  app: new App({ database: db }),
  hostname: config().apps.rssFeed.host,
  port: config().apps.rssFeed.port,
});

await shutdownPromise;
