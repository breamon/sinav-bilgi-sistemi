package app

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"

	apphttp "github.com/breamon/sinav-bilgi-sistemi/internal/delivery/http"
)

type App struct {
	Router *gin.Engine
	DB     *sqlx.DB
	Redis  *redis.Client
	Logger *zap.Logger
	Port   string
}

func New(router *gin.Engine, db *sqlx.DB, redisClient *redis.Client, logger *zap.Logger, port string) *App {
	return &App{
		Router: router,
		DB:     db,
		Redis:  redisClient,
		Logger: logger,
		Port:   port,
	}
}

func (a *App) Run() error {
	a.Logger.Info("server starting", zap.String("port", a.Port))
	return a.Router.Run(fmt.Sprintf(":%s", a.Port))
}

func BuildRouter(db *sqlx.DB) *gin.Engine {
	return apphttp.NewRouter(db)
}
