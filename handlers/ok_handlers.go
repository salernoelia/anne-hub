package handlers

import (
	"net/http"

	"anne-hub/pkg/uuid"

	"github.com/labstack/echo/v4"
)



func OkHandler(c echo.Context) error {
	// return c.String(http.StatusOK, "OK")
	return c.JSON(http.StatusOK, map[string]string{
		"message": "OK",
	})
}

func UUIDHandler(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{
		"uuid": uuid.CreateUUID(),
	})
}