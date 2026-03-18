package provider

import "github.com/breamon/sinav-bilgi-sistemi/internal/domain"

type ExamProvider interface {
	FetchExams() ([]domain.Exam, error)
}
