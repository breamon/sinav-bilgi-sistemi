package service

import (
	"github.com/breamon/sinav-bilgi-sistemi/internal/domain"
	"github.com/breamon/sinav-bilgi-sistemi/internal/repository/postgres"
)

type ImportLogService struct {
	importLogRepo *postgres.ImportLogRepository
}

func NewImportLogService(importLogRepo *postgres.ImportLogRepository) *ImportLogService {
	return &ImportLogService{importLogRepo: importLogRepo}
}

func (s *ImportLogService) Create(provider, status string, importedCount int, errMsg *string) error {
	entry := &domain.ImportLog{
		Provider:      provider,
		Status:        status,
		ImportedCount: importedCount,
		ErrorMessage:  errMsg,
	}

	return s.importLogRepo.Create(entry)
}

func (s *ImportLogService) List(limit int) ([]domain.ImportLog, error) {
	if limit <= 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}

	return s.importLogRepo.List(limit)
}
