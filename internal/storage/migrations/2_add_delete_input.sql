-- migrate:up

-- Add a default input to be used for the outputs of a given input, when they are deleted.
insert into inputs (id, internal_provider_id, provider, name)
values ('__internal_deleted_input__', '__internal_deleted_input__', '__internal_deleted_input__', 'deleted');

-- migrate:down
