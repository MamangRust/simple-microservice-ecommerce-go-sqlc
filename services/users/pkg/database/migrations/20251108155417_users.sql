-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS "users" (
    "user_id" serial PRIMARY KEY,
    "firstname" VARCHAR(100) NOT NULL,
    "lastname" varchar(100) NOT NULL,
    "email" varchar(100) UNIQUE NOT NULL,
    "password" varchar(100) NOT NULL,
    "verification_code" varchar(100) NOT NULL,
    "is_verified" bool DEFAULT false,
    "created_at" timestamp DEFAULT current_timestamp,
    "updated_at" timestamp DEFAULT current_timestamp,
    "deleted_at" TIMESTAMP DEFAULT NULL
);

CREATE INDEX IF NOT EXISTS idx_users_active_email ON users (email)
WHERE
    deleted_at IS NULL;

CREATE INDEX IF NOT EXISTS idx_users_active_created_at ON users (created_at DESC)
WHERE
    deleted_at IS NULL;

CREATE INDEX IF NOT EXISTS idx_users_active_firstname_lastname ON users (firstname, lastname)
WHERE
    deleted_at IS NULL;

CREATE INDEX IF NOT EXISTS idx_users_trashed_created_at ON users (created_at DESC)
WHERE
    deleted_at IS NOT NULL;

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_users_active_email;

DROP INDEX IF EXISTS idx_users_active_created_at;

DROP INDEX IF EXISTS idx_users_active_firstname_lastname;

DROP INDEX IF EXISTS idx_users_trashed_created_at;

DROP TABLE IF EXISTS "users";

-- +goose StatementEnd