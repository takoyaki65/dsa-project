package problem

import (
	"archive/zip"
	"context"
	"dsa-backend/handler/auth"
	"dsa-backend/handler/middleware"
	requeststatus "dsa-backend/handler/problem/requestStatus"
	"dsa-backend/handler/response"
	"dsa-backend/storage/model"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
)

// TODO: Discuss file size limits
// TODO: Make this configuration be configurable via env file, or admin API.
const (
	// Maximum uncompressed size for uploaded files (10MB)
	maxUncompressedSize = 10 * 1024 * 1024
	// Maximum size for a single uploaded file (5MB)
	maxFileSize = 5 * 1024 * 1024
	// Maximum number of uploaded files
	maxFiles = 500
	// Maximum size for an uploaded zip file (5MB)
	maxZipSize = 5 * 1024 * 1024
)

// RequestValidation godoc
//
//	@Summary		Request validation
//	@Description	request a validation request, which is just compiling program codes, and executes some simple test cases.
//	@Tags			problem
//	@Accept			multipart/form-data
//	@Produce		json
//	@Param			lectureid	path		int					true	"Lecture ID"
//	@Param			problemid	path		int					true	"Problem ID"
//	@Param			files		formData	[]file				true	"Files to validate"
//	@Success		200			{object}	response.Success	"Validation request registered successfully"
//	@Failure		400			{object}	response.Error		"Invalid request"
//	@Failure		500			{object}	response.Error		"Internal server error"
//	@Security		OAuth2Password[me]
//	@Router			/problem/validate/{lectureid}/{problemid} [post]
func (h *Handler) RequestValidation(c echo.Context) error {
	var req LectureIDProblemID
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, response.NewError("Invalid request payload"))
	}

	ctx := context.Background()

	// Check the existence of problem entry
	exists, err := h.problemStore.CheckProblemExists(&ctx, req.LectureID, req.ProblemID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, response.NewError("Failed to check problem existence"))
	}
	if !exists {
		return echo.NewHTTPError(http.StatusNotFound, response.NewError("Problem not found"))
	}

	//-------------
	// Read files
	//------------

	// Multipart form
	form, err := c.MultipartForm()
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, response.NewError("Invalid multipart form"))
	}
	files := form.File["files"]
	if len(files) == 0 {
		return echo.NewHTTPError(http.StatusBadRequest, response.NewError("No files uploaded"))
	}

	// File size validation
	var totalFileSize int64
	for _, file := range files {
		totalFileSize += file.Size
	}
	if totalFileSize > maxUncompressedSize {
		return echo.NewHTTPError(http.StatusBadRequest, response.NewError(fmt.Sprintf("Total file size exceeds the maximum limit (%d MB)", maxUncompressedSize/(1024*1024))))
	}

	// get user info from jwt
	claim, err := auth.GetJWTClaims(&c)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, response.NewError("Failed to get user info"))
	}

	userCode := claim.ID
	userID := claim.UserID

	requestTime := time.Now()

	// ---------------------------------------------------------------------------------------------
	// store files at dir: upload/validation/{userID}/{lectureID}/{problemID}/{YYYY-MM-DD-HH-mm-ss}/file
	// ---------------------------------------------------------------------------------------------
	basePath := fmt.Sprintf("upload/validation/%s/%d/%d/%s", userID, req.LectureID, req.ProblemID, requestTime.Format("2006-01-02-15-04-05"))
	uploadDir := fmt.Sprintf("%s/file", basePath)

	// Check the existence of directory
	if info, err := os.Stat(uploadDir); err != nil {
		if os.IsNotExist(err) {
			// Directory does not exist
		} else if info.IsDir() {
			// Directory exists
			return echo.NewHTTPError(http.StatusConflict, response.NewError("your must not request validation for the same problem twice at the same time, please try again later"))
		} else {
			return echo.NewHTTPError(http.StatusInternalServerError, response.NewError("Failed to check upload directory"))
		}
	}

	// Make directory
	if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, response.NewError("Failed to create upload directory"))
	}

	// Store files
	for _, file := range files {
		if file.Size > maxFileSize {
			// delete directory
			_ = os.RemoveAll(uploadDir)
			return echo.NewHTTPError(http.StatusBadRequest, response.NewError(fmt.Sprintf("File size exceeds the maximum limit (%d MB)", maxFileSize/(1024*1024))))
		}

		// Source
		src, err := file.Open()
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, response.NewError("Failed to open uploaded file"))
		}
		defer src.Close()

		// Destination
		// TODO: Check if this operation is safe and there are no risks like path-traversal attacks
		dst, err := os.Create(fmt.Sprintf("%s/%s", uploadDir, filepath.Clean(file.Filename)))
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, response.NewError("Failed to create destination file"))
		}
		defer dst.Close()

		// Copy
		if _, err := io.Copy(dst, src); err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, response.NewError("Failed to copy uploaded file"))
		}
	}

	// register file location
	fileLocation := model.FileLocation{
		Path: uploadDir,
		Ts:   requestTime,
	}
	err = h.fileStore.RegisterFileLocation(&ctx, &fileLocation)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, response.NewError("Failed to register file location"))
	}

	// Make request entry
	request := model.ValidationRequest{
		TS:          requestTime,
		UserCode:    userCode,
		LectureID:   req.LectureID,
		ProblemID:   req.ProblemID,
		UploadDirID: fileLocation.ID,
		ResultID:    int64(requeststatus.WJ),
		TimeMS:      0,
		MemoryKB:    0,
	}

	// Register request
	err = h.requestStore.RegisterValidationRequest(&ctx, &request)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, response.NewError("Failed to register validation request"))
	}

	// --------------------------------------
	// Submit this request to job queue.
	// --------------------------------------

	// Get Problem info
	problem, err := h.problemStore.GetProblem(&ctx, req.LectureID, req.ProblemID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, response.NewError("Failed to get problem info"))
	}
	resultDir := fmt.Sprintf("%s/result", basePath)

	filteredBuildTasks := make([]model.TestCase, 0)
	filteredJudgeTasks := make([]model.TestCase, 0)

	// Filter tasks if the role of requested user is not manager or admin.
	if middleware.HasAllScopes(claim.Scopes, []string{"manager"}) || middleware.HasAllScopes(claim.Scopes, []string{"admin"}) {
		// Do nothing, keep all tasks
		filteredBuildTasks = problem.Detail.BuildTasks
		filteredJudgeTasks = problem.Detail.JudgeTasks
	} else {
		// Filter out tasks if "Evaluation" flat is true
		for _, task := range problem.Detail.BuildTasks {
			if !task.Evaluation {
				filteredBuildTasks = append(filteredBuildTasks, task)
			}
		}
		for _, task := range problem.Detail.JudgeTasks {
			if !task.Evaluation {
				filteredJudgeTasks = append(filteredJudgeTasks, task)
			}
		}
	}

	// Create JobQueue Entry
	job := model.JobQueue{
		RequestType: "validation",
		RequestID:   request.ID,
		Status:      "pending",
		CreatedAt:   time.Now(),
		Detail: model.JobDetail{
			TimeMS:     problem.Detail.TimeMS,
			MemoryMB:   problem.Detail.MemoryMB,
			TestFiles:  problem.Detail.TestFiles,
			FileDir:    uploadDir,
			ResultDir:  resultDir,
			BuildTasks: filteredBuildTasks,
			JudgeTasks: filteredJudgeTasks,
		},
	}

	// Register job
	err = h.jobQueueStore.InsertJob(&ctx, &job)
	if err != nil {

		// update status of ValidationRequest to "IE (Internal Error)"
		err = h.requestStore.UpdateValidationRequestStatus(&ctx, request.ID, requeststatus.IE)
		if err != nil {
			// TODO: log this with FATAL Level, because this should not happen.
			return echo.NewHTTPError(http.StatusInternalServerError, response.NewError("Failed to update validation request status"))
		}

		return echo.NewHTTPError(http.StatusInternalServerError, response.NewError("Failed to register job"))
	}

	return c.JSON(http.StatusOK, response.NewSuccess("Validation request registered successfully"))
}

