import type { SQLTagStore } from "node:sqlite";

export interface IDatabase {
  sql: SQLTagStore;
  transaction<T>(fn: () => T | Promise<T>): Promise<T>;
}
