package http

import (
	"os"

	"github.com/breamon/sinav-bilgi-sistemi/internal/delivery/http/handler"
	"github.com/breamon/sinav-bilgi-sistemi/internal/provider/osym"
	"github.com/breamon/sinav-bilgi-sistemi/internal/repository/postgres"
	"github.com/breamon/sinav-bilgi-sistemi/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
)

func NewRouter(db *sqlx.DB, redisClient *redis.Client) *gin.Engine {
	r := gin.Default()

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "secret"
	}

	userRepo := postgres.NewUserRepository(db)
	examRepo := postgres.NewExamRepository(db)

	authService := service.NewAuthService(userRepo)
	examService := service.NewExamService(examRepo, redisClient)

	osymProvider := osym.NewOSYMProvider()
	examImportService := service.NewExamImportService(
		examRepo,
		osymProvider,
		"osym",
		redisClient,
	)

	healthHandler := handler.NewHealthHandler()
	authHandler := handler.NewAuthHandler(authService, jwtSecret)
	examHandler := handler.NewExamHandler(examService)
	examImportHandler := handler.NewExamImportHandler(examImportService)

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
			exams.GET("", examHandler.List)
			exams.GET("/:id", examHandler.GetByID)
			exams.POST("", examHandler.Create)
			exams.PUT("/:id", examHandler.Update)
			exams.DELETE("/:id", examHandler.Delete)

			exams.POST("/import/osym", examImportHandler.ImportOSYM)
		}
	}

	return r
}
