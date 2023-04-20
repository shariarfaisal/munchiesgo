package db

import (
	"context"

	"github.com/lib/pq"
)

type BrandsWithIdsRow struct {
	ID       int64  `json:"id"`
	Name     string `json:"name"`
	VendorID int64  `json:"vendor_id"`
}

func (store *SqlStore) BrandsWithIds(ctx context.Context, ids []int64) ([]BrandsWithIdsRow, error) {
	const brandsWithIds = `SELECT id, name, vendor_id FROM brands WHERE id = ANY($1)`
	rows, err := store.db.QueryContext(ctx, brandsWithIds, pq.Array(ids))
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var items []BrandsWithIdsRow
	for rows.Next() {
		var item BrandsWithIdsRow
		err := rows.Scan(&item.ID, &item.Name, &item.VendorID)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	return items, nil
}
