package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type MockOSYMHandler struct{}

func NewMockOSYMHandler() *MockOSYMHandler {
	return &MockOSYMHandler{}
}

func (h *MockOSYMHandler) GetExams(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"items": []gin.H{
			{
				"external_id": "osym-2026-yks",
				"title":       "2026 YKS",
				"status":      "published",
			},
			{
				"external_id": "osym-2026-kpss",
				"title":       "2026 KPSS",
				"status":      "published",
			},
			{
				"external_id": "osym-2026-ales",
				"title":       "2026 ALES",
				"status":      "draft",
			},
		},
	})
}
