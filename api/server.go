package api

import (
	db "github.com/Munchies-Engineering/backend/db/sqlc"
	"github.com/Munchies-Engineering/backend/util"
	"github.com/gin-gonic/gin"
)

type Server struct {
	config util.Config
	store  db.Store
	router *gin.Engine
}

func NewServer(config util.Config, store db.Store) *Server {
	server := &Server{
		config: config,
		store:  store,
	}

	server.routerSetup()

	return server
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func (server *Server) routerSetup() {
	router := gin.Default()

	// vendor routes
	vendor := router.Group("/api/vendor")
	{
		vendor.POST("/", server.createVendor)
		vendor.GET("/", server.searchVendors)
		vendor.GET("/search", server.searchVendors)
		vendor.GET("/:id", server.getVendor)
		vendor.DELETE("/:id", server.deleteVendor)
		vendor.PUT("/:id", server.updateVendor)
	}

	// vendor user routes
	vendorUser := router.Group("/api/vendor_user")
	{
		vendorUser.POST("/login", server.loginVendorUser)
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
