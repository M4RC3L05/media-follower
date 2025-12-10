-- migrate:up
CREATE TABLE release_sources (
  id text NOT NULL,
  provider text NOT NULL,
  raw blob NOT NULL CHECK (json_valid(raw, 4)),
  PRIMARY KEY(id, provider)
) strict,
without rowid;

-- migrate:down
