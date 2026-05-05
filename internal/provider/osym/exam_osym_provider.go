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

	text := cleanTextWithNewLines(doc.Text())
	exams := parseCalendarText(text)

	if len(exams) == 0 {
		return fallbackExams(), nil
	}

	return exams, nil
}

func parseCalendarText(text string) []domain.Exam {
	lines := strings.Split(text, "\n")
	var exams []domain.Exam

	for i := 0; i < len(lines); i++ {
		line := cleanText(lines[i])

		if !isExamTitle(line) {
			continue
		}

		title := line
		description := previousUsefulLine(lines, i)

		blockEnd := i + 18
		if blockEnd > len(lines) {
			blockEnd = len(lines)
		}

		block := lines[i:blockEnd]

		examDate := findDateAfterLabel(block, "Sınav Tarihi:")
		appDates := findDatesAfterLabel(block, "Başvuru Tarihleri:", 2)
		resultDate := findDateAfterLabel(block, "Sonuç Tarihi:")

		var appStart *time.Time
		var appEnd *time.Time

		if len(appDates) > 0 {
			appStart = appDates[0]
		}

		if len(appDates) > 1 {
			appEnd = appDates[1]
		}

		code := extractExamCode(title)

		exams = append(exams, domain.Exam{
			Source:               "osym",
			ExternalID:           stringPtr("osym-" + slugify(title)),
			Title:                title,
			Description:          stringPtr(description),
			Category:             stringPtr(code),
			Status:               "published",
			ApplicationStartDate: appStart,
			ApplicationEndDate:   appEnd,
			ExamDate:             examDate,
			ResultDate:           resultDate,
		})
	}

	return exams
}

func isExamTitle(value string) bool {
	if value == "" {
		return false
	}

	value = strings.TrimSpace(value)

	re := regexp.MustCompile(`^20\d{2}[- ][A-ZÇĞİÖŞÜa-zçğıöşü0-9/()., -]+$`)
	if !re.MatchString(value) {
		return false
	}

	ignored := []string{
		"Sınav Tarihi",
		"Başvuru Tarihleri",
		"Geç Başvuru Günü",
		"Sonuç Tarihi",
	}

	for _, item := range ignored {
		if strings.Contains(value, item) {
			return false
		}
	}

	return true
}

func previousUsefulLine(lines []string, index int) string {
	for i := index - 1; i >= 0 && i >= index-4; i-- {
		line := cleanText(lines[i])
		if line == "" {
			continue
		}

		if strings.Contains(line, "Sınav Tarihi") ||
			strings.Contains(line, "Başvuru Tarihleri") ||
			strings.Contains(line, "Sonuç Tarihi") {
			continue
		}

		if isExamTitle(line) {
			continue
		}

		return line
	}

	return ""
}

func findDateAfterLabel(block []string, label string) *time.Time {
	dates := findDatesAfterLabel(block, label, 1)
	if len(dates) == 0 {
		return nil
	}

	return dates[0]
}

func findDatesAfterLabel(block []string, label string, limit int) []*time.Time {
	var dates []*time.Time
	found := false

	for _, rawLine := range block {
		line := cleanText(rawLine)

		if strings.Contains(line, label) {
			found = true
			continue
		}

		if !found {
			continue
		}

		if isStopLabel(line) {
			break
		}

		date := parseFirstDate(line)
		if date == nil {
			continue
		}

		dates = append(dates, date)

		if len(dates) >= limit {
			break
		}
	}

	return dates
}

func isStopLabel(value string) bool {
	return strings.Contains(value, "Sınav Tarihi:") ||
		strings.Contains(value, "Başvuru Tarihleri:") ||
		strings.Contains(value, "Geç Başvuru Günü:") ||
		strings.Contains(value, "Sonuç Tarihi:")
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

func cleanTextWithNewLines(value string) string {
	value = strings.ReplaceAll(value, "\u00a0", " ")
	value = strings.ReplaceAll(value, "\r\n", "\n")
	value = strings.ReplaceAll(value, "\r", "\n")

	lines := strings.Split(value, "\n")
	cleanedLines := make([]string, 0, len(lines))

	for _, line := range lines {
		line = cleanText(line)
		if line == "" {
			continue
		}

		cleanedLines = append(cleanedLines, line)
	}

	return strings.Join(cleanedLines, "\n")
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

func extractExamCode(value string) string {
	knownCodes := []string{
		"YKS", "KPSS", "ALES", "DGS", "YDS", "e-YDS", "YÖKDİL", "e-YÖKDİL",
		"TUS", "DUS", "STS", "YDUS", "MSÜ", "MEB-AGS", "TR-YÖS", "HMGS",
		"EKPSS", "DİB-MBSTS", "İYÖS", "ÖZYES", "BKUBTS", "EUS",
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
