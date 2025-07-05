-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS "order_queue"
(
    id             UUID        DEFAULT gen_random_uuid() NOT NULL PRIMARY KEY,
    created_at     TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL,
    started_at     TIMESTAMPTZ                           NULL,
    locked_until   TIMESTAMPTZ                           NULL,
    scheduled_for  TIMESTAMPTZ                           NULL,
    processed_at   TIMESTAMPTZ                           NULL,
    consumed_count INTEGER     DEFAULT 0                 NOT NULL,
    error_detail   TEXT                                  NULL,
    payload        JSONB                                 NOT NULL,
    metadata       JSONB                                 NULL
);
CREATE INDEX IF NOT EXISTS "order_queue_created_at_idx" ON "order_queue" (created_at);
CREATE INDEX IF NOT EXISTS "order_queue_processed_at_null_idx" ON "order_queue" (processed_at) WHERE (processed_at IS NULL);
CREATE INDEX IF NOT EXISTS "order_queue_scheduled_for_idx" ON "order_queue" (scheduled_for ASC NULLS LAST) WHERE (processed_at IS NULL);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS "order_queue"
-- +goose StatementEnd
