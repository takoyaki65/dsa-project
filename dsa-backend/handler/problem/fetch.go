package problem

import (
	"context"
	"dsa-backend/handler/auth"
	"dsa-backend/handler/response"
	"net/http"
	"os"
	"time"

	"github.com/labstack/echo/v4"
)

type LectureEntryResponse struct {
	ID        int64     `json:"id"`
	Title     string    `json:"title"`
	StartDate time.Time `json:"start_date"`
	Deadline  time.Time `json:"deadline"`

	Problems []ProblemEntryResponse `json:"problems"`
}

type ProblemEntryResponse struct {
	LectureID int64  `json:"lecture_id"`
	ProblemID int64  `json:"problem_id"`
	Title     string `json:"title"`
}

// ListProblems godoc
//
//	@Summary		list all problem entry, nested in lecture entry.
//	@Description	get all lecture entries, each containing its problem entries. When you don't have scopes "grading" or "admin", you will only see lecture entries that are published.
//	@Tags			problem
//	@Produce		json
//	@Success		200	{array}		LectureEntryResponse
//	@Failure		500	{object}	response.Error	"failed to get lecture entries"
//	@Security		OAuth2Password[me]
//	@Router			/problem/fetch/list [get]
func (h *Handler) ListProblems(c echo.Context) error {
	ctx := context.Background()

	// Get all Lecture entries
	lectureList, err := h.problemStore.GetAllLectureAndProblems(ctx)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, response.NewError("failed to get lecture entries"))
	}

	// Check your role
	jwtClaim, err := auth.GetJWTClaims(&c)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, response.NewError("failed to get JWT claims"))
	}
	rightsToSeeAll := jwtClaim.HasAllScopes(auth.ScopeGrading) || jwtClaim.HasAllScopes(auth.ScopeAdmin)

	var responseList []LectureEntryResponse
	for _, lecture := range *lectureList {
		if !rightsToSeeAll && lecture.StartDate.After(time.Now()) {
			// If you have no rights to see unpublished lectures, just skip those.
			continue
		}
		newLectureEntryResponse := LectureEntryResponse{
			ID:        lecture.ID,
			Title:     lecture.Title,
			StartDate: lecture.StartDate,
			Deadline:  lecture.Deadline,
			Problems:  make([]ProblemEntryResponse, len(lecture.Problems)),
		}
		for j, problem := range lecture.Problems {
			newLectureEntryResponse.Problems[j] = ProblemEntryResponse{
				LectureID: problem.LectureID,
				ProblemID: problem.ProblemID,
				Title:     problem.Title,
			}
		}
		responseList = append(responseList, newLectureEntryResponse)
	}
	return c.JSON(200, responseList)
}

type ProblemDetailRequest struct {
	LectureID int64 `json:"lecture_id" param:"lectureid"`
	ProblemID int64 `json:"problem_id" param:"problemid"`
}

type ProblemDetailResponse struct {
	LectureID     int64    `json:"lecture_id"`
	ProblemID     int64    `json:"problem_id"`
	Title         string   `json:"title"`
	Description   string   `json:"description"`
	TimeMS        int64    `json:"time_ms"`
	MemoryMB      int64    `json:"memory_mb"`
	RequiredFiles []string `json:"required_files"`
}

// GetProblemInfo godoc
//
//	@Summary		Get problem detail
//	@Description	Get detailed information about a specific problem within a lecture.
//	@Tags			problem
//	@Produce		json
//	@Param			lectureid	path		int64	true	"Lecture ID"
//	@Param			problemid	path		int64	true	"Problem ID"
//	@Success		200			{object}	ProblemDetailResponse
//	@Failure		400			{object}	response.Error	"invalid request"
//	@Failure		404			{object}	response.Error	"problem not found"
//	@Failure		500			{object}	response.Error	"failed to get problem"
//	@Security		OAuth2Password[me]
//	@Router			/problem/fetch/detail/{lectureid}/{problemid} [get]
func (h *Handler) GetProblemInfo(c echo.Context) error {
	ctx := context.Background()
	var req ProblemDetailRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, response.NewError("invalid request"))
	}

	// Check your role
	jwtClaim, err := auth.GetJWTClaims(&c)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, response.NewError("failed to get JWT claims"))
	}
	rightsToSeeAll := jwtClaim.HasAllScopes(auth.ScopeGrading) || jwtClaim.HasAllScopes(auth.ScopeAdmin)

	// fetch lecture
	lecture, err := h.problemStore.GetLectureByID(ctx, req.LectureID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, response.NewError("failed to get lecture"))
	}

	// If you have no rights to see all, and the lecture is unpublished, return 403
	if !rightsToSeeAll && lecture.StartDate.After(time.Now()) {
		return echo.NewHTTPError(http.StatusForbidden, response.NewError("you do not have permission to view this problem"))
	}

	// Fetch problem entry in DB
	problem, err := h.problemStore.GetProblemByID(ctx, req.LectureID, req.ProblemID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, response.NewError("failed to get problem"))
	}

	// read description md file
	fileLocation, err := h.fileStore.GetFileLocation(ctx, problem.ResourceLocationID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, response.NewError("failed to get file location"))
	}

	filePath := fileLocation.Path + "/" + problem.Detail.DescriptionPath

	// read file
	mdContent, err := os.ReadFile(filePath)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, response.NewError("failed to read description file"))
	}

	response := ProblemDetailResponse{
		LectureID:     problem.LectureID,
		ProblemID:     problem.ProblemID,
		Title:         problem.Title,
		Description:   string(mdContent),
		TimeMS:        problem.Detail.TimeMS,
		MemoryMB:      problem.Detail.MemoryMB,
		RequiredFiles: problem.Detail.RequiredFiles,
	}

	return c.JSON(http.StatusOK, response)
}
