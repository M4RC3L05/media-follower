CREATE TABLE IF NOT EXISTS "schema_migrations" (version varchar(128) primary key);
CREATE TABLE release_sources (
  id text NOT NULL,
  provider text NOT NULL,
  raw blob NOT NULL CHECK (json_valid(raw, 4)),
  PRIMARY KEY(id, provider)
) strict,
without rowid;
CREATE TABLE releases (
  id text NOT NULL,
  provider text NOT NULL,
  "type" text NOT NULL,
  "releasedAt" text NOT NULL,
  raw blob NOT NULL CHECK (json_valid(raw, 4)),
  PRIMARY KEY(id, provider, "type")
) strict;
-- Dbmate schema migrations
INSERT INTO "schema_migrations" (version) VALUES
  ('20251206205917'),
  ('20251206205918');
