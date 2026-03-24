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
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	source := c.Query("source")
	status := c.Query("status")

	exams, err := h.examService.List(page, limit, source, status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list exams"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"items": exams,
		"pagination": gin.H{
			"page":   page,
			"limit":  limit,
			"source": source,
			"status": status,
		},
	})
}

func (h *ExamHandler) Upcoming(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	exams, err := h.examService.GetUpcoming(limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list upcoming exams"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"items": exams,
		"meta": gin.H{
			"limit": limit,
		},
	})
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

func (h *ExamHandler) Update(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid exam id"})
		return
	}

	var exam domain.Exam
	if err := c.ShouldBindJSON(&exam); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid body"})
		return
	}

	exam.ID = id

	if err := h.examService.Update(&exam); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"exam": exam})
}

func (h *ExamHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid exam id"})
		return
	}

	if err := h.examService.Delete(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete exam"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "exam deleted"})
}
