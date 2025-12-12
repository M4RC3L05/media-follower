import z from "@zod/zod";
import { parseArgs } from "@std/cli";
import { config, initConfig } from "#src/common/config/mod.ts";
import { BluRayComService } from "#src/common/services/blu-ray-com-service.ts";
import { HttpFetch } from "#src/common/http/mod.ts";
import { gracefulShutdown } from "#src/common/process/mod.ts";
import { ItunesService } from "#src/common/services/itunes-service.ts";
import { ReleaseSourceProvider } from "#src/common/database/enums/release-source-provider.ts";
import { CustomDatabase } from "#src/common/database/mod.ts";
import { App } from "#src/entrypoints/jobs/sync-releases/app.ts";

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
