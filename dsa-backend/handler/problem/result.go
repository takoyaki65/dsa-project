package problem

import (
	"context"
	"dsa-backend/handler/auth"
	"dsa-backend/handler/problem/util"
	"dsa-backend/handler/response"
	"net/http"

	"github.com/labstack/echo/v4"
)

type ListingProps struct {
	Last int64 `query:"last" validate:"min=-1" default:"-1"`
}

type ListingOutput struct {
	LastID      int64               `json:"last_id"`
	Results     []ValidationResult  `json:"results"`
	LectureInfo []util.LectureEntry `json:"lecture_info"`
}

type ValidationResult struct {
	ID        int64  `json:"id"`
	TS        int64  `json:"ts"`
	UserID    string `json:"user_id"`
	LectureID int64  `json:"lecture_id"`
	ProblemID int64  `json:"problem_id"`
	Result    int64  `json:"result"`
}

// ListValidationResults lists validation results (not detailed, just summary) for the current user.
//
//	@Summary		List Validation Results for Current User
//	@Description	List validation results (not detailed, just summary) for the current user.
//	@Tags			Result
//	@Produce		json
//	@Param			last	query		int64	false	"The last ID received in the previous request. Use -1 to get the most recent results."	default(-1)	minimum(-1)
//	@Success		200		{object}	ListingOutput
//	@Failure		400		{object}	response.Error	"Invalid request"
//	@Failure		401		{object}	response.Error	"Failed to get user info"
//	@Failure		500		{object}	response.Error	"Failed to get lecture entries"	"Failed to get validation results"
//	@Security		OAuth2Password[grading]
//	@Router			/problem/result/validation/list [get]
func (h *Handler) ListValidationResults(c echo.Context) error {
	var props ListingProps

	if err := c.Bind(&props); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, response.NewError("Invalid request"))
	}
	if err := c.Validate(&props); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, response.NewError("Invalid request"))
	}

	if props.Last < 0 {
		props.Last = 1 << 62 // large number
	}

	ctx := context.Background()

	// get user info from jwt
	claim, err := auth.GetJWTClaims(&c)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, response.NewError("Failed to get user info"))
	}

	userCode := claim.ID
	userID := claim.UserID

	filter := false
	if !claim.HasAllScopes(auth.ScopeGrading) && !claim.HasAllScopes(auth.ScopeAdmin) {
		// if the user is not manager or admin, filter out unpublished lectures
		filter = true
	}

	// get allowed lecture ids
	lectureEntries, err := util.FetchLectureEntry(ctx, h.problemStore, filter)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, response.NewError("Failed to get lecture entries"))
	}

	var allowedLectureIDs []int64
	for _, entry := range lectureEntries {
		allowedLectureIDs = append(allowedLectureIDs, entry.LectureID)
	}

	// get validation results
	results, err := h.requestStore.GetValidationResults(ctx, userCode, allowedLectureIDs, props.Last, 20)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, response.NewError("Failed to get validation results"))
	}

	var outputResults []ValidationResult
	// calculate last id , which is most recent id (minimum)
	lastID := int64(1 << 62)
	for _, result := range results {
		if result.ID < lastID {
			lastID = result.ID
		}
		outputResults = append(outputResults, ValidationResult{
			ID:        result.ID,
			TS:        result.TS.Unix(),
			UserID:    userID,
			LectureID: result.LectureID,
			ProblemID: result.ProblemID,
			Result:    int64(result.ResultID),
		})
	}

	lastID = max(lastID-1, int64(0)) // if there is no result, set lastID to 0

	return c.JSON(http.StatusOK, ListingOutput{
		LastID:      lastID,
		Results:     outputResults,
		LectureInfo: lectureEntries,
	})
}

type ValidationDetailProps struct {
	ID int64 `param:"id" validate:"required,min=1"`
}

type DetailOutput struct {
	ID            int64             `json:"id"`
	TS            int64             `json:"ts"`
	UserID        string            `json:"user_id"`
	LectureID     int64             `json:"lecture_id"`
	ProblemID     int64             `json:"problem_id"`
	SubmissionTS  int64             `json:"submission_ts"`
	ResultID      int64             `json:"result_id"`
	TimeMS        int64             `json:"time_ms"`
	MemoryKB      int64             `json:"memory_kb"`
	UploadedFiles []util.FileData   `json:"uploaded_files"`
	BuildLogs     []DetailedTaskLog `json:"build_logs"`
	JudgeLogs     []DetailedTaskLog `json:"judge_logs"`
}

