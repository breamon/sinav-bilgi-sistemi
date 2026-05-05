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
	result, err := h.examImportService.ImportOSYM()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "OSYM exam import failed",
			"error":   err.Error(),
			"result":  result,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "OSYM exams imported successfully",
		"result":  result,
	})
}
