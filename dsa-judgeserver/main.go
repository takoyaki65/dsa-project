package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"github.com/uptrace/bun/extra/bundebug"
)

func main() {
	db_user := "dsa_app"
	db_password := read_db_password()
	// TODO: modify this for production
	db_host := "127.0.0.1:5432"
	dsn := fmt.Sprintf("postgres://%s:%s@%s/dsa_db?sslmode=disable", db_user, db_password, db_host)

	// initialize connection
	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dsn)))

	// create db instance
	db := bun.NewDB(sqldb, pgdialect.New())
	// For debugging purpose, print all queries to stdout.
	db.AddQueryHook(bundebug.NewQueryHook(bundebug.WithVerbose(true)))

	// TODO launch go rountines for workers, and monitoring job queues.
}

func read_db_password() string {
	data, err := os.ReadFile("../config/db_app_password.txt")
	if err != nil {
		log.Fatalf("failed to read db password: %v", err)
	}
	return strings.TrimSpace(string(data))
}
