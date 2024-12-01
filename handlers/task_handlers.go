package handlers

import (
	"anne-hub/models"
	"anne-hub/pkg/db"
	_ "anne-hub/pkg/db"
	"net/http"

	_ "database/sql"

	"github.com/labstack/echo/v4"
)

func GetAllTasks(c echo.Context) error {
	rows, err := db.DB.Query("SELECT * FROM tasks")
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Error getting tasks:" + err.Error(),
		})
	}
	defer rows.Close()

	tasks := []models.Task{}
	for rows.Next() {
		task := models.Task{}
		err := rows.Scan(&task.ID, &task.UserID, &task.Title, &task.Description, &task.DueDate)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{
				"error": "Error scanning task:" + err.Error(),
			})
		}
		tasks = append(tasks, task)
	}

	return c.JSON(http.StatusOK, tasks)
}

func CreateTaskHandler(c echo.Context) error {
	task := new(models.Task)
	if err := c.Bind(task); err != nil {
		 return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request body:" + err.Error(),
		}) 
	}

	_, err := db.DB.Exec(`
		INSERT INTO tasks (
		user_id,
		title,
		description,
		due_date)
		VALUES ($1, $2, $3, $4),
		task.UserID,
		task.Title,
		task.Description,
		task.DueDate,
	`, task);



	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Error inserting task:" + err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, task)
}