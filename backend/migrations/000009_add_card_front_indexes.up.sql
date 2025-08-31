CREATE INDEX IF NOT EXISTS card_front_idx ON cards USING GIN (to_tsvector('english', front));
