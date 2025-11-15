-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS "roles" (
    "role_id" SERIAL PRIMARY KEY,
    "name" VARCHAR(50) UNIQUE NOT NULL,
    "created_at" timestamp DEFAULT current_timestamp,
    "updated_at" timestamp DEFAULT current_timestamp,
    "deleted_at" TIMESTAMP DEFAULT NULL
);

CREATE INDEX "idx_roles_name" ON "roles" ("name");

CREATE INDEX "idx_roles_created_at" ON "roles" ("created_at");

CREATE INDEX "idx_roles_updated_at" ON "roles" ("updated_at");
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS "idx_roles_name";

DROP INDEX IF EXISTS "idx_roles_created_at";

DROP INDEX IF EXISTS "idx_roles_updated_at";

DROP TABLE IF EXISTS "roles";

-- +goose StatementEnd