// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.17.0
// source: category.sql

package db

import (
	"context"
	"database/sql"
)

const countCategories = `-- name: CountCategories :one
SELECT COUNT(*) FROM categories
`

func (q *Queries) CountCategories(ctx context.Context) (int64, error) {
	row := q.queryRow(ctx, q.countCategoriesStmt, countCategories)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const countCategoriesByBrandID = `-- name: CountCategoriesByBrandID :one
SELECT COUNT(*) FROM brand_categories WHERE brand_id = $1
`

func (q *Queries) CountCategoriesByBrandID(ctx context.Context, brandID int64) (int64, error) {
	row := q.queryRow(ctx, q.countCategoriesByBrandIDStmt, countCategoriesByBrandID, brandID)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const createBrandCategory = `-- name: CreateBrandCategory :one
INSERT INTO brand_categories (
  brand_id,
  name,
  image
) VALUES ($1, $2, $3)
RETURNING id, brand_id, name, image, created_at
`

type CreateBrandCategoryParams struct {
	BrandID int64  `json:"brand_id"`
	Name    string `json:"name"`
	Image   string `json:"image"`
}

func (q *Queries) CreateBrandCategory(ctx context.Context, arg CreateBrandCategoryParams) (BrandCategory, error) {
	row := q.queryRow(ctx, q.createBrandCategoryStmt, createBrandCategory, arg.BrandID, arg.Name, arg.Image)
	var i BrandCategory
	err := row.Scan(
		&i.ID,
		&i.BrandID,
		&i.Name,
		&i.Image,
		&i.CreatedAt,
	)
	return i, err
}

const createCategory = `-- name: CreateCategory :one
INSERT INTO categories (
  name,
  image
) VALUES ($1, $2)
RETURNING id, name, image, created_at
`

type CreateCategoryParams struct {
	Name  string `json:"name"`
	Image string `json:"image"`
}

func (q *Queries) CreateCategory(ctx context.Context, arg CreateCategoryParams) (Category, error) {
	row := q.queryRow(ctx, q.createCategoryStmt, createCategory, arg.Name, arg.Image)
	var i Category
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Image,
		&i.CreatedAt,
	)
	return i, err
}

const deleteBrandCategory = `-- name: DeleteBrandCategory :exec
DELETE FROM brand_categories WHERE brand_id = $1 AND id = $2
`

type DeleteBrandCategoryParams struct {
	BrandID int64 `json:"brand_id"`
	ID      int64 `json:"id"`
}

func (q *Queries) DeleteBrandCategory(ctx context.Context, arg DeleteBrandCategoryParams) error {
	_, err := q.exec(ctx, q.deleteBrandCategoryStmt, deleteBrandCategory, arg.BrandID, arg.ID)
	return err
}

const deleteCategory = `-- name: DeleteCategory :exec
DELETE FROM categories WHERE id = $1
`

func (q *Queries) DeleteCategory(ctx context.Context, id int64) error {
	_, err := q.exec(ctx, q.deleteCategoryStmt, deleteCategory, id)
	return err
}

const getBrandCategory = `-- name: GetBrandCategory :one
SELECT id, brand_id, name, image, created_at FROM brand_categories WHERE brand_id = $1 AND id = $2
`

type GetBrandCategoryParams struct {
	BrandID int64 `json:"brand_id"`
	ID      int64 `json:"id"`
}

func (q *Queries) GetBrandCategory(ctx context.Context, arg GetBrandCategoryParams) (BrandCategory, error) {
	row := q.queryRow(ctx, q.getBrandCategoryStmt, getBrandCategory, arg.BrandID, arg.ID)
	var i BrandCategory
	err := row.Scan(
		&i.ID,
		&i.BrandID,
		&i.Name,
		&i.Image,
		&i.CreatedAt,
	)
	return i, err
}

const getCategory = `-- name: GetCategory :one
SELECT id, name, image, created_at FROM categories WHERE id = $1
`

func (q *Queries) GetCategory(ctx context.Context, id int64) (Category, error) {
	row := q.queryRow(ctx, q.getCategoryStmt, getCategory, id)
	var i Category
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Image,
		&i.CreatedAt,
	)
	return i, err
}

const listBrandCategories = `-- name: ListBrandCategories :many
SELECT id, brand_id, name, image, created_at FROM brand_categories WHERE brand_id = $1 ORDER BY id LIMIT $2 OFFSET $3
`

