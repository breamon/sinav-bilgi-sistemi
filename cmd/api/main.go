package main

import (
	"log"
	"time"

	"github.com/breamon/sinav-bilgi-sistemi/internal/app"
	"github.com/breamon/sinav-bilgi-sistemi/internal/config"
	"github.com/breamon/sinav-bilgi-sistemi/internal/infrastructure/cache"
	"github.com/breamon/sinav-bilgi-sistemi/internal/infrastructure/database"
	"github.com/breamon/sinav-bilgi-sistemi/internal/infrastructure/logger"
	"github.com/breamon/sinav-bilgi-sistemi/internal/repository/postgres"
	"github.com/breamon/sinav-bilgi-sistemi/internal/service"
	"go.uber.org/zap"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config load error: %v", err)
	}

	appLogger, err := logger.New()
	if err != nil {
		log.Fatalf("logger init error: %v", err)
	}
	defer appLogger.Sync()

	db, err := database.NewPostgres(cfg.Postgres)
	if err != nil {
		appLogger.Fatal("postgres connection error", zap.Error(err))
	}
	defer db.Close()

	redisClient, err := cache.NewRedis(cfg.Redis)
	if err != nil {
		appLogger.Fatal("redis connection error", zap.Error(err))
	}
	defer redisClient.Close()

	router := app.BuildRouter(db)

	interval, err := time.ParseDuration(cfg.ExamImportInterval)
	if err != nil {
		appLogger.Fatal("invalid exam import interval", zap.Error(err))
	}

	examRepo := postgres.NewExamRepository(db)
	importLogRepo := postgres.NewImportLogRepository(db)

	examScheduler := service.NewExamSchedulerService(
		examRepo,
		importLogRepo,
		appLogger,
		interval,
	)
	examScheduler.Start()

	application := app.New(router, db, redisClient, appLogger, cfg.AppPort)
	if err := application.Run(); err != nil {
		appLogger.Fatal("server run error", zap.Error(err))
	}
}
