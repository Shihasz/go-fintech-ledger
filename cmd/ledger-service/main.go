package main

import (
	"database/sql"
	"log"

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
	defer conn.Close()

	// Initialize the sqlc store
	store := db.New(conn)

	// Initialize and start the Gin server
	server := api.NewServer(store)
	log.Printf("🚀 Starting Ledger Service on port %s", cfg.ServerAddress)

	err = server.Start(cfg.ServerAddress)
	if err != nil {
		log.Fatal("cannot start server:", err)
	}
}
