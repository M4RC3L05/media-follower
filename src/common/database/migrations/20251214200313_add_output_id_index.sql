-- migrate:up
CREATE INDEX idx_outputs_id ON outputs(id);

-- migrate:down
