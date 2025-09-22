package problem

import (
	"context"
	"dsa-backend/handler/auth"
	"dsa-backend/handler/problem/util"
	"dsa-backend/handler/response"
	"net/http"
	"slices"

	"github.com/labstack/echo/v4"
	"github.com/takoyaki65/dsa-project/database/model"
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
	ResultID  int64  `json:"result_id"`
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
			ResultID:  int64(result.ResultID),
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

type GradingListProps struct {
	LectureID int64 `param:"lectureid"`
}

type GradingListOutput struct {
	LectureInfo util.LectureEntry   `json:"lecture_info"`
	Detail      []UserGradingResult `json:"detail"`
}

type UserGradingResult struct {
	UserID   string                    `json:"user_id"`
	UserName string                    `json:"user_name"`
	Results  []GradingResultPerProblem `json:"results"`
}

type GradingResultPerProblem struct {
	ID           int64 `json:"id"`
	ProblemID    int64 `json:"problem_id"`
	ResultID     int64 `json:"result_id"`
	SubmissionTS int64 `json:"submission_ts"`
	TimeMS       int64 `json:"time_ms"`
	MemoryKB     int64 `json:"memory_kb"`
}

// ListGradingResults lists grading results for a specific lecture.
//
//	@Summary		List Grading Results for a Specific Lecture
//	@Description	List grading results for a specific lecture.
//	@Tags			Result
//	@Produce		json
//	@Param			lectureid	path		int64	true	"The ID of the lecture to retrieve grading results for."
//	@Success		200			{object}	GradingListOutput
//	@Failure		400			{object}	response.Error	"Invalid request"	"No users found"
//	@Failure		401			{object}	response.Error	"Failed to get user info"
//	@Failure		500			{object}	response.Error	"Failed to get user list"	"Failed to get grading results"	"Failed to get lecture info"	"Inconsistent data: user not found"	"Inconsistent data: user not found in detail assembly"
//	@Security		OAuth2Password[grading]
//	@Router			/problem/result/grading/list/{lectureid} [get]
func (h *Handler) ListGradingResults(c echo.Context) error {
	var props GradingListProps
	ctx := context.Background()
	if err := c.Bind(&props); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, response.NewError("Invalid request"))
	}

	if err := c.Validate(&props); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, response.NewError("Invalid request"))
	}

	// Fetch all user lists
	userList, err := h.userStore.GetAllUserList(ctx)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, response.NewError("Failed to get user list"))
	}
	if userList == nil {
		return echo.NewHTTPError(http.StatusBadRequest, response.NewError("No users found"))
	}

	results, err := h.requestStore.GetGradingResults(ctx, props.LectureID)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, response.NewError("Failed to get grading results"))
	}

	// get lecture info
	lectureEntry, err := util.FetchLectureByID(ctx, h.problemStore, props.LectureID, false)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, response.NewError("Failed to get lecture info"))
	}

	if lectureEntry == nil {
		return echo.NewHTTPError(http.StatusBadRequest, response.NewError("Lecture not found"))
	}

	output := GradingListOutput{
		LectureInfo: *lectureEntry,
	}

	// TODO: implement logic to fill in output.Detail
	gradingResultDict := make(map[int64]UserGradingResult)
	for _, user := range *userList {
		gradingResultDict[user.ID] = UserGradingResult{
			UserID:   user.UserID,
			UserName: user.Name,
			Results:  []GradingResultPerProblem{},
		}
	}

	for _, result := range results {
		userResult, exists := gradingResultDict[result.UserCode]
		if !exists {
			// This must not happen, so return error
			return echo.NewHTTPError(http.StatusInternalServerError, response.NewError("Inconsistent data: user not found"))
		}

		userResult.Results = append(userResult.Results, GradingResultPerProblem{
			ID:           result.ID,
			ProblemID:    result.ProblemID,
			ResultID:     int64(result.ResultID),
			SubmissionTS: result.SubmissionTS.Unix(),
			TimeMS:       result.Log.TimeMS,
			MemoryKB:     result.Log.MemoryKB,
		})

		gradingResultDict[result.UserCode] = userResult
	}

	keys := make([]int64, 0, len(*userList))
	for _, user := range *userList {
		keys = append(keys, user.ID)
	}

	// sort keys to make output deterministic
	slices.Sort(keys)

	detail := make([]UserGradingResult, 0, len(gradingResultDict))

	for _, k := range keys {
		if userResult, exists := gradingResultDict[k]; exists {
			// sort results by ProblemID and submission timestamp to make output deterministic
			slices.SortFunc(userResult.Results, func(a, b GradingResultPerProblem) int {
				if a.ProblemID != b.ProblemID {
					return int(a.ProblemID - b.ProblemID)
				}
				return int(a.SubmissionTS - b.SubmissionTS)
			})

			detail = append(detail, userResult)
		} else {
			// This must not happen, so return error
			return echo.NewHTTPError(http.StatusInternalServerError, response.NewError("Inconsistent data: user not found in detail assembly"))
		}
	}

	output.Detail = detail

	return c.JSON(http.StatusOK, output)
}

type GradingDetailProps struct {
	LectureID int64  `param:"lectureid"`
	UserID    string `param:"userid"`
}

type GradingDetailOutput struct {
	LectureID    int64                     `json:"lecture_id"`
	LectureTitle string                    `json:"lecture_title"`
	Deadline     int64                     `json:"deadline"`
	UserID       string                    `json:"user_id"`
	UserName     string                    `json:"user_name"`
	FileGroups   []FileGroup               `json:"file_groups"`
	DetailList   []GradingDetailPerProblem `json:"detail_list"`
}

