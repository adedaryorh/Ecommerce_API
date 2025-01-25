package db_test

import (
	"database/sql"
	db "github.com/adedaryorh/ecommerceapi/db/sqlc"
	"github.com/adedaryorh/ecommerceapi/utils"
	_ "github.com/lib/pq"
	"log"
	"os"
	"testing"
)

var testQuery *db.Queries

func TestMain(m *testing.M) {
	config, err := utils.LoadConfig("../..")
	if err != nil {
		log.Fatal("cannot load config file")
	}
	conn, err := sql.Open(config.DBdriver, config.DB_source)

	if err != nil {
		log.Fatal("error connecting to postgres:", err)
	}

	testQuery = db.New(conn)
	os.Exit(m.Run())
}
