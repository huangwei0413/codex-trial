//go:build integration
// +build integration

package integration

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/huang/codex-trial/internal/handlers"
	"github.com/huang/codex-trial/internal/models"
	"github.com/huang/codex-trial/pkg/database"
	"gorm.io/driver/sqlite"
)

func setupIntegrationTest(t *testing.T) (*gin.Engine, *gorm.DB) {
	// Use in-memory SQLite for testing
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to create test database: %v", err)
	}

	// Auto-migrate
	err = db.AutoMigrate(&models.Student{})
	if err != nil {
		t.Fatalf("Failed to migrate test database: %v", err)
	}

	// Setup Gin router
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	handler := handlers.NewStudentHandler(db)

	v1 := r.Group("/api/v1")
	{
		v1.GET("/students", handler.GetStudents)
	v1.GET("/students/:id", handler.GetStudent)
		v1.POST("/students", handler.CreateStudent)
		v1.PUT("/students/:id", handler.UpdateStudent)
		v1.DELETE("/students/:id", handler.DeleteStudent)
	}

	return r, db
}

func TestAPI_StudentCRUD(t *testing.T) {
	router, _ := setupIntegrationTest(t)

	// Test 1: Create student
	student := models.Student{
		Name:  "Jane Doe",
		Email: "jane@example.com",
		Age:   22,
		Major: "Data Science",
		GPA:   3.9,
	}

	jsonData, _ := json.Marshal(student)
	req, _ := http.NewRequest("POST", "/api/v1/students", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected status %d, got %d", http.StatusCreated, w.Code)
	}

	// Test 2: Get all students
	req, _ = http.NewRequest("GET", "/api/v1/students", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	// Test 3: Health check
	req, _ = http.NewRequest("GET", "/health", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Health check expected status %d, got %d", http.StatusOK, w.Code)
	}

	t.Logf("Integration tests completed successfully!")
}
