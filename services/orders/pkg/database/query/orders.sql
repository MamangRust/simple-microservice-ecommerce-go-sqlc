-- GetOrders: Retrieves paginated list of active orders with search capability
-- Purpose: List all active orders for management UI
-- Parameters:
--   $1: search_term - Optional text to filter orders by ID or total price (NULL for no filter)
--   $2: limit - Maximum number of records to return (pagination limit)
--   $3: offset - Number of records to skip (pagination offset)
-- Returns:
--   All order fields plus total_count of matching records
-- Business Logic:
--   - Excludes soft-deleted orders (deleted_at IS NULL)
--   - Supports partial text matching on order_id and total_price fields (case-insensitive)
--   - Returns newest orders first (created_at DESC)
--   - Provides total_count for client-side pagination
--   - Uses window function COUNT(*) OVER() for efficient total count
-- name: GetOrders :many
SELECT *, COUNT(*) OVER () AS total_count
FROM orders
WHERE
    deleted_at IS NULL
    AND (
        $1::TEXT IS NULL
        OR order_id::TEXT ILIKE '%' || $1 || '%'
        OR total_price::TEXT ILIKE '%' || $1 || '%'
    )
ORDER BY created_at DESC
LIMIT $2
OFFSET
    $3;

-- GetOrdersActive: Retrieves paginated list of active orders (identical to GetOrders)
-- Purpose: Maintains consistent API pattern with other active/trashed endpoints
-- Parameters:
--   $1: search_term - Optional filter text for order ID or price
--   $2: limit - Pagination limit
--   $3: offset - Pagination offset
-- Returns:
--   Active order records with total_count
-- Business Logic:
--   - Same functionality as GetOrders
--   - Exists for consistency in API design patterns
-- Note: Could be consolidated with GetOrders if duplicate functionality is undesired
-- name: GetOrdersActive :many
SELECT *, COUNT(*) OVER () AS total_count
FROM orders
WHERE
    deleted_at IS NULL
    AND (
        $1::TEXT IS NULL
        OR order_id::TEXT ILIKE '%' || $1 || '%'
        OR total_price::TEXT ILIKE '%' || $1 || '%'
    )
ORDER BY created_at DESC
LIMIT $2
OFFSET
    $3;

-- GetOrdersTrashed: Retrieves paginated list of soft-deleted orders
-- Purpose: View and manage deleted orders for potential restoration
-- Parameters:
--   $1: search_term - Optional text to filter trashed orders
--   $2: limit - Maximum records per page
--   $3: offset - Records to skip
-- Returns:
--   Trashed order records with total_count
-- Business Logic:
--   - Only returns soft-deleted records (deleted_at IS NOT NULL)
--   - Maintains same search functionality as active order queries
--   - Preserves chronological sorting (newest first)
--   - Used in order recovery/audit interfaces
--   - Includes total_count for pagination in trash management UI
-- name: GetOrdersTrashed :many
SELECT *, COUNT(*) OVER () AS total_count
FROM orders
WHERE
    deleted_at IS NOT NULL
    AND (
        $1::TEXT IS NULL
        OR order_id::TEXT ILIKE '%' || $1 || '%'
        OR total_price::TEXT ILIKE '%' || $1 || '%'
    )
ORDER BY created_at DESC
LIMIT $2
OFFSET
    $3;

-- CreateOrder: Creates a new order record
-- Purpose: Register a new transaction in the system
-- Parameters:
--   $2: user_id - ID of the user processing the order
--   $3: total_price - Numeric total amount of the order
-- Returns: The complete created order record
-- Business Logic:
--   - Automatically sets created_at timestamp
--   - Requires merchant_id, user_id and total_price
--   - Typically followed by order item creation
-- name: CreateOrder :one
INSERT INTO
    orders (user_id, total_price)
VALUES ($1, $2)
RETURNING
    *;

-- GetOrderByID: Retrieves an active order by ID
-- Purpose: Fetch order details for display/processing
-- Parameters:
--   $1: order_id - UUID of the order to retrieve
-- Returns: Full order record if found and active
-- Business Logic:
--   - Excludes soft-deleted orders
--   - Used for order viewing, receipts, and processing
--   - Typically joined with order_items in application
-- name: GetOrderByID :one
SELECT * FROM orders WHERE order_id = $1 AND deleted_at IS NULL;

-- UpdateOrder: Modifies order information
-- Purpose: Update order details (primarily total price)
-- Parameters:
--   $1: order_id - UUID of order to update
--   $2: total_price - New total amount
-- Returns: Updated order record
-- Business Logic:
--   - Auto-updates updated_at timestamp
--   - Only modifies active (non-deleted) orders
--   - Used when order items change
--   - Should trigger recalculation of total_price
-- name: UpdateOrder :one
UPDATE orders
SET
    total_price = $2,
    updated_at = CURRENT_TIMESTAMP
WHERE
    order_id = $1
    AND deleted_at IS NULL
RETURNING
    *;

-- TrashedOrder: Soft-deletes an order
-- Purpose: Cancel/void an order without permanent deletion
-- Parameters:
--   $1: order_id - UUID of order to cancel
-- Returns: The soft-deleted order record
-- Business Logic:
--   - Sets deleted_at to current timestamp
--   - Preserves order data for reporting
--   - Only processes active orders
--   - Can be restored via RestoreOrder
-- name: TrashedOrder :one
UPDATE orders
SET
    deleted_at = current_timestamp
WHERE
    order_id = $1
    AND deleted_at IS NULL
RETURNING
    *;

-- RestoreOrder: Recovers a soft-deleted order
-- Purpose: Reactivate a cancelled order
-- Parameters:
--   $1: order_id - UUID of order to restore
-- Returns: The restored order record
-- Business Logic:
--   - Nullifies deleted_at field
--   - Only works on previously cancelled orders
--   - Maintains all original order data
-- name: RestoreOrder :one
UPDATE orders
SET
    deleted_at = NULL
WHERE
    order_id = $1
    AND deleted_at IS NOT NULL
RETURNING
    *;

-- DeleteOrderPermanently: Hard-deletes an order
-- Purpose: Completely remove order from database
-- Parameters:
--   $1: order_id - UUID of order to delete
-- Business Logic:
--   - Permanent deletion of already cancelled orders
--   - No return value (exec-only operation)
--   - Irreversible action - use with caution
--   - Should trigger deletion of related order_items
-- name: DeleteOrderPermanently :exec
DELETE FROM orders WHERE order_id = $1 AND deleted_at IS NOT NULL;

-- RestoreAllOrders: Mass restoration of cancelled orders
-- Purpose: Recover all trashed orders at once
-- Business Logic:
--   - Reactivates all soft-deleted orders
--   - No parameters needed (bulk operation)
--   - Typically used during system recovery
-- name: RestoreAllOrders :exec
UPDATE orders
SET
    deleted_at = NULL
WHERE
    deleted_at IS NOT NULL;

-- DeleteAllPermanentOrders: Purges all cancelled orders
-- Purpose: Clean up all soft-deleted order records
-- Business Logic:
--   - Irreversible bulk deletion operation
--   - Only affects already cancelled orders
--   - Typically used during database maintenance
--   - Should be restricted to admin users
-- name: DeleteAllPermanentOrders :exec
DELETE FROM orders WHERE deleted_at IS NOT NULL;