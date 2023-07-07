package db

import (
	"context"
	"database/sql"
	"fmt"
)

type Store interface {
	Querier
	BulkProductUpload(ctx context.Context, args []ProductUploadParams) ([]Product, error)
	FullProduct(ctx context.Context, id int64) (*FullProduct, error)
	BrandsWithIds(ctx context.Context, ids []int64) ([]Brand, error)
	BrandCategoriesByIds(ctx context.Context, ids []int64) ([]BrandCategory, error)
	ProductsByIds(ctx context.Context, ids []int64) ([]Product, error)
	ProductsTypeByIds(ctx context.Context, ids []int64) ([]ProductsTypeByIdsRow, error)
	ProductVariantsByProductIds(ctx context.Context, ids []int64) ([]ProductVariant, error)
	ProductVariantsItemsByProductIds(ctx context.Context, ids []int64) ([]ProductVariantItem, error)
	OrderItemProducts(ctx context.Context, orderItems []PlaceOrderItemParams) ([]*OrderItemProduct, error)
	ProductVariantsByIds(ctx context.Context, ids []int64) ([]ProductVariant, error)
	PlaceOrder(ctx context.Context, arg PlaceOrderParams) (*Order, error)
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

func (store *SqlStore) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	q := New(tx)
	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx error: %v, rb error: %v", err, rbErr)
		}
		return err
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}
