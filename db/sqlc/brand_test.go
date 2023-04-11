package db

import (
	"context"
	"database/sql/driver"
	"fmt"
	"testing"

	"github.com/Munchies-Engineering/backend/util"
	"github.com/stretchr/testify/require"
)

func TestCreateBrand(t *testing.T) {
	randomBrand(t)
}

type Point struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}

func (p Point) Value() driver.Value {
	return fmt.Sprintf("(%f,%f)", p.Lat, p.Lng)
}

func randomBrand(t *testing.T) Brand {
	vendor := randomVendor(t)

	arg := CreateBrandParams{
		Name:          util.RandomString(10),
		MetaTags:      util.RandomString(20),
		Slug:          util.RandomString(10),
		Type:          util.RandomBrandType(),
		Phone:         util.RandomString(10),
		Email:         util.RandomEmail(),
		EmailVerified: false,
		Logo:          util.RandomString(10),
		Banner:        util.RandomString(10),
		Rating:        5,
		VendorID:      vendor.ID,
		Prefix:        util.RandomString(3),
		Status:        util.RandomString(10), // FIXME: Status should be a enum type
		Availability:  true,
		Location:      Point{Lat: 1.0, Lng: 1.0}.Value(),
		Address:       util.RandomString(10),
	}

	brand, err := testQueries.CreateBrand(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, brand)
	require.NotEmpty(t, brand.ID)
	require.NotEmpty(t, brand.CreatedAt)

	require.Equal(t, arg.Name, brand.Name)
	require.Equal(t, arg.MetaTags, brand.MetaTags)
	require.Equal(t, arg.Slug, brand.Slug)
	require.Equal(t, arg.Type, brand.Type)
	require.Equal(t, arg.Phone, brand.Phone)
	require.Equal(t, arg.Email, brand.Email)
	require.Equal(t, arg.EmailVerified, brand.EmailVerified)
	require.Equal(t, arg.Logo, brand.Logo)
	require.Equal(t, arg.Banner, brand.Banner)
	require.Equal(t, arg.Rating, brand.Rating)
	require.Equal(t, arg.VendorID, brand.VendorID)
	require.Equal(t, arg.Prefix, brand.Prefix)
	require.Equal(t, arg.Status, brand.Status)
	require.Equal(t, arg.Availability, brand.Availability)
	require.NotEmpty(t, brand.Location)
	require.Equal(t, arg.Address, brand.Address)

	return brand
}
