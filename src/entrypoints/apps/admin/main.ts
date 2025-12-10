import { config, initConfig } from "../../../common/config/mod.ts";
import { HttpFetch } from "../../../common/http/mod.ts";
import { ItunesService } from "../../../common/services/itunes-service.ts";
import { gracefulShutdown } from "../../../common/process/mod.ts";
import { pageToHtmlResponse } from "./pages/page.tsx";
import { BluRayComService } from "../../../common/services/blu-ray-com-service.ts";
import { CustomDatabase } from "../../../common/database/mod.ts";
import { Server } from "#src/common/server/mod.ts";
import { errorPage } from "./pages/error.tsx";
import { App } from "./app.ts";
import { makeLogger } from "../../../common/logger/mod.ts";

initConfig();

const { promise: shutdownPromise, signal: shutdownSignal } = gracefulShutdown();

const log = makeLogger("admin-app");
log.critical;
using db = new CustomDatabase(config().database.path);
const fetchClient = new HttpFetch({ signal: shutdownSignal });

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
    blurayComService: new BluRayComService({ httpClient: fetchClient }),
    itunesService: new ItunesService({ httpClient: fetchClient }),
    database: db,
  }),
});

await shutdownPromise;
