package handlers

import (
	"anne-hub/models"
	"anne-hub/pkg/db"
	"database/sql"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/lib/pq"
)

// GetAllUsersHandler retrieves all users from the database
func GetAllUsersHandler(c echo.Context) error {
	var users []models.User

	query := "SELECT id, username, email, password_hash, created_at, age, interests FROM users"

	rows, err := db.DB.Query(query)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to fetch users: " + err.Error(),
		})
	}
	defer rows.Close()

	for rows.Next() {
		var user models.User
		var interests []string // Temp variable to handle array scanning

		err := rows.Scan(
			&user.ID,
			&user.Username,
			&user.Email,
			&user.PasswordHash,
			&user.CreatedAt,
			&user.Age,
			pq.Array(&interests), // Use pq.Array for PostgreSQL arrays
		)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{
				"error": "Failed to parse user data: " + err.Error(),
			})
		}
		user.Interests = interests // Assign parsed interests to the user struct
		users = append(users, user)
	}

	return c.JSON(http.StatusOK, users)
}

// get user by id parameter
func GetUserHandler(c echo.Context) error {
	idParam := c.Param("id") // UUID from the URL path parameter

	// Validate UUID format
	userID, err := uuid.Parse(idParam)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid user ID format.",
		})
	}

	var user models.User
	var interests []string // Temp variable to handle array scanning

	query := `
		SELECT id, username, email, password_hash, created_at, age, interests
		FROM users
		WHERE id = $1
	`
	err = db.DB.QueryRow(query, userID).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.PasswordHash,
		&user.CreatedAt,
		&user.Age,
		pq.Array(&interests), // Use pq.Array for PostgreSQL arrays
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.JSON(http.StatusNotFound, map[string]string{
				"error": "User not found.",
			})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to fetch user: " + err.Error(),
		})
	}

	user.Interests = interests // Assign parsed interests to the user struct

	return c.JSON(http.StatusOK, user)
}


// CreateUserHandler creates a new user in the database
func CreateUserHandler(c echo.Context) error {
    user := new(models.User)

    if err := c.Bind(user); err != nil {
        return c.JSON(http.StatusBadRequest, map[string]string{
            "error": "Invalid request body: " + err.Error(),
        })
    }

    // Validate required fields
    if user.Username == "" || user.Email == "" || user.PasswordHash == "" {
        return c.JSON(http.StatusBadRequest, map[string]string{
            "error": "Username, email, and password_hash are required fields.",
        })
    }

    query := `
        INSERT INTO users (id, username, email, password_hash, age, interests)
        VALUES (gen_random_uuid(), $1, $2, $3, $4, $5)
        RETURNING id, created_at
    `
    err := db.DB.QueryRow(query,
        user.Username,
        user.Email,
        user.PasswordHash,
        user.Age,
        pq.Array(user.Interests), // Use pq.Array for PostgreSQL arrays
    ).Scan(&user.ID, &user.CreatedAt)
    if err != nil {
        return c.JSON(http.StatusInternalServerError, map[string]string{
            "error": "Failed to create user: " + err.Error(),
        })
    }

    return c.JSON(http.StatusCreated, user)
}



// UpdateUserHandler updates an existing user in the database
func UpdateUserHandler(c echo.Context) error {
    idParam := c.Param("id") // UUID from the URL path parameter

    // Validate UUID format
    userID, err := uuid.Parse(idParam)
    if err != nil {
        return c.JSON(http.StatusBadRequest, map[string]string{
            "error": "Invalid user ID format.",
        })
    }

    user := new(models.User)

    if err := c.Bind(user); err != nil {
        return c.JSON(http.StatusBadRequest, map[string]string{
            "error": "Invalid request body: " + err.Error(),
        })
    }


    query := `
        UPDATE users
        SET username = $1, email = $2, password_hash = $3, age = $4, interests = $5
        WHERE id = $6
        RETURNING id, created_at
    `
    err = db.DB.QueryRow(query,
        user.Username,
        user.Email,
        user.PasswordHash,
        user.Age,
        pq.Array(user.Interests), // Use pq.Array for PostgreSQL arrays
        userID,
    ).Scan(&user.ID, &user.CreatedAt)
    if err != nil {
        if err == sql.ErrNoRows {
            return c.JSON(http.StatusNotFound, map[string]string{
                "error": "User not found.",
            })
        }
        return c.JSON(http.StatusInternalServerError, map[string]string{
            "error": "Failed to update user: " + err.Error(),
        })
    }

    return c.JSON(http.StatusOK, user)
}


// delete user endpoint
func DeleteUserHandler(c echo.Context) error {
	idParam := c.Param("id") // UUID from the URL path parameter

	// Validate UUID format
	userID, err := uuid.Parse(idParam)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid user ID format.",
		})
	}

	query := `
		DELETE FROM users
		WHERE id = $1
	`
	_, err = db.DB.Exec(query, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.JSON(http.StatusNotFound, map[string]string{
				"error": "User not found.",
			})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to delete user: " + err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]string{
		"message": "User deleted successfully.",
	})
}