type DetailedTaskLog struct {
	TestCaseID int64  `json:"test_case_id"`
	ResultID   int64  `json:"result_id"`
	TimeMS     int64  `json:"time_ms"`
	MemoryKB   int64  `json:"memory_kb"`
	ExitCode   int64  `json:"exit_code"`
	Stdout     string `json:"stdout"` // base64 encoded
	Stderr     string `json:"stderr"` // base64 encoded
}

// GetValidationDetail gets detailed information about a specific validation result.
//
//	@Summary		Get Validation Result Detail
//	@Description	Get detailed information about a specific validation result.
//	@Tags			Result
//	@Produce		json
//	@Param			id	path		int64	true	"Validation Result ID"
//	@Success		200	{object}	DetailOutput
//	@Failure		400	{object}	response.Error	"Invalid request"
//	@Failure		401	{object}	response.Error	"Failed to get user info"
//	@Failure		404	{object}	response.Error	"Validation result not found"
//	@Failure		500	{object}	response.Error	"Failed to get validation result"	"File location not found"	"Failed to fetch uploaded files"	"Failed to read build stdout"	"Failed to read build stderr"	"Failed to read judge stdout"	"Failed to read judge stderr"
//	@Security		OAuth2Password[me]
//	@Router			/problem/result/validation/detail/{id} [get]
func (h *Handler) GetValidationDetail(c echo.Context) error {
	var props ValidationDetailProps
	if err := c.Bind(&props); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, response.NewError("Invalid request"))
	}

	if err := c.Validate(&props); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, response.NewError("Invalid request"))
	}

	// Get user info from jwt
	claim, err := auth.GetJWTClaims(&c)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, response.NewError("Failed to get user info"))
	}

	ctx := context.Background()
	userID := claim.UserID

	validationRequest, err := h.requestStore.GetValidationResultByID(ctx, props.ID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, response.NewError("Failed to get validation result"))
	}
	if validationRequest == nil {
		return echo.NewHTTPError(http.StatusNotFound, response.NewError("Validation result not found"))
	}

	detail := DetailOutput{
		ID:           validationRequest.ID,
		TS:           validationRequest.TS.Unix(),
		UserID:       userID,
		LectureID:    validationRequest.LectureID,
		ProblemID:    validationRequest.ProblemID,
		SubmissionTS: validationRequest.TS.Unix(), // for validation request, submission ts is same as request ts
		ResultID:     int64(validationRequest.Log.ResultID),
		TimeMS:       validationRequest.Log.TimeMS,
		MemoryKB:     validationRequest.Log.MemoryKB,
	}

	// Fill in uploaded files
	if validationRequest.FileLocation == nil {
		return echo.NewHTTPError(http.StatusInternalServerError, response.NewError("File location not found"))
	}

	uploadDirPath := validationRequest.FileLocation.Path

	fileDataList, err := util.FetchAllFilesInDirectory(uploadDirPath)

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, response.NewError("Failed to fetch uploaded files"))
	}

	detail.UploadedFiles = fileDataList

	// Fill in build logs
	for _, buildResult := range validationRequest.Log.BuildResults {
		stdout, err := util.FetchFile(buildResult.StdoutPath)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, response.NewError("Failed to read build stdout"))
		}
		stderr, err := util.FetchFile(buildResult.StderrPath)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, response.NewError("Failed to read build stderr"))
		}

		detail.BuildLogs = append(detail.BuildLogs, DetailedTaskLog{
			TestCaseID: buildResult.TestCaseID,
			ResultID:   int64(buildResult.ResultID),
			TimeMS:     buildResult.TimeMS,
			MemoryKB:   buildResult.MemoryKB,
			ExitCode:   buildResult.ExitCode,
			Stdout:     stdout.Data,
			Stderr:     stderr.Data,
		})
	}

	// Fill in judge logs
	for _, judgeResult := range validationRequest.Log.JudgeResults {
		stdout, err := util.FetchFile(judgeResult.StdoutPath)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, response.NewError("Failed to read judge stdout"))
		}
		stderr, err := util.FetchFile(judgeResult.StderrPath)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, response.NewError("Failed to read judge stderr"))
		}

		detail.JudgeLogs = append(detail.JudgeLogs, DetailedTaskLog{
			TestCaseID: judgeResult.TestCaseID,
			ResultID:   int64(judgeResult.ResultID),
			TimeMS:     judgeResult.TimeMS,
			MemoryKB:   judgeResult.MemoryKB,
			ExitCode:   judgeResult.ExitCode,
			Stdout:     stdout.Data,
			Stderr:     stderr.Data,
		})
	}

	return c.JSON(http.StatusOK, detail)
}

func (h *Handler) ListGradingResults(c echo.Context) error {

	panic("not implemented")
}

func (h *Handler) GetGradingResult(c echo.Context) error {
	panic("not implemented")
}
