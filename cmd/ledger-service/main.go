package main

import (
	"database/sql"
	"log"

	"github.com/Shihasz/go-fintech-ledger/internal/common/config"
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

	// Ping the database to ensure the connection is actually valid
	if err = conn.Ping(); err != nil {
		log.Fatal("cannot ping db:", err)
	}

	log.Println("✅ Successfully connected to the FinTech Ledger Database!")
}
