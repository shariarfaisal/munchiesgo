package api

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	db "github.com/Munchies-Engineering/backend/db/sqlc"
	"github.com/Munchies-Engineering/backend/token"
	"github.com/gin-gonic/gin"
)

type createProductVariantItemRequest struct {
	Name         string  `json:"name" binding:"required"`
	Image        string  `json:"image" binding:"required"`
	Details      string  `json:"details" binding:"required"`
	Price        float64 `json:"price" binding:"min=0"`
	Availability bool    `json:"availability" binding:"required"`
}

type createProductVariantRequest struct {
	Title     string                            `json:"title" binding:"required"`
	MinSelect int32                             `json:"minSelect" binding:"required,min=0"`
	MaxSelect int32                             `json:"maxSelect" binding:"required,min=0"`
	Items     []createProductVariantItemRequest `json:"items" binding:"omitempty,dive"`
}

type createProductRequest struct {
	Name         string                        `json:"name" binding:"required"`
	CategoryID   int64                         `json:"categoryId" binding:"required,min=1"`
	Image        string                        `json:"image" binding:"required"`
	Details      string                        `json:"details" binding:"required"`
	Price        float64                       `json:"price" binding:"min=0"`
	BrandID      int64                         `json:"brandId" binding:"required,min=1"`
	Availability bool                          `json:"availability" binding:"required"`
	Variant      []createProductVariantRequest `json:"variant,omitempty" binding:"omitempty"`
}

func (server *Server) isAllBrandsAreValid(ctx context.Context, products []createProductRequest, brands map[int64]bool, vendorId int64) error {
	invalidBrandIds := []int64{}
	for _, product := range products {
		if !brands[product.BrandID] {
			invalidBrandIds = append(invalidBrandIds, product.BrandID)
		}
	}

	if len(invalidBrandIds) != 0 {
		return errors.New(fmt.Sprintf("Invalid brand id %v", invalidBrandIds))
	}

	return nil
}

func (server *Server) isAllCategoriesAreValid(ctx *gin.Context, products []createProductRequest, brands map[int64]bool) error {
	catIdGroup := map[int64]bool{}
	catIds := []int64{}

	for _, product := range products {
		if !catIdGroup[product.CategoryID] {
			catIds = append(catIds, product.CategoryID)
		}
		catIdGroup[product.CategoryID] = true
	}

	// Get all categories by ids
	categories, err := server.store.BrandCategoriesByIds(ctx, catIds)
	if err != nil {
		return err
	}

	catResultGroup := map[int64]bool{}
	invalidCatIds := []int64{}

	for _, cat := range categories {
		catResultGroup[cat.ID] = true

		// Check if this category belongs to this vendor
		if !brands[cat.BrandID] {
			invalidCatIds = append(invalidCatIds, cat.ID)
		}
	}

	// Check if all categories are valid
	for _, product := range products {
		if !catResultGroup[product.CategoryID] {
			invalidCatIds = append(invalidCatIds, product.CategoryID)
		}
	}

	if len(invalidCatIds) != 0 {
		return errors.New(fmt.Sprintf("Invalid category id %v", invalidCatIds))
	}

	return nil
}

