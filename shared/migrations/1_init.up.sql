CREATE TABLE IF NOT EXISTS users
(
    id           SERIAL  PRIMARY KEY,
    email        TEXT    NOT NULL UNIQUE,
    pass_hash    BYTEA   NOT NULL,
    phone        TEXT    NOT NULL,
    photo        BYTEA   NULL,
    active      TEXT NULL,
);
CREATE INDEX IF NOT EXISTS idx_email ON users (email);

CREATE TABLE IF NOT EXISTS roles (
    id   SERIAL PRIMARY KEY,
    name TEXT NOT NULL UNIQUE   
);

CREATE TABLE IF NOT EXISTS user_roles (
    user_id INT REFERENCES users(id) ON DELETE CASCADE,
    role_id INT REFERENCES roles(id) ON DELETE CASCADE,
    PRIMARY KEY (user_id, role_id)
);

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
    chat_type INT NOT NULL,
    pin TEXT NULL,
    mute TEXT NULL
)

CREATE TABLE IF NOT EXISTS messages
(
    id SERIAL PRIMARY KEY,
    content TEXT NOT NULL,
    created_at TIME NOT NULL,
    created_by INT NOT NULL,
    reply_to_id INT NOT NULL,
    received_by INT NOT NULL,
    received_at TIME NOT NULL,
    pin TEXT NULL,
    is_read TEXT NULL
)
