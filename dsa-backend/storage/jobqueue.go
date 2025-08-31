package storage

import "github.com/uptrace/bun"

type JobQueueStore struct {
	db *bun.DB
}

func NewJobQueueStore(db *bun.DB) *JobQueueStore {
	return &JobQueueStore{
		db: db,
	}
}
