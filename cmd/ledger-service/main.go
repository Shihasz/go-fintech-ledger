package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Shihasz/go-fintech-ledger/internal/common/config"
	"github.com/Shihasz/go-fintech-ledger/internal/ledger/api"
	"github.com/Shihasz/go-fintech-ledger/internal/ledger/db"
	_ "github.com/lib/pq"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}

	// Open connection to the database
	conn, err := sql.Open("postgres", cfg.DBSource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}
	// Manually close the connection during shutdown

	// Initialize the sqlc store and Gin server
	store := db.NewStore(conn)
	server, err := api.NewServer(cfg, store)
	if err != nil {
		log.Fatal("cannot create server:", err)
	}

	// Create a custom HTTP server
	srv := &http.Server{
		Addr:    cfg.ServerAddress,
		Handler: server.GetRouter(),
	}

	// Start the server in a separate Goroutine so it doesn't block the main thread
	go func() {
		log.Printf("🚀 Starting Ledger Service on port %s", cfg.ServerAddress)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// Wait for interrupt signal to gracefully shut down the server
	quit := make(chan os.Signal, 1)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT (Ctrl+C)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// This blocks until a signal is received
	<-quit
	log.Println("🛑 Shutting down server gracefully...")

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	// Clean up database connection
	log.Println("Closing database connection...")
	conn.Close()

	log.Println("✅ Server exiting cleanly")
}