type ProblemIDPathParam struct {
	LectureID int64 `param:"lectureid"`
}

// BatchValidation godoc
//
//	@Summary		Request validation for all problems in a specific lecture entry.
//	@Description	This endpoint allows users to request validation for all problems within a specific lecture.
//	@Tags			problem
//	@Accept			multipart/form-data
//	@Produce		json
//	@Param			lectureid	path		int					true	"Lecture ID"
//	@Param			zipfile		formData	file				true	"Zip file containing all program codes you're submitting"
//	@Success		200			{object}	response.Success	"Batch validation requests registered successfully"
//	@Failure		400			{object}	response.Error		"Invalid request payload"
//	@Failure		404			{object}	response.Error		"No problems found for the given lecture ID"
//	@Failure		500			{object}	response.Error		"Failed to register batch validation requests"
//	@Security		OAuth2Password[me]
//	@Router			/problem/validate/batch/{lectureid} [post]
func (h *Handler) BatchValidation(c echo.Context) error {
	var req ProblemIDPathParam
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, response.NewError("Invalid request payload"))
	}

	ctx := context.Background()

	// Check if the Lecture entry exists
	lecture, err := h.problemStore.GetLectureAndAllProblems(&ctx, req.LectureID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, response.NewError("Failed to check lecture existence"+err.Error()))
	}

	problems := lecture.Problems
	if len(problems) == 0 {
		return echo.NewHTTPError(http.StatusNotFound, response.NewError("No problems found for the given lecture ID"))
	}

	// get user info from jwt
	claim, err := auth.GetJWTClaims(&c)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, response.NewError("Failed to get user info"))
	}

	userCode := claim.ID
	userID := claim.UserID

	requestTime := time.Now()

	//----------------------------------------------------------------------------
	// Read a zip file from formData, and then unzip it.
	//
	// Before unzipping it, we have to check the size of uncompressed files to prevent zip bomb attacks.
	//
	//
	// store files at dir: upload/validation/{userID}/{lectureID}/{YYYY-MM-DD-HH-mm-ss}/file
	// ---------------------------------------------------------------------------------------------
	zipFile, err := c.FormFile("zipfile")
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, response.NewError("Invalid zip file"))
	}

	// check the size of zip file
	if zipFile.Size > maxZipSize {
		return echo.NewHTTPError(http.StatusBadRequest, response.NewError(fmt.Sprintf("Zip file size exceeds the limit of %d bytes", maxZipSize)))
	}

	basePath := fmt.Sprintf("upload/validation/%s/%d/%s", userID, req.LectureID, requestTime.Format("2006-01-02-15-04-05"))
	uploadDir := fmt.Sprintf("%s/file", basePath)

	// Check the existence of directory
	if info, err := os.Stat(uploadDir); err != nil {
		if os.IsNotExist(err) {
			// Directory does not exist
		} else if info.IsDir() {
			// Directory exists
			return echo.NewHTTPError(http.StatusConflict, response.NewError("you must not request grading for the same problem twice at the same time, please try again later"))
		} else {
			return echo.NewHTTPError(http.StatusInternalServerError, response.NewError("Failed to check upload directory"))
		}
	}

	// Open zip file
	src, err := zipFile.Open()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, response.NewError("Failed to open zip file"))
	}
	defer src.Close()

	// move to temporary directory
	tempFile, err := os.CreateTemp("", "upload-*.zip")
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, response.NewError("Failed to create temporary file"))
	}
	defer os.Remove(tempFile.Name())
	defer tempFile.Close()

	if _, err := io.Copy(tempFile, src); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, response.NewError("Failed to copy zip file to temporary file"))
	}

	// Extract zip file **safely**
	if err := safeExtractZip(tempFile.Name(), uploadDir); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, response.NewError("Failed to extract zip file: "+err.Error()))
	}

	// ---------------------------------------------------------------------------
	// Check if the first level of extracted dir contains only one folder.
	// In this case, unnest it.
	//
	// e.g.,
	// class1.zip
	//   |- class1/
	//        |- main.c
	//        |- Makefile
	//        |- Report.pdf
	// ---------------------------------------------------------------------------
	files, err := os.ReadDir(uploadDir)
	if err != nil {
		defer os.RemoveAll(uploadDir)
		return echo.NewHTTPError(http.StatusInternalServerError, response.NewError("Failed to read upload directory"))
	}
	if len(files) == 1 && files[0].IsDir() {
		// Unnest the folder
		unnestDir := fmt.Sprintf("%s/%s", uploadDir, files[0].Name())
		uploadDir = unnestDir
	}

	// Register file location
	fileLocation := model.FileLocation{
		Path: uploadDir,
		Ts:   requestTime,
	}
	err = h.fileStore.RegisterFileLocation(&ctx, &fileLocation)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, response.NewError("Failed to register file location"))
	}

	RequiringFilter := !middleware.HasAllScopes(claim.Scopes, []string{"manager"}) && !middleware.HasAllScopes(claim.Scopes, []string{"admin"})

	for _, problem := range problems {
		// Make request entry
		request := model.ValidationRequest{
			TS:          requestTime,
			UserCode:    userCode,
			LectureID:   req.LectureID,
			ProblemID:   problem.ProblemID,
			UploadDirID: fileLocation.ID,
			ResultID:    int64(requeststatus.WJ),
			TimeMS:      0,
			MemoryKB:    0,
		}

		// -------------------------------------------------------------
		// Register request
		// -------------------------------------------------------------
		err = h.requestStore.RegisterValidationRequest(&ctx, &request)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, response.NewError("Failed to register validation request"))
		}

		// -------------------------------------------------------------
		// Submit this request to job queue.
		// -------------------------------------------------------------
		resultDir := fmt.Sprintf("%s/result/%d", basePath, problem.ProblemID)
		filteredBuildTasks := make([]model.TestCase, 0)
		filteredJudgeTasks := make([]model.TestCase, 0)

		if RequiringFilter {
			// Filter out tasks if "Evaluation" flat is true
			for _, task := range problem.Detail.BuildTasks {
				if !task.Evaluation {
					filteredBuildTasks = append(filteredBuildTasks, task)
				}
			}
			for _, task := range problem.Detail.JudgeTasks {
				if !task.Evaluation {
					filteredJudgeTasks = append(filteredJudgeTasks, task)
				}
			}
		} else {
			// Do nothing, keep all tasks.
			filteredBuildTasks = problem.Detail.BuildTasks
			filteredJudgeTasks = problem.Detail.JudgeTasks
		}

		// Make an entity pushing to job queue.
		job := model.JobQueue{
			RequestType: "validation",
			RequestID:   request.ID,
			Status:      "pending",
			CreatedAt:   time.Now(),
			Detail: model.JobDetail{
				TimeMS:     problem.Detail.TimeMS,
				MemoryMB:   problem.Detail.MemoryMB,
				TestFiles:  problem.Detail.TestFiles,
				FileDir:    uploadDir,
				ResultDir:  resultDir,
				BuildTasks: filteredBuildTasks,
				JudgeTasks: filteredJudgeTasks,
			},
		}

		// Register job
		err = h.jobQueueStore.InsertJob(&ctx, &job)
		if err != nil {
			// update status of ValidationRequest to "IE (Internal Error)"
			err = h.requestStore.UpdateValidationRequestStatus(&ctx, request.ID, requeststatus.IE)
			if err != nil {
				// TODO: log this with FATAL Level, because this should not happen.
				return echo.NewHTTPError(http.StatusInternalServerError, response.NewError("Failed to update validation request status"))
			}

			return echo.NewHTTPError(http.StatusInternalServerError, response.NewError("Failed to register job"))
		}
	}

	return c.JSON(http.StatusOK, response.NewSuccess("Batch validation requests registered successfully"))
}

