package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"
)

const (
	dbDriver = "postgres"
	dbSource = "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable"
)

var testQueries *Queries
var testDBConn *sql.DB

func TestMain(m *testing.M) {
	var err error
	testDBConn, err = sql.Open(dbDriver, dbSource)

	if err != nil {
		log.Fatal("Error connecting to the database: ", err)
	}

	testQueries = New(testDBConn)

	os.Exit(m.Run())
}
