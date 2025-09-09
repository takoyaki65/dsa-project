package model

import (
	"context"
	"time"

	"github.com/takoyaki65/dsa-project/database/model/queuestatus"
	"github.com/takoyaki65/dsa-project/database/model/queuetype"
	"github.com/takoyaki65/dsa-project/database/model/requeststatus"
	"github.com/uptrace/bun"
)

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

type ResultQueue struct {
	bun.BaseModel `bun:"table:resultqueue"`

	ID        int64        `bun:"id,pk,autoincrement" json:"id"`
	JobID     int64        `bun:"job_id,notnull" json:"job_id"`
	CreatedAt time.Time    `bun:"created_at,notnull" json:"created_at"`
	Result    ResultDetail `bun:"result,notnull,type:jsonb" json:"result"`
}

type ResultDetail struct {
	TimeMS   int64               `json:"time_ms"`
	MemoryKB int64               `json:"memory_kb"`
	ResultID requeststatus.State `json:"result_id"`
	BuildLog []ResultLog         `json:"log"`
	JudgeLog []ResultLog         `json:"judge_log"`
}

type ResultLog struct {
	TestCaseID int64               `json:"test_case_id"`
	ResultID   requeststatus.State `json:"result_id"`
	TimeMS     int64               `json:"timeMS"`
	MemoryKB   int64               `json:"memoryKB"`
	ExitCode   int64               `json:"exitCode"`
	StdoutPath string              `json:"stdoutPath"`
	StderrPath string              `json:"stderrPath"`
}

func (r *ResultDetail) ConstructFromLogs(buildLogs []ResultLog, judgeLogs []ResultLog) {
	r.BuildLog = buildLogs
	r.JudgeLog = judgeLogs

	// calculate total time and memory
	var maxTimeMS int64 = 0
	var maxMemoryKB int64 = 0
	var maxResultState requeststatus.State = requeststatus.AC

	for _, log := range buildLogs {
		if log.TimeMS > maxTimeMS {
			maxTimeMS = log.TimeMS
		}

		if log.MemoryKB > maxMemoryKB {
			maxMemoryKB = log.MemoryKB
		}
		maxResultState = maxResultState.Max(log.ResultID)
	}

	r.TimeMS = maxTimeMS
	r.MemoryKB = maxMemoryKB
	r.ResultID = maxResultState
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