type GradingRequestParam struct {
	LectureID    int64  `param:"lectureid" validate:"required"`
	ProblemID    int64  `param:"problemid" validate:"required"`
	TargetUserID string `form:"userid" validate:"required"`
	SubmissionTS int64  `form:"ts" default:"1764464361" validate:"required"`
}

func (grp *GradingRequestParam) bind(c echo.Context) error {
	if err := c.Bind(grp); err != nil {
		return err
	}
	if err := c.Validate(grp); err != nil {
		return err
	}
	return nil
}

// RequestGrading godoc
//
//	@Summary		Request grading
//	@Description	request a grading request, which is compiling program codes, and executes all test cases. note that the submission timestamp is specified by the user (manager or admin), and the target user ID (e.g., student) is also specified by the user.
//	@Tags			problem
//	@Accept			multipart/form-data
//	@Produce		json
//	@Param			lectureid	path		int					true	"Lecture ID"
//	@Param			problemid	path		int					true	"Problem ID"
//	@Param			userid		formData	string				true	"User ID targeted for grading"
//	@Param			ts			formData	int					true	"Submission Timestamp, epoch seconds (e.g., 1764464361)"
//	@Param			files		formData	[]file				true	"Files to be graded"
//	@Success		200			{object}	response.Success	"Grading request registered successfully"
//	@Failure		400			{object}	response.Error		"Invalid request payload"
//	@Failure		404			{object}	response.Error		"Problem not found"
//	@Failure		500			{object}	response.Error		"Internal server error"
//	@Security		OAuth2Password[grading]
//	@Router			/problem/judge/{lectureid}/{problemid} [post]
func (h *Handler) RequestGrading(c echo.Context) error {
	req := &GradingRequestParam{}
	if err := req.bind(c); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, response.NewError("Invalid request payload"+err.Error()))
	}

	// convert epoch seconds to time.Time
	submissionTS := time.Unix(req.SubmissionTS, 0)

	ctx := context.Background()

	// Check the existence of problem entry
	exists, err := h.problemStore.CheckProblemExists(&ctx, req.LectureID, req.ProblemID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, response.NewError("Failed to check problem existence"))
	}
	if !exists {
		return echo.NewHTTPError(http.StatusNotFound, response.NewError("Problem not found"))
	}

	// Get user code of user subjected to grading
	userCodeOfSubject, err := h.userStore.GetIDByUserID(&ctx, req.TargetUserID)
	if err != nil || userCodeOfSubject == nil {
		return echo.NewHTTPError(http.StatusBadRequest, response.NewError("Targeted user does not exist"))
	}

	//-------------
	// Read files
	//------------

	// Multipart form
	form, err := c.MultipartForm()
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, response.NewError("Invalid multipart form"))
	}
	files := form.File["files"]
	if len(files) == 0 {
		return echo.NewHTTPError(http.StatusBadRequest, response.NewError("No files uploaded"))
	}

	// File size validation
	var totalFileSize int64
	for _, file := range files {
		totalFileSize += file.Size
	}
	if totalFileSize > maxUncompressedSize {
		return echo.NewHTTPError(http.StatusBadRequest, response.NewError(fmt.Sprintf("Total file size exceeds the maximum limit (%d MB)", maxUncompressedSize/(1024*1024))))
	}

	// get user info from jwt
	claim, err := auth.GetJWTClaims(&c)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, response.NewError("Failed to get user info"))
	}

	userCodeOfRequester := claim.ID

	requestTime := time.Now()

	// ---------------------------------------------------------------------------------------------
	// store files at dir: upload/grading/{userID}/{lectureID}/{problemID}/{YYYY-MM-DD-HH-mm-ss}/{YYYY-MM-DD-HH-mm-ss}/file
	//
	// Note:
	// {userID} is the ID of the user subjected to grading.
	// First {YYYY-MM-DD-HH-mm-ss} is the submission timestamp specified by the user, which is the actual timestamp when target user submit
	// program codes for grading.
	// Second {YYYY-MM-DD-HH-mm-ss} is the request timestamp.
	// ---------------------------------------------------------------------------------------------
	basePath := fmt.Sprintf("upload/grading/%s/%d/%d/%s/%s", req.TargetUserID, req.LectureID, req.ProblemID, submissionTS.Format("2006-01-02-15-04-05"), requestTime.Format("2006-01-02-15-04-05"))
	uploadDir := fmt.Sprintf("%s/file", basePath)

	// Check the existence of directory
	if info, err := os.Stat(uploadDir); err != nil {
		if os.IsNotExist(err) {
			// Directory does not exist
		} else if info.IsDir() {
			// Directory exists
			return echo.NewHTTPError(http.StatusConflict, response.NewError("you must not request grading for the same problem twice at the same time, please try again later"))
		} else {
			return echo.NewHTTPError(http.StatusInternalServerError, response.NewError("Failed to check upload directory"))
		}
	}

	// Make directory
	if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, response.NewError("Failed to create upload directory"))
	}

	// Store files
	for _, file := range files {
		if file.Size > maxFileSize {
			// delete directory
			_ = os.RemoveAll(uploadDir)
			return echo.NewHTTPError(http.StatusBadRequest, response.NewError(fmt.Sprintf("File size exceeds the maximum limit (%d MB)", maxFileSize/(1024*1024))))
		}

		// Source
		src, err := file.Open()
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, response.NewError("Failed to open uploaded file"))
		}
		defer src.Close()

		// Destination
		// TODO: Check if this operation is safe and there are no risks like path-traversal attacks
		dst, err := os.Create(fmt.Sprintf("%s/%s", uploadDir, filepath.Clean(file.Filename)))
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, response.NewError("Failed to create destination file"))
		}
		defer dst.Close()

		// Copy
		if _, err := io.Copy(dst, src); err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, response.NewError("Failed to copy uploaded file"))
		}
	}

	// register file location
	fileLocation := model.FileLocation{
		Path: uploadDir,
		Ts:   requestTime,
	}
	err = h.fileStore.RegisterFileLocation(&ctx, &fileLocation)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, response.NewError("Failed to register file location"))
	}

	// Make request entry
	request := model.GradingRequest{
		LectureID:       req.LectureID,
		ProblemID:       req.ProblemID,
		UserCode:        *userCodeOfSubject,
		SubmissionTS:    submissionTS,
		TS:              requestTime,
		RequestUserCode: userCodeOfRequester,
		UploadDirID:     fileLocation.ID,
		ResultID:        int64(requeststatus.WJ),
		TimeMS:          0,
		MemoryKB:        0,
	}

	// Register request
	err = h.requestStore.RegisterOrUpdateGradingRequest(&ctx, &request)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, response.NewError("Failed to register grading request"))
	}

	// --------------------------------------
	// Submit this request to job queue.
	// --------------------------------------

	// Get Problem info
	problem, err := h.problemStore.GetProblem(&ctx, req.LectureID, req.ProblemID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, response.NewError("Failed to get problem info"))
	}
	resultDir := fmt.Sprintf("%s/result", basePath)

	// Create JobQueue Entry
	job := model.JobQueue{
		RequestType: "grading",
		RequestID:   request.ID,
		Status:      "pending",
		CreatedAt:   time.Now(),
		Detail: model.JobDetail{
			TimeMS:     problem.Detail.TimeMS,
			MemoryMB:   problem.Detail.MemoryMB,
			TestFiles:  problem.Detail.TestFiles,
			FileDir:    uploadDir,
			ResultDir:  resultDir,
			BuildTasks: problem.Detail.BuildTasks, // We do not any filtering here, because only manager or admin can access this endpoint.
			JudgeTasks: problem.Detail.JudgeTasks,
		},
	}

	// Register job
	err = h.jobQueueStore.InsertJob(&ctx, &job)
	if err != nil {

		// update status of GradingRequest to "IE (Internal Error)"
		err = h.requestStore.UpdateGradingRequestStatus(&ctx, request.ID, requeststatus.IE)
		if err != nil {
			// TODO: log this with FATAL Level, because this should not happen.
			return echo.NewHTTPError(http.StatusInternalServerError, response.NewError("Failed to update grading request status"))
		}

		return echo.NewHTTPError(http.StatusInternalServerError, response.NewError("Failed to register job"))
	}

	return c.JSON(http.StatusOK, response.NewSuccess("Grading request registered successfully"))
}

