-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS users
(
    id       BIGSERIAL PRIMARY KEY,
    login    text NOT NULL,
    password bytea NOT NULL
);
-- +goose StatementEnd
-- +goose StatementBegin
CREATE UNIQUE INDEX IF NOT EXISTS users_login_uindex ON users (login);
-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS users;
-- +goose StatementEnd