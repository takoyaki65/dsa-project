package problem

import (
	"dsa-backend/handler/auth"
	"dsa-backend/handler/middleware"

	"github.com/labstack/echo/v4"
	"github.com/takoyaki65/dsa-project/database"
	"github.com/uptrace/bun"
)

type Handler struct {
	db            *bun.DB
	problemStore  database.ProblemStore
	requestStore  database.RequestStore
	fileStore     database.FileStore
	userStore     database.UserStore
	jobQueueStore database.JobQueueStore
	jwtSecret     string
}

func NewProblemHandler(jwtSecret string, db *bun.DB) *Handler {
	return &Handler{
		db:            db,
		problemStore:  *database.NewProblemStore(db),
		requestStore:  *database.NewRequestStore(db),
		fileStore:     *database.NewFileStore(db),
		userStore:     *database.NewUserStore(db),
		jobQueueStore: *database.NewJobQueueStore(db),
		jwtSecret:     jwtSecret,
	}
}

func (h *Handler) RegisterRoutes(r *echo.Group) {
	// require auth
	r.Use(middleware.JWTMiddleware(h.jwtSecret))
	r.Use(middleware.CheckValidityOfJWTMiddleware(h.db))

	fetchRouter := r.Group("/fetch")
	fetchRouter.GET("/list", h.ListProblems)
	fetchRouter.GET("/detail/:lectureid/:problemid", h.GetProblemInfo)

	validateRouter := r.Group("/validate")
	validateRouter.POST("/:lectureid/:problemid", h.RequestValidation)
	validateRouter.POST("/batch/:lectureid", h.BatchValidation)

	judgeRouter := r.Group("/judge", middleware.RequiredScopesMiddleware(auth.ScopeGrading))
	judgeRouter.POST("/:lectureid/:problemid", h.RequestGrading)
	judgeRouter.POST("/batch/:lectureid", h.BatchGrading)

	crudRouter := r.Group("/crud", middleware.RequiredScopesMiddleware(auth.ScopeGrading))
	crudRouter.PUT("/create", h.CreateLectureEntry)
	crudRouter.PATCH("/update/:lectureid", h.UpdateLectureEntry)
	crudRouter.DELETE("/delete/:lectureid", h.DeleteLectureEntry)
	crudRouter.POST("/create/:lectureid/:problemid", h.RegisterProblem)
	crudRouter.DELETE("/delete/:lectureid/:problemid", h.DeleteProblem)

	resultRouter := r.Group("/result")
	resultRouter.GET("/validation/list", h.ListValidationResults)
	resultRouter.GET("/validation/detail/:requestid", h.GetValidationResult)

	gradingResultRouter := resultRouter.Group("/grading", middleware.RequiredScopesMiddleware(auth.ScopeGrading))
	gradingResultRouter.GET("/list", h.ListGradingResults)
	gradingResultRouter.GET("/detail/:requestid", h.GetGradingResult)
}
