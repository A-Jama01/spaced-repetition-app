ALTER TABLE decks ADD CONSTRAINT unique_decks_user_id_name UNIQUE (user_id, name);
