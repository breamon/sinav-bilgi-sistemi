package service

import (
	"github.com/breamon/sinav-bilgi-sistemi/internal/domain"
	"github.com/breamon/sinav-bilgi-sistemi/internal/provider"
	"github.com/breamon/sinav-bilgi-sistemi/internal/repository/postgres"
)

type ExamImportService struct {
	examRepo     *postgres.ExamRepository
	provider     provider.ExamProvider
	providerName string
}

func NewExamImportService(
	examRepo *postgres.ExamRepository,
	provider provider.ExamProvider,
	providerName string,
) *ExamImportService {
	return &ExamImportService{
		examRepo:     examRepo,
		provider:     provider,
		providerName: providerName,
	}
}

func (s *ExamImportService) Import() ([]domain.Exam, error) {
	exams, err := s.provider.FetchExams()
	if err != nil {
		return nil, err
	}

	for i := range exams {
		if err := s.examRepo.UpsertBySourceAndExternalID(&exams[i]); err != nil {
			return nil, err
		}
	}

	return exams, nil
}

func (s *ExamImportService) ProviderName() string {
	return s.providerName
}
