package db

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/Munchies-Engineering/backend/util"
	"github.com/stretchr/testify/require"
)

func TestCreateVendor(t *testing.T) {
	randomVendor(t)
}

func randomVendor(t *testing.T) Vendor {

	paymentInfo, err := json.Marshal(map[string]interface{}{
		"account_number": util.RandomString(10),
	})
	require.NoError(t, err)

	socialLinks, err := json.Marshal(map[string]interface{}{
		"facebook": util.RandomString(10),
	})
	require.NoError(t, err)

	arg := CreateVendorParams{
		Name:        util.RandomString(6),
		Email:       util.RandomEmail(),
		Phone:       util.RandomString(10),
		PaymentInfo: paymentInfo,
		SocialLinks: socialLinks,
	}

	vendor, err := testQueries.CreateVendor(context.Background(), arg)
	fmt.Println(err)
	require.NoError(t, err)
	require.NotEmpty(t, vendor)

	require.Equal(t, arg.Name, vendor.Name)
	require.Equal(t, arg.Email, vendor.Email)
	require.Equal(t, arg.Phone, vendor.Phone)

	require.NotEmpty(t, vendor.PaymentInfo)
	require.NotEmpty(t, vendor.SocialLinks)

	var paymentInfoMap map[string]interface{}
	err = json.Unmarshal(vendor.PaymentInfo, &paymentInfoMap)
	require.NoError(t, err)
	require.NotEmpty(t, paymentInfoMap["account_number"])

	var socialLinksMap map[string]interface{}
	err = json.Unmarshal(vendor.SocialLinks, &socialLinksMap)
	require.NoError(t, err)
	require.NotEmpty(t, socialLinksMap["facebook"])

	arg2 := CreateVendorParams{
		Name:        util.RandomString(6),
		Email:       arg.Email,
		Phone:       util.RandomString(10),
		PaymentInfo: paymentInfo,
		SocialLinks: socialLinks,
	}
	vendor2, err := testQueries.CreateVendor(context.Background(), arg2)
	require.Error(t, err)
	require.Empty(t, vendor2)

	return vendor
}
