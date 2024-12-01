package handlers

import (
	"anne-hub/models"
	"anne-hub/pkg/db"
	"database/sql"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/lib/pq"
)

// GetAllTasks retrieves all tasks from the database with pagination and error handling
func GetAllTasks(c echo.Context) error {
	// Pagination parameters
	page, err := strconv.Atoi(c.QueryParam("page"))
	if err != nil || page < 1 {
		page = 1
	}
	limit, err := strconv.Atoi(c.QueryParam("limit"))
	if err != nil || limit < 1 {
		limit = 10 // default limit
	}
	offset := (page - 1) * limit

	query := `
		SELECT id, user_id, title, description, due_date, completed, created_at, interest_links 
		FROM tasks 
		ORDER BY created_at DESC 
		LIMIT $1 OFFSET $2
	`

	rows, err := db.DB.Queryx(query, limit, offset)
	if err != nil {
		// Log the error for internal monitoring
		c.Logger().Errorf("Error querying tasks: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to retrieve tasks.",
		})
	}
	defer rows.Close()

	tasks := []models.Task{}
	for rows.Next() {
		task := models.Task{}
		var interestLinks pq.StringArray

		err := rows.Scan(
			&task.ID,
			&task.UserID,
			&task.Title,
			&task.Description,
			&task.DueDate,
			&task.Completed,
			&task.CreatedAt,
			&interestLinks,
		)
		if err != nil {
			c.Logger().Errorf("Error scanning task: %v", err)
			return c.JSON(http.StatusInternalServerError, map[string]string{
				"error": "Failed to process tasks.",
			})
		}

		task.InterestLinks = []string(interestLinks)
		tasks = append(tasks, task)
	}

	// Check for errors from iterating over rows
	if err = rows.Err(); err != nil {
		c.Logger().Errorf("Row iteration error: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Error retrieving tasks.",
		})
	}

	return c.JSON(http.StatusOK, tasks)
}

// GetTaskByID retrieves a single task by its ID with error handling
func GetTaskByID(c echo.Context) error {
	idParam := c.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil || id < 1 {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid task ID.",
		})
	}

	var task models.Task
	var interestLinks pq.StringArray

	query := `
		SELECT id, user_id, title, description, due_date, completed, created_at, interest_links 
		FROM tasks 
		WHERE id = $1
	`

	err = db.DB.QueryRowx(query, id).Scan(
		&task.ID,
		&task.UserID,
		&task.Title,
		&task.Description,
		&task.DueDate,
		&task.Completed,
		&task.CreatedAt,
		&interestLinks,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.JSON(http.StatusNotFound, map[string]string{
				"error": "Task not found.",
			})
		}
		c.Logger().Errorf("Error retrieving task: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to retrieve task.",
		})
	}

	task.InterestLinks = []string(interestLinks)
	return c.JSON(http.StatusOK, task)
}

// CreateTaskHandler creates a new task with validation and error handling
func CreateTaskHandler(c echo.Context) error {
	task := new(models.Task)
	if err := c.Bind(task); err != nil {
		c.Logger().Warnf("Bind error: %v", err)
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request payload.",
		})
	}

	// Input Validation
	if err := validateTaskInput(task); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
	}

	// Set default values if necessary
	if task.CreatedAt.IsZero() {
		task.CreatedAt = time.Now()
	}
	if task.InterestLinks == nil {
		task.InterestLinks = []string{}
	}

	query := `
		INSERT INTO tasks (
			user_id,
			title,
			description,
			due_date,
			completed,
			created_at,
			interest_links
		) VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id
	`

	// Use pq.Array to handle the []string for interest_links
	err := db.DB.QueryRow(query,
		task.UserID,
		task.Title,
		task.Description,
		task.DueDate,
		task.Completed,
		task.CreatedAt,
		pq.Array(task.InterestLinks),
	).Scan(&task.ID)
	if err != nil {
		c.Logger().Errorf("Error inserting task: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to create task.",
		})
	}

	return c.JSON(http.StatusCreated, task)
}

// UpdateTaskHandler updates an existing task by its ID with validation and error handling
func UpdateTaskHandler(c echo.Context) error {
	idParam := c.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil || id < 1 {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid task ID.",
		})
	}

	task := new(models.Task)
	if err := c.Bind(task); err != nil {
		c.Logger().Warnf("Bind error: %v", err)
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request payload.",
		})
	}

	// Input Validation
	if err := validateTaskInput(task); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
	}

	query := `
		UPDATE tasks SET
			user_id = $1,
			title = $2,
			description = $3,
			due_date = $4,
			completed = $5,
			interest_links = $6
		WHERE id = $7
	`

	res, err := db.DB.Exec(query,
		task.UserID,
		task.Title,
		task.Description,
		task.DueDate,
		task.Completed,
		pq.Array(task.InterestLinks),
		id,
	)
	if err != nil {
		c.Logger().Errorf("Error updating task: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to update task.",
		})
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		c.Logger().Errorf("Error fetching update result: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to process update.",
		})
	}
	if rowsAffected == 0 {
		return c.JSON(http.StatusNotFound, map[string]string{
			"error": "Task not found.",
		})
	}

	// Fetch the updated task to return
	updatedTask, err := fetchTaskByID(id)
	if err != nil {
		c.Logger().Errorf("Error fetching updated task: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Task updated but failed to retrieve.",
		})
	}

	return c.JSON(http.StatusOK, updatedTask)
}

// DeleteTaskHandler deletes a task by its ID with comprehensive error handling
func DeleteTaskHandler(c echo.Context) error {
	idParam := c.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil || id < 1 {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid task ID.",
		})
	}

	query := `DELETE FROM tasks WHERE id = $1`

	res, err := db.DB.Exec(query, id)
	if err != nil {
		c.Logger().Errorf("Error deleting task: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to delete task.",
		})
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		c.Logger().Errorf("Error fetching delete result: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to process deletion.",
		})
	}
	if rowsAffected == 0 {
		return c.JSON(http.StatusNotFound, map[string]string{
			"error": "Task not found.",
		})
	}

	return c.NoContent(http.StatusNoContent)
}

// fetchTaskByID is a helper function to retrieve a task after update
func fetchTaskByID(id int64) (*models.Task, error) {
	var task models.Task
	var interestLinks pq.StringArray

	query := `
		SELECT id, user_id, title, description, due_date, completed, created_at, interest_links 
		FROM tasks 
		WHERE id = $1
	`

	err := db.DB.QueryRowx(query, id).Scan(
		&task.ID,
		&task.UserID,
		&task.Title,
		&task.Description,
		&task.DueDate,
		&task.Completed,
		&task.CreatedAt,
		&interestLinks,
	)
	if err != nil {
		return nil, err
	}

	task.InterestLinks = []string(interestLinks)
	return &task, nil
}

// validateTaskInput performs basic validation on the Task input
func validateTaskInput(task *models.Task) error {
	if len(task.Title) == 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "Title is required.")
	}
	if len(task.Title) > 255 {
		return echo.NewHTTPError(http.StatusBadRequest, "Title cannot exceed 255 characters.")
	}
	if task.DueDate != nil && task.DueDate.Before(time.Now()) {
		return echo.NewHTTPError(http.StatusBadRequest, "Due date cannot be in the past.")
	}
	// Add more validation rules as needed
	return nil
}
