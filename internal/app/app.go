package app

import (
	"fmt"

	httpDelivery "github.com/breamon/sinav-bilgi-sistemi/internal/delivery/http"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

type App struct {
	router      *gin.Engine
	db          *sqlx.DB
	redisClient *redis.Client
	logger      *zap.Logger
	port        string
}

func New(router *gin.Engine, db *sqlx.DB, redisClient *redis.Client, logger *zap.Logger, port string) *App {
	return &App{
		router:      router,
		db:          db,
		redisClient: redisClient,
		logger:      logger,
		port:        port,
	}
}

func BuildRouter(db *sqlx.DB, redisClient *redis.Client) *gin.Engine {
	return httpDelivery.NewRouter(db, redisClient)
}

func (a *App) Run() error {
	a.logger.Info("server starting", zap.String("port", a.port))
	return a.router.Run(fmt.Sprintf(":%s", a.port))
}
