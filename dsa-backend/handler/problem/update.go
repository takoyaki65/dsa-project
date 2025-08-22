package problem

import (
	"archive/zip"
	"context"
	"dsa-backend/handler/response"
	"dsa-backend/storage/model"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
)

type LectureEntryRequest struct {
	ID        int64     `json:"id" validate:"required" default:"0"`
	Title     string    `json:"title" validate:"required"`
	StartDate time.Time `json:"start_date" validate:"required" default:"2025-10-01T10:00:00+09:00"`
	Deadline  time.Time `json:"deadline" validate:"required" default:"2025-12-01T10:00:00+09:00"`
}

func (le *LectureEntryRequest) bind(c echo.Context) error {
	if err := c.Bind(le); err != nil {
		return err
	}
	if err := c.Validate(le); err != nil {
		return err
	}
	return nil
}

type LectureResponse struct {
	ID        int64     `json:"id"`
	Title     string    `json:"title"`
	StartDate time.Time `json:"start_date"`
	Deadline  time.Time `json:"deadline"`
}

type ProblemResponse struct {
	LectureID int64  `json:"lecture_id"`
	ProblemID int64  `json:"problem_id"`
	Title     string `json:"title"`
}

// CreateLectureEntry godoc
//
//	@Summary		Create a new lecture entry
//	@Description	Create a new lecture entry, accessible by manager and admin.
//	@Tags			problem
//	@Accept			json
//	@Produce		json
//	@param			lectureEntry	body		LectureEntryRequest	true	"Lecture entry details"
//	@Success		200				{object}	response.Success	"Lecture entry created successfully"
//	@Failure		400				{object}	response.Error		"Invalid request"
//	@Failure		500				{object}	response.Error		"Internal server error"
//	@Security		OAuth2Password[grading]
//	@Router			/problem/crud/create [put]
func (h *Handler) CreateLectureEntry(c echo.Context) error {
	lectureEntry := &LectureEntryRequest{}
	if err := lectureEntry.bind(c); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, response.NewError("invalid request: "+err.Error()))
	}
	ctx := context.Background()

	// Check existence of Lecture entry
	{
		_, err := h.problemStore.GetLectureByID(&ctx, lectureEntry.ID)
		if err == nil {
			return echo.NewHTTPError(http.StatusInternalServerError, response.NewError("lecture entry already exists"))
		}
	}

	err := h.problemStore.CreateLectureEntry(&ctx, &model.Lecture{
		ID:        lectureEntry.ID,
		Title:     lectureEntry.Title,
		StartDate: lectureEntry.StartDate,
		Deadline:  lectureEntry.Deadline,
	})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, response.NewError("failed to create lecture entry: "+err.Error()))
	}

	return c.JSON(http.StatusOK, response.NewSuccess("Lecture entry created successfully"))
}

// UpdateLectureEntry godoc
//
//	@Summary		Update an existing lecture entry
//	@Description	Update an existing lecture entry, accessible by manager and admin.
//	@Tags			problem
//	@Accept			json
//	@Produce		json
//	@Param			lectureid		path		int					true	"Lecture ID"
//	@Param			lectureEntry	body		LectureEntryRequest	true	"Lecture entry details"
//	@Success		200				{object}	response.Success	"Lecture entry updated successfully"
//	@Failure		400				{object}	response.Error		"Invalid request"
//	@Failure		500				{object}	response.Error		"Internal server error"
//	@Security		OAuth2Password[grading]
//	@Router			/problem/crud/update/{lectureid} [patch]
func (h *Handler) UpdateLectureEntry(c echo.Context) error {
	// Get lectureId from path param :lectureid, then convert to int
	lectureId, err := strconv.Atoi(c.Param("lectureid"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, response.NewError("invalid lecture ID"))
	}

	// Check the existence of lecture entry
	ctx := context.Background()
	lectureEntryInDB, err := h.problemStore.GetLectureByID(&ctx, int64(lectureId))
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, response.NewError("failed to get lecture entry: "+err.Error()))
	}

	lectureEntryRequest := &LectureEntryRequest{}
	if err := lectureEntryRequest.bind(c); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, response.NewError("invalid request: "+err.Error()))
	}

	lectureEntryInDB.Title = lectureEntryRequest.Title
	lectureEntryInDB.StartDate = lectureEntryRequest.StartDate
	lectureEntryInDB.Deadline = lectureEntryRequest.Deadline

	err = h.problemStore.UpdateLectureEntry(&ctx, &lectureEntryInDB)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, response.NewError("failed to update lecture entry: "+err.Error()))
	}

	return c.JSON(http.StatusOK, response.NewSuccess("Lecture entry updated successfully"))
}

