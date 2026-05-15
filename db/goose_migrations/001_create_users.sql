-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS users (
    username VARCHAR(45) PRIMARY KEY,
    password VARCHAR(45) NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS users;
-- +goose StatementEnd