type BatchGradingParam struct {
	LectureID    int64  `param:"lectureid" validate:"required"`
	TargetUserID string `form:"userid" validate:"required"`
	SubmissionTS int64  `form:"ts" default:"1764464361" validate:"required"`
}

func (bgp *BatchGradingParam) bind(c echo.Context) error {
	if err := c.Bind(bgp); err != nil {
		return err
	}
	if err := c.Validate(bgp); err != nil {
		return err
	}
	return nil
}

// BatchGrading godoc
//
//	@Summary		Request batched grading requests for all problems in a specific lecture entry.
//	@Description	This endpoint allows instructors to request grading for all problems in a specific lecture entry in a single request.
//	@Tags			problem
//	@Accept			multipart/form-data
//	@Produce		json
//	@Param			lectureid	path		int					true	"Lecture ID"
//	@Param			userid		formData	string				true	"User ID"
//	@Param			ts			formData	int64				true	"Submission Timestamp, epoch seconds (e.g., 1764464361)"
//	@Param			zipfile		formData	file				true	"Zip file containing all files user submitted"
//	@Success		200			{object}	response.Success	"Batched grading requests registered successfully"
//	@Failure		400			{object}	response.Error		"Invalid request payload"
//	@Failure		404			{object}	response.Error		"No problems found for the given lecture ID"
//	@Failure		500			{object}	response.Error		"Failed to register grading request"
//	@Security		OAuth2Password[grading]
//	@Router			/problem/judge/batch/{lectureid} [post]
func (h *Handler) BatchGrading(c echo.Context) error {
	var req BatchGradingParam
	if err := req.bind(c); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, response.NewError("Invalid request payload"))
	}

	// convert epoch seconds to time.Time
	submissionTS := time.Unix(req.SubmissionTS, 0)

	ctx := context.Background()

	// Check if the Lecture entry exists
	lecture, err := h.problemStore.GetLectureAndAllProblems(&ctx, req.LectureID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, response.NewError("Failed to check lecture existence"+err.Error()))
	}

	problems := lecture.Problems
	if len(problems) == 0 {
		return echo.NewHTTPError(http.StatusNotFound, response.NewError("No problems found for the given lecture ID"))
	}

	// Get user code of user subjected to grading
	userCodeOfSubject, err := h.userStore.GetIDByUserID(&ctx, req.TargetUserID)
	if err != nil || userCodeOfSubject == nil {
		return echo.NewHTTPError(http.StatusBadRequest, response.NewError("Targeted user does not exist"))
	}

	// get user info from jwt
	claim, err := auth.GetJWTClaims(&c)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, response.NewError("Failed to get user info"))
	}

	userCodeOfRequester := claim.ID
	requestTime := time.Now()

	//----------------------------------------------------------------------------
	// Read a zip file from formData, and then unzip it.
	//
	// Before unzipping it, we have to check the size of uncompressed files to prevent zip bomb attacks.
	//
	//
	// store files at dir: upload/grading/{userID}/{lectureID}/{YYYY-MM-DD-HH-mm-ss}/{YYYY-MM-DD-HH-mm-ss}/file
	// Note:
	// {userID} is the ID of the user being graded.
	// First {YYYY-MM-DD-HH-mm-ss} is the submission timestamp specified by the user, which is the
	// actual timestamp when target user submitted program codes for grading.
	// Second {YYYY-MM-DD-HH-mm-ss} is the timestamp when the grading request was created.
	// ---------------------------------------------------------------------------------------------
	zipFile, err := c.FormFile("zipfile")
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, response.NewError("Invalid zip file"))
	}

	// check the size of zip file
	if zipFile.Size > maxZipSize {
		return echo.NewHTTPError(http.StatusBadRequest, response.NewError(fmt.Sprintf("Zip file size exceeds the limit of %d bytes", maxZipSize)))
	}

	basePath := fmt.Sprintf("upload/grading/%s/%d/%s/%s", req.TargetUserID, req.LectureID, submissionTS.Format("2006-01-02-15-04-05"), requestTime.Format("2006-01-02-15-04-05"))
	uploadDir := fmt.Sprintf("%s/file", basePath)

	// Check the existence of directory
	if info, err := os.Stat(uploadDir); err != nil {
		if os.IsNotExist(err) {
			// Directory does not exist
		} else if info.IsDir() {
			// Directory exists
			return echo.NewHTTPError(http.StatusConflict, response.NewError("you must not request grading for the same problem twice at the same time, please try again later"))
		} else {
			return echo.NewHTTPError(http.StatusInternalServerError, response.NewError("Failed to check upload directory"))
		}
	}

	// Open zip file
	src, err := zipFile.Open()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, response.NewError("Failed to open zip file"))
	}
	defer src.Close()

	// move to temporary directory
	tempFile, err := os.CreateTemp("", "upload-*.zip")
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, response.NewError("Failed to create temporary file"))
	}
	defer os.Remove(tempFile.Name())
	defer tempFile.Close()

	if _, err := io.Copy(tempFile, src); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, response.NewError("Failed to copy zip file to temporary file"))
	}

	// Extract zip file **safely**
	if err := safeExtractZip(tempFile.Name(), uploadDir); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, response.NewError("Failed to extract zip file: "+err.Error()))
	}

	// ---------------------------------------------------------------------------
	// Check if the first level of extracted dir contains only one folder.
	// In this case, unnest it.
	//
	// e.g.,
	// class1.zip
	//   |- class1/
	//        |- main.c
	//        |- Makefile
	//        |- Report.pdf
	// ---------------------------------------------------------------------------
	files, err := os.ReadDir(uploadDir)
	if err != nil {
		defer os.RemoveAll(uploadDir)
		return echo.NewHTTPError(http.StatusInternalServerError, response.NewError("Failed to read upload directory"))
	}
	if len(files) == 1 && files[0].IsDir() {
		// Unnest the folder
		unnestDir := fmt.Sprintf("%s/%s", uploadDir, files[0].Name())
		uploadDir = unnestDir
	}

	// Register file location
	fileLocation := model.FileLocation{
		Path: uploadDir,
		Ts:   requestTime,
	}
	err = h.fileStore.RegisterFileLocation(&ctx, &fileLocation)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, response.NewError("Failed to register file location"))
	}

	// Make request body for each problem entries

	for _, problem := range problems {
		// Make request entry
		request := model.GradingRequest{
			LectureID:       req.LectureID,
			ProblemID:       problem.ProblemID,
			UserCode:        *userCodeOfSubject,
			SubmissionTS:    submissionTS,
			TS:              requestTime,
			RequestUserCode: userCodeOfRequester,
			UploadDirID:     fileLocation.ID,
			ResultID:        int64(requeststatus.WJ),
			TimeMS:          0,
			MemoryKB:        0,
		}

		// -------------------------------------------------------------
		// Register request
		// -------------------------------------------------------------
		err := h.requestStore.RegisterOrUpdateGradingRequest(&ctx, &request)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, response.NewError("Failed to register grading request"))
		}

		// -------------------------------------------------------------
		// Submit this request to job queue.
		// -------------------------------------------------------------
		resultDir := fmt.Sprintf("%s/result/%d", basePath, problem.ProblemID)

		// Make an entity pushing to job queue.
		job := model.JobQueue{
			RequestType: "grading",
			RequestID:   request.ID,
			Status:      "pending",
			CreatedAt:   time.Now(),
			Detail: model.JobDetail{
				TimeMS:     problem.Detail.TimeMS,
				MemoryMB:   problem.Detail.MemoryMB,
				TestFiles:  problem.Detail.TestFiles,
				FileDir:    uploadDir,
				ResultDir:  resultDir,
				BuildTasks: problem.Detail.BuildTasks, // We do not any filtering here, because only manager or admin can access this endpoint.
				JudgeTasks: problem.Detail.JudgeTasks,
			},
		}

		// Register job
		err = h.jobQueueStore.InsertJob(&ctx, &job)
		if err != nil {
			// update status of GradingRequest to "IE (Internal Error)"
			err = h.requestStore.UpdateGradingRequestStatus(&ctx, request.ID, requeststatus.IE)
			if err != nil {
				// TODO: log this with FATAL Level, because this should not happen.
				return echo.NewHTTPError(http.StatusInternalServerError, response.NewError("Failed to update grading request status"))
			}

			return echo.NewHTTPError(http.StatusInternalServerError, response.NewError("Failed to register job"))
		}
	}

	return c.JSON(http.StatusOK, response.NewSuccess("Batched grading requests registered successfully"))
}

