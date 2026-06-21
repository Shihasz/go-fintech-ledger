package api

import (
	"errors"
	"net/http"
	"strings"

	"github.com/Shihasz/go-fintech-ledger/internal/common/token"
	"github.com/gin-gonic/gin"
)

const (
	authorizationHeaderKey  = "authorization"
	authorizationTypeBearer = "bearer"
	authorizationPayloadKey = "authorization_payload"
)

// authMiddleware creates a gin middleware for authorization
func authMiddleware(tokenMaker *token.JWTMaker) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Get the Authorization header
		authorizationHeader := ctx.GetHeader(authorizationHeaderKey)
		if len(authorizationHeader) == 0 {
			err := errors.New("authorization header is not provided")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}

		// Split the header to get the token (Format: "Bearer <token>")
		fields := strings.Fields(authorizationHeader)
		if len(fields) < 2 {
			err := errors.New("invalid authorization header format")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}

		authorizationType := strings.ToLower(fields[0])
		if authorizationType != authorizationTypeBearer {
			err := errors.New("unsupported authorization type")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}

		accessToken := fields[1]

		// Verify the token using our JWTMaker
		payload, err := tokenMaker.VerifyToken(accessToken)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}

		// Store the payload in the context so the next handler can access the username
		ctx.Set(authorizationPayloadKey, payload)
		ctx.Next()
	}
}
