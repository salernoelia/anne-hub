package router

import (
	"anne-hub/handlers"

	"github.com/labstack/echo/v4"
)

func NewRouter() *echo.Echo {
	e := echo.New()


	// General routes
	e.GET("/ok", handlers.OkHandler)
	e.GET("/gh-actions-test", handlers.GitHubActionsTestHandler)
	e.GET("/uuid", handlers.UUIDHandler)


	// Task routes
	e.GET("/tasks", handlers.GetAllTasks)
	// e.GET("/tasks/:id", handlers.GetTaskByID)
	e.GET("/tasks/:id", handlers.GetAllTasksByUserID)
	e.POST("/tasks", handlers.CreateTaskHandler)
	e.PUT("/tasks/:id", handlers.UpdateTaskHandler)
	e.DELETE("/tasks/:id", handlers.DeleteTaskHandler)

	// Interest routes
	e.GET("/interests", handlers.GetAllInterests)
	e.GET("/interests/:id", handlers.GetInterestByID)
	e.POST("/interests", handlers.CreateInterestHandler)
	e.PUT("/interests/:id", handlers.UpdateInterestHandler)
	e.DELETE("/interests/:id", handlers.DeleteInterestHandler)


	// User routes
	e.GET("/users", handlers.GetAllUsersHandler)        
	e.GET("/users/:id", handlers.GetUserHandler)           // Fetch a specific user by ID
	e.POST("/users", handlers.CreateUserHandler)          // Create a new user
	e.PUT("/users/:id", handlers.UpdateUserHandler)       // Update a specific user by ID
	e.DELETE("/users/:id", handlers.DeleteUserHandler)    // Delete a specific user by ID

	// Conversation routes
	e.POST("/ConversationHandler", handlers.ConversationHandler)
	e.POST("/transcribe", handlers.TranscribeAudio)

    // e.GET("/ws", handlers.WebSocketTestHandler)
    e.GET("/ws", handlers.WebSocketConversationHandler)


	return e
}