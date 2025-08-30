CREATE INDEX IF NOT EXISTS decks_name_idx ON decks USING GIN (to_tsvector('simple', name));
