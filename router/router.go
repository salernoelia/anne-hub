package router

import (
	"anne-hub/handlers"

	"github.com/labstack/echo/v4"
)

func NewRouter() *echo.Echo {
	e := echo.New()


	// General routes
	e.GET("/", handlers.OkHandler)
	e.GET("/uuid", handlers.UUIDHandler)


	// Task routes
	e.GET("/tasks", handlers.GetAllTasks)
	e.POST("/tasks", handlers.CreateTaskHandler)

	// User routes
	e.GET("/users", handlers.GetAllUsersHandler)          // Fetch all users
	e.POST("/users", handlers.CreateUserHandler)          // Create a new user
	e.PUT("/users/:id", handlers.UpdateUserHandler)       // Update a specific user by ID


	// Cors middleware, in case we need it
	// e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
	// 	AllowOrigins: []string{"*"}, // You can specify allowed origins here
	// 	AllowMethods: []string{echo.GET, echo.POST, echo.PUT, echo.DELETE},
	// }))

	return e
}