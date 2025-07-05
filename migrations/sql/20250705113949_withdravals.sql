-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS withdrawns
(
    user_id  bigint NOT NULL,
    withdrawn  DECIMAL(10,5)
);
CREATE INDEX IF NOT EXISTS "withdrawns_user_id_index" ON "withdrawns" (user_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS withdrawns;
-- +goose StatementEnd