// DeleteLectureEntry godoc
//
//	@Summary		Delete an existing lecture entry
//	@Description	Delete an existing lecture entry, accessible by manager and admin.
//	@Tags			problem
//	@Accept			json
//	@Produce		json
//	@Param			lectureid	path		int					true	"Lecture ID"
//	@Success		200			{object}	response.Success	"Lecture entry deleted successfully"
//	@Failure		400			{object}	response.Error		"Invalid request"
//	@Failure		500			{object}	response.Error		"Internal server error"
//	@Security		OAuth2Password[grading]
//	@Router			/problem/crud/delete/{lectureid} [delete]
func (h *Handler) DeleteLectureEntry(c echo.Context) error {
	// Get lectureId from path param :lectureid, then convert to int
	lectureId, err := strconv.Atoi(c.Param("lectureid"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, response.NewError("invalid lecture ID"))
	}

	// Check the existence of lecture entry
	ctx := context.Background()
	_, err = h.problemStore.GetLectureByID(&ctx, int64(lectureId))
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, response.NewError("failed to get lecture entry: "+err.Error()))
	}

	err = h.problemStore.DeleteLectureEntry(&ctx, int64(lectureId))
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, response.NewError("failed to delete lecture entry: "+err.Error()))
	}

	return c.JSON(http.StatusOK, response.NewSuccess("Lecture entry deleted successfully"))
}

type LectureIDProblemID struct {
	LectureID int64 `param:"lectureid"`
	ProblemID int64 `param:"problemid"`
}

