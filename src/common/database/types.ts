/* GENERATED FILE CONTENT DO NOT EDIT */

import { EInputProvider } from "#src/common/database/enums/mod.ts";

export type DbInputsTable = {
  id: string;
  provider: EInputProvider;
  raw: string;
};

export type DbOutputsTable = {
  id: string;
  input_id: string;
  provider: EInputProvider;
  raw: string;
};

export type DbUsersTable = {
  id: string;
  username: string;
  password: string;
};
