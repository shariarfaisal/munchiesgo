package db

import "database/sql"

type Store interface {
	Querier
}

type SqlStore struct {
	*Queries
	db *sql.DB
}

func NewStore(db *sql.DB) *SqlStore {
	return &SqlStore{
		Queries: New(db),
		db:      db,
	}
}
