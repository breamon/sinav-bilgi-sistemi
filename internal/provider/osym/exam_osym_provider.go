package osym

import (
	"time"

	"github.com/breamon/sinav-bilgi-sistemi/internal/domain"
)

type OSYMProvider struct{}

func NewOSYMProvider() *OSYMProvider {
	return &OSYMProvider{}
}

func (p *OSYMProvider) FetchExams() ([]domain.Exam, error) {
	now := time.Now()

	return []domain.Exam{
		{
			Source:               "osym",
			ExternalID:           stringPtr("osym-2026-yks"),
			Title:                "2026 YKS",
			Description:          stringPtr("Yükseköğretim Kurumları Sınavı"),
			Category:             stringPtr("Üniversite"),
			Status:               "upcoming",
			ApplicationStartDate: timePtr(now.AddDate(0, 1, 0)),
			ApplicationEndDate:   timePtr(now.AddDate(0, 1, 10)),
			ExamDate:             timePtr(now.AddDate(0, 3, 0)),
		},
		{
			Source:               "osym",
			ExternalID:           stringPtr("osym-2026-kpss"),
			Title:                "2026 KPSS",
			Description:          stringPtr("Kamu Personeli Seçme Sınavı"),
			Category:             stringPtr("Kamu"),
			Status:               "upcoming",
			ApplicationStartDate: timePtr(now.AddDate(0, 2, 0)),
			ApplicationEndDate:   timePtr(now.AddDate(0, 2, 10)),
			ExamDate:             timePtr(now.AddDate(0, 4, 0)),
		},
	}, nil
}

func stringPtr(s string) *string {
	return &s
}

func timePtr(t time.Time) *time.Time {
	return &t
}
