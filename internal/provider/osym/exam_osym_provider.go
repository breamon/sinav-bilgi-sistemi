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

	return p.extractExamsFromCalendarPage(doc), nil
}

func (p *ExamOSYMProvider) extractExamsFromCalendarPage(doc *goquery.Document) []domain.Exam {
	pageText := cleanText(doc.Text())

	seen := make(map[string]bool)
	exams := make([]domain.Exam, 0)

	patterns := []string{
		`2026-MSÜ`,
		`2026-GUY $begin:math:text$ÖN BAŞVURU$end:math:text$`,
		`2026-GUY`,
		`2026-YÖKDİL/1`,
		`2026-MEB-EKYS`,
		`2026-TUS 1\. Dönem`,
		`2026-STS Tıp Doktorluğu 1\. Dönem`,
		`2026-DİB-MBSTS`,
		`2026-YDS/1`,
		`2026-TR-YÖS/1`,
		`2026-EKPSS`,
		`2026-EKPSS/KURA`,
		`2026-HMGS/1`,
		`2026-DUS 1\.Dönem`,
		`2026-STS Diş Hekimliği 1\.Dönem`,
		`2026-YDUS`,
		`2026-ALES/1`,
		`2026-STS Öğretmenlik`,
		`2026-YKS 1\. Oturum $begin:math:text$TYT$end:math:text$`,
		`2026-YKS 2\. Oturum $begin:math:text$AYT$end:math:text$`,
		`2026-YKS 3\. Oturum $begin:math:text$YDT$end:math:text$`,
		`2026-MEB-AGS $begin:math:text$Akademi Giriş Sınavı \\\(AGS$end:math:text$, Öğretmenlik Alan Bilgisi Testi $begin:math:text$ÖABT$end:math:text$\)`,
		`2026-DGS`,
		`2026-ALES/2`,
		`2026-YÖKDİL/2`,
		`2026-ÖZYES`,
		`2026-TUS 2\. Dönem`,
		`2026-STS Tıp Doktorluğu 2\. Dönem`,
		`2026-KPSS Lisans $begin:math:text$Genel Yetenek\-Genel Kültür$end:math:text$`,
		`2026-KPSS Lisans $begin:math:text$Alan Bilgisi$end:math:text$ 1\. gün`,
		`2026-KPSS Lisans $begin:math:text$Alan Bilgisi$end:math:text$ 2\. gün`,
		`2026-HMGS/2`,
		`2026-İYÖS`,
		`2026-KPSS Ön Lisans`,
		`2026-TR-YÖS/2`,
		`2026-BKUBTS`,
		`2026-KPSS Ortaöğretim`,
		`2026-DUS 2\.Dönem`,
		`2026-STS Diş Hekimliği 2\.Dönem`,
		`2026-KPSS Din Hizmetleri Alan Bilgisi Testi $begin:math:text$DHBT$end:math:text$`,
		`2026-EUS`,
		`2026-STS Eczacılık`,
		`2026-YDS/2`,
		`2026-ALES/3`,
	}

	for _, pattern := range patterns {
		re := regexp.MustCompile(pattern)
		matches := re.FindAllString(pageText, -1)

		for _, match := range matches {
			title := strings.TrimSpace(match)
			externalID := slugify(title)

			if seen[externalID] {
				continue
			}
			seen[externalID] = true

			exams = append(exams, domain.Exam{
				Source:     "osym",
				ExternalID: stringPtr(externalID),
				Title:      title,
				Status:     "published",
			})
		}
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

func cleanText(s string) string {
	s = strings.ReplaceAll(s, "\n", " ")
	s = strings.ReplaceAll(s, "\t", " ")
	s = strings.ReplaceAll(s, "\r", " ")

	for strings.Contains(s, "  ") {
		s = strings.ReplaceAll(s, "  ", " ")
	}

	return strings.TrimSpace(s)
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
