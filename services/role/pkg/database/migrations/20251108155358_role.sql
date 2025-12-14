-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS "roles" (
    "role_id" SERIAL PRIMARY KEY,
    "name" VARCHAR(50) UNIQUE NOT NULL,
    "created_at" timestamp DEFAULT current_timestamp,
    "updated_at" timestamp DEFAULT current_timestamp,
    "deleted_at" TIMESTAMP DEFAULT NULL
);

CREATE INDEX IF NOT EXISTS idx_roles_active_name ON roles (role_name)
WHERE
    deleted_at IS NULL;

CREATE INDEX IF NOT EXISTS idx_roles_active_created_at ON roles (created_at DESC)
WHERE
    deleted_at IS NULL;

CREATE INDEX IF NOT EXISTS idx_roles_trashed_created_at ON roles (created_at DESC)
WHERE
    deleted_at IS NOT NULL;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS "idx_roles_active_name";

DROP INDEX IF EXISTS "idx_roles_active_created_at";

DROP INDEX IF EXISTS "idx_roles_trashed_created_at";

DROP TABLE IF EXISTS "roles";

-- +goose StatementEnd