-- migrate:up

create table releases_new (
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

insert into releases_new
select * from releases
where strftime('%Y-%m-%dT%H:%M:%fZ', released_at) is not null
  and released_at = strftime('%Y-%m-%dT%H:%M:%fZ', released_at);

drop table releases;
alter table releases_new rename to releases;

-- migrate:down
