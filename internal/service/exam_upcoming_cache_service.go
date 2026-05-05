package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/breamon/sinav-bilgi-sistemi/internal/domain"
)

func (s *ExamService) GetUpcomingCached(limit int) ([]domain.Exam, error) {
	if limit <= 0 {
		limit = 10
	}

	if s.redisClient == nil {
		return s.GetUpcoming(limit)
	}

	ctx := context.Background()
	cacheKey := fmt.Sprintf("exams:upcoming:%d", limit)

	cachedValue, err := s.redisClient.Get(ctx, cacheKey).Result()
	if err == nil {
		var exams []domain.Exam
		if jsonErr := json.Unmarshal([]byte(cachedValue), &exams); jsonErr == nil {
			return exams, nil
		}
	}

	exams, err := s.GetUpcoming(limit)
	if err != nil {
		return nil, err
	}

	payload, err := json.Marshal(exams)
	if err == nil {
		_ = s.redisClient.Set(ctx, cacheKey, payload, 10*time.Minute).Err()
	}

	return exams, nil
}
