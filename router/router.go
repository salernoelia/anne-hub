package router

import (
	"anne-hub/handlers"

	"github.com/labstack/echo/v4"
)

func NewRouter() *echo.Echo {
	e := echo.New()


	// Routes
	e.GET("/", handlers.OkHandler)
	e.GET("/uuid", handlers.UUIDHandler)


	// Use CORS middleware
	// e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
	// 	AllowOrigins: []string{"*"}, // You can specify allowed origins here
	// 	AllowMethods: []string{echo.GET, echo.POST, echo.PUT, echo.DELETE},
	// }))

	return e
}