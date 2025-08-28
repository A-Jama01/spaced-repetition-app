CREATE TABLE IF NOT EXISTS decks (
    id bigserial PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE, 
    name varchar(70) NOT NULL,
    created_at timestamp(0) with time zone NOT NULL DEFAULT NOW()
);
