// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.17.0
// source: product.sql

package db

import (
	"context"
	"encoding/json"
)

const countProductVariantItems = `-- name: CountProductVariantItems :one
SELECT COUNT(*) FROM product_variant_items WHERE variant_id = $1
`

func (q *Queries) CountProductVariantItems(ctx context.Context, variantID int64) (int64, error) {
	row := q.queryRow(ctx, q.countProductVariantItemsStmt, countProductVariantItems, variantID)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const countProductVariants = `-- name: CountProductVariants :one
SELECT COUNT(*) FROM product_variants WHERE product_id = $1
`

func (q *Queries) CountProductVariants(ctx context.Context, productID int64) (int64, error) {
	row := q.queryRow(ctx, q.countProductVariantsStmt, countProductVariants, productID)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const countProducts = `-- name: CountProducts :one
SELECT COUNT(*) FROM products
`

func (q *Queries) CountProducts(ctx context.Context) (int64, error) {
	row := q.queryRow(ctx, q.countProductsStmt, countProducts)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const countProductsByBrandID = `-- name: CountProductsByBrandID :one
SELECT COUNT(*) FROM products WHERE brand_id = $1
`

func (q *Queries) CountProductsByBrandID(ctx context.Context, brandID int64) (int64, error) {
	row := q.queryRow(ctx, q.countProductsByBrandIDStmt, countProductsByBrandID, brandID)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const countProductsByCategoryID = `-- name: CountProductsByCategoryID :one
SELECT COUNT(*) FROM products WHERE category_id = $1
`

func (q *Queries) CountProductsByCategoryID(ctx context.Context, categoryID int64) (int64, error) {
	row := q.queryRow(ctx, q.countProductsByCategoryIDStmt, countProductsByCategoryID, categoryID)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const countProductsByVendorID = `-- name: CountProductsByVendorID :one
SELECT COUNT(*) FROM vendors
INNER JOIN brands ON vendors.id = brands.vendor_id
INNER JOIN products ON brands.id = products.brand_id
WHERE vendors.id = $1
`

func (q *Queries) CountProductsByVendorID(ctx context.Context, id int64) (int64, error) {
	row := q.queryRow(ctx, q.countProductsByVendorIDStmt, countProductsByVendorID, id)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const createProduct = `-- name: CreateProduct :one
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
RETURNING id, type, name, category_id, slug, image, details, price, status, brand_id, availability, use_inventory, created_at, updated_at
`

type CreateProductParams struct {
	Type         string  `json:"type"`
	Name         string  `json:"name"`
	CategoryID   int64   `json:"category_id"`
	Slug         string  `json:"slug"`
	Image        string  `json:"image"`
	Details      string  `json:"details"`
	Price        float64 `json:"price"`
	Status       string  `json:"status"`
	BrandID      int64   `json:"brand_id"`
	Availability bool    `json:"availability"`
}

func (q *Queries) CreateProduct(ctx context.Context, arg CreateProductParams) (Product, error) {
	row := q.queryRow(ctx, q.createProductStmt, createProduct,
		arg.Type,
		arg.Name,
		arg.CategoryID,
		arg.Slug,
		arg.Image,
		arg.Details,
		arg.Price,
		arg.Status,
		arg.BrandID,
		arg.Availability,
	)
	var i Product
	err := row.Scan(
		&i.ID,
		&i.Type,
		&i.Name,
		&i.CategoryID,
		&i.Slug,
		&i.Image,
		&i.Details,
		&i.Price,
		&i.Status,
		&i.BrandID,
		&i.Availability,
		&i.UseInventory,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const createProductVariant = `-- name: CreateProductVariant :one
INSERT INTO product_variants (
  product_id,
  title,
  min_select,
  max_select
) VALUES ($1, $2, $3, $4)
RETURNING id, product_id, title, min_select, max_select, created_at, updated_at
`

type CreateProductVariantParams struct {
	ProductID int64  `json:"product_id"`
	Title     string `json:"title"`
	MinSelect int32  `json:"min_select"`
	MaxSelect int32  `json:"max_select"`
}

func (q *Queries) CreateProductVariant(ctx context.Context, arg CreateProductVariantParams) (ProductVariant, error) {
	row := q.queryRow(ctx, q.createProductVariantStmt, createProductVariant,
		arg.ProductID,
		arg.Title,
		arg.MinSelect,
		arg.MaxSelect,
	)
	var i ProductVariant
	err := row.Scan(
		&i.ID,
		&i.ProductID,
		&i.Title,
		&i.MinSelect,
		&i.MaxSelect,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const createProductVariantItem = `-- name: CreateProductVariantItem :one
INSERT INTO product_variant_items (
  variant_id,
  product_id
) VALUES ($1, $2)
RETURNING id, variant_id, product_id, created_at
`

type CreateProductVariantItemParams struct {
	VariantID int64 `json:"variant_id"`
	ProductID int64 `json:"product_id"`
}

func (q *Queries) CreateProductVariantItem(ctx context.Context, arg CreateProductVariantItemParams) (ProductVariantItem, error) {
	row := q.queryRow(ctx, q.createProductVariantItemStmt, createProductVariantItem, arg.VariantID, arg.ProductID)
	var i ProductVariantItem
	err := row.Scan(
		&i.ID,
		&i.VariantID,
		&i.ProductID,
		&i.CreatedAt,
	)
	return i, err
}

const deleteProduct = `-- name: DeleteProduct :exec
DELETE FROM products WHERE id = $1
`

func (q *Queries) DeleteProduct(ctx context.Context, id int64) error {
	_, err := q.exec(ctx, q.deleteProductStmt, deleteProduct, id)
	return err
}

const deleteProductVariant = `-- name: DeleteProductVariant :exec
DELETE FROM product_variants WHERE id = $1
`

func (q *Queries) DeleteProductVariant(ctx context.Context, id int64) error {
	_, err := q.exec(ctx, q.deleteProductVariantStmt, deleteProductVariant, id)
	return err
}

const deleteProductVariantItem = `-- name: DeleteProductVariantItem :exec
DELETE FROM product_variant_items WHERE id = $1
`

func (q *Queries) DeleteProductVariantItem(ctx context.Context, id int64) error {
	_, err := q.exec(ctx, q.deleteProductVariantItemStmt, deleteProductVariantItem, id)
	return err
}

const getProduct = `-- name: GetProduct :one
SELECT id, type, name, category_id, slug, image, details, price, status, brand_id, availability, use_inventory, created_at, updated_at FROM products WHERE id = $1
`

func (q *Queries) GetProduct(ctx context.Context, id int64) (Product, error) {
	row := q.queryRow(ctx, q.getProductStmt, getProduct, id)
	var i Product
	err := row.Scan(
		&i.ID,
		&i.Type,
		&i.Name,
		&i.CategoryID,
		&i.Slug,
		&i.Image,
		&i.Details,
		&i.Price,
		&i.Status,
		&i.BrandID,
		&i.Availability,
		&i.UseInventory,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getProductBySlug = `-- name: GetProductBySlug :one
SELECT id, type, name, category_id, slug, image, details, price, status, brand_id, availability, use_inventory, created_at, updated_at FROM products WHERE slug = $1
`

func (q *Queries) GetProductBySlug(ctx context.Context, slug string) (Product, error) {
	row := q.queryRow(ctx, q.getProductBySlugStmt, getProductBySlug, slug)
	var i Product
	err := row.Scan(
		&i.ID,
		&i.Type,
		&i.Name,
		&i.CategoryID,
		&i.Slug,
		&i.Image,
		&i.Details,
		&i.Price,
		&i.Status,
		&i.BrandID,
		&i.Availability,
		&i.UseInventory,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getProductVariant = `-- name: GetProductVariant :one
SELECT id, product_id, title, min_select, max_select, created_at, updated_at FROM product_variants WHERE id = $1
`

func (q *Queries) GetProductVariant(ctx context.Context, id int64) (ProductVariant, error) {
	row := q.queryRow(ctx, q.getProductVariantStmt, getProductVariant, id)
	var i ProductVariant
	err := row.Scan(
		&i.ID,
		&i.ProductID,
		&i.Title,
		&i.MinSelect,
		&i.MaxSelect,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getProductVariantItem = `-- name: GetProductVariantItem :one
SELECT id, variant_id, product_id, created_at FROM product_variant_items WHERE id = $1
`

func (q *Queries) GetProductVariantItem(ctx context.Context, id int64) (ProductVariantItem, error) {
	row := q.queryRow(ctx, q.getProductVariantItemStmt, getProductVariantItem, id)
	var i ProductVariantItem
	err := row.Scan(
		&i.ID,
		&i.VariantID,
		&i.ProductID,
		&i.CreatedAt,
	)
	return i, err
}

const getProductWithVariants = `-- name: GetProductWithVariants :one
SELECT p.id AS product_id, p.name AS product_name, 
       v.id AS variant_id, v.title AS variant_title,
       json_agg(json_build_object('id', vi.id, 'name', vi.name)) AS variant_items
FROM products p
JOIN product_variants v ON v.product_id = p.id
JOIN product_variant_items vi ON vi.variant_id = v.id
WHERE p.id = $1
GROUP BY p.id, v.id
`

type GetProductWithVariantsRow struct {
	ProductID    int64           `json:"product_id"`
	ProductName  string          `json:"product_name"`
	VariantID    int64           `json:"variant_id"`
	VariantTitle string          `json:"variant_title"`
	VariantItems json.RawMessage `json:"variant_items"`
}

func (q *Queries) GetProductWithVariants(ctx context.Context, id int64) (GetProductWithVariantsRow, error) {
	row := q.queryRow(ctx, q.getProductWithVariantsStmt, getProductWithVariants, id)
	var i GetProductWithVariantsRow
	err := row.Scan(
		&i.ProductID,
		&i.ProductName,
		&i.VariantID,
		&i.VariantTitle,
		&i.VariantItems,
	)
	return i, err
}

const listProductVariantItems = `-- name: ListProductVariantItems :many
SELECT id, variant_id, product_id, created_at FROM product_variant_items WHERE variant_id = $1 ORDER BY id LIMIT $2 OFFSET $3
`

type ListProductVariantItemsParams struct {
	VariantID int64 `json:"variant_id"`
	Limit     int32 `json:"limit"`
	Offset    int32 `json:"offset"`
}

func (q *Queries) ListProductVariantItems(ctx context.Context, arg ListProductVariantItemsParams) ([]ProductVariantItem, error) {
	rows, err := q.query(ctx, q.listProductVariantItemsStmt, listProductVariantItems, arg.VariantID, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []ProductVariantItem{}
	for rows.Next() {
		var i ProductVariantItem
		if err := rows.Scan(
			&i.ID,
			&i.VariantID,
			&i.ProductID,
			&i.CreatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listProductVariants = `-- name: ListProductVariants :many
SELECT id, product_id, title, min_select, max_select, created_at, updated_at FROM product_variants WHERE product_id = $1 ORDER BY id LIMIT $2 OFFSET $3
`

type ListProductVariantsParams struct {
	ProductID int64 `json:"product_id"`
	Limit     int32 `json:"limit"`
	Offset    int32 `json:"offset"`
}

func (q *Queries) ListProductVariants(ctx context.Context, arg ListProductVariantsParams) ([]ProductVariant, error) {
	rows, err := q.query(ctx, q.listProductVariantsStmt, listProductVariants, arg.ProductID, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []ProductVariant{}
	for rows.Next() {
		var i ProductVariant
		if err := rows.Scan(
			&i.ID,
			&i.ProductID,
			&i.Title,
			&i.MinSelect,
			&i.MaxSelect,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listProducts = `-- name: ListProducts :many
SELECT id, type, name, category_id, slug, image, details, price, status, brand_id, availability, use_inventory, created_at, updated_at FROM products ORDER BY id LIMIT $1 OFFSET $2
`

type ListProductsParams struct {
	Limit  int32 `json:"limit"`
	Offset int32 `json:"offset"`
}

func (q *Queries) ListProducts(ctx context.Context, arg ListProductsParams) ([]Product, error) {
	rows, err := q.query(ctx, q.listProductsStmt, listProducts, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Product{}
	for rows.Next() {
		var i Product
		if err := rows.Scan(
			&i.ID,
			&i.Type,
			&i.Name,
			&i.CategoryID,
			&i.Slug,
			&i.Image,
			&i.Details,
			&i.Price,
			&i.Status,
			&i.BrandID,
			&i.Availability,
			&i.UseInventory,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listProductsByBrandID = `-- name: ListProductsByBrandID :many
SELECT id, type, name, category_id, slug, image, details, price, status, brand_id, availability, use_inventory, created_at, updated_at FROM products WHERE brand_id = $1 ORDER BY id LIMIT $2 OFFSET $3
`

type ListProductsByBrandIDParams struct {
	BrandID int64 `json:"brand_id"`
	Limit   int32 `json:"limit"`
	Offset  int32 `json:"offset"`
}

func (q *Queries) ListProductsByBrandID(ctx context.Context, arg ListProductsByBrandIDParams) ([]Product, error) {
	rows, err := q.query(ctx, q.listProductsByBrandIDStmt, listProductsByBrandID, arg.BrandID, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Product{}
	for rows.Next() {
		var i Product
		if err := rows.Scan(
			&i.ID,
			&i.Type,
			&i.Name,
			&i.CategoryID,
			&i.Slug,
			&i.Image,
			&i.Details,
			&i.Price,
			&i.Status,
			&i.BrandID,
			&i.Availability,
			&i.UseInventory,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listProductsByCategoryID = `-- name: ListProductsByCategoryID :many
SELECT id, type, name, category_id, slug, image, details, price, status, brand_id, availability, use_inventory, created_at, updated_at FROM products WHERE category_id = $1 ORDER BY id LIMIT $2 OFFSET $3
`

type ListProductsByCategoryIDParams struct {
	CategoryID int64 `json:"category_id"`
	Limit      int32 `json:"limit"`
	Offset     int32 `json:"offset"`
}

func (q *Queries) ListProductsByCategoryID(ctx context.Context, arg ListProductsByCategoryIDParams) ([]Product, error) {
	rows, err := q.query(ctx, q.listProductsByCategoryIDStmt, listProductsByCategoryID, arg.CategoryID, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Product{}
	for rows.Next() {
		var i Product
		if err := rows.Scan(
			&i.ID,
			&i.Type,
			&i.Name,
			&i.CategoryID,
			&i.Slug,
			&i.Image,
			&i.Details,
			&i.Price,
			&i.Status,
			&i.BrandID,
			&i.Availability,
			&i.UseInventory,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listProductsByVendorID = `-- name: ListProductsByVendorID :many
SELECT products.id, products.type, products.name, products.category_id, products.slug, products.image, products.details, products.price, products.status, products.brand_id, products.availability, products.use_inventory, products.created_at, products.updated_at FROM vendors
INNER JOIN brands ON vendors.id = brands.vendor_id
INNER JOIN products ON brands.id = products.brand_id
WHERE vendors.id = $1 ORDER BY products.id LIMIT $2 OFFSET $3
`

type ListProductsByVendorIDParams struct {
	ID     int64 `json:"id"`
	Limit  int32 `json:"limit"`
	Offset int32 `json:"offset"`
}

func (q *Queries) ListProductsByVendorID(ctx context.Context, arg ListProductsByVendorIDParams) ([]Product, error) {
	rows, err := q.query(ctx, q.listProductsByVendorIDStmt, listProductsByVendorID, arg.ID, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Product{}
	for rows.Next() {
		var i Product
		if err := rows.Scan(
			&i.ID,
			&i.Type,
			&i.Name,
			&i.CategoryID,
			&i.Slug,
			&i.Image,
			&i.Details,
			&i.Price,
			&i.Status,
			&i.BrandID,
			&i.Availability,
			&i.UseInventory,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const searchProducts = `-- name: SearchProducts :many
SELECT id, type, name, category_id, slug, image, details, price, status, brand_id, availability, use_inventory, created_at, updated_at FROM products
WHERE name ILIKE $1 ORDER BY id LIMIT $2 OFFSET $3
`

type SearchProductsParams struct {
	Name   string `json:"name"`
	Limit  int32  `json:"limit"`
	Offset int32  `json:"offset"`
}

func (q *Queries) SearchProducts(ctx context.Context, arg SearchProductsParams) ([]Product, error) {
	rows, err := q.query(ctx, q.searchProductsStmt, searchProducts, arg.Name, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Product{}
	for rows.Next() {
		var i Product
		if err := rows.Scan(
			&i.ID,
			&i.Type,
			&i.Name,
			&i.CategoryID,
			&i.Slug,
			&i.Image,
			&i.Details,
			&i.Price,
			&i.Status,
			&i.BrandID,
			&i.Availability,
			&i.UseInventory,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const searchProductsByBrandID = `-- name: SearchProductsByBrandID :many
SELECT id, type, name, category_id, slug, image, details, price, status, brand_id, availability, use_inventory, created_at, updated_at FROM products
WHERE brand_id = $2 AND name ILIKE $1 ORDER BY id LIMIT $3 OFFSET $4
`

type SearchProductsByBrandIDParams struct {
	Name    string `json:"name"`
	BrandID int64  `json:"brand_id"`
	Limit   int32  `json:"limit"`
	Offset  int32  `json:"offset"`
}

func (q *Queries) SearchProductsByBrandID(ctx context.Context, arg SearchProductsByBrandIDParams) ([]Product, error) {
	rows, err := q.query(ctx, q.searchProductsByBrandIDStmt, searchProductsByBrandID,
		arg.Name,
		arg.BrandID,
		arg.Limit,
		arg.Offset,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Product{}
	for rows.Next() {
		var i Product
		if err := rows.Scan(
			&i.ID,
			&i.Type,
			&i.Name,
			&i.CategoryID,
			&i.Slug,
			&i.Image,
			&i.Details,
			&i.Price,
			&i.Status,
			&i.BrandID,
			&i.Availability,
			&i.UseInventory,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const searchProductsByCategoryID = `-- name: SearchProductsByCategoryID :many
SELECT id, type, name, category_id, slug, image, details, price, status, brand_id, availability, use_inventory, created_at, updated_at FROM products
WHERE category_id = $2 AND name ILIKE $1 ORDER BY id LIMIT $3 OFFSET $4
`

type SearchProductsByCategoryIDParams struct {
	Name       string `json:"name"`
	CategoryID int64  `json:"category_id"`
	Limit      int32  `json:"limit"`
	Offset     int32  `json:"offset"`
}

func (q *Queries) SearchProductsByCategoryID(ctx context.Context, arg SearchProductsByCategoryIDParams) ([]Product, error) {
	rows, err := q.query(ctx, q.searchProductsByCategoryIDStmt, searchProductsByCategoryID,
		arg.Name,
		arg.CategoryID,
		arg.Limit,
		arg.Offset,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Product{}
	for rows.Next() {
		var i Product
		if err := rows.Scan(
			&i.ID,
			&i.Type,
			&i.Name,
			&i.CategoryID,
			&i.Slug,
			&i.Image,
			&i.Details,
			&i.Price,
			&i.Status,
			&i.BrandID,
			&i.Availability,
			&i.UseInventory,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateProduct = `-- name: UpdateProduct :one
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
RETURNING id, type, name, category_id, slug, image, details, price, status, brand_id, availability, use_inventory, created_at, updated_at
`

type UpdateProductParams struct {
	ID           int64   `json:"id"`
	Type         string  `json:"type"`
	Name         string  `json:"name"`
	CategoryID   int64   `json:"category_id"`
	Slug         string  `json:"slug"`
	Image        string  `json:"image"`
	Details      string  `json:"details"`
	Price        float64 `json:"price"`
	Status       string  `json:"status"`
	BrandID      int64   `json:"brand_id"`
	Availability bool    `json:"availability"`
}

func (q *Queries) UpdateProduct(ctx context.Context, arg UpdateProductParams) (Product, error) {
	row := q.queryRow(ctx, q.updateProductStmt, updateProduct,
		arg.ID,
		arg.Type,
		arg.Name,
		arg.CategoryID,
		arg.Slug,
		arg.Image,
		arg.Details,
		arg.Price,
		arg.Status,
		arg.BrandID,
		arg.Availability,
	)
	var i Product
	err := row.Scan(
		&i.ID,
		&i.Type,
		&i.Name,
		&i.CategoryID,
		&i.Slug,
		&i.Image,
		&i.Details,
		&i.Price,
		&i.Status,
		&i.BrandID,
		&i.Availability,
		&i.UseInventory,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const updateProductVariant = `-- name: UpdateProductVariant :one
UPDATE product_variants SET
  product_id = $2,
  title = $3,
  min_select = $4,
  max_select = $5
WHERE id = $1
RETURNING id, product_id, title, min_select, max_select, created_at, updated_at
`

type UpdateProductVariantParams struct {
	ID        int64  `json:"id"`
	ProductID int64  `json:"product_id"`
	Title     string `json:"title"`
	MinSelect int32  `json:"min_select"`
	MaxSelect int32  `json:"max_select"`
}

func (q *Queries) UpdateProductVariant(ctx context.Context, arg UpdateProductVariantParams) (ProductVariant, error) {
	row := q.queryRow(ctx, q.updateProductVariantStmt, updateProductVariant,
		arg.ID,
		arg.ProductID,
		arg.Title,
		arg.MinSelect,
		arg.MaxSelect,
	)
	var i ProductVariant
	err := row.Scan(
		&i.ID,
		&i.ProductID,
		&i.Title,
		&i.MinSelect,
		&i.MaxSelect,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}
