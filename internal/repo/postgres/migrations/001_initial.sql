CREATE TABLE users (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(10) UNIQUE,
    time_zone VARCHAR(6) DEFAULT 'UTC',
    password_hash bytea NOT NULL
);

CREATE TABLE notes (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT REFERENCES users(id) ON DELETE CASCADE,
    title VARCHAR(50),
    text TEXT,
    date TIMESTAMP DEFAULT current_timestamp,
    is_finished BOOLEAN DEFAULT false
);

---- create above / drop below ----

DROP TABLE notes;

DROP TABLE users;
