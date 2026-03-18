package handler

import (
	"net/http"

	"github.com/breamon/sinav-bilgi-sistemi/internal/service"
	"github.com/gin-gonic/gin"
)

type ExamImportHandler struct {
	importService *service.ExamImportService
}

func NewExamImportHandler(importService *service.ExamImportService) *ExamImportHandler {
	return &ExamImportHandler{importService: importService}
}

func (h *ExamImportHandler) Import(c *gin.Context) {
	exams, err := h.importService.Import()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "failed to import exams",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "exams imported",
		"items":   exams,
	})
}
