package database

import (
	"context"

	"github.com/takoyaki65/dsa-project/database/model"
	"github.com/takoyaki65/dsa-project/database/model/requeststatus"
	"github.com/uptrace/bun"
)

type RequestStore struct {
	db *bun.DB
}

func (r RequestStore) UpdateGradingRequestStatus(ctx context.Context, id int64, e requeststatus.State) error {
	_, err := r.db.NewUpdate().Model(&model.GradingRequest{}).Set("result = ?", int64(e)).Where("id = ?", id).Exec(ctx)
	return err
}

func (r RequestStore) UpdateValidationRequestStatus(ctx context.Context, id int64, status_id requeststatus.State) error {
	_, err := r.db.NewUpdate().Model(&model.ValidationRequest{}).Set("result = ?", int64(status_id)).Where("id = ?", id).Exec(ctx)
	return err
}

func (r *RequestStore) RegisterValidationRequest(ctx context.Context, request *model.ValidationRequest) error {
	_, err := r.db.NewInsert().Model(request).Returning("id"). // Return auto-incremented ID
									Exec(ctx)
	return err
}

func (r *RequestStore) RegisterOrUpdateGradingRequest(ctx context.Context, request *model.GradingRequest) error {
	_, err := r.db.NewInsert().Model(request).On("CONFLICT (lecture_id,problem_id,usercode,submission_ts) DO UPDATE"). // Upsert
																Returning("id"). // Return auto-incremented ID
																Exec(ctx)
	return err
}

func (r *RequestStore) UpdateResultOfValidationRequest(ctx context.Context, id int64, result_id requeststatus.State, Log model.RequestLog) error {
	_, err := r.db.NewUpdate().Model(&model.ValidationRequest{}).
		Set("result = ?", int64(result_id)).
		Set("log = ?", Log).
		Where("id = ?", id).
		Exec(ctx)
	return err
}

func (r *RequestStore) UpdateResultOfGradingRequest(ctx context.Context, id int64, result_id requeststatus.State, Log model.RequestLog) error {
	_, err := r.db.NewUpdate().Model(&model.GradingRequest{}).
		Set("result = ?", int64(result_id)).
		Set("log = ?", Log).
		Where("id = ?", id).
		Exec(ctx)
	return err
}

// GetValidationResults retrieves validation results for a user filtered by allowed lecture IDs, with pagination support.
// NOTE: This function does not utilize OFFSET, because OFFSET can be inefficient for large datasets.
// Instead, it uses the "last" parameter to fetch results with IDs less than the provided value.
func (r *RequestStore) GetValidationResults(ctx context.Context, usercode int64, lecture_ids []int64, last int64, limit int) ([]model.ValidationRequest, error) {
	var results []model.ValidationRequest
	err := r.db.NewSelect().Model(&results).
		Where("usercode = ?", usercode).
		Where("lecture_id IN (?)", bun.In(lecture_ids)).
		Where("id <= ?", last).
		Order("id DESC").
		Limit(limit).
		Scan(ctx)

	return results, err
}

func (r *RequestStore) GetValidationResultByID(ctx context.Context, id int64) (*model.ValidationRequest, error) {
	var result model.ValidationRequest
	err := r.db.NewSelect().Model(&result).Relation("FileLocation").Where("validation_request.id = ?", id).Scan(ctx)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func NewRequestStore(db *bun.DB) *RequestStore {
	return &RequestStore{
		db: db,
	}
}
