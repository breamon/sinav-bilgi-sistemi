package osym

import (
	"github.com/breamon/sinav-bilgi-sistemi/internal/domain"
)

type ExamOSYMProvider struct{}

func NewExamOSYMProvider() *ExamOSYMProvider {
	return &ExamOSYMProvider{}
}

func (p *ExamOSYMProvider) FetchExams() ([]domain.Exam, error) {
	externalID1 := "osym-2026-yks"
	externalID2 := "osym-2026-kpss"
	externalID3 := "osym-2026-ales"

	exams := []domain.Exam{
		{
			Source:     "osym",
			ExternalID: &externalID1,
			Title:      "2026 YKS",
			Status:     "published",
		},
		{
			Source:     "osym",
			ExternalID: &externalID2,
			Title:      "2026 KPSS",
			Status:     "published",
		},
		{
			Source:     "osym",
			ExternalID: &externalID3,
			Title:      "2026 ALES",
			Status:     "draft",
		},
	}

	return exams, nil
}
