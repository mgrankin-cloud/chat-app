CREATE TABLE IF NOT EXISTS users
(
    id           SERIAL  PRIMARY KEY,
    email        TEXT    NOT NULL UNIQUE,
    pass_hash    BYTEA   NOT NULL,
    phone        TEXT    NOT NULL,
    photo        BYTEA   NULL
);
CREATE INDEX IF NOT EXISTS idx_email ON users (email);

CREATE TABLE IF NOT EXISTS apps
(
    id     SERIAL PRIMARY KEY,
    name   TEXT NOT NULL UNIQUE,
    secret TEXT NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS chats
(
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL UNIQUE,
    photo        BYTEA   NULL,
    chat_type INT NOT NULL
)

CREATE TABLE IF NOT EXISTS messages
(
    id SERIAL PRIMARY KEY,
    content TEXT NOT NULL,
    created_at TIME NOT NULL,
    created_by INT NOT NULL,
    reply_to_id INT NOT NULL,
    received_by INT NOT NULL,
    received_at TIME NOT NULL
)
