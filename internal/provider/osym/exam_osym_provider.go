package osym

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/breamon/sinav-bilgi-sistemi/internal/domain"
)

type ExamOSYMProvider struct {
	baseURL string
	client  *http.Client
}

type osymExamItem struct {
	ExternalID string `json:"external_id"`
	Title      string `json:"title"`
	Status     string `json:"status"`
}

type osymExamResponse struct {
	Items []osymExamItem `json:"items"`
}

func NewExamOSYMProvider() *ExamOSYMProvider {
	return &ExamOSYMProvider{
		baseURL: "http://localhost:8080",
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (p *ExamOSYMProvider) FetchExams() ([]domain.Exam, error) {
	return p.fetchFromHTTP()
}

func (p *ExamOSYMProvider) fetchFromHTTP() ([]domain.Exam, error) {
	req, err := http.NewRequest(http.MethodGet, p.baseURL+"/mock/osym/exams", nil)
	if err != nil {
		return nil, err
	}

	resp, err := p.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("osym provider returned non-200 status")
	}

	var payload osymExamResponse
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return nil, err
	}

	return p.mapToDomain(payload)
}

func (p *ExamOSYMProvider) mapToDomain(payload osymExamResponse) ([]domain.Exam, error) {
	exams := make([]domain.Exam, 0, len(payload.Items))

	for _, item := range payload.Items {
		externalID := item.ExternalID
		status := item.Status
		if status == "" {
			status = "draft"
		}

		exams = append(exams, domain.Exam{
			Source:     "osym",
			ExternalID: &externalID,
			Title:      item.Title,
			Status:     status,
		})
	}

	return exams, nil
}
