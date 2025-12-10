import z from "@zod/zod";
import { delay } from "@std/async";
import { parseArgs } from "@std/cli";
import { initConfig } from "../../../common/config/mod.ts";
import { makeLogger } from "../../../common/logger/mod.ts";
import {
  bluRayComCountrySchema,
  itunesLookupArtistModelWithExtraSchema,
  ITunesLookupEntityType,
} from "../../../common/services/service.ts";
import { BluRayComService } from "../../../common/services/blu-ray-com-service.ts";
import { HttpFetch } from "../../../common/http/mod.ts";
import { gracefulShutdown } from "../../../common/process/mod.ts";
import { ItunesService } from "../../../common/services/itunes-service.ts";
import { ReleaseSourceProvider } from "../../../common/database/enums/release-source-provider.ts";
import { makeDatabase } from "../../../common/database/mod.ts";
import type {
  DbReleaseSourcesTable,
  DbReleasesTable,
} from "../../../common/database/types.ts";

initConfig();

const { done, signal: shutdownSignal } = gracefulShutdown();

const { provider } = z.object({ provider: z.enum(ReleaseSourceProvider) })
  .parse(parseArgs(Deno.args));

const now = new Date();
const log = makeLogger("sync-releases");
using db = makeDatabase();
const bluRayComService = new BluRayComService({
  httpClient: new HttpFetch({ signal: shutdownSignal }),
});
const itunesService = new ItunesService({
  httpClient: new HttpFetch({ signal: shutdownSignal }),
});

const sources = db.sql<DbReleaseSourcesTable>`
  select *, json(raw) as raw from release_sources
  where provider = ${provider}
`;

for (const source of sources) {
  log.info("Syncting releases from source", { raw: JSON.parse(source.raw) });

  let releases: (DbReleasesTable)[] = [];

  switch (provider) {
    case ReleaseSourceProvider.BLU_RAY_COM: {
      const parsed = bluRayComCountrySchema.parse(JSON.parse(source.raw));
      const remote = await bluRayComService
        .getBlurayReleasesByCountryForMonth(
          parsed.code,
          now.getFullYear(),
          now.getMonth() + 1,
        );
      releases = remote.map((item) =>
        BluRayComService.toReleasePersistance(item)
      );
      break;
    }
    case ReleaseSourceProvider.ITUNES: {
      const parsed = itunesLookupArtistModelWithExtraSchema.parse(
        JSON.parse(source.raw),
      );
      const [albums, songs] = await Promise.all([
        await itunesService.lookupLatestReleasesByArtist(
          String(parsed.artistId),
          ITunesLookupEntityType.ALBUM,
          30,
        ),
        await itunesService.lookupLatestReleasesByArtist(
          String(parsed.artistId),
          ITunesLookupEntityType.SONG,
          30,
        ),
      ]);

      releases.push(
        ...albums.map((item) => ItunesService.toReleasePersistance(item)),
      );
      releases.push(
        ...songs.map((item) => ItunesService.toReleasePersistance(item)),
      );
    }
  }

  await db.transaction(async () => {
    await Promise.all(
      releases.map((item) =>
        db.sql<DbReleasesTable>`
          insert or replace into releases
            (id,         "type",       provider,         "releasedAt",       raw)
          values
            (${item.id}, ${item.type}, ${item.provider}, ${item.releasedAt}, jsonb(${item.raw}))
          returning *;
        `
      ),
    );
  });

  log.info(`Synced ${releases.length} releases`);

  await delay(5000).catch(() => {});
}

// Truncate wal file as to not grow to mutch
db.sql`PRAGMA wal_checkpoint(TRUNCATE);`;

await done();
