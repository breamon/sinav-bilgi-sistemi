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

const osymCalendarURL = "https://www.osym.gov.tr/tr,8797/takvim.html"

type OSYMProvider struct {
	client *http.Client
}

func NewOSYMProvider() *OSYMProvider {
	return &OSYMProvider{
		client: &http.Client{
			Timeout: 15 * time.Second,
		},
	}
}

func (p *OSYMProvider) FetchExams() ([]domain.Exam, error) {
	req, err := http.NewRequest(http.MethodGet, osymCalendarURL, nil)
	if err != nil {
		return fallbackExams(), err
	}

	req.Header.Set("User-Agent", "sinav-bilgi-sistemi/1.0")

	resp, err := p.client.Do(req)
	if err != nil {
		return fallbackExams(), err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fallbackExams(), fmt.Errorf("osym returned status code: %d", resp.StatusCode)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return fallbackExams(), err
	}

	var exams []domain.Exam

	doc.Find("tr").Each(func(i int, row *goquery.Selection) {
		cells := row.Find("td")
		if cells.Length() < 2 {
			return
		}

		titleText := cleanText(cells.Eq(0).Text())
		if titleText == "" || strings.Contains(strings.ToLower(titleText), "sınav") && strings.Contains(strings.ToLower(titleText), "tarihi") {
			return
		}

		examDateText := cleanText(cells.Eq(1).Text())
		appDateText := ""
		if cells.Length() > 2 {
			appDateText = cleanText(cells.Eq(2).Text())
		}

		examDate := parseFirstDate(examDateText)
		appStart, appEnd := parseDateRange(appDateText)

		code := extractExamCode(titleText)
		title := buildTitle(titleText, code)

		exams = append(exams, domain.Exam{
			Source:               "osym",
			ExternalID:           stringPtr("osym-" + slugify(title)),
			Title:                title,
			Description:          stringPtr(titleText),
			Category:             stringPtr(code),
			Status:               "upcoming",
			ApplicationStartDate: appStart,
			ApplicationEndDate:   appEnd,
			ExamDate:             examDate,
		})
	})

	if len(exams) == 0 {
		return fallbackExams(), nil
	}

	return exams, nil
}

func fallbackExams() []domain.Exam {
	now := time.Now()

	return []domain.Exam{
		{
			Source:               "osym",
			ExternalID:           stringPtr("osym-2026-yks"),
			Title:                "2026 YKS",
			Description:          stringPtr("Yükseköğretim Kurumları Sınavı"),
			Category:             stringPtr("YKS"),
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
			Category:             stringPtr("KPSS"),
			Status:               "upcoming",
			ApplicationStartDate: timePtr(now.AddDate(0, 2, 0)),
			ApplicationEndDate:   timePtr(now.AddDate(0, 2, 10)),
			ExamDate:             timePtr(now.AddDate(0, 4, 0)),
		},
	}
}

func cleanText(value string) string {
	value = strings.ReplaceAll(value, "\u00a0", " ")
	value = strings.Join(strings.Fields(value), " ")
	return strings.TrimSpace(value)
}

func parseFirstDate(value string) *time.Time {
	dateRegex := regexp.MustCompile(`\d{2}\.\d{2}\.\d{4}`)
	match := dateRegex.FindString(value)
	if match == "" {
		return nil
	}

	parsed, err := time.Parse("02.01.2006", match)
	if err != nil {
		return nil
	}

	return &parsed
}

func parseDateRange(value string) (*time.Time, *time.Time) {
	dateRegex := regexp.MustCompile(`\d{2}\.\d{2}\.\d{4}`)
	matches := dateRegex.FindAllString(value, -1)

	if len(matches) == 0 {
		return nil, nil
	}

	start, err := time.Parse("02.01.2006", matches[0])
	if err != nil {
		return nil, nil
	}

	if len(matches) == 1 {
		return &start, nil
	}

	end, err := time.Parse("02.01.2006", matches[1])
	if err != nil {
		return &start, nil
	}

	return &start, &end
}

func extractExamCode(value string) string {
	knownCodes := []string{
		"YKS", "KPSS", "ALES", "DGS", "YDS", "e-YDS", "YÖKDİL", "e-YÖKDİL",
		"TUS", "DUS", "STS", "YDUS", "MSÜ", "MEB-AGS", "TR-YÖS", "HMGS",
	}

	upperValue := strings.ToUpper(value)

	for _, code := range knownCodes {
		if strings.Contains(upperValue, strings.ToUpper(code)) {
			return code
		}
	}

	parts := strings.Fields(value)
	if len(parts) > 0 {
		return strings.Trim(parts[0], ":-")
	}

	return "OSYM"
}

func buildTitle(raw string, code string) string {
	yearRegex := regexp.MustCompile(`20\d{2}`)
	year := yearRegex.FindString(raw)

	if year != "" && code != "" {
		return year + " " + code
	}

	if code != "" {
		return code
	}

	return raw
}

func slugify(value string) string {
	value = strings.ToLower(value)

	replacer := strings.NewReplacer(
		"ç", "c",
		"ğ", "g",
		"ı", "i",
		"i̇", "i",
		"ö", "o",
		"ş", "s",
		"ü", "u",
	)

	value = replacer.Replace(value)

	reg := regexp.MustCompile(`[^a-z0-9]+`)
	value = reg.ReplaceAllString(value, "-")
	value = strings.Trim(value, "-")

	if value == "" {
		return fmt.Sprintf("%d", time.Now().UnixNano())
	}

	return value
}

func stringPtr(s string) *string {
	return &s
}

func timePtr(t time.Time) *time.Time {
	return &t
}
