package db

import (
	"context"
	"testing"

	"github.com/Munchies-Engineering/backend/util"
	"github.com/stretchr/testify/require"
)

func TestVendorUserCreate(t *testing.T) {
	randomVendorUser(t)
}

func randomVendorUser(t *testing.T) VendorUser {
	vendor := randomVendor(t)

	arg := CreateVendorUserParams{
		Username:       util.RandomString(6),
		Email:          util.RandomEmail(),
		HashedPassword: util.RandomString(6),
		Role:           util.RandomString(6),
		VendorID:       vendor.ID,
	}

	vendorUser, err := testQueries.CreateVendorUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, vendorUser)
	require.NotEmpty(t, vendorUser.ID)
	require.NotEmpty(t, vendorUser.CreatedAt)

	require.Equal(t, arg.Username, vendorUser.Username)
	require.Equal(t, arg.Email, vendorUser.Email)
	require.Equal(t, arg.HashedPassword, vendorUser.HashedPassword)
	require.Equal(t, arg.Role, vendorUser.Role)
	require.Equal(t, arg.VendorID, vendorUser.VendorID)

	return vendorUser
}
