CREATE TABLE servers (
    id SERIAL PRIMARY KEY,
    addr varchar(16) NOT NULL,
    port integer,
    nick varchar(80),
    votes_up integer DEFAULT 0,
    votes_down integer DEFAULT 0,
    last_modified timestamp NOT NULL DEFAULT NOW(),
    date_added timestamp NOT NULL DEFAULT NOW(),
    CONSTRAINT validate_port CHECK (port > 0 AND port < 65535)
);


CREATE TABLE tags (
    id serial PRIMARY KEY,
    tag varchar(30) UNIQUE NOT NULL
);


CREATE TABLE server_tags (
    server_id serial NOT NULL,
    tag_id serial NOT NULL,
    FOREIGN KEY (server_id) REFERENCES servers (id),
    FOREIGN KEY (tag_id) REFERENCES tags(id),
    CONSTRAINT no_duplicates UNIQUE(server_id, tag_id)
)