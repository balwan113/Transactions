-- +goose Up
-- +goose StatementBegin
CREATE TABLE transactions (
    id SERIAL PRIMARY KEY,
    to_userid INT NOT NULL,
    at_userid INT NOT NULL,
    amount DECIMAL(10,2),
    created_at  TIMESTAMP DEFAULT now(),
    FOREIGN KEY(to_userid) REFERENCES users(id) ON DELETE CASCADE,
     FOREIGN KEY(at_userid) REFERENCES users(id) ON DELETE CASCADE
);
-- +goose Down
DELETE TABLE transactions;
-- +goose StatementEnd
