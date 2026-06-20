-- migrate:up

create table inputs (
  -- Input id.
  id text not null primary key,
  -- Input provider id to be used in the provider.
  internal_provider_id text not null,
  -- Input provier.
  provider text not null,
  -- Input name.
  name text not null,
  -- Input description.
  description text,
  -- Input image.
  image_url text,

  unique(internal_provider_id, provider)
-- Strict table, maintain rowid to untie during sorting.
) strict;

-- migrate:down
