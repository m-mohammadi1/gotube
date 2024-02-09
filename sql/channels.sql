CREATE TABLE channels
(
    id         SERIAL PRIMARY KEY,
    user_id    INTEGER                 NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    youtube_id VARCHAR(255)            NOT NULL UNIQUE,
    title      VARCHAR(255)            NOT NULL,
    token      JSONB                   NOT NULL,
    added_at   TIMESTAMP DEFAULT NOW() NOT NULL
);