type GradingDetailPerProblem struct {
	ID              int64             `json:"id"`
	ProblemID       int64             `json:"problem_id"`
	ProblemTitle    string            `json:"problem_title"`
	RequestUserID   string            `json:"request_user_id"`
	RequestUserName string            `json:"request_user_name"`
	TS              int64             `json:"ts"`
	SubmissionTS    int64             `json:"submission_ts"`
	ResultID        int64             `json:"result_id"`
	FileGroupID     int64             `json:"file_group_id"`
	TimeMS          int64             `json:"time_ms"`
	MemoryKB        int64             `json:"memory_kb"`
	BuildLogs       []DetailedTaskLog `json:"build_logs"`
	JudgeLogs       []DetailedTaskLog `json:"judge_logs"`
}

type FileGroup struct {
	ID    int64           `json:"id"`
	Files []util.FileData `json:"files"`
}

// GetGradingResult gets detailed information about a specific grading result.
//
//	@Summary		Get Grading Result Detail
//	@Description	Get detailed information about a specific grading result.
//	@Tags			Result
//	@Produce		json
//	@Param			lectureid	path		int64	true	"Lecture ID"
//	@Param			userid		path		string	true	"User ID"
//	@Success		200			{object}	GradingDetailOutput
//	@Failure		400			{object}	response.Error	"Invalid request"
//	@Failure		401			{object}	response.Error	"Failed to get user info"
//	@Failure		404			{object}	response.Error	"Grading result not found"
//	@Failure		500			{object}	response.Error	"Failed to get grading result"	"Failed to get lecture info"	"Failed to get problem info"	"Inconsistent data: user info missing"	"File location not found"	"Failed to fetch uploaded files"	"Failed to read build stdout"	"Failed to read build stderr"	"Failed to read judge stdout"	"Failed to read judge stderr"
//	@Security		OAuth2Password[grading]
//	@Router			/problem/result/grading/summary/{lectureid}/{userid} [get]
func (h *Handler) GetGradingResult(c echo.Context) error {
	var props GradingDetailProps
	ctx := context.Background()
	if err := c.Bind(&props); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, response.NewError("Invalid request"))
	}

	if err := c.Validate(&props); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, response.NewError("Invalid request"))
	}

	// Get user code from user ID
	user, err := h.userStore.GetUserByUserID(ctx, props.UserID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, response.NewError("Failed to get user info"))
	}
	userCode := user.ID

	// Get lecture and problem info
	lectureData, err := h.problemStore.GetLectureAndAllProblems(ctx, props.LectureID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, response.NewError("Failed to get lecture info"))
	}

	problemDict := make(map[int64]model.Problem)
	for _, problem := range lectureData.Problems {
		problemDict[problem.ProblemID] = *problem
	}

	grResults, err := h.requestStore.GetGradingResultsByLectureIDAndUserCode(ctx, props.LectureID, userCode)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, response.NewError("Failed to get grading results"))
	}

	output := GradingDetailOutput{
		LectureID:    props.LectureID,
		LectureTitle: lectureData.Title,
		Deadline:     lectureData.Deadline.Unix(),
		UserID:       props.UserID,
		UserName:     user.Name,
		// FileGroups to be filled later,
		// DetailList to be filled later
	}

	fileGroupDict := make(map[int64]FileGroup)

	for _, grResult := range grResults {
		upload_dir_id := grResult.UploadDirID
		if _, exists := fileGroupDict[upload_dir_id]; exists {
			continue
		}

		fileDataList, err := util.FetchAllFilesInDirectory(grResult.FileLocation.Path)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, response.NewError("Failed to fetch uploaded files"))
		}

		fileGroupDict[upload_dir_id] = FileGroup{
			ID:    upload_dir_id,
			Files: fileDataList,
		}
	}

	// Fill in FileGroups
	file_groups := make([]FileGroup, 0, len(fileGroupDict))
	for _, fg := range fileGroupDict {
		file_groups = append(file_groups, fg)
	}
	output.FileGroups = file_groups

	// Fill in DetailList
	for _, grResult := range grResults {
		var problemData model.Problem
		problemData, exists := problemDict[grResult.ProblemID]
		if !exists {
			return echo.NewHTTPError(http.StatusInternalServerError, response.NewError("Failed to get problem info"))
		}

		if grResult.RequestUser == nil {
			return echo.NewHTTPError(http.StatusInternalServerError, response.NewError("Inconsistent data: user info missing"))
		}

		detail := GradingDetailPerProblem{
			ID:              grResult.ID,
			ProblemID:       grResult.ProblemID,
			ProblemTitle:    problemData.Title,
			RequestUserID:   grResult.RequestUser.UserID,
			RequestUserName: grResult.RequestUser.Name,
			TS:              grResult.TS.Unix(),
			SubmissionTS:    grResult.SubmissionTS.Unix(),
			ResultID:        int64(grResult.Log.ResultID),
			FileGroupID:     grResult.UploadDirID,
			TimeMS:          grResult.Log.TimeMS,
			MemoryKB:        grResult.Log.MemoryKB,
			// BuildLogs to be filled later
			// JudgeLogs to be filled later
		}

		for _, buildResult := range grResult.Log.BuildResults {
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

		for _, judgeResult := range grResult.Log.JudgeResults {
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

		output.DetailList = append(output.DetailList, detail)
	}

	return c.JSON(http.StatusOK, output)
}
