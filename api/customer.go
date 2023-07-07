package api

import (
	"net/http"

	db "github.com/Munchies-Engineering/backend/db/sqlc"
	"github.com/gin-gonic/gin"
)

const (
	CustomerStatusActive   = "ACTIVE"
	CustomerStatusInactive = "INACTIVE"
	CustomerStatusBlocked  = "BLOCKED"
	CustomerStatusDeleted  = "DELETED"
)

type customerSignUpRequest struct {
	Name  string `json:"name" binding:"required"`
	Phone string `json:"phone" binding:"required"`
	Email string `json:"email" binding:"omitempty,email"`
	Image string `json:"image"`
	Nid   string `json:"nid"`
}

func (server *Server) customerSignUp(ctx *gin.Context) {
	var req customerSignUpRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(400, errorResponse(err))
		return
	}

	arg := db.CreateCustomerParams{
		Name:          req.Name,
		Phone:         req.Phone,
		Email:         req.Email,
		Image:         req.Image,
		EmailVerified: false,
		Nid:           req.Nid,
		Status:        CustomerStatusActive,
	}

	customer, err := server.store.CreateCustomer(ctx, arg)
	if err != nil {
		ctx.JSON(500, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusCreated, customer)
}

type addCustomerAddressRequest struct {
	CustomerID int64  `json:"customer_id" binding:"required"`
	Label      string `json:"label" binding:"required"`
	Address    string `json:"address" binding:"required"`
	GeoPoint   Point  `json:"geo_point" binding:"required"`
	Apartment  string `json:"apartment"`
	Area       string `json:"area"`
	Floor      string `json:"floor"`
}

func (server *Server) addCustomerAddress(ctx *gin.Context) {
	var req addCustomerAddressRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(400, errorResponse(err))
		return
	}

	arg := db.CreateCustomerAddressParams{
		CustomerID: req.CustomerID,
		Label:      req.Label,
		Address:    req.Address,
		GeoPoint:   req.GeoPoint,
		Apartment:  req.Apartment,
		Area:       req.Area,
		Floor:      req.Floor,
	}

	customerAddress, err := server.store.CreateCustomerAddress(ctx, arg)
	if err != nil {
		ctx.JSON(500, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusCreated, customerAddress)
}

func (server *Server) customerGetOtp(ctx *gin.Context) {

}

func (server *Server) customerValidateOtp(ctx *gin.Context) {

}

func (server *Server) customerLogin(ctx *gin.Context) {

}

func (server *Server) customerProfile(ctx *gin.Context) {

}

func (server *Server) GetCustomer(ctx *gin.Context) {

}

func (server *Server) ListCustomer(ctx *gin.Context) {

}
