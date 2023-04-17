package api

import (
	"database/sql"
	"errors"
	"net/http"
	"reflect"

	db "github.com/Munchies-Engineering/backend/db/sqlc"
	"github.com/Munchies-Engineering/backend/token"
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

const (
	BrandTypeRestaurant string = "restaurant"
	BrandTypeStore      string = "store"
)

const (
	BrandStatusActive   string = "active"
	BrandStatusInactive string = "inactive"
	BrandStatusPending  string = "pending"
	BrandStatusRejected string = "rejected"
)

/*
	CreateBrand
	UpdateBrand
	DeleteBrand
	GetBrand
	ListBrands
	ListBrandsByVendorID
	CreateOperationTime
	UpdateOperationTime
	DeleteOperationTime
	GetOperationTime
	ListOperationTimesByBrandId
*/

type createBrandRequest struct {
	Name     string `json:"name" binding:"required,min=2"`
	MetaTags string `json:"metaTags" binding:"required"`
	Type     string `json:"type" binding:"required,oneof=restaurant store"`
	Phone    string `json:"phone" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Logo     string `json:"logo" binding:"required"`
	Banner   string `json:"banner" binding:"required"`
	Location Point  `json:"location" binding:"required"`
	Address  string `json:"address" binding:"required"`
}

type brandResponse struct {
	ID           int64  `json:"id"`
	Name         string `json:"name"`
	MetaTags     string `json:"metaTags"`
	Slug         string `json:"slug"`
	Type         string `json:"type"`
	Phone        string `json:"phone"`
	Email        string `json:"email"`
	Logo         string `json:"logo"`
	Banner       string `json:"banner"`
	Rating       int32  `json:"rating"`
	Status       string `json:"status"`
	Availability bool   `json:"availability"`
	Location     *Point `json:"location"`
	Prefix       string `json:"prefix"`
	VendorID     int64  `json:"vendorID"`
	Address      string `json:"address"`
}

func getBrandResponse(brand db.Brand) *brandResponse {
	return &brandResponse{
		ID:           brand.ID,
		Name:         brand.Name,
		MetaTags:     brand.MetaTags,
		Slug:         brand.Slug,
		Type:         brand.Type,
		Phone:        brand.Phone,
		Email:        brand.Email,
		Logo:         brand.Logo,
		Banner:       brand.Banner,
		Rating:       brand.Rating,
		Status:       brand.Status,
		Availability: brand.Availability,
		Location:     getPoint(brand.Location),
		Prefix:       brand.Prefix,
		VendorID:     brand.VendorID,
		Address:      brand.Address,
	}
}

func (server *Server) createBrand(ctx *gin.Context) {
	var req createBrandRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	authPayload := ctx.MustGet(authorizationPayload).(*token.Payload)
	vendor, err := server.store.GetVendor(ctx, authPayload.VendorID)

	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	location := req.Location.Value()

	arg := db.CreateBrandParams{
		Name:          req.Name,
		MetaTags:      req.MetaTags,
		Slug:          getBrandsSlug(req.Name, 1),
		Type:          req.Type,
		Phone:         req.Phone,
		Email:         req.Email,
		EmailVerified: false,
		Logo:          req.Logo,
		Banner:        req.Banner,
		Rating:        0,
		Status:        BrandStatusPending,
		Availability:  false,
		Location:      location,
		Prefix:        getPrefix(req.Name),
		VendorID:      vendor.ID,
		Address:       req.Address,
	}

	brand, err := server.store.CreateBrand(ctx, arg)

	if err != nil {
		if err := err.(*pq.Error); err != nil {
			switch err.Code.Name() {
			case "unique_violation":
				ctx.JSON(http.StatusConflict, errorResponse(err))
				return
			}
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
	}

	ctx.JSON(http.StatusOK, getBrandResponse(brand))
}

type updateBrandRequest struct {
	Name     string `json:"name" binding:"min=2"`
	MetaTags string `json:"metaTags"`
	Type     string `json:"type" binding:"omitempty,oneof=restaurant store"`
	Phone    string `json:"phone"`
	Email    string `json:"email" binding:"omitempty,email"`
	Logo     string `json:"logo"`
	Banner   string `json:"banner"`
	Location Point  `json:"location"`
	Address  string `json:"address"`
}

type updateBrandRequestParams struct {
	ID int64 `uri:"id" binding:"required"`
}

func (server *Server) updateBrand(ctx *gin.Context) {
	var req updateBrandRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	var params updateBrandRequestParams
	if err := ctx.ShouldBindUri(&params); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	brand, err := server.store.GetBrand(ctx, params.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	arg := db.UpdateBrandParams{
		ID:           brand.ID,
		Name:         brand.Name,
		MetaTags:     brand.MetaTags,
		Slug:         brand.Slug,
		Type:         brand.Type,
		Phone:        brand.Phone,
		Email:        brand.Email,
		Logo:         brand.Logo,
		Banner:       brand.Banner,
		Rating:       brand.Rating,
		Status:       brand.Status,
		Availability: brand.Availability,
		Location:     getPoint(brand.Location).Value(),
		Prefix:       brand.Prefix,
		VendorID:     brand.VendorID,
		Address:      brand.Address,
	}

	reflectValue := reflect.ValueOf(req)
	reflectType := reflect.TypeOf(req)

	for i := 0; i < reflectValue.NumField(); i++ {
		if reflectValue.Field(i).IsZero() {
			continue
		}

		switch reflectType.Field(i).Name {
		case "Name":
			arg.Name = req.Name
			arg.Slug = getBrandsSlug(req.Name, 1)
			arg.Prefix = getPrefix(req.Name)
		case "MetaTags":
			arg.MetaTags = req.MetaTags
		case "Type":
			arg.Type = req.Type
		case "Phone":
			arg.Phone = req.Phone
		case "Email":
			arg.Email = req.Email
		case "Logo":
			arg.Logo = req.Logo
		case "Banner":
			arg.Banner = req.Banner
		case "Location":
			arg.Location = req.Location.Value()
		case "Address":
			arg.Address = req.Address
		}
	}

	brand, err = server.store.UpdateBrand(ctx, arg)
	if err != nil {
		if err := err.(*pq.Error); err != nil {
			switch err.Code.Name() {
			case "unique_violation":
				ctx.JSON(http.StatusConflict, errorResponse(err))
				return
			}
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, getBrandResponse(brand))
}

type deleteBrandRequestParams struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (server *Server) deleteBrand(ctx *gin.Context) {
	var req deleteBrandRequestParams
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	brand, err := server.store.GetBrand(ctx, req.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	err = server.store.DeleteBrand(ctx, brand.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, getBrandResponse(brand))
}

type getBrandRequestParams struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (server *Server) getBrand(ctx *gin.Context) {
	var req getBrandRequestParams
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	authPayload := ctx.MustGet(authorizationPayload).(*token.Payload)

	brand, err := server.store.GetBrand(ctx, req.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	if brand.VendorID != authPayload.VendorID {
		ctx.JSON(http.StatusUnauthorized, errorResponse(errors.New("Unauthorized")))
		return
	}

	ctx.JSON(http.StatusOK, getBrandResponse(brand))
}

type listBrandsRequest struct {
	PageID   int32 `form:"pageId" binding:"required,min=1"`
	PageSize int32 `form:"pageSize" binding:"required,min=1,max=10000"`
}

func (server *Server) listBrands(ctx *gin.Context) {
	var req listBrandsRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	authPayload := ctx.MustGet(authorizationPayload).(*token.Payload)
	arg := db.ListBrandsByVendorIDParams{
		Limit:    req.PageSize,
		Offset:   (req.PageID - 1) * req.PageSize,
		VendorID: authPayload.VendorID,
	}

	brands, err := server.store.ListBrandsByVendorID(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	result := make([]*brandResponse, 0, len(brands))
	for i := range brands {
		result = append(result, getBrandResponse(brands[i]))
	}

	ctx.JSON(http.StatusOK, result)
}
