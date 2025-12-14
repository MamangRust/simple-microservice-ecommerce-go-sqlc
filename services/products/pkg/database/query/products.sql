-- name: GetProducts :many
SELECT *, COUNT(*) OVER () AS total_count
FROM products AS p
WHERE
    p.deleted_at IS NULL
    AND (
        $1::TEXT IS NULL
        OR p.name ILIKE '%' || $1 || '%'
    )
ORDER BY p.created_at DESC
LIMIT $2
OFFSET
    $3;

-- name: GetActiveProducts :many
SELECT *, COUNT(*) OVER () AS total_count
FROM products AS p
WHERE
    p.deleted_at IS NULL
    AND (
        $1::TEXT IS NULL
        OR p.name ILIKE '%' || $1 || '%'
    )
ORDER BY p.created_at DESC
LIMIT $2
OFFSET
    $3;

-- name: GetTrashedProducts :many
SELECT *, COUNT(*) OVER () AS total_count
FROM products AS p
WHERE
    p.deleted_at IS NOT NULL
    AND (
        $1::TEXT IS NULL
        OR p.name ILIKE '%' || $1 || '%'
    )
ORDER BY p.created_at DESC
LIMIT $2
OFFSET
    $3;

-- name: CreateProduct :one
INSERT INTO
    products (name, price, stock)
VALUES ($1, $2, $3)
RETURNING
    *;

-- name: GetProductByID :one
SELECT *
FROM products AS p
WHERE
    p.product_id = $1
    AND p.deleted_at IS NULL;

-- name: GetProductByIdTrashed :one
SELECT * FROM products AS p WHERE p.product_id = $1;

-- name: UpdateProduct :one
UPDATE products
SET
    name = $2,
    price = $3,
    stock = $4,
    updated_at = CURRENT_TIMESTAMP
WHERE
    product_id = $1
    AND deleted_at IS NULL
RETURNING
    *;

-- name: UpdateProductCountStock :one
UPDATE products
SET
    stock = $2
WHERE
    product_id = $1
    AND deleted_at IS NULL
RETURNING
    *;

-- name: TrashProduct :one
UPDATE products
SET
    deleted_at = current_timestamp
WHERE
    product_id = $1
    AND deleted_at IS NULL
RETURNING
    *;

-- name: RestoreProduct :one
UPDATE products
SET
    deleted_at = NULL
WHERE
    product_id = $1
    AND deleted_at IS NOT NULL
RETURNING
    *;

-- name: DeleteProductPermanently :exec
DELETE FROM products
WHERE
    product_id = $1
    AND deleted_at IS NOT NULL;

-- name: RestoreAllProducts :exec
UPDATE products
SET
    deleted_at = NULL
WHERE
    deleted_at IS NOT NULL;

-- name: DeleteAllPermanentProducts :exec
DELETE FROM products WHERE deleted_at IS NOT NULL;