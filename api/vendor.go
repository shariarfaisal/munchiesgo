package api

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	db "github.com/Munchies-Engineering/backend/db/sqlc"
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

type createVendorRequest struct {
	Name        string      `json:"name" binding:"required"`
	Email       string      `json:"email" binding:"required,email"`
	Phone       string      `json:"phone" binding:"required"`
	PaymentInfo interface{} `json:"paymentInfo" binding:"required"`
	SocialLinks interface{} `json:"socialLinks" binding:"required"`
}

func (server *Server) createVendor(ctx *gin.Context) {
	var req createVendorRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	paymentInfo, err := json.Marshal(req.PaymentInfo)
	if err != nil {

		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	socialLinks, err := json.Marshal(req.SocialLinks)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.CreateVendorParams{
		Name:        req.Name,
		Email:       req.Email,
		Phone:       req.Phone,
		PaymentInfo: paymentInfo,
		SocialLinks: socialLinks,
	}

	vendor, err := server.store.CreateVendor(ctx, arg)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			fmt.Println(pqErr.Code.Name())
			switch pqErr.Code.Name() {
			case "unique_violation":
				ctx.JSON(http.StatusForbidden, errorResponse(err))
				return
			}
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, vendor)
}

type updateVendorRequest struct {
	Name        string      `json:"name" binding:"required"`
	Email       string      `json:"email" binding:"required,email"`
	Phone       string      `json:"phone" binding:"required"`
	PaymentInfo interface{} `json:"paymentInfo"`
	SocialLinks interface{} `json:"socialLinks"`
}

type updateVendorId struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (server *Server) updateVendor(ctx *gin.Context) {
	var req updateVendorRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	reqURI := updateVendorId{}
	if err := ctx.ShouldBindUri(&reqURI); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	paymentInfo, err := json.Marshal(req.PaymentInfo)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	socialLinks, err := json.Marshal(req.SocialLinks)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.UpdateVendorParams{
		ID:          reqURI.ID,
		Name:        req.Name,
		Email:       req.Email,
		Phone:       req.Phone,
		PaymentInfo: paymentInfo,
		SocialLinks: socialLinks,
	}
	vendor, err := server.store.UpdateVendor(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, vendor)
}

type getVendorRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (server *Server) getVendor(ctx *gin.Context) {
	var req getVendorRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	vendor, err := server.store.GetVendor(ctx, req.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, vendor)
}

type listVendorsRequest struct {
	PageSize int32 `form:"pageSize" binding:"required,min=1,max=100"`
	PageID   int32 `form:"pageId" binding:"required,min=1"`
}

func (server *Server) listVendors(ctx *gin.Context) {
	var req listVendorsRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.ListVendorsParams{
		Limit:  req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	}

	list, err := server.store.ListVendors(ctx, arg)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, list)
}

type searchVendorsRequest struct {
	PageSize int32  `form:"pageSize" binding:"required,min=1,max=100"`
	PageID   int32  `form:"pageId" binding:"required,min=1"`
	Name     string `form:"name" binding:"required"`
}

func (server *Server) searchVendors(ctx *gin.Context) {
	var req searchVendorsRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.SearchVendorsParams{
		Name:   req.Name,
		Limit:  req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	}
	vendors, err := server.store.SearchVendors(ctx, arg)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, vendors)
}

type deleteVendorRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (server *Server) deleteVendor(ctx *gin.Context) {
	var req deleteVendorRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	err := server.store.DeleteVendor(ctx, req.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "deleted"})
}
