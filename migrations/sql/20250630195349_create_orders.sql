-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS orders
(
    number   text PRIMARY KEY,
    status   text NOT NULL,
    accrual  DECIMAL(10,5) DEFAULT 0,
    uploaded_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    user_id  bigint NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS orders;
-- +goose StatementEnd
