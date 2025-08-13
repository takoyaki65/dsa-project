package main

import (
	"database/sql"
	"dsa-backend/handler"
	"dsa-backend/router"
	"log"
	"os"
	"strings"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
)

func main() {
	r := router.New()

	db_user := "dsa_app"
	db_password := read_db_password()
	dsn := "postgres://" + db_user + ":" + db_password + "@db:5432/dsa_db"

	// initialize connection
	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dsn)))

	// create db instance
	db := bun.NewDB(sqldb, pgdialect.New())

	// create handler
	h := handler.NewHandler(db)

	// register api routes to router
	h.RegisterRoutes(r)

	r.Logger.Fatal(r.Start(":15050")) // Start the Echo server on port 15050
}

func read_db_password() string {
	data, err := os.ReadFile("/run/secrets/db_app_password")
	if err != nil {
		log.Fatalf("failed to read db password: %v", err)
	}
	return strings.TrimSpace(string(data))
}
