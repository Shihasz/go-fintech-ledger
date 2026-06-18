package api

import (
	"github.com/Shihasz/go-fintech-ledger/internal/common/config"
	"github.com/Shihasz/go-fintech-ledger/internal/common/token"
	"github.com/Shihasz/go-fintech-ledger/internal/ledger/db"
	"github.com/gin-gonic/gin"
)

// Server serves HTTP requests for our auth service.
type Server struct {
	config     config.Config
	store      *db.Queries
	tokenMaker *token.JWTMaker
	router     *gin.Engine
}

// NewServer creates a new HTTP server and sets up routing.
func NewServer(config config.Config, store *db.Queries) (*Server, error) {
	tokenMaker, err := token.NewJWTMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, err
	}

	server := &Server{
		config:     config,
		store:      store,
		tokenMaker: tokenMaker,
	}

	router := gin.Default()

	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "UP", "service": "auth"})
	})

	router.POST("/users", server.createUser)
	router.POST("/users/login", server.loginUser)

	server.router = router
	return server, nil
}

// GetRouter returns the underlying Gin engine.
func (server *Server) GetRouter() *gin.Engine {
	return server.router
}
