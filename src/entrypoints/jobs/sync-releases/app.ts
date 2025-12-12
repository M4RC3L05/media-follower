import type { IDatabase } from "#src/common/database/database.ts";
import type { ReleaseSourceProvider } from "#src/common/database/enums/release-source-provider.ts";
import type {
  DbReleaseSourcesTable,
  DbReleasesTable,
} from "#src/common/database/types.ts";
import { makeLogger } from "#src/common/logger/mod.ts";
import { delay } from "@std/async";

const log = makeLogger("sync-releases-app");

export type AppProps = {
  database: IDatabase;
  provider: ReleaseSourceProvider;
  service: {
    fetchReleasesFromSource: (
      source: DbReleaseSourcesTable,
    ) => Promise<Array<DbReleasesTable>>;
  };
  signal: AbortSignal;
};

const delayIf = async (cond: () => boolean, signal: AbortSignal) => {
  if (!cond()) return;

  await delay(5000, { signal }).catch(() => {});
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

    for (const [index, source] of sources.entries()) {
      if (this.#props.signal.aborted) break;

      try {
        log.info("Syncting releases from source", {
          source: { id: source.id },
        });

        const releases = await this.#props.service.fetchReleasesFromSource(
          source,
        );

        if (releases.length <= 0) {
          await delayIf(
            () => index < (sources.length - 1),
            this.#props.signal,
          );
          continue;
        }

        await this.#props.database.transaction(() => {
          releases.map((item) =>
            this.#props.database.sql<DbReleasesTable>`
              insert into releases
                (id,         "type",       provider,         "releasedAt",       raw)
              values
                (${item.id}, ${item.type}, ${item.provider}, ${item.releasedAt}, jsonb(${item.raw}))
              on conflict (id, provider, "type")
                do update
                  set raw = jsonb(${item.raw})
            `
          );
        });

        log.info(`Synced ${releases.length} releases`);

        await delayIf(
          () => index < (sources.length - 1),
          this.#props.signal,
        );
      } catch (error) {
        log.error("Could not sync releases for source successfully", {
          source: { id: source.id },
          error,
        });

        await delayIf(
          () => index < (sources.length - 1),
          this.#props.signal,
        );
      }
    }

    // Truncate wal file as to not grow to mutch
    this.#props.database.sql`PRAGMA wal_checkpoint(TRUNCATE);`;
  }
}
