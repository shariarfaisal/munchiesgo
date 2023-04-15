package api

import (
	db "github.com/Munchies-Engineering/backend/db/sqlc"
	"github.com/Munchies-Engineering/backend/token"
	"github.com/Munchies-Engineering/backend/util"
	"github.com/gin-gonic/gin"
)

type Server struct {
	config     util.Config
	store      db.Store
	router     *gin.Engine
	tokenMaker token.Maker
}

func NewServer(config util.Config, store db.Store) (*Server, error) {

	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, err
	}

	server := &Server{
		config:     config,
		store:      store,
		tokenMaker: tokenMaker,
	}

	server.routerSetup()

	return server, nil
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func (server *Server) routerSetup() {
	router := gin.Default()

	// vendor routes
	vendor := router.Group("/api/vendor", authMiddleware(server.tokenMaker))
	{
		vendor.POST("/", server.createVendor)
		vendor.GET("/", server.listVendors)
		vendor.GET("/search", server.searchVendors)
		vendor.GET("/:id", server.getVendor)
		vendor.DELETE("/:id", server.deleteVendor)
		vendor.PUT("/:id", server.updateVendor)
	}

	// vendor user routes
	router.POST("/api/vendor_user/login", server.loginVendorUser)
	vendorUser := router.Group("/api/vendor_user", authMiddleware(server.tokenMaker))
	{
		vendorUser.POST("/", server.createVendorUser)
		vendorUser.GET("/", server.listVendorUsers)
		vendorUser.GET("/:id", server.getVendorUser)
		vendorUser.DELETE("/:id", server.deleteVendorUser)
		vendorUser.PUT("/:id", server.updateVendorUser)
	}

	server.router = router
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
