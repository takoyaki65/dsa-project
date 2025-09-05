package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"log/slog"
	"os"
	"strings"
	"time"

	"github.com/takoyaki65/dsa-project/database"
	"github.com/takoyaki65/dsa-project/database/model/queuestatus"
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

	jobQueueStore := database.NewJobQueueStore(db)

	ctx := context.Background()
	textHandler := slog.NewTextHandler(os.Stdout, nil)
	logger := slog.New(textHandler)

	// Create Docker Client
	jobExecutor, err := NewJobExecutor()
	if err != nil {
		logger.Error("Failed to create job executor", slog.String("error", err.Error()))
		return
	}

	for {
		// Infinite loop to keep the program running

		// Fetch Pending tasks from JobQueue
		jobs, err := jobQueueStore.FetchJobs(&ctx, queuestatus.Pending, 1)

		if err != nil {
			logger.Error("Failed to fetch jobs", slog.String("error", err.Error()))
			continue
		}

		if len(jobs) == 0 {
			logger.Info("No pending jobs found. Sleeping for a while...")
			// Sleep for a while before checking again
			time.Sleep(5 * time.Second) // Uncomment this line to add a delay
			continue
		}

		// Pick the first job
		job := jobs[0]
		logger.Info("Processing job", slog.Int64("job_id", job.ID))

		// Execute the job
		jobExecutor.execute_job(&job)
	}
}

func read_db_password() string {
	data, err := os.ReadFile("../config/db_app_password.txt")
	if err != nil {
		log.Fatalf("failed to read db password: %v", err)
	}
	return strings.TrimSpace(string(data))
}
