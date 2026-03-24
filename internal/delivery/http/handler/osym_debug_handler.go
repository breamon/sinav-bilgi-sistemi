package handler

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/gin-gonic/gin"
	"golang.org/x/net/html"
)

type OSYMDebugHandler struct{}

func NewOSYMDebugHandler() *OSYMDebugHandler {
	return &OSYMDebugHandler{}
}

func (h *OSYMDebugHandler) RawSegments(c *gin.Context) {
	target := strings.TrimSpace(c.Query("q"))
	if target == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "q query param is required",
		})
		return
	}

	req, err := http.NewRequest(http.MethodGet, "https://www.osym.gov.tr/tr,8797/takvim.html", nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	req.Header.Set("User-Agent", "sinav-bilgi-sistemi/1.0")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	segments := extractDebugDOMSegments(doc)

	matchIndexes := make([]int, 0)
	for i, seg := range segments {
		if strings.Contains(strings.ToLower(seg), strings.ToLower(target)) {
			matchIndexes = append(matchIndexes, i)
		}
	}

	if len(matchIndexes) == 0 {
		c.JSON(http.StatusOK, gin.H{
			"query":   target,
			"message": "no matching segment found",
			"items":   []gin.H{},
		})
		return
	}

	items := make([]gin.H, 0)
	seen := make(map[string]bool)

	for _, idx := range matchIndexes {
		start := idx - 8
		if start < 0 {
			start = 0
		}
		end := idx + 20
		if end > len(segments) {
			end = len(segments)
		}

		key := fmt.Sprintf("%d-%d", start, end)
		if seen[key] {
			continue
		}
		seen[key] = true

		block := make([]gin.H, 0, end-start)
		for i := start; i < end; i++ {
			block = append(block, gin.H{
				"index": i,
				"text":  segments[i],
			})
		}

		items = append(items, gin.H{
			"match_index": idx,
			"window":      block,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"query": target,
		"items": items,
	})
}

func extractDebugDOMSegments(doc *goquery.Document) []string {
	segments := make([]string, 0)
	last := ""

	doc.Find("body, body *").Each(func(_ int, s *goquery.Selection) {
		if len(s.Nodes) == 0 {
			return
		}

		text := ownDebugText(s.Nodes[0])
		text = collapseDebugSpaces(strings.TrimSpace(text))
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

func ownDebugText(n *html.Node) string {
	var b strings.Builder
	for child := n.FirstChild; child != nil; child = child.NextSibling {
		if child.Type == html.TextNode {
			b.WriteString(child.Data)
			b.WriteString(" ")
		}
	}
	return b.String()
}

func collapseDebugSpaces(s string) string {
	return strings.Join(strings.Fields(s), " ")
}
