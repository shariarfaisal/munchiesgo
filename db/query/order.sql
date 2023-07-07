
-- name: CreateOrder :one
INSERT INTO orders (customer_id, status, payment_method, payment_status, rider_note, dispatch_time, total, total_discount, service_charge, payable)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) RETURNING *;

-- name: GetOrder :one
SELECT * FROM orders WHERE id = $1;

-- name: ListOrder :many
SELECT * FROM orders ORDER BY id Limit $1 OFFSET $2;

-- name: ListOrderByCustomerId :many
SELECT * FROM orders WHERE customer_id = $1 ORDER BY id Limit $2 OFFSET $3;

-- name: CreateDeliveryAddress :one
INSERT INTO delivery_address (order_id, customer_id, address, geo_point, apartment, area, floor, phone)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING *;

-- name: GetDeliveryAddress :one
SELECT * FROM delivery_address WHERE order_id = $1;

-- name: CreateOrderItem :one
INSERT INTO order_items (order_id, product_id, brand_id, price, quantity, discount)
VALUES ($1, $2, $3, $4, $5, $6) RETURNING *;

-- name: GetOrderItem :one
SELECT * FROM order_items WHERE order_id = $1;

-- name: ListOrderItemsByOrderId :many
SELECT * FROM order_items WHERE order_id = $1;

-- name: ListOrderItemsByBrandId :many
SELECT * FROM order_items WHERE brand_id = $1 ORDER BY id Limit $2 OFFSET $3;

-- name: CreateBrandOrder :one
INSERT INTO brand_orders (brand_id, order_id, status, total, discount, note)
VALUES ($1, $2, $3, $4, $5, $6) RETURNING *;

-- name: GetBrandOrder :one
SELECT * FROM brand_orders WHERE order_id = $1;

-- name: ListBrandOrdersByOrderId :many
SELECT * FROM brand_orders WHERE order_id = $1;

-- name: ListBrandOrdersByBrandId :many
SELECT * FROM brand_orders WHERE brand_id = $1 ORDER BY id Limit $2 OFFSET $3;

-- name: CreateRiderAssign :one
INSERT INTO rider_assign (order_id, rider_id)
VALUES ($1, $2) RETURNING *;

-- name: GetRiderAssign :one
SELECT * FROM rider_assign WHERE order_id = $1;

-- name: ListRiderAssignByRiderId :many
SELECT * FROM rider_assign WHERE rider_id = $1 ORDER BY id Limit $2 OFFSET $3;