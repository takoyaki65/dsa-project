package model

import (
	"context"
	"time"

	"github.com/uptrace/bun"
)

type Lecture struct {
	bun.BaseModel `bun:"table:lecture"`

	ID        int64     `bun:"id,pk,autoincrement" json:"id"`
	Title     string    `bun:"title,notnull" json:"title"`
	StartDate time.Time `bun:"start_date,notnull" json:"start_date"`
	EndDate   time.Time `bun:"end_date,notnull" json:"end_date"`
	Deadline  time.Time `bun:"deadline,notnull" json:"deadline"`
}

type Problem struct {
	bun.BaseModel `bun:"table:problem"`

	LectureID          int64                  `bun:"lecture_id,pk,notnull" json:"lecture_id"`
	ProblemID          int64                  `bun:"problem_id,pk,notnull" json:"problem_id"`
	Title              string                 `bun:"title,notnull" json:"title"`
	ResourceLocationID int64                  `bun:"resource_location_id,notnull" json:"resource_location_id"`
	Detail             map[string]interface{} `bun:"detail,notnull,type:jsonb" json:"detail"`

	Lecture *Lecture `bun:"rel:belongs-to,join:lecture_id=id"`
}

var _ bun.BeforeAppendModelHook = (*Lecture)(nil)

func (l *Lecture) BeforeAppendModel(ctx context.Context, query bun.Query) error {
	switch query.(type) {
	case *bun.InsertQuery, *bun.UpdateQuery:
		// remove fraction less than seconds (milliseconds, microseconds, ...)
		l.StartDate = l.StartDate.Truncate(time.Second)
		l.EndDate = l.EndDate.Truncate(time.Second)
		l.Deadline = l.Deadline.Truncate(time.Second)
	}
	return nil
}
