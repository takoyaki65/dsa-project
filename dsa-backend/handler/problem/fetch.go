package problem

import (
	"context"
	"dsa-backend/handler/auth"
	"dsa-backend/handler/problem/util"
	"dsa-backend/handler/response"
	"net/http"

	"github.com/labstack/echo/v4"
)

// ListProblems godoc
//
//	@Summary		list all problem entry, nested in lecture entry.
//	@Description	get all lecture entries, each containing its problem entries. When you don't have scopes "grading" or "admin", you will only see lecture entries that are published.
//	@Tags			Fetch
//	@Produce		json
//	@Success		200	{array}		util.LectureEntry
//	@Failure		500	{object}	response.Error	"failed to get lecture entries"
//	@Security		OAuth2Password[me]
//	@Router			/problem/fetch/list [get]
func (h *Handler) ListProblems(c echo.Context) error {
	ctx := context.Background()

	// Check your role
	jwtClaim, err := auth.GetJWTClaims(&c)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, response.NewError("failed to get JWT claims"))
	}
	rightsToSeeAll := jwtClaim.HasAllScopes(auth.ScopeGrading) || jwtClaim.HasAllScopes(auth.ScopeAdmin)
	filter := !rightsToSeeAll

	responseList, err := util.FetchLectureEntry(ctx, h.problemStore, filter)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, response.NewError("failed to fetch lecture entries"))
	}
	return c.JSON(200, responseList)
}

type ProblemDetailRequest struct {
	LectureID int64 `json:"lecture_id" param:"lectureid"`
	ProblemID int64 `json:"problem_id" param:"problemid"`
}

// GetProblemInfo godoc
//
//	@Summary		Get problem detail
//	@Description	Get detailed information about a specific problem within a lecture.
//	@Tags			Fetch
//	@Produce		json
//	@Param			lectureid	path		int64	true	"Lecture ID"
//	@Param			problemid	path		int64	true	"Problem ID"
//	@Success		200			{object}	util.ProblemDetail
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
	filter := !rightsToSeeAll

	detail, err := util.FetchProblemDetail(ctx, h.problemStore, h.fileStore, req.LectureID, req.ProblemID, filter)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, response.NewError("failed to get problem"))
	}

	return c.JSON(http.StatusOK, *detail)
}
