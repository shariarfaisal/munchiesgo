

-- name: CreateCategory :one
INSERT INTO categories (
  name,
  image
) VALUES ($1, $2)
RETURNING *;

-- name: UpdateCategory :one
UPDATE categories SET
  name = $2,
  image = $3
WHERE id = $1
RETURNING *;

-- name: DeleteCategory :exec
DELETE FROM categories WHERE id = $1;

-- name: GetCategory :one
SELECT * FROM categories WHERE id = $1;

-- name: ListCategories :many
SELECT * FROM categories ORDER BY id LIMIT $1 OFFSET $2;

-- name: SearchCategories :many
SELECT * FROM categories WHERE name ILIKE '%' || $1 || '%' ORDER BY id LIMIT $2 OFFSET $3;

-- name: CreateBrandCategory :one
INSERT INTO brand_categories (
  brand_id,
  category_id, 
  name
) VALUES ($1, $2, $3)
RETURNING *;

-- name: UpdateBrandCategory :one
UPDATE brand_categories SET
  name = $2
WHERE brand_id = $1 AND id = $3
RETURNING *;

-- name: GetBrandCategory :one
SELECT * FROM brand_categories WHERE id = $1;

-- name: SearchBrandCategories :many
SELECT * FROM brand_categories WHERE brand_id = $1 AND name ILIKE '%' || $2 || '%' ORDER BY id LIMIT $3 OFFSET $4;

-- name: DeleteBrandCategory :exec
DELETE FROM brand_categories WHERE brand_id = $1 AND id = $2;

-- name: ListBrandCategories :many
SELECT * FROM brand_categories WHERE brand_id = $1 ORDER BY id LIMIT $2 OFFSET $3;

-- name: CountCategories :one
SELECT COUNT(*) FROM categories;

-- name: CountCategoriesByBrandID :one
SELECT COUNT(*) FROM brand_categories WHERE brand_id = $1;