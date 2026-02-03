import { delay } from "@std/async";
import type { IDatabase } from "#src/common/database/database.ts";
import { makeLogger } from "#src/common/logger/mod.ts";
import type { IProvider } from "#src/common/providers/interfaces.ts";
import type { DbInputsTable } from "#src/common/database/types.ts";

const log = makeLogger("sync-release-sources-app");

type AppProps = {
  provider: IProvider;
  database: IDatabase;
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

    const dbInputs = this.#props.database.sql.all`
      select id, provider, json(raw) as raw from inputs
      where provider = ${this.#props.provider.provider}
    ` as DbInputsTable[];

    for (const [index, dbInput] of dbInputs.entries()) {
      try {
        if (this.#props.signal.aborted) break;

        log.info({ input: dbInput }, "Syncing source");

        const releaseSourceFetched = await this.#props.provider.fetchInput(
          dbInput,
        );

        if (!releaseSourceFetched) {
          await delayIf(
            () => index < (dbInputs.length - 1),
            this.#props.signal,
          );
          continue;
        }

        const db = this.#props.provider.fromInputToPersistence(
          releaseSourceFetched,
        );

        this.#props.database.sql.run`
          update inputs
          set raw = jsonb(${db.raw})
          where id = ${dbInput.id}
        `;

        log.info({ input: dbInput }, "Synced input");

        await delayIf(
          () => index < (dbInputs.length - 1),
          this.#props.signal,
        );
      } catch (error) {
        log.error({
          input: dbInput,
          error,
        }, "Could not sync input successfully");

        await delayIf(
          () => index < (dbInputs.length - 1),
          this.#props.signal,
        );
      }
    }

    // Truncate wal file as to not grow to mutch
    this.#props.database.sql.all`PRAGMA wal_checkpoint(TRUNCATE);`;
  }
}
