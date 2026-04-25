package handler

import (
	"net/http"

	"github.com/breamon/sinav-bilgi-sistemi/internal/service"
	"github.com/gin-gonic/gin"
)

type ExamImportHandler struct {
	examImportService *service.ExamImportService
}

func NewExamImportHandler(examImportService *service.ExamImportService) *ExamImportHandler {
	return &ExamImportHandler{
		examImportService: examImportService,
	}
}

func (h *ExamImportHandler) ImportOSYM(c *gin.Context) {
	if err := h.examImportService.ImportOSYM(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "OSYM exams imported successfully",
	})
}
