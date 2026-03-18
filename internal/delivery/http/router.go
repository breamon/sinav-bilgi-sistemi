package http

import (
	"os"

	"github.com/breamon/sinav-bilgi-sistemi/internal/delivery/http/handler"
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

	r.GET("/health", healthHandler.HealthCheck)

	api := r.Group("/api/v1")
	{
		auth := api.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
			auth.GET("/me", authHandler.Me)
		}

		exams := api.Group("/exams")
		{
			exams.POST("", examHandler.Create)
			exams.GET("", examHandler.List)
			exams.GET("/:id", examHandler.GetByID)
			exams.PUT("/:id", examHandler.Update)
			exams.DELETE("/:id", examHandler.Delete)
		}
	}

	return r
}
