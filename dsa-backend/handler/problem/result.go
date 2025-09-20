package problem

import (
	"context"
	"dsa-backend/handler/auth"
	"dsa-backend/handler/problem/util"
	"dsa-backend/handler/response"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/takoyaki65/dsa-project/database/model/requeststatus"
)

type ListingProps struct {
	Last int64 `json:"last" query:"last" validate:"min=0"`
}

type ListingOutput struct {
	LastID      int64               `json:"last_id"`
	Results     []ValidationResult  `json:"results"`
	LectureInfo []util.LectureEntry `json:"lecture_info"`
}

type ValidationResult struct {
	ID        int64               `json:"id"`
	TS        int64               `json:"ts"`
	UserID    string              `json:"user_id"`
	LectureID int64               `json:"lecture_id"`
	ProblemID int64               `json:"problem_id"`
	Result    requeststatus.State `json:"result"`
}

func (h *Handler) ListValidationResults(c echo.Context) error {
	var props ListingProps

	if err := c.Bind(&props); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, response.NewError("Invalid request"))
	}
	if err := c.Validate(&props); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, response.NewError("Invalid request"))
	}

	ctx := context.Background()

	// get user info from jwt
	claim, err := auth.GetJWTClaims(&c)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, response.NewError("Failed to get user info"))
	}

	userCode := claim.ID
	userID := claim.UserID

	// get allowed lecture ids
	lectureEntries, err := util.FetchLectureEntry(c.Request().Context(), h.problemStore, true)
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
			Result:    result.ResultID,
		})
	}

	return c.JSON(http.StatusOK, ListingOutput{
		LastID:      lastID,
		Results:     outputResults,
		LectureInfo: lectureEntries,
	})
}

func (h *Handler) GetValidationResult(c echo.Context) error {
	panic("not implemented")
}

func (h *Handler) ListGradingResults(c echo.Context) error {
	panic("not implemented")
}

func (h *Handler) GetGradingResult(c echo.Context) error {
	panic("not implemented")
}
