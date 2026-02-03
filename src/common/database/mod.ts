import { DatabaseSync, type SQLTagStore } from "node:sqlite";
import type { IDatabase } from "./database.ts";

export class CustomDatabase extends DatabaseSync implements IDatabase {
  sql: SQLTagStore;

  constructor(path: string) {
    super(path);

    this.sql = this.createTagStore(100);
    this.exec("pragma journal_mode = WAL");
    this.exec("pragma busy_timeout = 5000");
    this.exec("pragma foreign_keys = ON");
    this.exec("pragma synchronous = NORMAL");
    this.exec("pragma temp_store = MEMORY");
    this.exec("pragma optimize = 0x10002");
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

  override [Symbol.dispose]() {
    this.exec("pragma optimize");

    this.close();
  }
}
