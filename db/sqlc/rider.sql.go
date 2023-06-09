// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.17.0
// source: rider.sql

package db

import (
	"context"
)

const createRider = `-- name: CreateRider :one
INSERT INTO riders (name, phone, hashed_password)
VALUES ($1, $2, $3) RETURNING id, name, phone, hashed_password
`

type CreateRiderParams struct {
	Name           string `json:"name"`
	Phone          string `json:"phone"`
	HashedPassword string `json:"hashed_password"`
}

func (q *Queries) CreateRider(ctx context.Context, arg CreateRiderParams) (Rider, error) {
	row := q.queryRow(ctx, q.createRiderStmt, createRider, arg.Name, arg.Phone, arg.HashedPassword)
	var i Rider
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Phone,
		&i.HashedPassword,
	)
	return i, err
}

const getRider = `-- name: GetRider :one
SELECT id, name, phone, hashed_password FROM riders WHERE id = $1
`

func (q *Queries) GetRider(ctx context.Context, id int64) (Rider, error) {
	row := q.queryRow(ctx, q.getRiderStmt, getRider, id)
	var i Rider
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Phone,
		&i.HashedPassword,
	)
	return i, err
}

const listRider = `-- name: ListRider :many
SELECT id, name, phone, hashed_password FROM riders ORDER BY id Limit $1 OFFSET $2
`

type ListRiderParams struct {
	Limit  int32 `json:"limit"`
	Offset int32 `json:"offset"`
}

func (q *Queries) ListRider(ctx context.Context, arg ListRiderParams) ([]Rider, error) {
	rows, err := q.query(ctx, q.listRiderStmt, listRider, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Rider{}
	for rows.Next() {
		var i Rider
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Phone,
			&i.HashedPassword,
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
