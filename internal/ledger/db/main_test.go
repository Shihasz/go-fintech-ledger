package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"
)

// testQueries is a global variable we can use in all our database tests
var testQueries *Queries
var testDB *sql.DB

func TestMain(m *testing.M) {
	// Connect to our local Docker database
	dbSource := "postgresql://root:secret@127.0.0.1:5432/fintech_ledger?sslmode=disable"
	var err error
	testDB, err = sql.Open("postgres", dbSource)
	if err != nil {
		log.Fatal("cannot connect to test db:", err)
	}

	testQueries = New(testDB)

	// Run the tests and exit
	os.Exit(m.Run())
}
