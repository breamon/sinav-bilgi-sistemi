package service

import (
	"time"

	"github.com/breamon/sinav-bilgi-sistemi/internal/provider/osym"
	"github.com/breamon/sinav-bilgi-sistemi/internal/repository/postgres"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

type ExamSchedulerService struct {
	importService *ExamImportService
	importLogRepo *postgres.ImportLogRepository
	logger        *zap.Logger
	interval      time.Duration
}

func NewExamSchedulerService(
	examRepo *postgres.ExamRepository,
	importLogRepo *postgres.ImportLogRepository,
	redisClient *redis.Client,
	logger *zap.Logger,
	interval time.Duration,
) *ExamSchedulerService {
	osymProvider := osym.NewOSYMProvider()

	importService := NewExamImportService(
		examRepo,
		osymProvider,
		"osym",
		redisClient,
	)

	return &ExamSchedulerService{
		importService: importService,
		importLogRepo: importLogRepo,
		logger:        logger,
		interval:      interval,
	}
}

func (s *ExamSchedulerService) Start() {
	go func() {
		ticker := time.NewTicker(s.interval)
		defer ticker.Stop()

		for {
			if err := s.RunOnce(); err != nil {
				s.logger.Error("exam import failed", zap.Error(err))
			}

			<-ticker.C
		}
	}()
}

func (s *ExamSchedulerService) RunOnce() error {
	s.logger.Info("exam import started",
		zap.String("provider", s.importService.ProviderName()),
	)

	if err := s.importService.Import(); err != nil {
		return err
	}

	s.logger.Info("exam import finished",
		zap.String("provider", s.importService.ProviderName()),
	)

	return nil
}
