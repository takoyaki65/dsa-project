package storage

import (
	"context"
	"dsa-backend/storage/model"

	"github.com/uptrace/bun"
)

type JobQueueStore struct {
	db *bun.DB
}

func (j *JobQueueStore) InsertJob(context *context.Context, job *model.JobQueue) error {
	_, err := j.db.NewInsert().Model(job).Exec(*context)
	return err
}

func NewJobQueueStore(db *bun.DB) *JobQueueStore {
	return &JobQueueStore{
		db: db,
	}
}