// Extracts zip file with validation of size constraints.
// We check those constraints before/during extracting.
//
//  1. Check the number of files is above `maxFiles`.
//  2. Check the total size of all uncompressed files is below `maxUncompressedSize`.
//  3. Check the individual file sizes are below `maxFileSize`.
//
// Also, this function takes care of path-traversal attacks by sanitizing file paths.
//
// When all checks pass, extract the zip file to the specified destination directory.
// Otherwise, remove the destination directory and return an error.
func safeExtractZip(zipPath, destDir string) error {
	reader, err := zip.OpenReader(zipPath)
	if err != nil {
		return fmt.Errorf("failed to open zip file: %w", err)
	}
	defer reader.Close()

	// Check the number of files in the zip file.
	if len(reader.File) > maxFiles {
		return fmt.Errorf("zip file contains too many files (max: %d)", maxFiles)
	}

	// Check the total expected size of uncompressed files, before extracting.
	var totalUncompressed uint64
	for _, file := range reader.File {
		totalUncompressed += file.UncompressedSize64
		if totalUncompressed > maxUncompressedSize {
			return fmt.Errorf("uncompressed size too large (max: %d MB)", maxUncompressedSize/(1024*1024))
		}
	}

	// Make destination directory.
	if err := os.MkdirAll(destDir, 0755); err != nil {
		return fmt.Errorf("failed to create destination directory: %w", err)
	}

	// Extract files.
	for _, file := range reader.File {
		if err := extractFile(file, destDir); err != nil {
			// When error occurs, remove destination directory
			os.RemoveAll(destDir)
			return err
		}
	}

	return nil
}

