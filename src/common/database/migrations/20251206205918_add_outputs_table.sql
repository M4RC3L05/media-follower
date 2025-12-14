-- migrate:up
CREATE TABLE outputs (
  id text NOT NULL,
  input_id text NOT NULL,
  provider text NOT NULL,
  raw blob NOT NULL CHECK (json_valid(raw, 4)),
  PRIMARY KEY(id, input_id, provider)
) strict;

-- migrate:down
