

-- name: CreateCustomer :one
INSERT INTO customers (name, phone, email, image, email_verified, nid, status) 
VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING *;

-- name: GetCustomer :one
SELECT * FROM customers WHERE id = $1;

-- name: CreateCustomerAddress :one
INSERT INTO customer_addresses (customer_id, label, address, geo_point, apartment, area, floor)
VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING *;

-- name: GetCustomerAddress :one
SELECT * FROM customer_addresses WHERE id = $1;

-- name: ListCustomerAddresses :many
SELECT * FROM customer_addresses WHERE customer_id = $1;