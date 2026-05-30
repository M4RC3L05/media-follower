import type { VNode } from "preact";
import type { EInputProvider } from "../database/enums/mod.ts";
import type { DbInputsTable, DbOutputsTable } from "../database/types.ts";
import type { Feed } from "feed";

export interface IProvider<
  T extends EInputProvider = EInputProvider,
  // deno-lint-ignore no-explicit-any
  I = Record<string, any>,
  // deno-lint-ignore no-explicit-any
  O = Record<string, any>,
> extends
  IProviderRender,
  IProviderMapper<I, O>,
  IProviderRepository,
  IProviderFeed {
  readonly provider: T;

  lookupInput(term: string): Promise<I | undefined>;

  fetchInput(row: DbInputsTable): Promise<I>;
  fetchOutputs(input: I): Promise<O[]>;
}

export interface IProviderRender {
  renderInputListItem(row: DbInputsTable): VNode;
  renderOutputListItem(row: DbOutputsTable): VNode;
}

export interface IProviderMapper<I, O> {
  fromInputToPersistence(item: I): DbInputsTable;
  fromPersistenceToInput(row: DbInputsTable): I;
  fromOutputToPersistence(row: DbInputsTable, item: O): DbOutputsTable;
  fromPersistenceToOutput(row: DbOutputsTable): O;
}

export type IProviderRepositoryQueryOutputsProps = {
  pagination: { limit: number; page: number };
  queries?: Record<string, unknown>;
};

export interface IProviderRepository {
  queryOutputs(
    props: IProviderRepositoryQueryOutputsProps,
  ): Promise<DbOutputsTable[]>;
}

export type IProviderFeedGetOutputsFeedProps = {
  queries?: Record<string, unknown>;
};

export interface IProviderFeed {
  getOutputsFeed(props: IProviderFeedGetOutputsFeedProps): Feed;
}
