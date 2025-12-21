-- +goose up
CREATE TABLE users(
    id UUID PRIMARY KEY,
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    email TEXT NOT NULL UNIQUE
);

-- +goose down
DROP TABLE users;
