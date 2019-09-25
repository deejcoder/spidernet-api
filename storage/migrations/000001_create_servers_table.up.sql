CREATE TABLE servers (
    id SERIAL PRIMARY KEY,
    addr varchar(16) NOT NULL,
    port integer,
    nick varchar(80),
    votes_up integer DEFAULT 0,
    votes_down integer DEFAULT 0,
    tags text[] NOT NULL DEFAULT '{}'::text[],
    last_modified timestamp NOT NULL DEFAULT NOW(),
    date_added timestamp NOT NULL DEFAULT NOW(),
    CONSTRAINT validate_port CHECK (port > 0 AND port < 65535)
);

CREATE INDEX tags_index ON servers USING GIN(tags);