package storage

import (
	"context"
	"dsa-backend/storage/model"

	"github.com/uptrace/bun"
)

type ProblemStore struct {
	db *bun.DB
}

func (ps ProblemStore) DeleteLectureEntry(context *context.Context, i int64) error {
	_, err := ps.db.NewDelete().Model(&model.Lecture{}).Where("id = ?", i).Exec(*context)
	if err != nil {
		return err
	}
	return nil
}

func (ps ProblemStore) UpdateLectureEntry(context *context.Context, lectureEntryInDB *model.Lecture) error {
	_, err := ps.db.NewUpdate().Model(lectureEntryInDB).Where("id = ?", lectureEntryInDB.ID).Exec(*context)
	if err != nil {
		return err
	}
	return nil
}

func NewProblemStore(db *bun.DB) *ProblemStore {
	return &ProblemStore{
		db: db,
	}
}

func (ps *ProblemStore) GetLectureByID(ctx *context.Context, id int64) (model.Lecture, error) {
	var lecture model.Lecture
	err := ps.db.NewSelect().Model(&lecture).Where("id = ?", id).Scan(*ctx)
	if err != nil {
		return model.Lecture{}, err
	}
	return lecture, nil
}

func (ps *ProblemStore) CreateLectureEntry(ctx *context.Context, lec *model.Lecture) error {
	_, err := ps.db.NewInsert().Model(lec).Exec(*ctx)
	if err != nil {
		return err
	}
	return nil
}
