-- migrate:up
CREATE INDEX idx_outputs_provider ON outputs(provider);

CREATE INDEX idx_outputs_input_id ON outputs(input_id);

CREATE INDEX idx_outputs_provider_input_id ON outputs(provider, input_id);

-- migrate:down
