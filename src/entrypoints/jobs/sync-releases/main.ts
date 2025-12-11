import z from "@zod/zod";
import { parseArgs } from "@std/cli";
import { config, initConfig } from "../../../common/config/mod.ts";
import { BluRayComService } from "../../../common/services/blu-ray-com-service.ts";
import { HttpFetch } from "../../../common/http/mod.ts";
import { gracefulShutdown } from "../../../common/process/mod.ts";
import { ItunesService } from "../../../common/services/itunes-service.ts";
import { ReleaseSourceProvider } from "../../../common/database/enums/release-source-provider.ts";
import { CustomDatabase } from "../../../common/database/mod.ts";
import { App } from "./app.ts";

initConfig();

const { done, signal: shutdownSignal } = gracefulShutdown();

const { provider } = z.object({ provider: z.enum(ReleaseSourceProvider) })
  .parse(parseArgs(Deno.args));

using db = new CustomDatabase(config().database.path);
const httpClient = new HttpFetch({ signal: shutdownSignal });
const bluRayComService = new BluRayComService({ httpClient: httpClient });
const itunesService = new ItunesService({ httpClient: httpClient });

await new App({
  bluRayComService,
  database: db,
  itunesService,
  provider,
  signal: shutdownSignal,
}).execute();

await done();
