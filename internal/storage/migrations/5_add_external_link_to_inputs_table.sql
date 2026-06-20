-- migrate:up

alter table inputs
  -- Input external link where it lives.
add column external_link text

-- migrate:down