// RegisterProblem godoc
//
//	@Summary		Register a new problem
//	@Description	Register a new problem associated with a lecture
//	@Tags			problem
//	@Accept			multipart/form-data
//	@Produce		json
//	@Param			lectureid	path		int					true	"Lecture ID"
//	@Param			problemid	path		int					true	"Problem ID"
//	@Param			file		formData	file				true	"Zip file contains problem resources"
//	@Success		200			{object}	response.Success	"Problem registered successfully"
//	@Failure		400			{object}	response.Error		"Invalid request"
//	@Failure		500			{object}	response.Error		"Internal server error"
//	@Security		OAuth2Password[grading]
//	@Router			/problem/crud/create/{lectureid}/{problemid} [post]
func (h *Handler) RegisterProblem(c echo.Context) error {
	var req LectureIDProblemID
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, response.NewError("invalid request: "+err.Error()))
	}

	// Read zip file
	file, err := c.FormFile("file")
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, response.NewError("failed to read file: "+err.Error()))
	}
	src, err := file.Open()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, response.NewError("failed to open file: "+err.Error()))
	}
	defer src.Close()

	// Create a temporary directory to extract the zip file
	tempDir, err := os.MkdirTemp("", "upload-*")
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, response.NewError("failed to create temp dir: "+err.Error()))
	}
	defer os.RemoveAll(tempDir)

	// Save the uploaded file to the temp directory
	zipPath := filepath.Join(tempDir, "uploaded.zip")
	out, err := os.Create(zipPath)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, response.NewError("failed to create temp file: "+err.Error()))
	}
	defer out.Close()

	if _, err := io.Copy(out, src); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, response.NewError("failed to copy file: "+err.Error()))
	}

	// Extract the zip file
	extractedDir := filepath.Join(tempDir, "extracted")
	err = unzip(zipPath, extractedDir)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, response.NewError("failed to unzip file: "+err.Error()))
	}

	// Check if the first level contains only one folder
	files, err := os.ReadDir(extractedDir)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, response.NewError("failed to read extracted directory: "+err.Error()))
	}

	if len(files) == 1 && files[0].IsDir() {
		// Unnest the folder
		unnestedDir := filepath.Join(extractedDir, files[0].Name())
		extractedDir = unnestedDir
	}

	// Check if init.json exists
	initPath := filepath.Join(extractedDir, "init.json")
	if _, err := os.Stat(initPath); os.IsNotExist(err) {
		return echo.NewHTTPError(http.StatusBadRequest, response.NewError("init.json not found"))
	}

	// Parse init.json into AssignmentConfig
	var config AssignmentConfig
	initFile, err := os.Open(initPath)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, response.NewError("failed to open init.json: "+err.Error()))
	}
	defer initFile.Close()

	err = json.NewDecoder(initFile).Decode(&config)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, response.NewError("failed to parse init.json: "+err.Error()))
	}

	// Set default value in config
	config.SetDefaults()

	// TODO: parse readme, and capture every link referencing image file in this zip file.
	// e.g., [image1.png](image1.png)
	// After that, we need to do those tasks below:
	// 1. Register referenced file into FileReference Table.
	// 2. Convert image links to use the new file reference URLs.
	// e.g., Convert ![image1.png](image1.png) to ![image1.png](<base_url>/fileref/<id>)

	// Set the destination directory
	timestamp := time.Now().Format("2006-01-02T15-04-05")
	lectureIDstr := strconv.FormatInt(req.LectureID, 10)
	problemIDstr := strconv.FormatInt(req.ProblemID, 10)
	destDir := filepath.Join("upload/resource", lectureIDstr, problemIDstr, timestamp)

	// Check if the destination directory already exists
	if _, err := os.Stat(destDir); !os.IsNotExist(err) {
		return echo.NewHTTPError(http.StatusConflict, response.NewError("destination directory already exists"))
	}

	// Move to the extracted directory to the destination
	err = os.Rename(extractedDir, destDir)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, response.NewError("failed to move directory: "+err.Error()))
	}

	context := context.Background()

	// Register file location
	fileLocation := model.FileLocation{
		Path: destDir,
		Ts:   time.Now(),
	}
	err = h.fileStore.RegisterFileLocation(&context, &fileLocation)

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, response.NewError("failed to register file location: "+err.Error()))
	}

	var buildtasks []model.TestCase
	var judgeTasks []model.TestCase

	for _, t := range config.Build {
		testcase := model.TestCase{
			Title:       t.Title,
			Description: t.Description,
			Command:     t.Command,
			Evaluation:  *t.EvalOnly,
			StdinPath:   t.Stdin,
			StdoutPath:  t.Stdout,
			StderrPath:  t.Stderr,
			ExitCode:    *t.ExitCode,
		}
		buildtasks = append(buildtasks, testcase)
	}

	for _, t := range config.Judge {
		testcase := model.TestCase{
			Title:       t.Title,
			Description: t.Description,
			Command:     t.Command,
			Evaluation:  *t.EvalOnly,
			StdinPath:   t.Stdin,
			StdoutPath:  t.Stdout,
			StderrPath:  t.Stderr,
			ExitCode:    *t.ExitCode,
		}
		judgeTasks = append(judgeTasks, testcase)
	}

	detail := model.Detail{
		DescriptionPath: config.MDfile,
		TimeMS:          *config.TimeMS,
		MemoryMB:        *config.MemoryMB,
		TestFiles:       config.TestFiles,
		RequiredFiles:   config.RequiredFiles,
		BuildTasks:      buildtasks,
		JudgeTasks:      judgeTasks,
	}

	problem := &model.Problem{
		LectureID:          req.LectureID,
		ProblemID:          req.ProblemID,
		Title:              config.Title,
		ResourceLocationID: fileLocation.ID,
		Detail:             &detail,
	}

	err = h.problemStore.RegisterProblem(&context, problem)

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, response.NewError("failed to register problem: "+err.Error()))
	}

	return c.JSON(http.StatusOK, response.NewSuccess("Problem registered successfully"))
}

func (h *Handler) DeleteProblem(c echo.Context) error {
	panic("DeleteProblem handler not implemented yet")
}

// Helper function to unzip a file
func unzip(src, dest string) error {
	r, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer r.Close()

	for _, f := range r.File {
		fpath := filepath.Join(dest, f.Name)
		if !strings.HasPrefix(fpath, filepath.Clean(dest)+string(os.PathSeparator)) {
			return errors.New("illegal file path: " + fpath)
		}

		if f.FileInfo().IsDir() {
			if err := os.MkdirAll(fpath, os.ModePerm); err != nil {
				return err
			}
		} else {
			if err := os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
				return err
			}

			outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
			if err != nil {
				return err
			}

			rc, err := f.Open()
			if err != nil {
				return err
			}

			_, err = io.Copy(outFile, rc)

			outFile.Close()
			rc.Close()

			if err != nil {
				return err
			}
		}
	}
	return nil
}
