package problem

import (
	"context"
	"dsa-backend/handler/auth"
	"dsa-backend/handler/middleware"
	"dsa-backend/handler/response"
	"dsa-backend/storage"
	"dsa-backend/storage/model"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/uptrace/bun"
)

type Handler struct {
	db           *bun.DB
	problemStore storage.ProblemStore
	jwtSecret    string
}

func NewProblemHandler(jwtSecret string, db *bun.DB) *Handler {
	return &Handler{
		db:           db,
		problemStore: *storage.NewProblemStore(db),
		jwtSecret:    jwtSecret,
	}
}

func (h *Handler) RegisterRoutes(r *echo.Group) {
	// require auth
	r.Use(middleware.JWTMiddleware(h.jwtSecret))

	fetchRouter := r.Group("/fetch")
	fetchRouter.GET("/list", h.ListProblems)
	fetchRouter.GET("/detail/:lectureid/:problemid", h.GetProblemInfo)

	validateRouter := r.Group("/validate")
	validateRouter.POST("/:lectureid/:problemid", h.ValidateSubmission)

	judgeRouter := r.Group("/judge", middleware.RequiredScopesMiddleware(auth.ScopeGrading))
	judgeRouter.POST("/:lectureid/:problemid", h.JudgeSubmission)

	crudRouter := r.Group("/crud", middleware.RequiredScopesMiddleware(auth.ScopeGrading))
	crudRouter.PUT("/create", h.CreateLectureEntry)
	crudRouter.PATCH("/update/:lectureid", h.UpdateLectureEntry)
	crudRouter.DELETE("/delete/:lectureid", h.DeleteLectureEntry)
	crudRouter.POST("/create/:lectureid", h.RegisterProblem)
	crudRouter.DELETE("/delete/:lectureid/:problemid", h.DeleteProblem)
}

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

func (h *Handler) ListProblems(c echo.Context) error {
	panic("ListProblems handler not implemented yet")
}

func (h *Handler) GetProblemInfo(c echo.Context) error {
	panic("GetProblemInfo handler not implemented yet")
}

func (h *Handler) ValidateSubmission(c echo.Context) error {
	panic("ValidateSubmission handler not implemented yet")
}

func (h *Handler) JudgeSubmission(c echo.Context) error {
	panic("JudgeSubmission handler not implemented yet")
}

// CreateLectureEntry godoc
//
//	@Summary		Create a new lecture entry
//	@Description	Create a new lecture entry, accessible by manager and admin.
//	@Tags			problem
//	@Accept			json
//	@Produce		json
//	@param			lectureEntry	body		LectureEntryRequest	true	"Lecture entry details"
//	@Success		200				{object}	response.Success		"Lecture entry created successfully"
//	@Failure		400				{object}	response.Error		"Invalid request"
//	@Failure		500				{object}	response.Error		"Internal server error"
//	@Security		OAuth2Password[grading]
//	@Router			/problem/create [put]
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
//	@Success		200				{object}	response.Success		"Lecture entry updated successfully"
//	@Failure		400				{object}	response.Error		"Invalid request"
//	@Failure		500				{object}	response.Error		"Internal server error"
//	@Security		OAuth2Password[grading]
//	@Router			/problem/update/{lectureid} [patch]
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
//	@Param			lectureid	path		int				true	"Lecture ID"
//	@Success		200			{object}	response.Success	"Lecture entry deleted successfully"
//	@Failure		400			{object}	response.Error	"Invalid request"
//	@Failure		500			{object}	response.Error	"Internal server error"
//	@Security		OAuth2Password[grading]
//	@Router			/problem/delete/{lectureid} [delete]
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

func (h *Handler) RegisterProblem(c echo.Context) error {
	panic("RegisterProblem handler not implemented yet")
}

func (h *Handler) DeleteProblem(c echo.Context) error {
	panic("DeleteProblem handler not implemented yet")
}
