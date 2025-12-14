-- +goose Up
-- +goose StatementBegin
CREATE TABLE "orders" (
    "order_id" SERIAL PRIMARY KEY,
    "user_id" INT NOT NULL,
    "total_price" INT NOT NULL,
    "created_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    "deleted_at" TIMESTAMP DEFAULT NULL
);

CREATE INDEX IF NOT EXISTS idx_orders_user_active_created_at ON orders (user_id, created_at DESC)
WHERE
    deleted_at IS NULL;

CREATE INDEX IF NOT EXISTS idx_orders_active_created_at ON orders (created_at DESC)
WHERE
    deleted_at IS NULL;

CREATE INDEX IF NOT EXISTS idx_orders_trashed_created_at ON orders (created_at DESC)
WHERE
    deleted_at IS NOT NULL;

CREATE INDEX IF NOT EXISTS idx_orders_active_total_price ON orders (total_price)
WHERE
    deleted_at IS NULL;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS "idx_orders_user_active_created_at";

DROP INDEX IF EXISTS "idx_orders_active_created_at";

DROP INDEX IF EXISTS "idx_orders_trashed_created_at";

DROP INDEX IF EXISTS "idx_orders_active_total_price";

DROP TABLE IF EXISTS "orders";
-- +goose StatementEnd