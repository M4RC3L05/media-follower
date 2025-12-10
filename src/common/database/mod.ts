import { config } from "#src/common/config/mod.ts";
import { DatabaseSync } from "node:sqlite";
import type { SQLInputValue, StatementSync } from "node:sqlite";
import type { IDatabase } from "./database.ts";

export class CustomDatabase extends DatabaseSync implements IDatabase {
  #cache = new Map<string, StatementSync>();

  #ensureInCache(query: string) {
    if (!this.#cache.has(query)) {
      this.#cache.set(query, super.prepare(query));
    }

    return this.#cache.get(query)!;
  }

  async transaction<T>(fn: () => T | Promise<T>) {
    try {
      this.exec("begin immediate");
      const result = await fn();
      this.exec("commit");

      return result;
    } catch (error) {
      this.exec("rollback");

      throw error;
    }
  }

  sql<T extends Record<string, unknown> = Record<string, unknown>>(
    strings: TemplateStringsArray,
    ...parameters: (SQLInputValue | SQLInputValue[])[]
  ): T[] {
    const sql = strings.reduce((acc, str, i) => {
      return acc + str +
        (i < parameters.length
          ? (Array.isArray(parameters[i])
            ? parameters[i].length <= 0
              ? ""
              : Array.from({ length: parameters[i].length }, () => "?").join(
                ",",
              )
            : "?")
          : "");
    }, "").trim();

    try {
      const stmt = this.prepare(sql);
      return stmt.all(...parameters.flat()) as T[];
    } catch (error) {
      this.#cache.delete(sql);
      throw error;
    }
  }

  override prepare(sql: string, cache = true): StatementSync {
    if (!cache) return super.prepare(sql);
    return this.#ensureInCache(sql);
  }

  override close(): void {
    super.close();

    this.#cache.clear();
  }

  override [Symbol.dispose]() {
    this.exec("pragma optimize");

    this.close();
  }
}

export const makeDatabase = () => {
  const db = new CustomDatabase(config().database.path);

  db.exec("pragma journal_mode = WAL");
  db.exec("pragma busy_timeout = 5000");
  db.exec("pragma foreign_keys = ON");
  db.exec("pragma synchronous = NORMAL");
  db.exec("pragma temp_store = MEMORY");
  db.exec("pragma optimize = 0x10002");

  return db;
};
