package model

import (
	"context"
	"dsa-judgeserver/storage/model/queuestatus"
	"dsa-judgeserver/storage/model/queuetype"
	"dsa-judgeserver/storage/model/requeststatus"
	"time"

	"github.com/uptrace/bun"
)

/** Those codes are almost same as [here](dsa-backend/storage/model/jobqueue.go) **/

type JobQueue struct {
	bun.BaseModel `bun:"table:jobqueue"`

	ID          int64              `bun:"id,pk,autoincrement" json:"id"`
	RequestType queuetype.Type     `bun:"request_type,notnull" json:"request_type"` // "validation" or "grading"
	RequestID   int64              `bun:"request_id,notnull" json:"request_id"`
	Status      queuestatus.Status `bun:"status,notnull" json:"status"` // "pending", "processing", "done"
	CreatedAt   time.Time          `bun:"created_at,notnull" json:"created_at"`
	Detail      JobDetail          `bun:"detail,notnull,type:jsonb" json:"detail"`
}

type JobDetail struct {
	TimeMS     int64      `json:"time_ms"`
	MemoryMB   int64      `json:"memory_mb"`
	TestFiles  []string   `json:"test_files"`
	FileDir    string     `json:"file_dir"`   // directory that contain submitted codes
	ResultDir  string     `json:"result_dir"` // directory that outputs will be stored
	BuildTasks []TestCase `json:"build"`
	JudgeTasks []TestCase `json:"judge"`
}

type TestCase struct {
	ID          int64  `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Command     string `json:"command"`
	Evaluation  bool   `json:"eval_only"`
	StdinPath   string `json:"stdin,omitempty"`
	StdoutPath  string `json:"stdout,omitempty"`
	StderrPath  string `json:"stderr,omitempty"`
	ExitCode    int64  `json:"exit,omitempty"`
}

type ResultQueue struct {
	bun.BaseModel `bun:"table:resultqueue"`

	ID        int64        `bun:"id,pk,autoincrement" json:"id"`
	JobID     int64        `bun:"job_id,notnull" json:"job_id"`
	CreatedAt time.Time    `bun:"created_at,notnull" json:"created_at"`
	Result    ResultDetail `bun:"result,notnull,type:jsonb" json:"result"`
}

type ResultDetail struct {
	TimeMS   int64               `json:"time_ms"`
	MemoryMB int64               `json:"memory_mb"`
	ResultID requeststatus.State `json:"result_id"`
	Log      RequestLog          `json:"log"`
}

type RequestLog struct {
	BuildResults []TaskLog `json:"build_results"`
	JudgeResults []TaskLog `json:"judge_results"`
}

type TaskLog struct {
	TestCaseID int64 `json:"test_case_id"`
	ResultID   int64 `json:"result_id"`
	TimeMS     int64 `json:"timeMS"`
	MemoryKB   int64 `json:"memoryKB"`
	ExitCode   int64 `json:"exitCode"`
	StdoutPath int64 `json:"stdoutPath"`
	StdErrPath int64 `json:"stderrPath"`
}

var _ bun.BeforeAppendModelHook = (*JobQueue)(nil)

func (jq *JobQueue) BeforeAppendModel(ctx context.Context, query bun.Query) error {
	switch query.(type) {
	case *bun.InsertQuery, *bun.UpdateQuery:
		// remove fraction less than seconds (milliseconds, microseconds, ...)
		jq.CreatedAt = jq.CreatedAt.Truncate(time.Second)
	}
	return nil
}

var _ bun.BeforeAppendModelHook = (*ResultQueue)(nil)

func (rq *ResultQueue) BeforeAppendModel(ctx context.Context, query bun.Query) error {
	switch query.(type) {
	case *bun.InsertQuery, *bun.UpdateQuery:
		// remove fraction less than seconds (milliseconds, microseconds, ...)
		rq.CreatedAt = rq.CreatedAt.Truncate(time.Second)
	}
	return nil
}
