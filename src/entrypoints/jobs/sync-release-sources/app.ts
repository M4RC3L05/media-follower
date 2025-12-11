import { delay } from "@std/async";
import type { DbReleaseSourcesTable } from "../../../common/database/types.ts";
import {
  type IItunesService,
  itunesLookupArtistModelWithExtraSchema,
} from "../../../common/services/service.ts";
import * as itunesMappers from "#src/common/mappers/itunes-mappers.ts";
import type { IDatabase } from "../../../common/database/database.ts";
import type { ReleaseSourceProvider } from "../../../common/database/enums/release-source-provider.ts";
import { makeLogger } from "../../../common/logger/mod.ts";

const log = makeLogger("sync-release-sources-app");

type AppProps = {
  itunesService: IItunesService;
  database: IDatabase;
  provider: ReleaseSourceProvider.ITUNES;
  signal: AbortSignal;
};

export class App {
  #props: AppProps;

  constructor(props: AppProps) {
    this.#props = props;
  }

  async execute() {
    if (this.#props.signal.aborted) return;

    const sources = this.#props.database.sql<DbReleaseSourcesTable>`
      select *, json(raw) as raw from release_sources
      where provider = ${this.#props.provider}
    `;

    for (const source of sources) {
      try {
        if (this.#props.signal.aborted) break;

        log.info("Syncing source", { source });

        const parsed = itunesLookupArtistModelWithExtraSchema.parse(
          JSON.parse(source.raw),
        );

        const releaseSourceFetched = await this.#props.itunesService
          .lookupArtistById(
            parsed.artistId,
          );

        if (!releaseSourceFetched) {
          await delay(5000, { signal: this.#props.signal }).catch(() => {});
          continue;
        }

        const fetchedMapped = itunesMappers.fromReleaseSourceToPersistance(
          releaseSourceFetched,
        );

        this.#props.database.sql`
          update release_sources
          set raw = jsonb(${fetchedMapped.raw})
          where id = ${source.id}
        `;

        log.info("Synced source", { source });
        await delay(5000, { signal: this.#props.signal }).catch(() => {});
      } catch (error) {
        log.error("Could not sync release source successfully", {
          releaseSource: source,
          error,
        });
        await delay(5000, { signal: this.#props.signal }).catch(() => {});
      }
    }

    // Truncate wal file as to not grow to mutch
    this.#props.database.sql`PRAGMA wal_checkpoint(TRUNCATE);`;
  }
}
