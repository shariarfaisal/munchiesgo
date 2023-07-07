package db

import (
	"context"
	"fmt"
	"time"

	"github.com/lib/pq"
)

var (
	ProductStatusPending string = "pending"
	ProductStatusActive  string = "active"
)

var (
	ProductTypeProduct string = "product"
	ProductTypeVariant string = "variant"
)

type ProductUploadVariantParams struct {
	Title     string
	MinSelect int32
	MaxSelect int32
	Items     []ProductUploadParams
}

type ProductUploadParams struct {
	Type         string
	Name         string
	Slug         string
	Status       string
	CategoryID   int64
	Image        string
	Details      string
	Price        float64
	BrandID      int64
	Availability bool
	Variant      []ProductUploadVariantParams
}

func SingleProductUpload(ctx context.Context, q *Queries, prod ProductUploadParams) (*Product, error) {
	arg := CreateProductParams{
		Type:         prod.Type,
		Name:         prod.Name,
		Slug:         prod.Slug,
		Status:       prod.Status,
		CategoryID:   prod.CategoryID,
		Image:        prod.Image,
		Details:      prod.Details,
		Price:        prod.Price,
		BrandID:      prod.BrandID,
		Availability: prod.Availability,
	}

	product, err := q.CreateProduct(ctx, arg)
	if err != nil {
		return nil, err
	}

	if len(prod.Variant) > 0 {
		for _, variant := range prod.Variant {
			arg := CreateProductVariantParams{
				ProductID: product.ID,
				Title:     variant.Title,
				MinSelect: variant.MinSelect,
				MaxSelect: variant.MaxSelect,
			}

			productVariant, err := q.CreateProductVariant(ctx, arg)
			if err != nil {
				return nil, err
			}

			for _, item := range variant.Items {

				// add category and brand id to variant item from parent product
				item.CategoryID = product.CategoryID
				item.BrandID = product.BrandID

				product, err := SingleProductUpload(ctx, q, item)
				if err != nil {
					return nil, err
				}

				arg := CreateProductVariantItemParams{
					VariantID: productVariant.ID,
					ProductID: product.ID,
				}

				_, err = q.CreateProductVariantItem(ctx, arg)
				if err != nil {
					return nil, err
				}
			}
		}
	}

	return &product, nil
}

