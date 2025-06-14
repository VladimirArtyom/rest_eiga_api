CREATE TABLE IF NOT EXISTS users (
    id bigserial PRIMARY KEY,
    created_at timestamp(0) with time zone NOT NULL DEFAULT now(),
    name text NOT NULL, 
    email citext NOT NULL UNIQUE,  -- case insensitive
    password_hash bytea NOT NULL, -- byte type data
    activated bool NOT NULL, -- No default, NULL est la default.
    version integer NOT NULL DEFAULT 1
)
