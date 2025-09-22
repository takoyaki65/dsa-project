package database

import (
	"context"

	"github.com/takoyaki65/dsa-project/database/model"
	"github.com/uptrace/bun"
)

type ProblemStore struct {
	db *bun.DB
}

func (ps *ProblemStore) FetchResourcePath(ctx context.Context, lectureID int64, problemID int64) (string, error) {
	// Fetch the resource path for the given lectureID and problemID
	// only select FileLocation.path
	var path string
	err := ps.db.NewSelect().
		Model((*model.Problem)(nil)).
		Column("filelocation.path").
		Join("JOIN filelocation ON problem.resource_location_id = filelocation.id").
		Where("problem.lecture_id = ? AND problem.problem_id = ?", lectureID, problemID).
		Scan(ctx, &path)
	if err != nil {
		return "", err
	}
	return path, nil
}

func (ps *ProblemStore) GetLectureAndAllProblems(ctx context.Context, d int64) (model.Lecture, error) {
	var lecture model.Lecture
	err := ps.db.NewSelect().Model(&lecture).
		Relation("Problems").
		Where("id = ?", d).
		Scan(ctx)
	if err != nil {
		return model.Lecture{}, err
	}
	return lecture, nil
}

func (ps *ProblemStore) GetProblemByID(ctx context.Context, lectureID int64, problemID int64) (model.Problem, error) {
	var problem model.Problem
	err := ps.db.NewSelect().Model(&problem).
		Where("lecture_id = ? AND problem_id = ?", lectureID, problemID).
		Scan(ctx)
	if err != nil {
		return model.Problem{}, err
	}
	return problem, nil
}

func (ps *ProblemStore) GetAllLectureAndProblems(ctx context.Context) ([]model.Lecture, error) {
	var lectures []model.Lecture
	err := ps.db.NewSelect().Model(&lectures).Relation("Problems", func(q *bun.SelectQuery) *bun.SelectQuery {
		return q.Order("problem.problem_id")
	}).Order("id").Scan(ctx)
	if err != nil {
		return nil, err
	}
	return lectures, nil
}

func (ps *ProblemStore) DeleteLectureEntry(ctx context.Context, i int64) error {
	_, err := ps.db.NewDelete().Model(&model.Lecture{}).Where("id = ?", i).Exec(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (ps *ProblemStore) UpdateLectureEntry(ctx context.Context, lectureEntryInDB *model.Lecture) error {
	_, err := ps.db.NewUpdate().Model(lectureEntryInDB).Where("id = ?", lectureEntryInDB.ID).Exec(ctx)
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

func (ps *ProblemStore) GetLectureByID(ctx context.Context, id int64) (model.Lecture, error) {
	var lecture model.Lecture
	err := ps.db.NewSelect().Model(&lecture).Where("id = ?", id).Scan(ctx)
	if err != nil {
		return model.Lecture{}, err
	}
	return lecture, nil
}

func (ps *ProblemStore) CreateLectureEntry(ctx context.Context, lec *model.Lecture) error {
	_, err := ps.db.NewInsert().Model(lec).Exec(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (ps *ProblemStore) RegisterProblem(ctx context.Context, problem *model.Problem) error {
	_, err := ps.db.NewInsert().Model(problem).Exec(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (ps *ProblemStore) CheckProblemExists(ctx context.Context, lectureID, problemID int64) (bool, error) {
	count, err := ps.db.NewSelect().Model(&model.Problem{}).
		Where("lecture_id = ? AND problem_id = ?", lectureID, problemID).
		Count(ctx)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (ps *ProblemStore) DeleteProblem(ctx context.Context, lectureID, problemID int64) error {
	_, err := ps.db.NewDelete().Model(&model.Problem{}).Where("lecture_id = ? AND problem_id = ?", lectureID, problemID).Exec(ctx)
	if err != nil {
		return err
	}
	return nil
}
