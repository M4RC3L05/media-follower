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

    const dbInputs = this.#props.database.sql.all`
      select id, provider, json(raw) as raw from inputs
      where provider = ${this.#props.provider.provider}
    ` as DbInputsTable[];

    for (const [index, dbInput] of dbInputs.entries()) {
      if (this.#props.signal.aborted) break;

      try {
        log.info({ input: dbInput }, "Syncting outputs from input");

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

        const toDb = releases.map((item) => ({
          toDb: this.#props.provider.fromOutputToPersistence(dbInput, item),
          toJsonPatchDb: this.#props.provider.fromOutputToJsonPatchPersistance(
            dbInput,
            item,
          ),
        }));

        await this.#props.database.transaction(() => {
          toDb.map(({ toDb: item, toJsonPatchDb: jsonPatch }) =>
            this.#props.database.sql.run`
              insert into outputs
                (id,         input_id,      provider,            raw)
              values
                (${item.id}, ${item.input_id}, ${item.provider}, jsonb(${item.raw}))
              on conflict (id, input_id, provider)
                do update
                  set raw = jsonb_patch(raw, jsonb(${jsonPatch.raw}))
            `
          );
        });

        log.info(`Synced ${releases.length} outputs`);

        await delayIf(
          () => index < (dbInputs.length - 1),
          this.#props.signal,
        );
      } catch (error) {
        log.error({
          input: dbInput,
          error,
        }, "Could not sync outputs for input successfully");

        await delayIf(
          () => index < (dbInputs.length - 1),
          this.#props.signal,
        );
      }
    }

    // Truncate wal file as to not grow to mutch
    this.#props.database.sql.run`PRAGMA wal_checkpoint(TRUNCATE);`;
  }
}
