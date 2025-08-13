package model

import (
	"context"
	"time"

	"github.com/uptrace/bun"
)

type ResultValues struct {
	bun.BaseModel `bun:"table:resultvalues"`

	Value int64  `bun:"value,notnull" json:"value"`
	Name  string `bun:"name,notnull" json:"name"`
}

type Request struct {
	bun.BaseModel `bun:"table:request"`

	ID            int64                  `bun:"id,pk,autoincrement" json:"id"`
	TS            time.Time              `bun:"ts,notnull" json:"ts"`
	UserID        int64                  `bun:"user_id,notnull" json:"user_id"`
	SubmissionTS  time.Time              `bun:"submission_ts,notnull" json:"submission_ts"`
	RequestUserID int64                  `bun:"request_user_id,notnull" json:"request_user_id"`
	Eval          bool                   `bun:"eval,notnull" json:"eval"`
	LectureID     int64                  `bun:"lecture_id,notnull" json:"lecture_id"`
	ProblemID     int64                  `bun:"problem_id,notnull" json:"problem"`
	UploadDirID   int64                  `bun:"upload_dir_id,notnull" json:"upload_dir_id"`
	ResultID      int64                  `bun:"result,notnull" json:"result_id"`
	Log           map[string]interface{} `bun:"log,notnull,type:jsonb" json:"log"`
	TimeMS        int64                  `bun:"time_ms,notnull" json:"time_ms"`
	MemoryKB      int64                  `bun:"memory_kb,notnull" json:"memory_kb"`

	Problem      *Problem      `bun:"rel:belongs-to,join:(lecture_id,problem_id)=(lecture_id,problem_id)"`
	Result       *ResultValues `bun:"rel:has-one,join:result=id"`
	FileLocation *FileLocation `bun:"rel:has-one,join:upload_dir_id=id"`
}

var _ bun.BeforeAppendModelHook = (*Request)(nil)

func (r *Request) BeforeAppendModel(ctx context.Context, query bun.Query) error {
	switch query.(type) {
	case *bun.InsertQuery, *bun.UpdateQuery:
		// remove fraction less than seconds (milliseconds, microeconds, ...)
		r.TS = r.TS.Truncate(time.Second)
		r.SubmissionTS = r.SubmissionTS.Truncate(time.Second)
	}
	return nil
}
