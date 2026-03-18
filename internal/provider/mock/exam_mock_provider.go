package mock

import (
	"github.com/breamon/sinav-bilgi-sistemi/internal/domain"
)

type ExamMockProvider struct{}

func NewExamMockProvider() *ExamMockProvider {
	return &ExamMockProvider{}
}

func (p *ExamMockProvider) FetchExams() ([]domain.Exam, error) {
	externalID1 := "osym-2026-yks"
	externalID2 := "osym-2026-kpss"
	externalID3 := "meb-2026-bursluluk"

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
			Source:     "meb",
			ExternalID: &externalID3,
			Title:      "2026 Bursluluk Sınavı",
			Status:     "draft",
		},
	}

	return exams, nil
}
