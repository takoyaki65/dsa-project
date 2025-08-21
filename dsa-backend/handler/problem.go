package handler

import (
	"context"
	"dsa-backend/storage/model"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

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
//	@Success		200				{object}	RequestSuccess		"Lecture entry created successfully"
//	@Failure		400				{object}	ErrorResponse		"Invalid request"
//	@Failure		500				{object}	ErrorResponse		"Internal server error"
//	@Security		OAuth2Password[grading]
//	@Router			/problem/create [put]
func (h *Handler) CreateLectureEntry(c echo.Context) error {
	lectureEntry := &LectureEntryRequest{}
	if err := lectureEntry.bind(c); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, newErrorResponse("invalid request: "+err.Error()))
	}
	ctx := context.Background()

	// Check existence of Lecture entry
	{
		_, err := h.problemStore.GetLectureByID(&ctx, lectureEntry.ID)
		if err == nil {
			return echo.NewHTTPError(http.StatusInternalServerError, newErrorResponse("lecture entry already exists"))
		}
	}

	err := h.problemStore.CreateLectureEntry(&ctx, &model.Lecture{
		ID:        lectureEntry.ID,
		Title:     lectureEntry.Title,
		StartDate: lectureEntry.StartDate,
		Deadline:  lectureEntry.Deadline,
	})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, newErrorResponse("failed to create lecture entry: "+err.Error()))
	}

	return c.JSON(http.StatusOK, RequestSuccess{
		Msg: "Lecture entry created successfully",
	})
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
//	@Success		200				{object}	RequestSuccess		"Lecture entry updated successfully"
//	@Failure		400				{object}	ErrorResponse		"Invalid request"
//	@Failure		500				{object}	ErrorResponse		"Internal server error"
//	@Security		OAuth2Password[grading]
//	@Router			/problem/update/{lectureid} [patch]
func (h *Handler) UpdateLectureEntry(c echo.Context) error {
	// Get lectureId from path param :lectureid, then convert to int
	lectureId, err := strconv.Atoi(c.Param("lectureid"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, newErrorResponse("invalid lecture ID"))
	}

	// Check the existence of lecture entry
	ctx := context.Background()
	lectureEntryInDB, err := h.problemStore.GetLectureByID(&ctx, int64(lectureId))
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, newErrorResponse("failed to get lecture entry: "+err.Error()))
	}

	lectureEntryRequest := &LectureEntryRequest{}
	if err := lectureEntryRequest.bind(c); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, newErrorResponse("invalid request: "+err.Error()))
	}

	lectureEntryInDB.Title = lectureEntryRequest.Title
	lectureEntryInDB.StartDate = lectureEntryRequest.StartDate
	lectureEntryInDB.Deadline = lectureEntryRequest.Deadline

	err = h.problemStore.UpdateLectureEntry(&ctx, &lectureEntryInDB)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, newErrorResponse("failed to update lecture entry: "+err.Error()))
	}

	return c.JSON(http.StatusOK, RequestSuccess{
		Msg: "Lecture entry updated successfully",
	})
}

// DeleteLectureEntry godoc
//
//	@Summary		Delete an existing lecture entry
//	@Description	Delete an existing lecture entry, accessible by manager and admin.
//	@Tags			problem
//	@Accept			json
//	@Produce		json
//	@Param			lectureid	path		int				true	"Lecture ID"
//	@Success		200			{object}	RequestSuccess	"Lecture entry deleted successfully"
//	@Failure		400			{object}	ErrorResponse	"Invalid request"
//	@Failure		500			{object}	ErrorResponse	"Internal server error"
//	@Security		OAuth2Password[grading]
//	@Router			/problem/delete/{lectureid} [delete]
func (h *Handler) DeleteLectureEntry(c echo.Context) error {
	// Get lectureId from path param :lectureid, then convert to int
	lectureId, err := strconv.Atoi(c.Param("lectureid"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, newErrorResponse("invalid lecture ID"))
	}

	// Check the existence of lecture entry
	ctx := context.Background()
	_, err = h.problemStore.GetLectureByID(&ctx, int64(lectureId))
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, newErrorResponse("failed to get lecture entry: "+err.Error()))
	}

	err = h.problemStore.DeleteLectureEntry(&ctx, int64(lectureId))
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, newErrorResponse("failed to delete lecture entry: "+err.Error()))
	}

	return c.JSON(http.StatusOK, RequestSuccess{
		Msg: "Lecture entry deleted successfully",
	})
}

func (h *Handler) RegisterProblem(c echo.Context) error {
	panic("RegisterProblem handler not implemented yet")
}

func (h *Handler) DeleteProblem(c echo.Context) error {
	panic("DeleteProblem handler not implemented yet")
}
