import { delay } from "@std/async";
import type { DbReleaseSourcesTable } from "#src/common/database/types.ts";
import type { IDatabase } from "#src/common/database/database.ts";
import { makeLogger } from "#src/common/logger/mod.ts";
import type { ReleaseSourceProvider } from "#src/common/database/enums/release-source-provider.ts";

const log = makeLogger("sync-release-sources-app");

type AppProps = {
  service: {
    fetchReleaseSource: (
      source: DbReleaseSourcesTable,
    ) => Promise<DbReleaseSourcesTable | undefined>;
  };
  database: IDatabase;
  provider: ReleaseSourceProvider;
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
      try {
        if (this.#props.signal.aborted) break;

        log.info("Syncing source", { source: { id: source.id } });

        const releaseSourceFetched = await this.#props.service
          .fetchReleaseSource(
            source,
          );

        if (!releaseSourceFetched) {
          await delayIf(
            () => index < (sources.length - 1),
            this.#props.signal,
          );
          continue;
        }

        this.#props.database.sql`
          update release_sources
          set raw = jsonb(${releaseSourceFetched.raw})
          where id = ${source.id}
        `;

        log.info("Synced source", { source: { id: source.id } });

        await delayIf(
          () => index < (sources.length - 1),
          this.#props.signal,
        );
      } catch (error) {
        log.error("Could not sync release source successfully", {
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