func (store *SqlStore) BulkProductUpload(ctx context.Context, args []ProductUploadParams) ([]Product, error) {

	products := []Product{}

	err := store.execTx(ctx, func(q *Queries) error {
		for _, prod := range args {
			product, err := SingleProductUpload(ctx, q, prod)
			if err != nil {
				return err
			}

			products = append(products, *product)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return products, nil
}

type FullProduct struct {
	*Product
	Brand    *Brand         `json:"brand,omitempty"`
	Category *BrandCategory `json:"category,omitempty"`
	Variants []struct {
		*ProductVariant
		Items []ListVariantItemsWithProductDetailsRow `json:"items,omitempty"`
	} `json:"variants,omitempty"`
}

func (store *SqlStore) FullProduct(ctx context.Context, id int64) (*FullProduct, error) {

	var result FullProduct

	product, err := store.GetProduct(ctx, id)
	if err != nil {
		return nil, err
	}
	result.Product = &product

	brand, err := store.GetBrand(ctx, product.BrandID)
	if err != nil {
		return nil, err
	}

	result.Brand = &brand

	category, err := store.GetBrandCategory(ctx, product.CategoryID)
	if err != nil {
		return nil, err
	}

	result.Category = &category

	variants, err := store.ListProductVariants(ctx, ListProductVariantsParams{
		ProductID: id,
		Limit:     100,
		Offset:    0,
	})
	if err != nil {
		return nil, err
	}

	for _, variant := range variants {
		items, err := store.ListVariantItemsWithProductDetails(ctx, variant.ID)
		if err != nil {
			return nil, err
		}

		result.Variants = append(result.Variants, struct {
			*ProductVariant
			Items []ListVariantItemsWithProductDetailsRow `json:"items,omitempty"`
		}{
			ProductVariant: &variant,
			Items:          items,
		})
	}

	return &result, nil
}

func (store *SqlStore) ProductsByIds(ctx context.Context, ids []int64) ([]Product, error) {
	const query = `
		SELECT * FROM products WHERE id = ANY($1);
	`

	rows, err := store.db.QueryContext(ctx, query, pq.Array(ids))
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var items []Product
	for rows.Next() {
		var item Product
		err := rows.Scan(
			&item.ID,
			&item.Type,
			&item.Name,
			&item.CategoryID,
			&item.Slug,
			&item.Image,
			&item.Details,
			&item.Price,
			&item.Status,
			&item.BrandID,
			&item.Availability,
			&item.UseInventory,
			&item.CreatedAt,
			&item.UpdatedAt,
		)

		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	return items, nil
}

type ProductsTypeByIdsRow struct {
	ID   int64  `json:"id"`
	Type string `json:"type"`
}

func (store *SqlStore) ProductsTypeByIds(ctx context.Context, ids []int64) ([]ProductsTypeByIdsRow, error) {
	const query = `
		SELECT id,type FROM products WHERE id = ANY($1)
	`

	rows, err := store.db.QueryContext(ctx, query, pq.Array(ids))
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var items []ProductsTypeByIdsRow
	for rows.Next() {
		var item ProductsTypeByIdsRow
		err := rows.Scan(
			&item.ID,
			&item.Type,
		)

		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	return items, nil
}

func (store *SqlStore) ProductVariantsByProductIds(ctx context.Context, ids []int64) ([]ProductVariant, error) {
	const query = `
		SELECT * FROM product_variants WHERE product_id = ANY($1)
	`

	rows, err := store.db.QueryContext(ctx, query, pq.Array(ids))
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var items []ProductVariant
	for rows.Next() {
		var item ProductVariant
		err := rows.Scan(
			&item.ID,
			&item.ProductID,
			&item.Title,
			&item.MinSelect,
			&item.MaxSelect,
		)

		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	return items, nil
}

func (store *SqlStore) ProductVariantsByIds(ctx context.Context, ids []int64) ([]ProductVariant, error) {
	const query = `
		SELECT * FROM product_variants WHERE id = ANY($1);
	`

	rows, err := store.db.QueryContext(ctx, query, pq.Array(ids))
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var items []ProductVariant
	for rows.Next() {
		var item ProductVariant
		err := rows.Scan(
			&item.ID,
			&item.ProductID,
			&item.Title,
			&item.MinSelect,
			&item.MaxSelect,
			&item.CreatedAt,
			&item.UpdatedAt,
		)

		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	return items, nil
}

func (store *SqlStore) ProductVariantsItemsByProductIds(ctx context.Context, ids []int64) ([]ProductVariantItem, error) {
	const query = `
		SELECT * FROM product_variant_items WHERE product_id = ANY($1)
	`

	rows, err := store.db.QueryContext(ctx, query, pq.Array(ids))
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var items []ProductVariantItem
	for rows.Next() {
		var item ProductVariantItem
		err := rows.Scan(
			&item.ID,
			&item.VariantID,
			&item.ProductID,
			&item.CreatedAt,
		)

		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	return items, nil
}

type OrderItemProduct struct {
	ID           int64
	Type         string
	Name         string
	CategoryID   int64
	Category     BrandCategory
	Slug         string
	Image        string
	Details      string
	Price        float64
	Status       string
	BrandID      int64
	Brand        Brand
	Availability bool
	UseInventory bool
	CreatedAt    time.Time
	UpdatedAt    time.Time
	Parent       Product
}

func (store *SqlStore) OrderItemProducts(ctx context.Context, orderItems []PlaceOrderItemParams) ([]*OrderItemProduct, error) {
	var products []*OrderItemProduct

	productIds := make([]int64, len(orderItems))
	for i, item := range orderItems {
		productIds[i] = item.ProductID
	}

	// get products
	productItems, err := store.ProductsByIds(ctx, productIds)
	if err != nil {
		return nil, err
	}

	brandIds := []int64{}
	categoryIds := []int64{}

	for _, item := range productItems {
		brandIds = append(brandIds, item.BrandID)
		categoryIds = append(categoryIds, item.CategoryID)
	}

	brands, err := store.BrandsWithIds(ctx, brandIds)
	if err != nil {
		return nil, err
	}

	brandsGroup := map[int64]Brand{}
	for _, brand := range brands {
		//
		brandsGroup[brand.ID] = brand
	}

	categories, err := store.BrandCategoriesByIds(ctx, categoryIds)
	if err != nil {
		return nil, err
	}

	categoriesGroup := map[int64]BrandCategory{}
	for _, category := range categories {
		categoriesGroup[category.ID] = category
	}

	for _, item := range productItems {
		brand, ok := brandsGroup[item.BrandID]
		if !ok {
			return nil, fmt.Errorf("brand not found")
		}

		category, ok := categoriesGroup[item.CategoryID]
		if !ok {
			return nil, fmt.Errorf("category not found")
		}

		products = append(products, &OrderItemProduct{
			ID:           item.ID,
			Type:         item.Type,
			Name:         item.Name,
			CategoryID:   item.CategoryID,
			Category:     category,
			Slug:         item.Slug,
			Image:        item.Image,
			Details:      item.Details,
			Price:        item.Price,
			Status:       item.Status,
			BrandID:      item.BrandID,
			Brand:        brand,
			Availability: item.Availability,
			UseInventory: item.UseInventory,
			CreatedAt:    item.CreatedAt,
		})
	}

	// check is all products are exists and available
	for _, id := range productIds {
		exists := false
		for _, item := range productItems {
			if item.ID == id {
				exists = true
			}
		}

		if !exists {
			return nil, fmt.Errorf("no product found for id: %d", id)
		}
	}

	variantProductIds := []int64{}
	for _, item := range productItems {
		if item.Type == ProductTypeVariant {
			variantProductIds = append(variantProductIds, item.ID)
		}
	}

	if len(variantProductIds) > 0 {
		variantItems, err := store.ProductVariantsItemsByProductIds(ctx, variantProductIds)
		if err != nil {
			return nil, err
		}

		// check is all variant items are exists
		for _, id := range variantProductIds {
			exists := false
			for _, item := range variantItems {
				if item.ProductID == id {
					exists = true
				}
			}

			if !exists {
				return nil, fmt.Errorf("no variant items found for product id: %d", id)
			}
		}

		// group variant items by variant id
		itemsByVariantId := map[int64][]ProductVariantItem{}
		for _, item := range variantItems {
			itemsByVariantId[item.VariantID] = append(itemsByVariantId[item.VariantID], item)
		}

		variantIds := []int64{}
		for _, item := range variantItems {
			variantIds = append(variantIds, item.VariantID)
		}

		variants, err := store.ProductVariantsByIds(ctx, variantIds)
		if err != nil {
			return nil, err
		}

		productIdsOfVariants := []int64{}
		for _, variant := range variants {
			productIdsOfVariants = append(productIdsOfVariants, variant.ProductID)
		}

		productsOfVariants, err := store.ProductsByIds(ctx, productIdsOfVariants)
		if err != nil {
			return nil, err
		}

		for _, product := range products {
			if product.Type == ProductTypeVariant {
				for _, variantItem := range variantItems {
					if variantItem.ProductID == product.ID {
						for _, variant := range variants {
							if variant.ID == variantItem.VariantID {
								for _, parentProduct := range productsOfVariants {
									if variant.ProductID == parentProduct.ID {
										product.Parent = parentProduct
										break
									}
								}
								break
							}
						}
						break
					}
				}
			}
		}

		// check is all variants are exists
		for _, id := range variantIds {
			exists := false
			for _, item := range variants {
				if item.ID == id {
					exists = true
				}
			}

			if !exists {
				return nil, fmt.Errorf("no variant found for variant id: %d", id)
			}
		}

		// check is min max select rules match for variants
		for _, variant := range variants {
			v, ok := itemsByVariantId[variant.ID]

			if !ok || variant.MinSelect > 0 && len(v) < int(variant.MinSelect) {
				return nil, fmt.Errorf("min select rule not match for variant id: %d", variant.ID)
			} else if variant.MaxSelect > 0 && len(v) > int(variant.MaxSelect) {
				return nil, fmt.Errorf("max select rule not match for variant id: %d", variant.ID)
			}
		}

	}

	return products, nil
}
