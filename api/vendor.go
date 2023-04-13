package api

import (
	"encoding/json"
	"net/http"

	db "github.com/Munchies-Engineering/backend/db/sqlc"
	"github.com/gin-gonic/gin"
)

type createVendorRequest struct {
	Name        string      `json:"name" binding:"required"`
	Email       string      `json:"email" binding:"required,email"`
	Phone       string      `json:"phone" binding:"required"`
	PaymentInfo interface{} `json:"paymentInfo"`
	SocialLinks interface{} `json:"socialLinks"`
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
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, vendor)
}

func (server *Server) updateVendor(ctx *gin.Context) {

}

func (server *Server) getVendor(ctx *gin.Context) {

}

func (server *Server) listVendors(ctx *gin.Context) {

}

func (server *Server) searchVendors(ctx *gin.Context) {

}

func (server *Server) deleteVendor(ctx *gin.Context) {

}
