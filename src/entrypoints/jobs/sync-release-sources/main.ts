import z from "@zod/zod";
import { parseArgs } from "@std/cli";
import { initConfig } from "../../../common/config/mod.ts";
import { delay } from "@std/async/delay";
import { HttpFetch } from "../../../common/http/mod.ts";
import { ItunesService } from "../../../common/services/itunes-service.ts";
import { gracefulShutdown } from "../../../common/process/mod.ts";
import { ReleaseSourceProvider } from "../../../common/database/enums/release-source-provider.ts";
import { makeDatabase } from "../../../common/database/mod.ts";
import type { DbReleaseSourcesTable } from "../../../common/database/types.ts";
import {
  itunesLookupArtistModelWithExtraSchema,
} from "../../../common/services/service.ts";
import { makeLogger } from "../../../common/logger/mod.ts";

initConfig();

const { done, signal: shutdownSignal } = gracefulShutdown();
const { provider } = z.object({
  provider: z.literal(ReleaseSourceProvider.ITUNES),
})
  .parse(parseArgs(Deno.args));

const log = makeLogger("sync-release-sources");
using db = makeDatabase();

const itunesService = new ItunesService({
  httpClient: new HttpFetch({ signal: shutdownSignal }),
});

const sources = db.sql<DbReleaseSourcesTable>`
  select *, json(raw) as raw from release_sources
  where provider = ${provider}
`;

for (const releaseSource of sources) {
  log.info("Syncing source", { raw: JSON.parse(releaseSource.raw) });

  const parsed = itunesLookupArtistModelWithExtraSchema.parse(
    JSON.parse(releaseSource.raw),
  );
  const releaseSourceFetched = await itunesService.lookupArtistById(
    parsed.artistId,
  );

  if (!releaseSourceFetched) {
    await delay(5000).catch(() => {});
    continue;
  }

  const fetchedMapped = ItunesService.toReleaseSourcePersistance(
    releaseSourceFetched,
  );

  db.sql`
    update release_sources
    set raw = jsonb(${fetchedMapped.raw})
    where id = ${releaseSource.id}
  `;

  log.info("Synced source", { raw: JSON.parse(fetchedMapped.raw) });
  await delay(5000).catch(() => {});
}

// Truncate wal file as to not grow to mutch
db.sql`PRAGMA wal_checkpoint(TRUNCATE);`;

await done();
