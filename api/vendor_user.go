package api

import (
	"net/http"

	db "github.com/Munchies-Engineering/backend/db/sqlc"
	"github.com/Munchies-Engineering/backend/token"
	"github.com/Munchies-Engineering/backend/util"
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
	"github.com/pkg/errors"
)

const (
	VendorTypeAdmin      = "admin"
	VendorTypeManager    = "manager"
	VendorTypeReporting  = "reporting"
	VendorTypeOperations = "operations"
	VendorTypeFinance    = "finance"
	VendorTypeSupport    = "support"
)

type createVendorUserRequest struct {
	Username string `json:"username" binding:"required"`
	FullName string `json:"fullName" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
	Role     string `json:"role" binding:"required,oneof=admin manager reporting operations finance support"`
}

type vendorUserResponseType struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	FullName string `json:"fullName"`
	Email    string `json:"email"`
	Role     string `json:"role"`
	VendorID int64  `json:"vendorId"`
}

func vendorUserResponse(user db.VendorUser) *vendorUserResponseType {
	return &vendorUserResponseType{
		ID:       user.ID,
		Username: user.Username,
		FullName: user.FullName,
		Email:    user.Email,
		Role:     user.Role,
		VendorID: user.VendorID,
	}
}

func (server *Server) getVendorUserByID(ctx *gin.Context, userId int64) (db.VendorUser, error) {
	user, err := server.store.GetVendorUser(ctx, userId)
	if err != nil {
		return db.VendorUser{}, err
	}

	return user, nil
}

func (server *Server) createVendorUser(ctx *gin.Context) {
	var req createVendorUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	authPayload := ctx.MustGet(authorizationPayload).(*token.Payload)
	user, err := server.getVendorUserByID(ctx, authPayload.UserID)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	hashedPassword, err := util.HashPassword(req.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(errors.New("error hashing password")))
		return
	}

	arg := db.CreateVendorUserParams{
		Username:       req.Username,
		FullName:       req.FullName,
		Email:          req.Email,
		HashedPassword: hashedPassword,
		Role:           req.Role,
		VendorID:       user.VendorID,
	}

	newUser, err := server.store.CreateVendorUser(ctx, arg)
	if err != nil {
		if err := err.(*pq.Error); err != nil {
			switch err.Code.Name() {
			case "unique_violation":
				ctx.JSON(http.StatusConflict, errorResponse(errors.New("username or email already exists")))
				return
			}
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, vendorUserResponse(newUser))
}

type loginVendorUser struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (server *Server) loginVendorUser(ctx *gin.Context) {
	var req loginVendorUser
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	user, err := server.store.GetVendorUserByUsername(ctx, req.Username)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	err = util.CheckPassword(req.Password, user.HashedPassword)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	accessToken, err := server.tokenMaker.CreateToken(user.ID, user.VendorID, user.Role, server.config.TokenDuration)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"accessToken": accessToken,
		"user":        vendorUserResponse(user),
	})
}

func (server *Server) updateVendorUser(ctx *gin.Context) {
	// TODO: Implement
}

func (server *Server) getVendorUser(ctx *gin.Context) {
	// TODO: Implement
}

func (server *Server) listVendorUsers(ctx *gin.Context) {
	// TODO: Implement
}

func (server *Server) deleteVendorUser(ctx *gin.Context) {
	// TODO: Implement
}
