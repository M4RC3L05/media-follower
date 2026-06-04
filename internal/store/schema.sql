CREATE TABLE schema_migrations (version uint64,dirty bool);
CREATE UNIQUE INDEX version_unique ON schema_migrations (version);
CREATE TABLE inputs (
  -- Input id.
  id text not null,
  -- Input provier.
  provider text not null,
  -- Input raw data, used to resolve outputs, must be valid json data.
  raw blob not null check (json_valid(raw, 4)),

  primary key(id, provider)
-- Strict table, maintain rowid to untie during sorting.
) strict;
CREATE TABLE outputs (
  -- Output id.
  id text not null,
  -- Output input id ref, default value will be used if input is deleted.
  input_id text not null default '__internal_deleted_input__',
  -- Output input provider ref, default value will be used if input is deleted.
  input_provider text not null default '__internal_deleted_input__',
  -- Output provider rused to get this output.
  provider text not null,
  -- Output raw data, must be valid json data.
  raw blob not null check (json_valid(raw, 4)),

  primary key(id, input_id, input_provider),
  foreign key(input_id, input_provider)
    references inputs(id, provider)
    on update cascade
    on delete set default
-- Strict table, maintain rowid to untie during sorting.
) strict;
CREATE TABLE users (
  -- User id.
  id text not null primary key,
  -- User username.
  username text not null,
  -- User password, stored in hash form.
  "password" text not null

-- Strict table, no rowid.
) strict, without rowid;
