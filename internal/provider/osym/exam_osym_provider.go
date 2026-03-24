package osym

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/breamon/sinav-bilgi-sistemi/internal/domain"
	"golang.org/x/net/html"
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

	segments := extractDOMSegments(doc)
	return p.extractExams(segments), nil
}

func (p *ExamOSYMProvider) extractExams(segments []string) []domain.Exam {
	exams := make([]domain.Exam, 0)
	seen := make(map[string]bool)

	for i := 0; i < len(segments); i++ {
		title, ok := extractExamTitleFromSegment(segments[i])
		if !ok {
			continue
		}

		// "(ÖN BAŞVURU)" gibi devam segmentlerini title'a ekle
		if i+1 < len(segments) && isParentheticalSegment(segments[i+1]) {
			title = strings.TrimSpace(title + " " + segments[i+1])
		}
		if i+2 < len(segments) && isParentheticalSegment(segments[i+2]) && !strings.Contains(title, segments[i+2]) {
			title = strings.TrimSpace(title + " " + segments[i+2])
		}

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

		end := i + 1
		for end < len(segments) {
			if nextTitle, ok := extractExamTitleFromSegment(segments[end]); ok && nextTitle != title {
				break
			}
			end++
		}

		block := segments[i:end]

		exam.ExamDate = findFirstDateInLabelLine(block, "Sınav Tarihi")
		appDates := findAllDatesInLabelLine(block, "Başvuru Tarihleri")
		if len(appDates) > 0 {
			exam.ApplicationStartDate = appDates[0]
		}
		if len(appDates) > 1 {
			exam.ApplicationEndDate = appDates[1]
		}
		exam.ResultDate = findFirstDateInLabelLine(block, "Sonuç Tarihi")

		exams = append(exams, exam)
		i = end - 1
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

func extractDOMSegments(doc *goquery.Document) []string {
	segments := make([]string, 0)
	last := ""

	doc.Find("body, body *").Each(func(_ int, s *goquery.Selection) {
		if len(s.Nodes) == 0 {
			return
		}

		text := ownText(s.Nodes[0])
		text = collapseSpaces(strings.TrimSpace(text))
		if text == "" {
			return
		}

		if text == last {
			return
		}

		segments = append(segments, text)
		last = text
	})

	return segments
}

func ownText(n *html.Node) string {
	var b strings.Builder
	for child := n.FirstChild; child != nil; child = child.NextSibling {
		if child.Type == html.TextNode {
			b.WriteString(child.Data)
			b.WriteString(" ")
		}
	}
	return b.String()
}

func extractExamTitleFromSegment(seg string) (string, bool) {
	seg = strings.TrimSpace(seg)
	if seg == "" {
		return "", false
	}

	idx := strings.Index(seg, "2026-")
	if idx == -1 {
		return "", false
	}

	title := strings.TrimSpace(seg[idx:])
	if title == "" {
		return "", false
	}

	return title, true
}

func isParentheticalSegment(seg string) bool {
	seg = strings.TrimSpace(seg)
	return strings.HasPrefix(seg, "(") && strings.HasSuffix(seg, ")")
}

func normalizeLabel(line string) string {
	line = strings.TrimSpace(line)
	line = strings.TrimSuffix(line, ":")
	return collapseSpaces(line)
}

func findFirstDateInLabelLine(block []string, label string) *time.Time {
	label = normalizeLabel(label)

	for _, line := range block {
		if !strings.HasPrefix(normalizeLabel(line), label) {
			continue
		}

		dates := extractDatesFromText(line)
		if len(dates) > 0 {
			return dates[0]
		}
	}

	return nil
}

func findAllDatesInLabelLine(block []string, label string) []*time.Time {
	label = normalizeLabel(label)

	for _, line := range block {
		if !strings.HasPrefix(normalizeLabel(line), label) {
			continue
		}

		return extractDatesFromText(line)
	}

	return nil
}

func extractDatesFromText(text string) []*time.Time {
	text = strings.TrimSpace(text)
	if text == "" {
		return nil
	}

	re := regexp.MustCompile(`\d{2}\.\d{2}\.\d{4}( \d{2}:\d{2})?`)
	matches := re.FindAllString(text, -1)
	if len(matches) == 0 {
		return nil
	}

	results := make([]*time.Time, 0, len(matches))
	for _, m := range matches {
		if t := parseOSYMDate(m); t != nil {
			results = append(results, t)
		}
	}

	return results
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

	return nil
}

func collapseSpaces(s string) string {
	return strings.Join(strings.Fields(s), " ")
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
