-- migrate:up
CREATE TABLE releases (
  id text NOT NULL,
  provider text NOT NULL,
  "type" text NOT NULL,
  "releasedAt" text NOT NULL,
  raw blob NOT NULL CHECK (json_valid(raw, 4)),
  PRIMARY KEY(id, provider, "type")
) strict;

-- migrate:down
