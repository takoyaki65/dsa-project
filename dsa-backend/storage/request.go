package storage

import (
	"context"
	"dsa-backend/storage/model"

	"github.com/uptrace/bun"
)

type RequestStore struct {
	db *bun.DB
}

func (r RequestStore) RegisterValidationRequest(ctx *context.Context, request *model.ValidationRequest) error {
	_, err := r.db.NewInsert().Model(request).Exec(*ctx)
	return err
}

func NewRequestStore(db *bun.DB) *RequestStore {
	return &RequestStore{
		db: db,
	}
}
