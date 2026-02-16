package main

import (
	"log"
	"net/http"

	"github.com/Shihasz/go-fintech-ledger/internal/common/config"
	"github.com/gin-gonic/gin"
)

func main() {
	// Load Configuration
	cfg, err := config.LoadConfig(".")
	if err != nil {
		log.Fatal("Cannot load config:", err)
	}

	// Initialize Gin router
	r := gin.Default()

	// Define a simple Health Check endpoint
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"service": "auth-service",
			"port":    cfg.ServerAddress,
		})
	})

	// Start the server on port 8080
	log.Printf("Starting Auth Service on %s...", cfg.ServerAddress)
	if err := r.Run(cfg.ServerAddress); err != nil {
		log.Fatal("Failed to start server: ", err)
	}
}
