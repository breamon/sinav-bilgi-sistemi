package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/breamon/sinav-bilgi-sistemi/internal/domain"
	"github.com/breamon/sinav-bilgi-sistemi/internal/repository/postgres"
	"github.com/redis/go-redis/v9"
)

type ExamService struct {
	examRepo    *postgres.ExamRepository
	redisClient *redis.Client
}

func NewExamService(examRepo *postgres.ExamRepository, redisClient *redis.Client) *ExamService {
	return &ExamService{
		examRepo:    examRepo,
		redisClient: redisClient,
	}
}

func (s *ExamService) Create(exam *domain.Exam) error {
	exam.Title = strings.TrimSpace(exam.Title)
	exam.Source = strings.TrimSpace(exam.Source)
	exam.Status = strings.TrimSpace(exam.Status)

	if exam.Title == "" {
		return errors.New("title is required")
	}

	if exam.Source == "" {
		return errors.New("source is required")
	}

	if exam.Status == "" {
		exam.Status = "draft"
	}

	if err := s.examRepo.Create(exam); err != nil {
		return err
	}

	s.invalidateListCache()
	return nil
}

func (s *ExamService) List(page, limit int, source, status string) ([]domain.Exam, error) {
	if page <= 0 {
		page = 1
	}

	if limit <= 0 {
		limit = 10
	}

	if limit > 100 {
		limit = 100
	}

	source = strings.TrimSpace(source)
	status = strings.TrimSpace(status)

	cacheKey := s.listCacheKey(page, limit, source, status)

	if exams, ok := s.getCachedList(cacheKey); ok {
		return exams, nil
	}

	exams, err := s.examRepo.List(page, limit, source, status)
	if err != nil {
		return nil, err
	}

	s.setCachedList(cacheKey, exams)

	return exams, nil
}

func (s *ExamService) GetByID(id int64) (*domain.Exam, error) {
	if id <= 0 {
		return nil, errors.New("invalid exam id")
	}

	return s.examRepo.GetByID(id)
}

func (s *ExamService) Update(exam *domain.Exam) error {
	if exam.ID <= 0 {
		return errors.New("invalid exam id")
	}

	exam.Title = strings.TrimSpace(exam.Title)
	exam.Source = strings.TrimSpace(exam.Source)
	exam.Status = strings.TrimSpace(exam.Status)

	if exam.Title == "" {
		return errors.New("title is required")
	}

	if exam.Source == "" {
		return errors.New("source is required")
	}

	if exam.Status == "" {
		exam.Status = "draft"
	}

	if err := s.examRepo.Update(exam); err != nil {
		return err
	}

	s.invalidateListCache()
	return nil
}

func (s *ExamService) Delete(id int64) error {
	if id <= 0 {
		return errors.New("invalid exam id")
	}

	if err := s.examRepo.Delete(id); err != nil {
		return err
	}

	s.invalidateListCache()
	return nil
}

func (s *ExamService) listCacheKey(page, limit int, source, status string) string {
	return fmt.Sprintf(
		"exams:list:page:%d:limit:%d:source:%s:status:%s",
		page,
		limit,
		source,
		status,
	)
}

func (s *ExamService) getCachedList(key string) ([]domain.Exam, bool) {
	if s.redisClient == nil {
		return nil, false
	}

	ctx := context.Background()

	value, err := s.redisClient.Get(ctx, key).Result()
	if err != nil {
		return nil, false
	}

	var exams []domain.Exam
	if err := json.Unmarshal([]byte(value), &exams); err != nil {
		return nil, false
	}

	return exams, true
}

func (s *ExamService) setCachedList(key string, exams []domain.Exam) {
	if s.redisClient == nil {
		return
	}

	ctx := context.Background()

	payload, err := json.Marshal(exams)
	if err != nil {
		return
	}

	_ = s.redisClient.Set(ctx, key, payload, 0).Err()
}

func (s *ExamService) invalidateListCache() {
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