func (server *Server) createProduct(ctx *gin.Context) {

	var req []createProductRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	authPayload := ctx.MustGet(authorizationPayload).(*token.Payload)

	brandIdGroup := map[int64]bool{}
	allBrandIds, err := server.store.ListBrandIdsByVendorID(ctx, authPayload.VendorID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	for _, brandId := range allBrandIds {
		brandIdGroup[brandId] = true
	}

	// Validate brand and category ids, is all provided brands and categories are belongs to this vendor
	err = server.isAllBrandsAreValid(ctx, req, brandIdGroup, authPayload.VendorID)
	if err != nil {
		ctx.JSON(http.StatusForbidden, errorResponse(err))
		return
	}

	err = server.isAllCategoriesAreValid(ctx, req, brandIdGroup)
	if err != nil {
		ctx.JSON(http.StatusForbidden, errorResponse(err))
		return
	}

	args := []db.ProductUploadParams{}
	for _, product := range req {
		arg := db.ProductUploadParams{
			Type:         db.ProductTypeProduct,
			Name:         product.Name,
			Slug:         getProductSlug(product.Name, product.BrandID, product.CategoryID),
			CategoryID:   product.CategoryID,
			Image:        product.Image,
			Details:      product.Details,
			Price:        product.Price,
			Status:       db.ProductStatusPending,
			BrandID:      product.BrandID,
			Availability: product.Availability,
		}

		if len(product.Variant) > 0 {
			arg.Variant = []db.ProductUploadVariantParams{}

			for _, variant := range product.Variant {
				variantArg := db.ProductUploadVariantParams{
					Title:     variant.Title,
					MinSelect: variant.MinSelect,
					MaxSelect: variant.MaxSelect,
				}

				// check variants min, max data is valid
				errText := ""
				if variantArg.MinSelect > variantArg.MaxSelect {
					errText = "minSelect must be less than maxSelect"
				} else if variantArg.MinSelect < 0 {
					errText = "minSelect must be greater than 0"
				} else if variantArg.MaxSelect < 0 {
					errText = "maxSelect must be greater than 0"
				} else if len(variant.Items) < int(variantArg.MinSelect) {
					errText = "minSelect must be less than or equal to the number of items"
				} else if len(variant.Items) < int(variantArg.MaxSelect) {
					errText = "maxSelect must be less than or equal to the number of items"
				} else if len(variant.Items) == 0 {
					errText = "variant must have at least 1 item"
				}

				if errText != "" {
					ctx.JSON(http.StatusBadRequest, errorResponse(errors.New(errText)))
					return
				}

				variantArg.Items = []db.ProductUploadParams{}

				for _, item := range variant.Items {
					variantArg.Items = append(variantArg.Items, db.ProductUploadParams{
						Type:         db.ProductTypeVariant,
						Status:       db.ProductStatusPending,
						Name:         item.Name,
						Slug:         getProductVariantSlug(product.Name, item.Name, product.BrandID, product.CategoryID),
						Image:        item.Image,
						Details:      item.Details,
						Price:        item.Price,
						Availability: item.Availability,
					})
				}

				arg.Variant = append(arg.Variant, variantArg)
			}
		}

		args = append(args, arg)
	}

	products, err := server.store.BulkProductUpload(ctx, args)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, products)
}

func (server *Server) getProductDetails(ctx *gin.Context) {

	productId, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	product, err := server.store.FullProduct(ctx, productId)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, product)
}

func (server *Server) getProduct(ctx *gin.Context)                 {}
func (server *Server) getVariantsByProductId(ctx *gin.Context)     {}
func (server *Server) getVariantItemsByVariantId(ctx *gin.Context) {}
func (server *Server) updateProduct(ctx *gin.Context)              {}
func (server *Server) deleteProduct(ctx *gin.Context)              {}
func (server *Server) getProductBySlug(ctx *gin.Context)           {}
func (server *Server) listProducts(ctx *gin.Context)               {}
func (server *Server) countProducts(ctx *gin.Context)              {}
func (server *Server) listProductsByCategoryID(ctx *gin.Context)   {}
func (server *Server) listProductsByBrandID(ctx *gin.Context)      {}
func (server *Server) updateProductVariant(ctx *gin.Context)       {}
func (server *Server) deleteProductVariant(ctx *gin.Context)       {}
func (server *Server) getProductVariant(ctx *gin.Context)          {}
func (server *Server) listProductVariants(ctx *gin.Context)        {}
func (server *Server) createProductVariantItem(ctx *gin.Context)   {}
func (server *Server) deleteProductVariantItem(ctx *gin.Context)   {}
func (server *Server) getProductVariantItem(ctx *gin.Context)      {}
func (server *Server) listProductVariantItems(ctx *gin.Context)    {}
