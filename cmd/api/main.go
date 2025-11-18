package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/huang/codex-trial/internal/handlers"
	"github.com/huang/codex-trial/internal/models"
	"github.com/huang/codex-trial/pkg/database"
)

func main() {
	// Initialize database
	db, err := database.InitDB()
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}

	// Auto-migrate database schema
	err = db.AutoMigrate(&models.Student{})
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	// Setup Gin router
	r := gin.Default()

	// Initialize handlers with database
	studentHandler := handlers.NewStudentHandler(db)

	// Setup routes
	v1 := r.Group("/api/v1")
	{
		v1.GET("/students", studentHandler.GetStudents)
		v1.GET("/students/:id", studentHandler.GetStudent)
		v1.POST("/students", studentHandler.CreateStudent)
		v1.PUT("/students/:id", studentHandler.UpdateStudent)
		v1.DELETE("/students/:id", studentHandler.DeleteStudent)
	}

	// Health check endpoint
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "healthy",
		})
	})

	log.Println("Server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
