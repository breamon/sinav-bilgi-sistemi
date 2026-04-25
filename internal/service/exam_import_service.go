package service

import (
	"github.com/breamon/sinav-bilgi-sistemi/internal/domain"
	"github.com/breamon/sinav-bilgi-sistemi/internal/repository/postgres"
	"github.com/redis/go-redis/v9"
)

type ExamProvider interface {
	FetchExams() ([]domain.Exam, error)
}

type ExamImportService struct {
	examRepo     *postgres.ExamRepository
	provider     ExamProvider
	providerName string
	redisClient  *redis.Client
}

func NewExamImportService(
	examRepo *postgres.ExamRepository,
	provider ExamProvider,
	providerName string,
	redisClient *redis.Client,
) *ExamImportService {
	return &ExamImportService{
		examRepo:     examRepo,
		provider:     provider,
		providerName: providerName,
		redisClient:  redisClient,
	}
}

func (s *ExamImportService) Import() error {
	exams, err := s.provider.FetchExams()
	if err != nil {
		return err
	}

	for _, exam := range exams {
		if err := s.examRepo.Create(&exam); err != nil {
			return err
		}
	}

	return nil
}

func (s *ExamImportService) ImportOSYM() error {
	return s.Import()
}

func (s *ExamImportService) ProviderName() string {
	return s.providerName
}
