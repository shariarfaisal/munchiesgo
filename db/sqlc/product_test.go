package db

import (
	"context"
	"testing"

	"github.com/Munchies-Engineering/backend/util"
	"github.com/stretchr/testify/require"
)

func TestCreateProduct(t *testing.T) {
	randomProduct(t)
}

func randomProduct(t *testing.T) Product {

	category := randomBrandCategory(t)

	arg := CreateProductParams{
		Type:         "product",
		Name:         util.RandomString(6),
		Slug:         util.RandomString(16),
		CategoryID:   category.ID,
		Image:        util.RandomString(6),
		Details:      util.RandomString(16),
		Price:        float64(util.RandomInt(1, 100)),
		Status:       util.RandomString(6),
		Availability: true,
		BrandID:      category.BrandID,
	}

	product, err := testQueries.CreateProduct(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, product)
	require.NotEmpty(t, product.ID)
	require.NotEmpty(t, product.CreatedAt)

	require.Equal(t, arg.Type, product.Type)
	require.Equal(t, arg.Name, product.Name)
	require.Equal(t, arg.Slug, product.Slug)
	require.Equal(t, arg.CategoryID, product.CategoryID)
	require.Equal(t, arg.Image, product.Image)
	require.Equal(t, arg.Details, product.Details)
	require.Equal(t, arg.Price, product.Price)
	require.Equal(t, arg.Status, product.Status)
	require.Equal(t, arg.Availability, product.Availability)
	require.Equal(t, arg.BrandID, product.BrandID)

	return product
}
