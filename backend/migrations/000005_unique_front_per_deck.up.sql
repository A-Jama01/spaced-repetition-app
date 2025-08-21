ALTER TABLE cards ADD CONSTRAINT unique_front_per_deck UNIQUE (deck_id, front);
