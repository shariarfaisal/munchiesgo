

-- name: CreateProduct :one
INSERT INTO products (
  type,
  name,
  category_id,
  slug,
  image,
  details,
  price,
  status,
  brand_id,
  availability
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
RETURNING *;

-- name: UpdateProduct :one
UPDATE products SET
  type = $2,
  name = $3,
  category_id = $4,
  slug = $5,
  image = $6,
  details = $7,
  price = $8,
  status = $9,
  brand_id = $10,
  availability = $11
WHERE id = $1
RETURNING *;

-- name: DeleteProduct :exec
DELETE FROM products WHERE id = $1;

-- name: GetProduct :one
SELECT * FROM products WHERE id = $1;

-- name: GetProductBySlug :one
SELECT * FROM products WHERE slug = $1;

-- name: GetProductWithVariants :one
SELECT p.id AS product_id, p.name AS product_name, 
       v.id AS variant_id, v.title AS variant_title,
       json_agg(json_build_object('id', vi.id, 'name', vi.name)) AS variant_items
FROM products p
JOIN product_variants v ON v.product_id = p.id
JOIN product_variant_items vi ON vi.variant_id = v.id
WHERE p.id = $1
GROUP BY p.id, v.id;

-- name: ListProducts :many
SELECT * FROM products ORDER BY id LIMIT $1 OFFSET $2;

-- name: CountProducts :one
SELECT COUNT(*) FROM products;

-- name: SearchProducts :many
SELECT * FROM products
WHERE name ILIKE $1 ORDER BY id LIMIT $2 OFFSET $3;

-- name: ListProductsByCategoryID :many
SELECT * FROM products WHERE category_id = $1 ORDER BY id LIMIT $2 OFFSET $3;

-- name: CountProductsByCategoryID :one
SELECT COUNT(*) FROM products WHERE category_id = $1;

-- name: SearchProductsByCategoryID :many
SELECT * FROM products
WHERE category_id = $2 AND name ILIKE $1 ORDER BY id LIMIT $3 OFFSET $4;

-- name: ListProductsByBrandID :many
SELECT * FROM products WHERE brand_id = $1 ORDER BY id LIMIT $2 OFFSET $3;

-- name: CountProductsByBrandID :one
SELECT COUNT(*) FROM products WHERE brand_id = $1;

-- name: SearchProductsByBrandID :many  
SELECT * FROM products
WHERE brand_id = $2 AND name ILIKE $1 ORDER BY id LIMIT $3 OFFSET $4;

-- name: ListProductsByVendorID :many
SELECT products.* FROM vendors
INNER JOIN brands ON vendors.id = brands.vendor_id
INNER JOIN products ON brands.id = products.brand_id
WHERE vendors.id = $1 ORDER BY products.id LIMIT $2 OFFSET $3; 

-- name: CountProductsByVendorID :one
SELECT COUNT(*) FROM vendors
INNER JOIN brands ON vendors.id = brands.vendor_id
INNER JOIN products ON brands.id = products.brand_id
WHERE vendors.id = $1;

-- name: CreateProductVariant :one
INSERT INTO product_variants (
  product_id,
  title,
  min_select,
  max_select
) VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: UpdateProductVariant :one
UPDATE product_variants SET
  product_id = $2,
  title = $3,
  min_select = $4,
  max_select = $5
WHERE id = $1
RETURNING *;

-- name: DeleteProductVariant :exec
DELETE FROM product_variants WHERE id = $1;

-- name: GetProductVariant :one
SELECT * FROM product_variants WHERE id = $1;

-- name: ListProductVariants :many
SELECT * FROM product_variants WHERE product_id = $1 ORDER BY id LIMIT $2 OFFSET $3;

-- name: CountProductVariants :one
SELECT COUNT(*) FROM product_variants WHERE product_id = $1;

-- name: CreateProductVariantItem :one
INSERT INTO product_variant_items (
  variant_id,
  product_id
) VALUES ($1, $2)
RETURNING *;

-- name: DeleteProductVariantItem :exec
DELETE FROM product_variant_items WHERE id = $1;

-- name: GetProductVariantItem :one
SELECT * FROM product_variant_items WHERE id = $1;

-- name: ListProductVariantItems :many
SELECT * FROM product_variant_items WHERE variant_id = $1 ORDER BY id LIMIT $2 OFFSET $3;

-- name: CountProductVariantItems :one
SELECT COUNT(*) FROM product_variant_items WHERE variant_id = $1;

