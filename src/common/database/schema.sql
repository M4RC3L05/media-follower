CREATE TABLE IF NOT EXISTS "schema_migrations" (version varchar(128) primary key);
CREATE TABLE inputs (
  id text NOT NULL,
  provider text NOT NULL,
  raw blob NOT NULL CHECK (json_valid(raw, 4)),
  PRIMARY KEY(id, provider)
) strict,
without rowid;
CREATE TABLE outputs (
  id text NOT NULL,
  input_id text NOT NULL,
  provider text NOT NULL,
  raw blob NOT NULL CHECK (json_valid(raw, 4)),
  PRIMARY KEY(id, input_id, provider)
) strict;
CREATE INDEX idx_outputs_provider ON outputs(provider);
CREATE INDEX idx_outputs_input_id ON outputs(input_id);
CREATE INDEX idx_outputs_provider_input_id ON outputs(provider, input_id);
-- Dbmate schema migrations
INSERT INTO "schema_migrations" (version) VALUES
  ('20251206205917'),
  ('20251206205918'),
  ('20251214193125');
