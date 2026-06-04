create table inputs (
  -- Input id.
  id text not null,
  -- Input provier.
  provider text not null,
  -- Input raw data, used to resolve outputs, must be valid json data.
  raw blob not null check (json_valid(raw, 4)),

  primary key(id, provider)
-- Strict table, maintain rowid to untie during sorting.
) strict;
