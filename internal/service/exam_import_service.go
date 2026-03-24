package service

import (
	"context"

	"github.com/breamon/sinav-bilgi-sistemi/internal/domain"
	"github.com/breamon/sinav-bilgi-sistemi/internal/provider"
	"github.com/breamon/sinav-bilgi-sistemi/internal/repository/postgres"
	"github.com/redis/go-redis/v9"
)

type ExamImportService struct {
	examRepo     *postgres.ExamRepository
	provider     provider.ExamProvider
	providerName string
	redisClient  *redis.Client
}

func NewExamImportService(
	examRepo *postgres.ExamRepository,
	provider provider.ExamProvider,
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

	s.invalidateListCache()

	return exams, nil
}

func (s *ExamImportService) ProviderName() string {
	return s.providerName
}

func (s *ExamImportService) invalidateListCache() {
	if s.redisClient == nil {
		return
	}

	ctx := context.Background()

	keys, err := s.redisClient.Keys(ctx, "exams:list:*").Result()
	if err != nil || len(keys) == 0 {
		return
	}

	_ = s.redisClient.Del(ctx, keys...).Err()
}
