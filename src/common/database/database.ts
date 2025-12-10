import type { SQLInputValue } from "node:sqlite";

export interface IDatabase {
  sql<T extends Record<string, unknown> = Record<string, unknown>>(
    strings: TemplateStringsArray,
    ...parameters: (SQLInputValue | SQLInputValue[])[]
  ): T[];
  transaction<T>(fn: () => T | Promise<T>): Promise<T>;
}
