
-- name: CreateProductInventory :one
INSERT INTO product_inventory (
    product_id,
    quantity,
    purchase_price,
    selling_price,
    expire_date
    ) VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: UpdateProductInventory :one
UPDATE product_inventory SET
    quantity = $2,
    purchase_price = $3,
    selling_price = $4,
    expire_date = $5
WHERE id = $1
RETURNING *;

-- name: DeleteProductInventory :exec
DELETE FROM product_inventory WHERE id = $1;

-- name: GetProductInventory :one
SELECT * FROM product_inventory WHERE id = $1;

-- name: ListProductInventory :many
SELECT * FROM product_inventory WHERE product_id = $1 ORDER BY id LIMIT $2 OFFSET $3;

-- name: CountProductInventory :one
SELECT COUNT(*) FROM product_inventory WHERE product_id = $1;

-- name: GetProductInventoryByProductID :one
SELECT * FROM product_inventory WHERE product_id = $1;

-- name: CountProductInventoryByProductId :one
SELECT COUNT(*) FROM product_inventory WHERE product_id = $1;

-- name: SearchProductInventory :many
SELECT * FROM product_inventory
WHERE product_id = $1 AND (quantity ILIKE $2 OR purchase_price ILIKE $2 OR selling_price ILIKE $2 OR expire_date ILIKE $2) ORDER BY id LIMIT $3 OFFSET $4;

-- name: CreateInventoryHistory :one
INSERT INTO inventory_history (
    product_id,
    quantity,
    type
    ) VALUES ($1, $2, $3)
RETURNING *;

-- name: ListInventoryHistory :many
SELECT * FROM inventory_history WHERE product_id = $1 ORDER BY id LIMIT $2 OFFSET $3;

-- name: CountInventoryHistory :one
SELECT COUNT(*) FROM inventory_history WHERE product_id = $1;

-- name: SearchInventoryHistory :many
SELECT * FROM inventory_history
WHERE product_id = $1 AND (quantity ILIKE $2 OR type ILIKE $2) ORDER BY id LIMIT $3 OFFSET $4;

-- name: GetInventoryHistory :one
SELECT * FROM inventory_history WHERE id = $1;

-- name: DeleteInventoryHistory :exec
DELETE FROM inventory_history WHERE id = $1;

-- name: UpdateInventoryHistory :one
UPDATE inventory_history SET
    product_id = $2,
    quantity = $3,
    type = $4
WHERE id = $1
RETURNING *;

-- name: GetInventoryHistoryByProductID :one
SELECT * FROM inventory_history WHERE product_id = $1;

-- name: CountInventoryHistoryByProductId :one
SELECT COUNT(*) FROM inventory_history WHERE product_id = $1;