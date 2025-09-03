package storage

import (
	"context"
	"dsa-backend/storage/model"
	"dsa-backend/storage/model/requeststatus"

	"github.com/uptrace/bun"
)

type RequestStore struct {
	db *bun.DB
}

func (r RequestStore) UpdateGradingRequestStatus(context *context.Context, id int64, e requeststatus.State) error {
	_, err := r.db.NewUpdate().Model(&model.GradingRequest{}).Set("result = ?", int64(e)).Where("id = ?", id).Exec(*context)
	return err
}

func (r RequestStore) UpdateValidationRequestStatus(context *context.Context, id int64, status_id requeststatus.State) error {
	_, err := r.db.NewUpdate().Model(&model.ValidationRequest{}).Set("result = ?", int64(status_id)).Where("id = ?", id).Exec(*context)
	return err
}

func (r *RequestStore) RegisterValidationRequest(ctx *context.Context, request *model.ValidationRequest) error {
	_, err := r.db.NewInsert().Model(request).Returning("id"). // Return auto-incremented ID
									Exec(*ctx)
	return err
}

func (r *RequestStore) RegisterOrUpdateGradingRequest(ctx *context.Context, request *model.GradingRequest) error {
	_, err := r.db.NewInsert().Model(request).On("CONFLICT (lecture_id,problem_id,usercode,submission_ts) DO UPDATE"). // Upsert
																Returning("id"). // Return auto-incremented ID
																Exec(*ctx)
	return err
}

func NewRequestStore(db *bun.DB) *RequestStore {
	return &RequestStore{
		db: db,
	}
}
