package api

import (
	"errors"
	"net/http"
	"strings"

	"github.com/Munchies-Engineering/backend/token"
	"github.com/gin-gonic/gin"
)

const (
	authorizationHeader     = "authorization"
	authorizationTypeBearer = "bearer"
	authorizationPayload    = "authorization_payload"
)

func authMiddleware(maker token.Maker) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader(authorizationHeader)
		if header == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(errors.New("authorization header is empty")))
			return
		}

		fields := strings.Split(header, " ")
		if len(fields) != 2 {
			c.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(errors.New("authorization header is invalid")))
			return
		}

		authorizationType := strings.ToLower(fields[0])
		if authorizationType != authorizationTypeBearer {
			c.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(errors.New("authorization type is not bearer")))
			return
		}

		accessToken := fields[1]
		payload, err := maker.VerifyToken(accessToken)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		c.Set(authorizationPayload, payload)
	}
}
