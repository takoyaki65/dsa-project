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
	Status      string    `bun:"status,notnull" json:"status"` // "pending", "in_progress"
	CreatedAt   time.Time `bun:"created_at,notnull" json:"created_at"`
	FullMode    bool      `bun:"fullmode,notnull" json:"full_mode"`
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
		// Check Status is "pending" or "in_progress"
		if jq.Status != "pending" && jq.Status != "in_progress" {
			return errors.New("invalid status: " + jq.Status)
		}
	}
	return nil
}
