package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/ryory2/test-go-app-todo-go/config"
	"github.com/ryory2/test-go-app-todo-go/internal/handler"
	"github.com/ryory2/test-go-app-todo-go/internal/repository"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// Load configuration
	cfg := config.LoadConfig()

	// Initialize database connection
	db, err := gorm.Open(postgres.Open(cfg.DSN()), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Run migrations (optional if using golang-migrate)
	// db.AutoMigrate(&model.Task{})

	// Initialize repository
	taskRepo := repository.NewTaskRepository(db)

	// Initialize validator
	validate := validator.New()

	// Initialize Gin router
	router := gin.Default()

	// Middleware to set Content-Type to application/json
	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Content-Type", "application/json")
		c.Next()
	})

	// Initialize handlers
	taskHandler := handler.NewTaskHandler(taskRepo, validate)

	// Define routes
	api := router.Group("/api/v1")
	{
		api.GET("/tasks", taskHandler.GetTasks)
		api.POST("/tasks", taskHandler.CreateTask)
		api.PUT("/tasks/:id", taskHandler.UpdateTask)
		api.DELETE("/tasks/:id", taskHandler.DeleteTask)
		api.PATCH("/tasks/:id/toggle", taskHandler.ToggleTask)
	}

	// Start server
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
