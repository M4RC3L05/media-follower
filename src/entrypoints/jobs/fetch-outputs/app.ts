import type { IDatabase } from "#src/common/database/database.ts";
import { makeLogger } from "#src/common/logger/mod.ts";
import { delay } from "@std/async";
import type { DbInputsTable } from "#src/common/database/types.ts";
import type { IProvider } from "#src/common/providers/interfaces.ts";

const log = makeLogger("sync-releases-app");

export type AppProps = {
  database: IDatabase;
  provider: IProvider;
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

    const dbInputs = this.#props.database.sql<DbInputsTable>`
      select id, provider, json(raw) as raw from inputs
      where provider = ${this.#props.provider.provider}
    `;

    for (const [index, dbInput] of dbInputs.entries()) {
      if (this.#props.signal.aborted) break;

      try {
        log.info("Syncting outputs from input", { input: dbInput });

        const input = this.#props.provider.fromPersistenceToInput(
          dbInput,
        );
        const releases = await this.#props.provider.fetchOutputs(input);

        if (releases.length <= 0) {
          await delayIf(
            () => index < (dbInputs.length - 1),
            this.#props.signal,
          );
          continue;
        }

        const toDb = releases.map((item) =>
          this.#props.provider.fromOutputToPersistence(dbInput, item)
        );

        await this.#props.database.transaction(() => {
          toDb.map((item) =>
            this.#props.database.sql`
              insert into outputs
                (id,         input_id,      provider,            raw)
              values
                (${item.id}, ${item.input_id}, ${item.provider}, jsonb(${item.raw}))
              on conflict (id, input_id, provider)
                do update
                  set raw = jsonb(${item.raw})
            `
          );
        });

        log.info(`Synced ${releases.length} outputs`);

        await delayIf(
          () => index < (dbInputs.length - 1),
          this.#props.signal,
        );
      } catch (error) {
        log.error("Could not sync outputs for input successfully", {
          input: dbInput,
          error,
        });

        await delayIf(
          () => index < (dbInputs.length - 1),
          this.#props.signal,
        );
      }
    }

    // Truncate wal file as to not grow to mutch
    this.#props.database.sql`PRAGMA wal_checkpoint(TRUNCATE);`;
  }
}
