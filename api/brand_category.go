package api

import (
	"database/sql"
	"errors"
	"net/http"
	"strconv"

	db "github.com/Munchies-Engineering/backend/db/sqlc"
	"github.com/Munchies-Engineering/backend/token"
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

/*
	CreateBrandCategory
	UpdateBrandCategory
	GetBrandCategory
	SearchBrandCategories
	DeleteBrandCategory
	ListBrandCategories
	CountCategories
	CountCategoriesByBrandID
*/

type createBrandCategoryRequest struct {
	CategoryID int64  `json:"categoryId" binding:"required,min=1"`
	Name       string `json:"name" binding:"required"`
}

func (server *Server) createBrandCategory(ctx *gin.Context) {
	var req createBrandCategoryRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	authPayload := ctx.MustGet(authorizationPayload).(*token.Payload)

	brandId, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(errors.New("Invalid ID in params.")))
		return
	}

	brand, err := server.store.GetBrand(ctx, brandId)
	if err != nil {
		if err := sql.ErrNoRows; err != nil {
			ctx.JSON(http.StatusNotFound, errorResponse(errors.New("Brand not found.")))
			return
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	if brand.VendorID != authPayload.VendorID {
		ctx.JSON(http.StatusUnauthorized, errorResponse(errors.New("Unauthorized.")))
		return
	}

	arg := db.CreateBrandCategoryParams{
		Name:       req.Name,
		BrandID:    brand.ID,
		CategoryID: req.CategoryID,
	}

	// TODO: add admin authentication

	category, err := server.store.CreateBrandCategory(ctx, arg)
	if err != nil {
		if err := err.(*pq.Error); err != nil {
			switch err.Code.Name() {
			case "unique_violation":
				ctx.JSON(http.StatusConflict, errorResponse(errors.New("category already exists")))
				return
			}
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, category)
}

type updateBrandCategoryRequest struct {
	Name string `json:"name" binding:"required"`
}

func (server *Server) updateBrandCategory(ctx *gin.Context) {
	var req updateBrandCategoryRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	authPayload := ctx.MustGet(authorizationPayload).(*token.Payload)

	brandId, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(errors.New("Invalid brand ID")))
		return
	}

	brand, err := server.store.GetBrand(ctx, brandId)
	if err != nil {
		if err := sql.ErrNoRows; err != nil {
			ctx.JSON(http.StatusNotFound, errorResponse(errors.New("Brand not found.")))
			return
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	if brand.VendorID != authPayload.VendorID {
		ctx.JSON(http.StatusUnauthorized, errorResponse(errors.New("Unauthorized.")))
		return
	}

	categoryId, err := strconv.ParseInt(ctx.Param("categoryId"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(errors.New("Invalid category ID")))
		return
	}

	category, err := server.store.GetBrandCategory(ctx, categoryId)
	if err != nil {
		if err := sql.ErrNoRows; err != nil {
			ctx.JSON(http.StatusNotFound, errorResponse(errors.New("Category not found.")))
			return
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	if category.BrandID != brand.ID {
		ctx.JSON(http.StatusUnauthorized, errorResponse(errors.New("Invalid brand id or category id.")))
		return
	}

	arg := db.UpdateBrandCategoryParams{
		ID:      category.ID,
		Name:    category.Name,
		BrandID: category.BrandID,
	}

	if req.Name != "" {
		arg.Name = req.Name
	}

	updatedCategory, err := server.store.UpdateBrandCategory(ctx, arg)
	if err != nil {
		if err := err.(*pq.Error); err != nil {
			switch err.Code.Name() {
			case "unique_violation":
				ctx.JSON(http.StatusConflict, errorResponse(errors.New("category already exists")))
				return
			}
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, updatedCategory)
}

func (server *Server) getBrandCategory(ctx *gin.Context) {

	brandId, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(errors.New("Invalid ID in params.")))
		return
	}
	authPayload := ctx.MustGet(authorizationPayload).(*token.Payload)

	brand, err := server.store.GetBrand(ctx, brandId)
	if err != nil {
		if err := sql.ErrNoRows; err != nil {
			ctx.JSON(http.StatusNotFound, errorResponse(errors.New("Brand not found.")))
			return
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	if brand.VendorID != authPayload.VendorID {
		ctx.JSON(http.StatusUnauthorized, errorResponse(errors.New("Unauthorized.")))
		return
	}

	id, err := strconv.ParseInt(ctx.Param("categoryId"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(errors.New("Invalid ID in params.")))
		return
	}

	category, err := server.store.GetBrandCategory(ctx, id)
	if err != nil {
		if err := sql.ErrNoRows; err != nil {
			ctx.JSON(http.StatusNotFound, errorResponse(errors.New("Category not found.")))
			return
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	if category.BrandID != brand.ID {
		ctx.JSON(http.StatusUnauthorized, errorResponse(errors.New("Invalid brand ID or category ID.")))
		return
	}

	ctx.JSON(http.StatusOK, category)
}

type listBrandCategoriesRequest struct {
	PageSize int32 `form:"pageSize" binding:"required,min=1,max=1000"`
	PageID   int32 `form:"pageId" binding:"required,min=1"`
}

func (server *Server) listBrandCategories(ctx *gin.Context) {
	var req listBrandCategoriesRequest
	if err := ctx.ShouldBind(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	brandId, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(errors.New("Invalid ID in params.")))
		return
	}

	brand, err := server.store.GetBrand(ctx, brandId)
	if err != nil {
		if err := sql.ErrNoRows; err != nil {
			ctx.JSON(http.StatusNotFound, errorResponse(errors.New("Brand not found.")))
			return
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	authPayload := ctx.MustGet(authorizationPayload).(*token.Payload)

	if brand.VendorID != authPayload.VendorID {
		ctx.JSON(http.StatusUnauthorized, errorResponse(errors.New("Unauthorized.")))
		return
	}

	arg := db.ListBrandCategoriesParams{
		BrandID: brand.ID,
		Limit:   req.PageSize,
		Offset:  (req.PageID - 1) * req.PageSize,
	}

	categories, err := server.store.ListBrandCategories(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, categories)
}

func (server *Server) deleteBrandCategory(ctx *gin.Context) {
	brandId, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(errors.New("Invalid ID in params.")))
		return
	}
	authPayload := ctx.MustGet(authorizationPayload).(*token.Payload)

	brand, err := server.store.GetBrand(ctx, brandId)
	if err != nil {
		if err := sql.ErrNoRows; err != nil {
			ctx.JSON(http.StatusNotFound, errorResponse(errors.New("Brand not found.")))
			return
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	if brand.VendorID != authPayload.VendorID {
		ctx.JSON(http.StatusUnauthorized, errorResponse(errors.New("Unauthorized.")))
		return
	}

	id, err := strconv.ParseInt(ctx.Param("categoryId"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(errors.New("Invalid ID in params.")))
		return
	}

	category, err := server.store.GetBrandCategory(ctx, id)
	if err != nil {
		if err := sql.ErrNoRows; err != nil {
			ctx.JSON(http.StatusNotFound, errorResponse(errors.New("Category not found.")))
			return
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	if category.BrandID != brand.ID {
		ctx.JSON(http.StatusUnauthorized, errorResponse(errors.New("Invalid brand ID or category ID.")))
		return
	}

	arg := db.DeleteBrandCategoryParams{
		ID:      id,
		BrandID: brand.ID,
	}

	err = server.store.DeleteBrandCategory(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, category)
}

func (server *Server) searchBrandCategories(ctx *gin.Context) {

}
