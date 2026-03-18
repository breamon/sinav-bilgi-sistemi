package http

import (
	"os"

	"github.com/breamon/sinav-bilgi-sistemi/internal/delivery/http/handler"
	"github.com/breamon/sinav-bilgi-sistemi/internal/delivery/http/middleware"
	"github.com/breamon/sinav-bilgi-sistemi/internal/provider/mock"
	"github.com/breamon/sinav-bilgi-sistemi/internal/provider/osym"
	"github.com/breamon/sinav-bilgi-sistemi/internal/repository/postgres"
	"github.com/breamon/sinav-bilgi-sistemi/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

func NewRouter(db *sqlx.DB) *gin.Engine {
	r := gin.Default()

	healthHandler := handler.NewHealthHandler()

	userRepo := postgres.NewUserRepository(db)
	authService := service.NewAuthService(userRepo)
	authHandler := handler.NewAuthHandler(authService, os.Getenv("JWT_SECRET"))

	examRepo := postgres.NewExamRepository(db)
	examService := service.NewExamService(examRepo)
	examHandler := handler.NewExamHandler(examService)

	examMockProvider := mock.NewExamMockProvider()
	examMockImportService := service.NewExamImportService(examRepo, examMockProvider)
	examMockImportHandler := handler.NewExamImportHandler(examMockImportService)

	examOSYMProvider := osym.NewExamOSYMProvider()
	examOSYMImportService := service.NewExamImportService(examRepo, examOSYMProvider)
	examOSYMImportHandler := handler.NewExamImportHandler(examOSYMImportService)

	authMiddleware := middleware.AuthMiddleware(os.Getenv("JWT_SECRET"))
	adminOnlyMiddleware := middleware.AdminOnlyMiddleware()

	r.GET("/health", healthHandler.HealthCheck)

	api := r.Group("/api/v1")
	{
		auth := api.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
			auth.GET("/me", authMiddleware, authHandler.Me)
		}

		exams := api.Group("/exams")
		{
			exams.GET("", examHandler.List)
			exams.GET("/:id", examHandler.GetByID)

			admin := exams.Group("")
			admin.Use(authMiddleware, adminOnlyMiddleware)
			{
				admin.POST("", examHandler.Create)
				admin.PUT("/:id", examHandler.Update)
				admin.DELETE("/:id", examHandler.Delete)
				admin.POST("/import/mock", examMockImportHandler.Import)
				admin.POST("/import/osym", examOSYMImportHandler.Import)
			}
		}
	}

	return r
}
