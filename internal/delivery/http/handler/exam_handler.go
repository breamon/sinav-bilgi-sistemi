package handler

import (
	"net/http"
	"strconv"

	"github.com/breamon/sinav-bilgi-sistemi/internal/domain"
	"github.com/breamon/sinav-bilgi-sistemi/internal/service"
	"github.com/gin-gonic/gin"
)

type ExamHandler struct {
	examService *service.ExamService
}

func NewExamHandler(examService *service.ExamService) *ExamHandler {
	return &ExamHandler{examService: examService}
}

func (h *ExamHandler) Create(c *gin.Context) {
	var exam domain.Exam

	if err := c.ShouldBindJSON(&exam); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid body"})
		return
	}

	if err := h.examService.Create(&exam); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"exam": exam})
}

func (h *ExamHandler) List(c *gin.Context) {
	exams, err := h.examService.List()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list exams"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"items": exams})
}

func (h *ExamHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid exam id"})
		return
	}

	exam, err := h.examService.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "exam not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"exam": exam})
}
