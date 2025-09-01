package model

import (
	"context"
	"errors"
	"time"

	"github.com/uptrace/bun"
)

type JobQueue struct {
	bun.BaseModel `bun:"table:jobqueue"`

	ID          int64     `bun:"id,pk,autoincrement" json:"id"`
	RequestType string    `bun:"request_type,notnull" json:"request_type"` // "validation" or "grading"
	RequestID   int64     `bun:"request_id,notnull" json:"request_id"`
	Status      string    `bun:"status,notnull" json:"status"` // "pending", "processing", "done"
	CreatedAt   time.Time `bun:"created_at,notnull" json:"created_at"`
	FileDir     string    `bun:"file_dir,notnull" json:"file_dir"`
	ResultDir   string    `bun:"result_dir,notnull" json:"result_dir"`
	Detail      JobDetail `bun:"detail,notnull,type:jsonb" json:"detail"`
}

type JobDetail struct {
	TimeMS     int64      `json:"time_ms"`
	MemoryMB   int64      `json:"memory_mb"`
	TestFiles  int64      `json:"test_files"`
	BuildTasks []TestCase `json:"build"`
	JudgeTasks []TestCase `json:"judge"`
}

type ResultQueue struct {
	bun.BaseModel `bun:"table:resultqueue"`

	ID        int64      `bun:"id,pk,autoincrement" json:"id"`
	JobID     int64      `bun:"job_id,notnull" json:"job_id"`
	CreatedAt time.Time  `bun:"created_at,notnull" json:"created_at"`
	Result    RequestLog `bun:"result,notnull,type:jsonb" json:"result"`
}

var _ bun.BeforeAppendModelHook = (*JobQueue)(nil)

func (jq *JobQueue) BeforeAppendModel(ctx context.Context, query bun.Query) error {
	switch query.(type) {
	case *bun.InsertQuery, *bun.UpdateQuery:
		// remove fraction less than seconds (milliseconds, microseconds, ...)
		jq.CreatedAt = jq.CreatedAt.Truncate(time.Second)
		// Check RequestType is "validation" or "grading"
		if jq.RequestType != "validation" && jq.RequestType != "grading" {
			return errors.New("invalid request type: " + jq.RequestType)
		}
		// Check Status is "pending" or "processing" or "done"
		if jq.Status != "pending" && jq.Status != "processing" && jq.Status != "done" {
			return errors.New("invalid status: " + jq.Status)
		}
	}
	return nil
}

var _ bun.BeforeAppendModelHook = (*ResultQueue)(nil)

func (rq *ResultQueue) BeforeAppendModel(ctx context.Context, query bun.Query) error {
	switch query.(type) {
	case *bun.InsertQuery, *bun.UpdateQuery:
		// remove fraction less than seconds (milliseconds, microseconds, ...)
		rq.CreatedAt = rq.CreatedAt.Truncate(time.Second)
	}
	return nil
}
