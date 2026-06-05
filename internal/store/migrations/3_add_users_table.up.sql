create table users (
  -- User id.
  id text not null primary key,
  -- User username.
  username text not null unique,
  -- User password, stored in hash form.
  "password" text not null

-- Strict table, no rowid.
) strict, without rowid;
