-- +gooser Up

CREATE TABLE posts (
    id SERIAL,
    title text NOT NULL,
    slug text NOT NULL UNIQUE,
    content text NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY(id)
);

CREATE TABLE tags (
    id SERIAL,
    tag text NOT NULL,
    post_id INTEGER,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY(id, post_id),
    FOREIGN KEY (post_id) REFERENCES posts (id)
);

-- +gooser Down

DROP TABLE tags;
DROP TABLE posts;


