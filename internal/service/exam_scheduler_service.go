package service

import (
	"time"

	"github.com/breamon/sinav-bilgi-sistemi/internal/provider/osym"
	"github.com/breamon/sinav-bilgi-sistemi/internal/repository/postgres"
	"go.uber.org/zap"
)

type ExamSchedulerService struct {
	importService    *ExamImportService
	importLogService *ImportLogService
	logger           *zap.Logger
	interval         time.Duration
}

func NewExamSchedulerService(
	examRepo *postgres.ExamRepository,
	importLogRepo *postgres.ImportLogRepository,
	logger *zap.Logger,
	interval time.Duration,
) *ExamSchedulerService {
	osymProvider := osym.NewExamOSYMProvider()
	importService := NewExamImportService(examRepo, osymProvider, "osym")
	importLogService := NewImportLogService(importLogRepo)

	return &ExamSchedulerService{
		importService:    importService,
		importLogService: importLogService,
		logger:           logger,
		interval:         interval,
	}
}

func (s *ExamSchedulerService) Start() {
	go func() {
		s.logger.Info("exam scheduler started", zap.Duration("interval", s.interval))

		time.Sleep(5 * time.Second)

		s.runImport()

		ticker := time.NewTicker(s.interval)
		defer ticker.Stop()

		for range ticker.C {
			s.runImport()
		}
	}()
}

func (s *ExamSchedulerService) runImport() {
	exams, err := s.importService.Import()
	if err != nil {
		errMsg := err.Error()
		_ = s.importLogService.Create(
			s.importService.ProviderName(),
			"failed",
			0,
			&errMsg,
		)

		s.logger.Error("scheduled exam import failed", zap.Error(err))
		return
	}

	_ = s.importLogService.Create(
		s.importService.ProviderName(),
		"success",
		len(exams),
		nil,
	)

	s.logger.Info("scheduled exam import completed", zap.Int("count", len(exams)))
}
