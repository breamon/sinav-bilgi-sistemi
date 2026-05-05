package service

import (
	"context"

	"github.com/breamon/sinav-bilgi-sistemi/internal/domain"
	"github.com/breamon/sinav-bilgi-sistemi/internal/repository/postgres"
	"github.com/redis/go-redis/v9"
)

type ExamProvider interface {
	FetchExams() ([]domain.Exam, error)
}

type ImportResult struct {
	Provider string `json:"provider"`
	Total    int    `json:"total"`
	Inserted int    `json:"inserted"`
	Updated  int    `json:"updated"`
	Failed   int    `json:"failed"`
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
	_, err := s.ImportWithResult()
	return err
}

func (s *ExamImportService) ImportOSYM() (*ImportResult, error) {
	return s.ImportWithResult()
}

func (s *ExamImportService) ImportWithResult() (*ImportResult, error) {
	exams, err := s.provider.FetchExams()
	if err != nil {
		return nil, err
	}

	result := &ImportResult{
		Provider: s.providerName,
		Total:    len(exams),
	}

	for _, exam := range exams {
		inserted, err := s.examRepo.UpsertImportedExam(&exam)
		if err != nil {
			result.Failed++
			return result, err
		}

		if inserted {
			result.Inserted++
		} else {
			result.Updated++
		}
	}

	s.clearExamCache()

	return result, nil
}

func (s *ExamImportService) ProviderName() string {
	return s.providerName
}

func (s *ExamImportService) clearExamCache() {
	if s.redisClient == nil {
		return
	}

	ctx := context.Background()

	keys, err := s.redisClient.Keys(ctx, "exams:*").Result()
	if err != nil || len(keys) == 0 {
		return
	}

	_ = s.redisClient.Del(ctx, keys...).Err()
}
