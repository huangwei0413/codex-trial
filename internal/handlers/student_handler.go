package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/huang/codex-trial/internal/models"
	"gorm.io/gorm"
)

type StudentHandler struct {
	db *gorm.DB
}

func NewStudentHandler(db *gorm.DB) *StudentHandler {
	return &StudentHandler{db: db}
}

func (h *StudentHandler) GetStudents(c *gin.Context) {
	var students []models.Student
	err := h.db.Find(&students).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve students"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": students})
}

func (h *StudentHandler) GetStudent(c *gin.Context) {
	id := c.Param("id")
	var student models.Student
	err := h.db.Where("id = ?", id).First(&student).Error
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Student not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": student})
}

func (h *StudentHandler) CreateStudent(c *gin.Context) {
	var student models.Student
	if err := c.ShouldBindJSON(&student); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.db.Create(&student).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create student"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": student})
}

func (h *StudentHandler) UpdateStudent(c *gin.Context) {
	id := c.Param("id")
	var student models.Student

	// Check if student exists
	if err := h.db.Where("id = ?", id).First(&student).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Student not found"})
		return
	}

	// Bind updated data
	var updateData models.Student
	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Update student
	err := h.db.Model(&student).Updates(updateData).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update student"})
		return
	}

	// Get updated student
	h.db.Where("id = ?", id).First(&student)
	c.JSON(http.StatusOK, gin.H{"data": student})
}

func (h *StudentHandler) DeleteStudent(c *gin.Context) {
	id := c.Param("id")
	err := h.db.Where("id = ?", id).Delete(&models.Student{}).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete student"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Student deleted successfully"})
}
