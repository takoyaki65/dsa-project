package problem

import (
	"context"
	"dsa-backend/handler/auth"
	requeststatus "dsa-backend/handler/problem/requestStatus"
	"dsa-backend/handler/response"
	"dsa-backend/storage/model"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/labstack/echo/v4"
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

	// get user info from jwt
	claim, err := auth.GetJWTClaims(&c)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, response.NewError("Failed to get user info"))
	}

	userCode := claim.ID
	userID := claim.UserID

	requestTime := time.Now()

	// ---------------------------------------------------------------------------------------------
	// store files at dir: upload/validation/{userID}/{lectureID}/{problemID}/{YYYY-MM-DD-HH-mm-ss}
	// ---------------------------------------------------------------------------------------------
	uploadDir := fmt.Sprintf("upload/validation/%s/%d/%d/%s", userID, req.LectureID, req.ProblemID, requestTime.Format("2006-01-02-15-04-05"))

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
		// Source
		src, err := file.Open()
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, response.NewError("Failed to open uploaded file"))
		}
		defer src.Close()

		// Destination
		dst, err := os.Create(fmt.Sprintf("%s/%s", uploadDir, file.Filename))
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

	return c.JSON(http.StatusOK, response.NewSuccess("Validation request registered successfully"))
}

type ProblemIDPathParam struct {
	LectureID int64 `param:"lectureid"`
}

func (h *Handler) BatchValidation(c echo.Context) error {
	var req ProblemIDPathParam
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, response.NewError("Invalid request payload"))
	}
	panic("BatchValidateSubmissions handler not implemented yet")
}

func (h *Handler) RequestGrading(c echo.Context) error {
	var req LectureIDProblemID
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, response.NewError("Invalid request payload"))
	}

	panic("JudgeSubmission handler not implemented yet")
}

func (h *Handler) BatchGrading(c echo.Context) error {
	var req ProblemIDPathParam
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, response.NewError("Invalid request payload"))
	}
	panic("BatchGrading handler not implemented yet")
}
