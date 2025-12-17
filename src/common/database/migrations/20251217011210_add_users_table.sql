-- migrate:up
CREATE TABLE users(
  id text NOT NULL PRIMARY KEY,
  username text NOT NULL,
  "password" text NOT NULL
) strict,
without rowid;

-- migrate:down
