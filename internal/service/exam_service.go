package service

import (
	"errors"
	"strings"

	"github.com/breamon/sinav-bilgi-sistemi/internal/domain"
	"github.com/breamon/sinav-bilgi-sistemi/internal/repository/postgres"
)

type ExamService struct {
	examRepo *postgres.ExamRepository
}

func NewExamService(examRepo *postgres.ExamRepository) *ExamService {
	return &ExamService{examRepo: examRepo}
}

func (s *ExamService) Create(exam *domain.Exam) error {
	exam.Title = strings.TrimSpace(exam.Title)
	exam.Source = strings.TrimSpace(exam.Source)

	if exam.Title == "" {
		return errors.New("title is required")
	}

	if exam.Source == "" {
		return errors.New("source is required")
	}

	if exam.Status == "" {
		exam.Status = "draft"
	}

	return s.examRepo.Create(exam)
}

func (s *ExamService) List() ([]domain.Exam, error) {
	return s.examRepo.List()
}

func (s *ExamService) GetByID(id int64) (*domain.Exam, error) {
	if id <= 0 {
		return nil, errors.New("invalid exam id")
	}

	return s.examRepo.GetByID(id)
}
