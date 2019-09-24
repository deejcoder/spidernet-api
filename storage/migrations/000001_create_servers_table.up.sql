CREATE TABLE servers (
    id integer PRIMARY KEY,
    addr varchar(16),
    port integer,
    nick varchar(80),
    votes_up integer DEFAULT 0,
    votes_down integer DEFAULT 0,
    tags varchar(80),
    last_modified timestamp,
    date_added timestamp,
    CONSTRAINT validate_port CHECK (port > 0 AND port < 65535)
)