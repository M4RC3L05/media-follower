-- migrate:up

create table releases (
  -- Release id.
  id text not null,
  -- Release input id ref, default value will be used if input is deleted.
  input_id text not null default '__internal_deleted_input__',
  -- Release input provider ref, default value will be used if input is deleted.
  input_provider text not null default '__internal_deleted_input__',
  -- Output provider used to get this output.
  provider text not null,
  -- Release release date.
  released_at text not null check (released_at is strftime("%Y-%m-%dT%H:%M:%fZ", released_at)),
  -- Release raw data, must be valid json data.
  raw blob not null check (json_valid(raw, 4)),

  primary key(id, input_id, input_provider, provider),
  foreign key(input_id, input_provider)
    references inputs(id, provider)
    on update cascade
    on delete set default
-- Strict table, maintain rowid to untie during sorting.
) strict;

-- migrate:down
