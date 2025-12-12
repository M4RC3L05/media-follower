import z from "@zod/zod";
import { parseArgs } from "@std/cli";
import { config, initConfig } from "#src/common/config/mod.ts";
import { HttpFetch } from "#src/common/http/mod.ts";
import { ItunesService } from "#src/common/services/itunes-service.ts";
import { gracefulShutdown } from "#src/common/process/mod.ts";
import { ReleaseSourceProvider } from "#src/common/database/enums/release-source-provider.ts";
import { CustomDatabase } from "#src/common/database/mod.ts";
import { App } from "#src/entrypoints/jobs/sync-release-sources/app.ts";
import { itunesMappers } from "../../../common/mappers/mod.ts";
import { itunesLookupArtistModelWithExtraSchema } from "../../../common/services/service.ts";

initConfig();

const { done, signal: shutdownSignal } = gracefulShutdown();
const { provider } = z.object({
  provider: z.literal(ReleaseSourceProvider.ITUNES),
})
  .parse(parseArgs(Deno.args));

using db = new CustomDatabase(config().database.path);

const itunesService = new ItunesService({
  httpClient: new HttpFetch({ signal: shutdownSignal }),
});

await new App({
  database: db,
  service: {
    fetchReleaseSource: async (source) => {
      const parsed = itunesLookupArtistModelWithExtraSchema.parse(source.raw);
      const fetched = await itunesService.lookupArtistById(parsed.artistId);

      if (!fetched) return;

      return itunesMappers.fromReleaseSourceToPersistance(fetched);
    },
  },
  provider,
  signal: shutdownSignal,
})
  .execute();

await done();
