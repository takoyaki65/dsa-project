package problem

import (
	"github.com/labstack/echo/v4"
)

type ListingProps struct {
	Last int64 `json:"last" query:"last" validate:"min=0"`
}

type ListingOutput struct {
	LectureInfo []LectureEntryResponse `json:"lecture_info"`
}

func (h *Handler) ListValidationResults(c echo.Context) error {
	// var props ListingProps

	// if err := c.Bind(&props); err != nil {
	// 	return echo.NewHTTPError(http.StatusBadRequest, response.NewError("Invalid request"))
	// }
	// if err := c.Validate(&props); err != nil {
	// 	return echo.NewHTTPError(http.StatusBadRequest, response.NewError("Invalid request"))
	// }

	// // get user info from jwt
	// claim, err := auth.GetJWTClaims(&c)
	// if err != nil {
	// 	return echo.NewHTTPError(http.StatusUnauthorized, response.NewError("Failed to get user info"))
	// }

	// userCode := claim.ID
	// userID := claim.UserID

	panic("not implemented")
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
