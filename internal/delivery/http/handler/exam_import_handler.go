package handler

import (
	"net/http"

	"github.com/breamon/sinav-bilgi-sistemi/internal/service"
	"github.com/gin-gonic/gin"
)

type ExamImportHandler struct {
	importService    *service.ExamImportService
	importLogService *service.ImportLogService
}

func NewExamImportHandler(
	importService *service.ExamImportService,
	importLogService *service.ImportLogService,
) *ExamImportHandler {
	return &ExamImportHandler{
		importService:    importService,
		importLogService: importLogService,
	}
}

func (h *ExamImportHandler) Import(c *gin.Context) {
	exams, err := h.importService.Import()
	if err != nil {
		errMsg := err.Error()
		_ = h.importLogService.Create(
			h.importService.ProviderName(),
			"failed",
			0,
			&errMsg,
		)

		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "failed to import exams",
			"details": err.Error(),
		})
		return
	}

	_ = h.importLogService.Create(
		h.importService.ProviderName(),
		"success",
		len(exams),
		nil,
	)

	c.JSON(http.StatusOK, gin.H{
		"message": "exams imported",
		"items":   exams,
	})
}
