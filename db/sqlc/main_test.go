package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"
	"github.com/santinofajardo/simpleBank/util"
)

var testQueries *Queries
var testDBConn *sql.DB

func TestMain(m *testing.M) {
	config, err := util.LoadConfig("../..")
	if err != nil {
		log.Fatal("Error loading config: ", err)
	}
	testDBConn, err = sql.Open(config.DBDriver, config.DBSource)

	if err != nil {
		log.Fatal("Error connecting to the database: ", err)
	}

	testQueries = New(testDBConn)

	os.Exit(m.Run())
}
