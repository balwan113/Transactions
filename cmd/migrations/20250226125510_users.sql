-- +goose Up
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    name TEXT UNIQUE NOT NULL,
    password TEXT NOT NULL,
    balance DECIMAL(10,2)
);
-- +goose Down
DELETE TABLE users;

