package storage

import (
	"context"
	"dsa-judgeserver/storage/model"
	"dsa-judgeserver/storage/model/queuestatus"

	"github.com/uptrace/bun"
)

type JobQueueStore struct {
	db *bun.DB
}

func (j *JobQueueStore) FetchPendingJobsAndUpdateToFetched(ctx *context.Context, limit int) ([]model.JobQueue, error) {
	var jobs []model.JobQueue
	err := j.db.NewSelect().Model(&jobs).Where("status = ?", queuestatus.Pending).Limit(limit).Scan(*ctx)
	if err != nil {
		return nil, err
	}

	_, err = j.db.NewUpdate().Model(&jobs).Set("status = ?", queuestatus.Fetched).WherePK().Exec(*ctx)
	if err != nil {
		return nil, err
	}

	return jobs, nil
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
