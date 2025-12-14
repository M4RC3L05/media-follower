-- migrate:up
CREATE INDEX idx_outputs_imr_provider_input_id_raw_wrapper_raw_collection_id ON outputs(
  provider,
  input_id,
  raw ->> 'wrapperType',
  raw ->> 'collectionId'
);

CREATE INDEX idx_outputs_imr_provider_raw_release_date_raw_wrapper_type ON outputs(
  provider,
  raw ->> 'releaseDate',
  raw ->> 'wrapperType'
);

-- migrate:down
