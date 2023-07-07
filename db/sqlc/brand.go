package db

import (
	"context"

	"github.com/lib/pq"
)

func (store *SqlStore) BrandsWithIds(ctx context.Context, ids []int64) ([]Brand, error) {
	const brandsWithIds = `SELECT * FROM brands WHERE id = ANY($1)`
	rows, err := store.db.QueryContext(ctx, brandsWithIds, pq.Array(ids))
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var items []Brand
	for rows.Next() {
		var item Brand
		err := rows.Scan(
			&item.ID,
			&item.Name,
			&item.MetaTags,
			&item.Slug,
			&item.Type,
			&item.Phone,
			&item.Email,
			&item.EmailVerified,
			&item.Logo,
			&item.Banner,
			&item.Rating,
			&item.VendorID,
			&item.Prefix,
			&item.Status,
			&item.Availability,
			&item.Location,
			&item.Address,
			&item.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	return items, nil
}