type ListBrandCategoriesParams struct {
	BrandID int64 `json:"brand_id"`
	Limit   int32 `json:"limit"`
	Offset  int32 `json:"offset"`
}

func (q *Queries) ListBrandCategories(ctx context.Context, arg ListBrandCategoriesParams) ([]BrandCategory, error) {
	rows, err := q.query(ctx, q.listBrandCategoriesStmt, listBrandCategories, arg.BrandID, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []BrandCategory{}
	for rows.Next() {
		var i BrandCategory
		if err := rows.Scan(
			&i.ID,
			&i.BrandID,
			&i.Name,
			&i.Image,
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

const listCategories = `-- name: ListCategories :many
SELECT id, name, image, created_at FROM categories ORDER BY id LIMIT $1 OFFSET $2
`

type ListCategoriesParams struct {
	Limit  int32 `json:"limit"`
	Offset int32 `json:"offset"`
}

func (q *Queries) ListCategories(ctx context.Context, arg ListCategoriesParams) ([]Category, error) {
	rows, err := q.query(ctx, q.listCategoriesStmt, listCategories, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Category{}
	for rows.Next() {
		var i Category
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Image,
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

const searchBrandCategories = `-- name: SearchBrandCategories :many
SELECT id, brand_id, name, image, created_at FROM brand_categories WHERE brand_id = $1 AND name ILIKE '%' || $2 || '%' ORDER BY id LIMIT $3 OFFSET $4
`

type SearchBrandCategoriesParams struct {
	BrandID int64          `json:"brand_id"`
	Column2 sql.NullString `json:"column_2"`
	Limit   int32          `json:"limit"`
	Offset  int32          `json:"offset"`
}

func (q *Queries) SearchBrandCategories(ctx context.Context, arg SearchBrandCategoriesParams) ([]BrandCategory, error) {
	rows, err := q.query(ctx, q.searchBrandCategoriesStmt, searchBrandCategories,
		arg.BrandID,
		arg.Column2,
		arg.Limit,
		arg.Offset,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []BrandCategory{}
	for rows.Next() {
		var i BrandCategory
		if err := rows.Scan(
			&i.ID,
			&i.BrandID,
			&i.Name,
			&i.Image,
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

const searchCategories = `-- name: SearchCategories :many
SELECT id, name, image, created_at FROM categories WHERE name ILIKE '%' || $1 || '%' ORDER BY id LIMIT $2 OFFSET $3
`

type SearchCategoriesParams struct {
	Column1 sql.NullString `json:"column_1"`
	Limit   int32          `json:"limit"`
	Offset  int32          `json:"offset"`
}

func (q *Queries) SearchCategories(ctx context.Context, arg SearchCategoriesParams) ([]Category, error) {
	rows, err := q.query(ctx, q.searchCategoriesStmt, searchCategories, arg.Column1, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Category{}
	for rows.Next() {
		var i Category
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Image,
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

const updateBrandCategory = `-- name: UpdateBrandCategory :one
UPDATE brand_categories SET
  name = $2,
  image = $3
WHERE brand_id = $1 AND id = $4
RETURNING id, brand_id, name, image, created_at
`

type UpdateBrandCategoryParams struct {
	BrandID int64  `json:"brand_id"`
	Name    string `json:"name"`
	Image   string `json:"image"`
	ID      int64  `json:"id"`
}

func (q *Queries) UpdateBrandCategory(ctx context.Context, arg UpdateBrandCategoryParams) (BrandCategory, error) {
	row := q.queryRow(ctx, q.updateBrandCategoryStmt, updateBrandCategory,
		arg.BrandID,
		arg.Name,
		arg.Image,
		arg.ID,
	)
	var i BrandCategory
	err := row.Scan(
		&i.ID,
		&i.BrandID,
		&i.Name,
		&i.Image,
		&i.CreatedAt,
	)
	return i, err
}

const updateCategory = `-- name: UpdateCategory :one
UPDATE categories SET
  name = $2,
  image = $3
WHERE id = $1
RETURNING id, name, image, created_at
`

type UpdateCategoryParams struct {
	ID    int64  `json:"id"`
	Name  string `json:"name"`
	Image string `json:"image"`
}

func (q *Queries) UpdateCategory(ctx context.Context, arg UpdateCategoryParams) (Category, error) {
	row := q.queryRow(ctx, q.updateCategoryStmt, updateCategory, arg.ID, arg.Name, arg.Image)
	var i Category
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Image,
		&i.CreatedAt,
	)
	return i, err
}
