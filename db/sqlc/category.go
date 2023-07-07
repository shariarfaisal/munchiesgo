package db

import (
	"context"

	"github.com/lib/pq"
)

func (store *SqlStore) BrandCategoriesByIds(ctx context.Context, ids []int64) ([]BrandCategory, error) {
	const categoriesByIds = `
		SELECT * FROM brand_categories WHERE id = ANY($1)
	`

	rows, err := store.db.QueryContext(ctx, categoriesByIds, pq.Array(ids))
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var items []BrandCategory
	for rows.Next() {
		var item BrandCategory
		err := rows.Scan(
			&item.ID,
			&item.BrandID,
			&item.CategoryID,
			&item.Name,
			&item.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	return items, nil
}
