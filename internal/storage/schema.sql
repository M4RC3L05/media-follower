CREATE TABLE "schema_migrations" (version varchar(128) primary key);
CREATE TABLE inputs (
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
  image_url text, external_link text,

  unique(internal_provider_id, provider)
-- Strict table, maintain rowid to untie during sorting.
) strict;
CREATE TABLE users (
  -- User id.
  id text not null primary key,
  -- User username.
  username text not null unique,
  -- User password, stored in hash form.
  "password" text not null

-- Strict table, no rowid.
) strict, without rowid;
CREATE TABLE "releases" (
  -- Release id.
  id text not null primary key,
  -- Release provier id to be used in the provider.
  internal_provider_id text not null,
  -- Release input id ref, default value will be used if input is deleted.
  input_id text not null default '__internal_deleted_input__',
  -- Release title.
  title text not null,
  -- Release description.
  description text,
  -- Release image.
  image_url text,
  -- Release external link where it lives.
  external_link text,
  -- Release release date.
  released_at text not null check (
    strftime('%Y-%m-%dT%H:%M:%fZ', released_at) is not null
    and released_at = strftime('%Y-%m-%dT%H:%M:%fZ', released_at)
  ),

  unique(internal_provider_id, input_id),

  foreign key(input_id)
    references inputs(id)
    on update cascade
    on delete set default
-- Strict table, maintain rowid to untie during sorting.
) strict;
-- Dbmate schema migrations
INSERT INTO "schema_migrations" (version) VALUES
  ('1'),
  ('2'),
  ('3'),
  ('4'),
  ('5'),
  ('6');
