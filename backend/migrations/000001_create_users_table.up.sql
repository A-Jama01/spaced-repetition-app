CREATE TABLE IF NOT EXISTS users (
    id bigserial PRIMARY KEY,    
    username varchar(40) UNIQUE NOT NULL,
    password_hash bytea NOT NULL,
    created_at timestamp(0) with time zone NOT NULL DEFAULT NOW()
);
