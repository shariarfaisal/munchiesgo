// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.17.0
// source: vendor.sql

package db

import (
	"context"
	"encoding/json"
)

const countVendors = `-- name: CountVendors :one
SELECT COUNT(*) FROM vendors
`

func (q *Queries) CountVendors(ctx context.Context) (int64, error) {
	row := q.queryRow(ctx, q.countVendorsStmt, countVendors)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const createVendor = `-- name: CreateVendor :one
INSERT INTO vendors (
  name,
  email,
  phone,
  payment_info,
  social_links
) VALUES ($1, $2, $3, $4, $5)
RETURNING id, name, email, phone, payment_info, social_links, created_at
`

type CreateVendorParams struct {
	Name        string          `json:"name"`
	Email       string          `json:"email"`
	Phone       string          `json:"phone"`
	PaymentInfo json.RawMessage `json:"payment_info"`
	SocialLinks json.RawMessage `json:"social_links"`
}

func (q *Queries) CreateVendor(ctx context.Context, arg CreateVendorParams) (Vendor, error) {
	row := q.queryRow(ctx, q.createVendorStmt, createVendor,
		arg.Name,
		arg.Email,
		arg.Phone,
		arg.PaymentInfo,
		arg.SocialLinks,
	)
	var i Vendor
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Email,
		&i.Phone,
		&i.PaymentInfo,
		&i.SocialLinks,
		&i.CreatedAt,
	)
	return i, err
}

const deleteVendor = `-- name: DeleteVendor :exec
DELETE FROM vendors WHERE id = $1
`

func (q *Queries) DeleteVendor(ctx context.Context, id int64) error {
	_, err := q.exec(ctx, q.deleteVendorStmt, deleteVendor, id)
	return err
}

const getVendor = `-- name: GetVendor :one
SELECT id, name, email, phone, payment_info, social_links, created_at FROM vendors WHERE id = $1
`

func (q *Queries) GetVendor(ctx context.Context, id int64) (Vendor, error) {
	row := q.queryRow(ctx, q.getVendorStmt, getVendor, id)
	var i Vendor
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Email,
		&i.Phone,
		&i.PaymentInfo,
		&i.SocialLinks,
		&i.CreatedAt,
	)
	return i, err
}

const listVendors = `-- name: ListVendors :many
SELECT id, name, email, phone, payment_info, social_links, created_at FROM vendors ORDER BY id LIMIT $1 OFFSET $2
`

type ListVendorsParams struct {
	Limit  int32 `json:"limit"`
	Offset int32 `json:"offset"`
}

func (q *Queries) ListVendors(ctx context.Context, arg ListVendorsParams) ([]Vendor, error) {
	rows, err := q.query(ctx, q.listVendorsStmt, listVendors, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Vendor{}
	for rows.Next() {
		var i Vendor
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Email,
			&i.Phone,
			&i.PaymentInfo,
			&i.SocialLinks,
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

const searchVendors = `-- name: SearchVendors :many
SELECT id, name, email, phone, payment_info, social_links, created_at FROM vendors
WHERE name ILIKE $1 OR email ILIKE $1 OR phone ILIKE $1 ORDER BY id LIMIT $2 OFFSET $3
`

type SearchVendorsParams struct {
	Name   string `json:"name"`
	Limit  int32  `json:"limit"`
	Offset int32  `json:"offset"`
}

func (q *Queries) SearchVendors(ctx context.Context, arg SearchVendorsParams) ([]Vendor, error) {
	rows, err := q.query(ctx, q.searchVendorsStmt, searchVendors, arg.Name, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Vendor{}
	for rows.Next() {
		var i Vendor
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Email,
			&i.Phone,
			&i.PaymentInfo,
			&i.SocialLinks,
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

const updateVendor = `-- name: UpdateVendor :one
UPDATE vendors SET
  name = $2,
  email = $3,
  phone = $4,
  payment_info = $5,
  social_links = $6
WHERE id = $1
RETURNING id, name, email, phone, payment_info, social_links, created_at
`

type UpdateVendorParams struct {
	ID          int64           `json:"id"`
	Name        string          `json:"name"`
	Email       string          `json:"email"`
	Phone       string          `json:"phone"`
	PaymentInfo json.RawMessage `json:"payment_info"`
	SocialLinks json.RawMessage `json:"social_links"`
}

func (q *Queries) UpdateVendor(ctx context.Context, arg UpdateVendorParams) (Vendor, error) {
	row := q.queryRow(ctx, q.updateVendorStmt, updateVendor,
		arg.ID,
		arg.Name,
		arg.Email,
		arg.Phone,
		arg.PaymentInfo,
		arg.SocialLinks,
	)
	var i Vendor
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Email,
		&i.Phone,
		&i.PaymentInfo,
		&i.SocialLinks,
		&i.CreatedAt,
	)
	return i, err
}
