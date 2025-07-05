-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS withdrawals
(
    user_id      bigint NOT NULL,
    order_number text   NOT NULL,
    withdraw    DECIMAL(10, 5),
    processed_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS "withdrawals_user_id_index" ON "withdrawals" (user_id);
CREATE UNIQUE INDEX IF NOT EXISTS withdrawals_user_id_withdrawn_uindex ON withdrawals (user_id, order_number);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS withdrawals;
-- +goose StatementEnd
