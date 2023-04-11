package db

import (
	"context"
	"testing"

	"github.com/Munchies-Engineering/backend/util"
	"github.com/stretchr/testify/require"
)

func TestCreateCategory(t *testing.T) {
	randomCategory(t)
}

func TestCreateBrandCategory(t *testing.T) {
	randomBrandCategory(t)
}

func randomCategory(t *testing.T) Category {
	arg := CreateCategoryParams{
		Name:  util.RandomString(6),
		Image: util.RandomString(6),
	}

	category, err := testQueries.CreateCategory(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, category)
	require.NotEmpty(t, category.ID)
	require.NotEmpty(t, category.CreatedAt)

	require.Equal(t, arg.Name, category.Name)
	require.Equal(t, arg.Image, category.Image)

	return category
}

func randomBrandCategory(t *testing.T) BrandCategory {
	brand := randomBrand(t)
	category := randomCategory(t)

	arg := CreateBrandCategoryParams{
		BrandID:    brand.ID,
		CategoryID: category.ID,
		Name:       util.RandomString(6),
	}

	brandCategory, err := testQueries.CreateBrandCategory(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, brandCategory)
	require.NotEmpty(t, brandCategory.ID)
	require.NotEmpty(t, brandCategory.CreatedAt)

	require.Equal(t, arg.BrandID, brandCategory.BrandID)
	require.Equal(t, arg.CategoryID, brandCategory.CategoryID)

	return brandCategory
}
