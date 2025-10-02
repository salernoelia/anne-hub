package router

import (
	"anne-hub/handlers"

	"github.com/labstack/echo/v4"
)

func NewRouter() *echo.Echo {
	e := echo.New()

	e.GET("/ok", handlers.OkHandler)
	e.GET("/gh-actions-test", handlers.GitHubActionsTestHandler)
	e.GET("/uuid", handlers.UUIDHandler)

	e.GET("/tasks", handlers.GetAllTasks)
	e.GET("/tasks/:id", handlers.GetAllTasksByUserID)
	e.POST("/tasks", handlers.CreateTaskHandler)
	e.PUT("/tasks/:id", handlers.UpdateTaskHandler)
	e.DELETE("/tasks/:id", handlers.DeleteTaskHandler)

	e.GET("/interests", handlers.GetAllInterests)
	e.GET("/interests/:id", handlers.GetInterestByID)
	e.POST("/interests", handlers.CreateInterestHandler)
	e.PUT("/interests/:id", handlers.UpdateInterestHandler)
	e.DELETE("/interests/:id", handlers.DeleteInterestHandler)

	e.GET("/users", handlers.GetAllUsersHandler)
	e.GET("/users/:id", handlers.GetUserHandler)
	e.POST("/users", handlers.CreateUserHandler)
	e.PUT("/users/:id", handlers.UpdateUserHandler)
	e.DELETE("/users/:id", handlers.DeleteUserHandler)

	e.POST("/ConversationHandler", handlers.ConversationHandler)
	e.POST("/transcribe", handlers.TranscribeAudio)

	e.GET("/ws", handlers.WebSocketConversationHandler)

	return e
}
