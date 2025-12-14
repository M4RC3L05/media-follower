import z from "@zod/zod";
import { parseArgs } from "@std/cli";
import { config, initConfig } from "#src/common/config/mod.ts";
import { HttpFetch } from "#src/common/http/mod.ts";
import { gracefulShutdown } from "#src/common/process/mod.ts";
import { EInputProvider } from "../../../common/database/enums/input-provider.ts";
import { CustomDatabase } from "#src/common/database/mod.ts";
import { providerFactory } from "../../../common/providers/provider.ts";
import { App } from "./app.ts";

initConfig();

const { done, signal: shutdownSignal } = gracefulShutdown();

const run = async () => {
  const { provider } = z.object({ provider: z.enum(EInputProvider) })
    .parse(parseArgs(Deno.args));

  using database = new CustomDatabase(config().database.path);
  const httpClient = new HttpFetch({ signal: shutdownSignal });

  await new App({
    database,
    provider: providerFactory(provider, { database, httpClient }),
    signal: shutdownSignal,
  }).execute();
};

if (!shutdownSignal.aborted) {
  await run();
}

await done();
