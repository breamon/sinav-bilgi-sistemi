package osym

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/breamon/sinav-bilgi-sistemi/internal/domain"
)

type ExamOSYMProvider struct {
	baseURL string
	client  *http.Client
}

func NewExamOSYMProvider() *ExamOSYMProvider {
	return &ExamOSYMProvider{
		baseURL: "https://www.osym.gov.tr/tr,8797/takvim.html",
		client: &http.Client{
			Timeout: 15 * time.Second,
		},
	}
}

func (p *ExamOSYMProvider) FetchExams() ([]domain.Exam, error) {
	req, err := http.NewRequest(http.MethodGet, p.baseURL, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", "sinav-bilgi-sistemi/1.0")

	resp, err := p.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("osym provider returned status: %d", resp.StatusCode)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}

	lines := normalizeLines(doc.Text())
	return p.extractExams(lines), nil
}

func (p *ExamOSYMProvider) extractExams(lines []string) []domain.Exam {
	exams := make([]domain.Exam, 0)
	seen := make(map[string]bool)

	for i := 0; i < len(lines); i++ {
		if !isOSYMExamTitle(lines[i]) {
			continue
		}

		title := cleanExamTitle(lines[i])
		externalID := slugify(title)

		if seen[externalID] {
			continue
		}
		seen[externalID] = true

		exam := domain.Exam{
			Source:     "osym",
			ExternalID: stringPtr(externalID),
			Title:      title,
			Status:     "published",
		}

		j := i + 1
		for j < len(lines) && !isOSYMExamTitle(lines[j]) {
			switch lines[j] {
			case "Sınav Tarihi:":
				if j+1 < len(lines) {
					exam.ExamDate = parseOSYMDate(lines[j+1])
					j++
				}
			case "Başvuru Tarihleri:":
				if j+1 < len(lines) {
					exam.ApplicationStartDate = parseOSYMDate(lines[j+1])
					j++
				}
				if j+1 < len(lines) && !isSectionLabel(lines[j+1]) {
					exam.ApplicationEndDate = parseOSYMDate(lines[j+1])
					j++
				}
			case "Sonuç Tarihi:":
				if j+1 < len(lines) {
					exam.ResultDate = parseOSYMDate(lines[j+1])
					j++
				}
			}

			j++
		}

		exams = append(exams, exam)
		i = j - 1
	}

	if len(exams) == 0 {
		fallback := []string{
			"2026-YKS",
			"2026-KPSS",
			"2026-ALES",
		}

		for _, title := range fallback {
			externalID := slugify(title)
			exams = append(exams, domain.Exam{
				Source:     "osym",
				ExternalID: stringPtr(externalID),
				Title:      title,
				Status:     "draft",
			})
		}
	}

	return exams
}

func normalizeLines(text string) []string {
	raw := strings.Split(text, "\n")
	lines := make([]string, 0, len(raw))

	for _, line := range raw {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		line = collapseSpaces(line)
		if line == "" {
			continue
		}

		lines = append(lines, line)
	}

	return lines
}

func collapseSpaces(s string) string {
	return strings.Join(strings.Fields(s), " ")
}

func isOSYMExamTitle(line string) bool {
	// İlk sürüm: ÖSYM takvimindeki ana başlık satırları
	// örn: 2026-MSÜ, 2026-YÖKDİL/1, 2026-TUS 1. Dönem
	return strings.HasPrefix(line, "2026-")
}

func cleanExamTitle(line string) string {
	return strings.TrimSpace(line)
}

func isSectionLabel(line string) bool {
	switch line {
	case "Sınav Tarihi:", "Başvuru Tarihleri:", "Geç Başvuru Günü:", "Sonuç Tarihi:":
		return true
	default:
		return false
	}
}

func parseOSYMDate(value string) *time.Time {
	value = strings.TrimSpace(value)
	if value == "" {
		return nil
	}

	layouts := []string{
		"02.01.2006 15:04",
		"02.01.2006",
	}

	for _, layout := range layouts {
		if t, err := time.Parse(layout, value); err == nil {
			return &t
		}
	}

	// "04.02.2026 23:59" gibi değerlerin yanında bazen fazladan metin olabilir.
	re := regexp.MustCompile(`\d{2}\.\d{2}\.\d{4}( \d{2}:\d{2})?`)
	match := re.FindString(value)
	if match == "" {
		return nil
	}

	for _, layout := range layouts {
		if t, err := time.Parse(layout, match); err == nil {
			return &t
		}
	}

	return nil
}

func slugify(s string) string {
	s = strings.ToLower(strings.TrimSpace(s))
	replacer := strings.NewReplacer(
		" ", "-",
		"/", "-",
		".", "",
		",", "",
		":", "",
		"(", "",
		")", "",
		"ı", "i",
		"İ", "i",
		"ç", "c",
		"Ç", "c",
		"ğ", "g",
		"Ğ", "g",
		"ö", "o",
		"Ö", "o",
		"ş", "s",
		"Ş", "s",
		"ü", "u",
		"Ü", "u",
	)
	s = replacer.Replace(s)
	return "osym-" + s
}

func stringPtr(v string) *string {
	return &v
}
