package db

import (
	"context"
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
