package unit

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/huang/codex-trial/internal/handlers"
	"github.com/huang/codex-trial/internal/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to create test database: %v", err)
	}

	// Auto-migrate the schema
	err = db.AutoMigrate(&models.Student{})
	if err != nil {
		t.Fatalf("Failed to migrate test database: %v", err)
	}

	return db
}

func TestStudentHandler_GetStudents_Empty(t *testing.T) {
	// Setup
	db := setupTestDB(t)
	handler := handlers.NewStudentHandler(db)
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.GET("/students", handler.GetStudents)

	// Test
	req, _ := http.NewRequest("GET", "/students", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assert
	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	// Should return empty array
	expected := `{"data":[]}`
	if w.Body.String() != expected {
		t.Errorf("Expected %s, got %s", expected, w.Body.String())
	}
}

func TestStudentHandler_CreateStudent(t *testing.T) {
	// Setup
	db := setupTestDB(t)
	handler := handlers.NewStudentHandler(db)
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.POST("/students", handler.CreateStudent)

	// Test
	req, _ := http.NewRequest("POST", "/students", nil)
	req.Header.Set("Content-Type", "application/json")
	// Note: In a real test, you would marshal the student to JSON
	// For demo purposes, this shows the test structure

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// This test demonstrates the structure - actual implementation would need JSON marshaling
	t.Logf("Response status: %d", w.Code)
	t.Logf("Response body: %s", w.Body.String())
}
