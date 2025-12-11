import type { IDatabase } from "../../../common/database/database.ts";
import { ReleaseSourceProvider } from "../../../common/database/enums/release-source-provider.ts";
import type {
  DbReleaseSourcesTable,
  DbReleasesTable,
} from "../../../common/database/types.ts";
import { makeLogger } from "../../../common/logger/mod.ts";
import {
  bluRayComCountrySchema,
  type IBlurayComService,
  type IItunesService,
  itunesLookupArtistModelWithExtraSchema,
  ITunesLookupEntityType,
} from "../../../common/services/service.ts";
import * as bluRayComMappers from "#src/common/mappers/blu-ray-com-mappers.ts";
import * as itunesMappers from "#src/common/mappers/itunes-mappers.ts";
import { delay } from "@std/async/delay";

const log = makeLogger("sync-releases-app");

export type AppProps = {
  database: IDatabase;
  provider: ReleaseSourceProvider;
  bluRayComService: IBlurayComService;
  itunesService: IItunesService;
  signal: AbortSignal;
};

export class App {
  #props: AppProps;

  constructor(props: AppProps) {
    this.#props = props;
  }

  async execute() {
    if (this.#props.signal.aborted) return;

    const now = new Date();
    const sources = this.#props.database.sql<DbReleaseSourcesTable>`
      select *, json(raw) as raw from release_sources
      where provider = ${this.#props.provider}
    `;

    for (const source of sources) {
      if (this.#props.signal.aborted) break;

      try {
        log.info("Syncting releases from source", { source });

        let releases: (DbReleasesTable)[] = [];

        switch (this.#props.provider) {
          case ReleaseSourceProvider.BLU_RAY_COM: {
            const parsed = bluRayComCountrySchema.parse(JSON.parse(source.raw));
            const remote = await this.#props.bluRayComService
              .getBlurayReleasesByCountryForMonth(
                parsed.code,
                now.getFullYear(),
                now.getMonth() + 1,
              );

            releases = remote.map((item) =>
              bluRayComMappers.fromReleaseToPersistance(item)
            );
            break;
          }
          case ReleaseSourceProvider.ITUNES: {
            const parsed = itunesLookupArtistModelWithExtraSchema.parse(
              JSON.parse(source.raw),
            );
            const [albums, songs] = await Promise.all([
              await this.#props.itunesService.lookupLatestReleasesByArtist(
                String(parsed.artistId),
                ITunesLookupEntityType.ALBUM,
                50,
              ),
              await this.#props.itunesService.lookupLatestReleasesByArtist(
                String(parsed.artistId),
                ITunesLookupEntityType.SONG,
                50,
              ),
            ]);

            releases.push(
              ...albums.map((item) =>
                itunesMappers.fromReleaseToPersistance(item)
              ),
            );
            releases.push(
              ...songs.map((item) =>
                itunesMappers.fromReleaseToPersistance(item)
              ),
            );
          }
        }

        await this.#props.database.transaction(async () => {
          await Promise.all(
            releases.map((item) =>
              this.#props.database.sql<DbReleasesTable>`
              insert into releases
                (id,         "type",       provider,         "releasedAt",       raw)
              values
                (${item.id}, ${item.type}, ${item.provider}, ${item.releasedAt}, jsonb(${item.raw}))
              on conflict (id, provider, "type")
                do update set
                  "releasedAt" = iif("releasedAt" is not null, "releasedAt", ${item.releasedAt}),
                  raw = jsonb(${item.raw})
              returning *;
            `
            ),
          );
        });

        log.info(`Synced ${releases.length} releases`);

        await delay(5000, { signal: this.#props.signal }).catch(() => {});
      } catch (error) {
        log.error("Could not sync releases for source successfully", {
          source,
          error,
        });
        await delay(5000, { signal: this.#props.signal }).catch(() => {});
      }
    }

    // Truncate wal file as to not grow to mutch
    this.#props.database.sql`PRAGMA wal_checkpoint(TRUNCATE);`;
  }
}