func extractFile(file *zip.File, destDir string) error {
	// Sanitize file name to prevent path traversal attacks.
	cleanPath := filepath.Clean(file.Name)
	if strings.Contains(cleanPath, "..") {
		return fmt.Errorf("invalid file path: %s", file.Name)
	}

	targetPath := filepath.Join(destDir, cleanPath)

	// Check if the target path is within the destination directory.
	if !strings.HasPrefix(filepath.Clean(targetPath), filepath.Clean(destDir)) {
		return fmt.Errorf("file path outside of destination directory: %s", file.Name)
	}

	// Check if this individual file exceeds the size limit.
	if file.UncompressedSize64 > maxFileSize {
		return fmt.Errorf("file %s too large (max %d MB)", file.Name, maxFileSize/(1024*1024))
	}

	// In the case of "file" is a directory
	if file.FileInfo().IsDir() {
		// make directory
		return os.MkdirAll(targetPath, file.Mode())
	}

	// In the case of "file" is a regular file
	// Make parent directory
	if err := os.MkdirAll(filepath.Dir(targetPath), 0755); err != nil {
		return fmt.Errorf("failed to create directory for %s: %w", targetPath, err)
	}

	// Open file
	rc, err := file.Open()
	if err != nil {
		return fmt.Errorf("failed to open file %s in zip: %w", file.Name, err)
	}
	defer rc.Close()

	// Create output file
	outFile, err := os.OpenFile(targetPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
	if err != nil {
		return fmt.Errorf("failed to create output file %s: %w", targetPath, err)
	}
	defer outFile.Close()

	// Copy it with size limit
	limitedReader := &io.LimitedReader{
		R: rc,
		N: int64(maxFileSize),
	}

	written, err := io.Copy(outFile, limitedReader)
	if err != nil {
		return fmt.Errorf("failed to extract file: %w", err)
	}

	// Check if the entire file was copied
	if uint64(written) != file.UncompressedSize64 {
		return fmt.Errorf("file size mismatch for %s", file.Name)
	}

	return nil
}
