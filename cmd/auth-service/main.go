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

	"github.com/Shihasz/go-fintech-ledger/internal/auth/api"
	"github.com/Shihasz/go-fintech-ledger/internal/common/config"
	"github.com/Shihasz/go-fintech-ledger/internal/ledger/db"
	_ "github.com/lib/pq"
)

func main() {
	cfg, err := config.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}

	conn, err := sql.Open("postgres", cfg.DBSource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	store := db.New(conn)
	server, err := api.NewServer(cfg, store)
	if err != nil {
		log.Fatal("cannot create server:", err)
	}

	srv := &http.Server{
		Addr:    cfg.ServerAddress,
		Handler: server.GetRouter(),
	}

	go func() {
		log.Printf("🚀 Starting Auth Service on port %s", cfg.ServerAddress)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("🛑 Shutting down server gracefully...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Closing database connection...")
	conn.Close()

	log.Println("✅ Server exiting cleanly")
}
