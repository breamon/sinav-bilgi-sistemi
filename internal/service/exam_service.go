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

	return s.examRepo.Create(exam)
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

	return s.examRepo.List(page, limit, source, status)
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

	return s.examRepo.Update(exam)
}

func (s *ExamService) Delete(id int64) error {
	if id <= 0 {
		return errors.New("invalid exam id")
	}

	return s.examRepo.Delete(id)
}
