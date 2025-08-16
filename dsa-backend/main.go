package main

import (
	"context"
	"database/sql"
	"dsa-backend/handler"
	"dsa-backend/router"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
)

func main() {
	r := router.New()

	db_user := "dsa_app"
	db_password := read_db_password()
	dsn := fmt.Sprintf("postgres://%s:%s@127.0.0.1:5432/dsa_db?sslmode=disable", db_user, db_password)

	// initialize connection
	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dsn)))

	// create db instance
	db := bun.NewDB(sqldb, pgdialect.New())

	// create handler
	h := handler.NewHandler(db)

	// register api routes to router
	h.RegisterRoutes(r)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	go func() {
		port := ":8000"

		if err := r.Start(port); err != nil && err != http.ErrServerClosed {
			r.Logger.Fatal("shutting down the server: ", err)
		}
	}()

	select {
	case <-ctx.Done():
		r.Logger.Info("Interrupt signal received, shutting down server...")
	case <-h.GetShutdownChan():
		r.Logger.Info("Admin user created, initiating graceful shutdown")
		r.Logger.Info("Please restart the server to enable all endpoints")
	}

	// Gracefully shutdown the server with a timeout of 10 seconds.
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := r.Shutdown(shutdownCtx); err != nil {
		r.Logger.Fatal("Server forced to shutdown: ", err)
	}

	r.Logger.Info("Server gracefully stopped")
}

func read_db_password() string {
	data, err := os.ReadFile("../config/db_app_password.txt")
	if err != nil {
		log.Fatalf("failed to read db password: %v", err)
	}
	return strings.TrimSpace(string(data))
}
