package handlers

import (
	"anne-hub/models"
	"anne-hub/pkg/db"
	"database/sql"
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

// GetAllInterests retrieves all interests from the database with pagination and error handling
func GetAllInterests(c echo.Context) error {
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
        SELECT id, user_id, created_at, updated_at, name, description, level, level_accuracy
        FROM interests 
        ORDER BY created_at DESC 
        LIMIT $1 OFFSET $2
    `

    rows, err := db.DB.Queryx(query, limit, offset)
    if err != nil {
        // Log the error for internal monitoring
        c.Logger().Errorf("Error querying interests: %v", err)
        return c.JSON(http.StatusInternalServerError, map[string]string{
            "error": "Failed to retrieve interests.",
        })
    }
    defer rows.Close()

    interests := []models.Interest{}
    for rows.Next() {
        interest := models.Interest{}

        err := rows.Scan(
            &interest.ID,
            &interest.UserID,
            &interest.CreatedAt,
            &interest.UpdatedAt,
            &interest.Name,
            &interest.Description,
            &interest.Level,
            &interest.LevelAccuracy,
        )
        if err != nil {
            c.Logger().Errorf("Error scanning interest: %v", err)
            return c.JSON(http.StatusInternalServerError, map[string]string{
                "error": "Failed to process interests.",
            })
        }

        interests = append(interests, interest)
    }

    // Check for errors from iterating over rows
    if err = rows.Err(); err != nil {
        c.Logger().Errorf("Row iteration error: %v", err)
        return c.JSON(http.StatusInternalServerError, map[string]string{
            "error": "Error retrieving interests.",
        })
    }

    return c.JSON(http.StatusOK, interests)
}

// GetInterestByID retrieves a single interest by its ID with error handling
func GetInterestByID(c echo.Context) error {
    idParam := c.Param("id")
    id, err := strconv.ParseInt(idParam, 10, 64)
    if err != nil || id < 1 {
        return c.JSON(http.StatusBadRequest, map[string]string{
            "error": "Invalid interest ID.",
        })
    }

    var interest models.Interest

    query := `
        SELECT id, user_id, created_at, updated_at, name, description, level, level_accuracy
        FROM interests 
        WHERE id = $1
    `

    err = db.DB.QueryRowx(query, id).Scan(
        &interest.ID,
        &interest.UserID,
        &interest.CreatedAt,
        &interest.UpdatedAt,
        &interest.Name,
        &interest.Description,
        &interest.Level,
        &interest.LevelAccuracy,
    )
    if err != nil {
        if err == sql.ErrNoRows {
            return c.JSON(http.StatusNotFound, map[string]string{
                "error": "Interest not found.",
            })
        }
        c.Logger().Errorf("Error retrieving interest: %v", err)
        return c.JSON(http.StatusInternalServerError, map[string]string{
            "error": "Failed to retrieve interest.",
        })
    }

    return c.JSON(http.StatusOK, interest)
}

// GetAllInterestsByUserID retrieves all interests for a specific user
func GetAllInterestsByUserID(c echo.Context) error {
    userIDParam := c.Param("id")
    userID, err := uuid.Parse(userIDParam)
    if err != nil {
        return c.JSON(http.StatusBadRequest, map[string]string{
            "error": "Invalid user ID format",
        })
    }

    query := `
        SELECT id, user_id, created_at, updated_at, name, description, level, level_accuracy
        FROM interests 
        WHERE user_id = $1
        ORDER BY created_at DESC 
    `

    interests := []models.Interest{}
    err = db.DB.Select(&interests, query, userID)
    if err != nil {
        c.Logger().Errorf("Error querying interests: %v", err)
        return c.JSON(http.StatusInternalServerError, map[string]string{
            "error": "Failed to retrieve interests",
        })
    }
    return c.JSON(http.StatusOK, interests)
}

// CreateInterestHandler creates a new interest with validation and error handling
func CreateInterestHandler(c echo.Context) error {
    interest := new(models.Interest)
    if err := c.Bind(interest); err != nil {
        c.Logger().Warnf("Bind error: %v", err)
        return c.JSON(http.StatusBadRequest, map[string]string{
            "error": "Invalid request payload.",
        })
    }

    // Input Validation
    if err := validateInterestInput(interest); err != nil {
        return c.JSON(http.StatusBadRequest, map[string]string{
            "error": err.Error(),
        })
    }

    // Set default values if necessary
    if interest.CreatedAt == "" {
        interest.CreatedAt = time.Now().Format(time.RFC3339)
    }
    if interest.UpdatedAt == "" {
        interest.UpdatedAt = time.Now().Format(time.RFC3339)
    }

    query := `
        INSERT INTO interests (
            user_id,
            created_at,
            updated_at,
            name,
            description,
            level,
            level_accuracy
        ) VALUES ($1, $2, $3, $4, $5, $6, $7)
        RETURNING id
    `

    err := db.DB.QueryRow(query,
        interest.UserID,
        interest.CreatedAt,
        interest.UpdatedAt,
        interest.Name,
        interest.Description,
        interest.Level,
        interest.LevelAccuracy,
    ).Scan(&interest.ID)
    if err != nil {
        c.Logger().Errorf("Error inserting interest: %v", err)
        return c.JSON(http.StatusInternalServerError, map[string]string{
            "error": "Failed to create interest.",
        })
    }

    return c.JSON(http.StatusCreated, interest)
}

// UpdateInterestHandler updates an existing interest by its ID with validation and error handling
func UpdateInterestHandler(c echo.Context) error {
    idParam := c.Param("id")
    id, err := strconv.ParseInt(idParam, 10, 64)
    if err != nil || id < 1 {
        return c.JSON(http.StatusBadRequest, map[string]string{
            "error": "Invalid interest ID.",
        })
    }

    interest := new(models.Interest)
    if err := c.Bind(interest); err != nil {
        c.Logger().Warnf("Bind error: %v", err)
        return c.JSON(http.StatusBadRequest, map[string]string{
            "error": "Invalid request payload.",
        })
    }

    // Input Validation
    if err := validateInterestInput(interest); err != nil {
        return c.JSON(http.StatusBadRequest, map[string]string{
            "error": err.Error(),
        })
    }

    // Update the UpdatedAt field
    interest.UpdatedAt = time.Now().Format(time.RFC3339)

    query := `
        UPDATE interests SET
            user_id = $1,
            updated_at = $2,
            name = $3,
            description = $4,
            level = $5,
            level_accuracy = $6
        WHERE id = $7
    `

    res, err := db.DB.Exec(query,
        interest.UserID,
        interest.UpdatedAt,
        interest.Name,
        interest.Description,
        interest.Level,
        interest.LevelAccuracy,
        id,
    )
    if err != nil {
        c.Logger().Errorf("Error updating interest: %v", err)
        return c.JSON(http.StatusInternalServerError, map[string]string{
            "error": "Failed to update interest.",
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
            "error": "Interest not found.",
        })
    }

    // Fetch the updated interest to return
    updatedInterest, err := fetchInterestByID(id)
    if err != nil {
        c.Logger().Errorf("Error fetching updated interest: %v", err)
        return c.JSON(http.StatusInternalServerError, map[string]string{
            "error": "Interest updated but failed to retrieve.",
        })
    }

    return c.JSON(http.StatusOK, updatedInterest)
}

// DeleteInterestHandler deletes an interest by its ID with comprehensive error handling
func DeleteInterestHandler(c echo.Context) error {
    idParam := c.Param("id")
    id, err := strconv.ParseInt(idParam, 10, 64)
    if err != nil || id < 1 {
        return c.JSON(http.StatusBadRequest, map[string]string{
            "error": "Invalid interest ID.",
        })
    }

    query := `DELETE FROM interests WHERE id = $1`

    res, err := db.DB.Exec(query, id)
    if err != nil {
        c.Logger().Errorf("Error deleting interest: %v", err)
        return c.JSON(http.StatusInternalServerError, map[string]string{
            "error": "Failed to delete interest.",
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
            "error": "Interest not found.",
        })
    }

    return c.NoContent(http.StatusNoContent)
}

// fetchInterestByID is a helper function to retrieve an interest after update
func fetchInterestByID(id int64) (*models.Interest, error) {
    var interest models.Interest

    query := `
        SELECT id, user_id, created_at, updated_at, name, description, level, level_accuracy
        FROM interests 
        WHERE id = $1
    `

    err := db.DB.QueryRowx(query, id).Scan(
        &interest.ID,
        &interest.UserID,
        &interest.CreatedAt,
        &interest.UpdatedAt,
        &interest.Name,
        &interest.Description,
        &interest.Level,
        &interest.LevelAccuracy,
    )
    if err != nil {
        return nil, err
    }

    return &interest, nil
}

// validateInterestInput performs basic validation on the Interest input
func validateInterestInput(interest *models.Interest) error {
    if len(interest.Name) == 0 {
        return echo.NewHTTPError(http.StatusBadRequest, "Name is required.")
    }
    if len(interest.Name) > 255 {
        return echo.NewHTTPError(http.StatusBadRequest, "Name cannot exceed 255 characters.")
    }
    // Add more validation rules as needed
    return nil
}