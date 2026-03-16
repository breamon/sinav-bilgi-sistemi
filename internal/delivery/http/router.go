package http

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"

	"github.com/breamon/sinav-bilgi-sistemi/internal/delivery/http/handler"
	"github.com/breamon/sinav-bilgi-sistemi/internal/repository/postgres"
	"github.com/breamon/sinav-bilgi-sistemi/internal/service"
)

func NewRouter(db *sqlx.DB) *gin.Engine {
	r := gin.Default()

	healthHandler := handler.NewHealthHandler()

	userRepo := postgres.NewUserRepository(db)
	authService := service.NewAuthService(userRepo)
	authHandler := handler.NewAuthHandler(authService, os.Getenv("JWT_SECRET"))

	r.GET("/health", healthHandler.HealthCheck)

	api := r.Group("/api/v1")
	{
		auth := api.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
			auth.GET("/me", authHandler.Me)
		}
	}

	return r
}
