package db

import (
	"context"

	"github.com/lib/pq"
)

type GetBrandCategoriesByIdsRow struct {
	ID      int64  `json:"id"`
	Name    string `json:"name"`
	BrandID int64  `json:"brand_id"`
}

func (store *SqlStore) BrandCategoriesByIds(ctx context.Context, ids []int64) ([]GetBrandCategoriesByIdsRow, error) {
	const categoriesByIds = `
		SELECT id, name, brand_id FROM brand_categories WHERE id = ANY($1)
	`

	rows, err := store.db.QueryContext(ctx, categoriesByIds, pq.Array(ids))
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var items []GetBrandCategoriesByIdsRow
	for rows.Next() {
		var item GetBrandCategoriesByIdsRow
		err := rows.Scan(&item.ID, &item.Name, &item.BrandID)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	return items, nil
}
