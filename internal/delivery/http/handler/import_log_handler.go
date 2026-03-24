package handler

import (
	"net/http"
	"strconv"

	"github.com/breamon/sinav-bilgi-sistemi/internal/service"
	"github.com/gin-gonic/gin"
)

type ImportLogHandler struct {
	importLogService *service.ImportLogService
}

func NewImportLogHandler(importLogService *service.ImportLogService) *ImportLogHandler {
	return &ImportLogHandler{importLogService: importLogService}
}

func (h *ImportLogHandler) List(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	logs, err := h.importLogService.List(limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list import logs"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"items": logs,
	})
}
