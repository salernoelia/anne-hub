package handlers

import (
	"anne-hub/models"
	"anne-hub/pkg/db"
	"database/sql"
	"net/http"

	"anne-hub/services"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

// GetAllUsersHandler retrieves all users along with their interests from the database
func GetAllUsersHandler(c echo.Context) error {
    var users []models.UserData

    query := "SELECT id, username, email, password_hash, created_at, age, country, city, first_name, last_name FROM users"

    rows, err := db.DB.Query(query)
    if err != nil {
        return c.JSON(http.StatusInternalServerError, map[string]string{
            "error": "Failed to fetch users: " + err.Error(),
        })
    }
    defer rows.Close()

    for rows.Next() {
        var user models.User

        err := rows.Scan(
            &user.ID,
            &user.Username,
            &user.Email,
            &user.PasswordHash,
            &user.CreatedAt,
            &user.Age,
            &user.Country,
            &user.City,
            &user.FirstName,
            &user.LastName,
        )
        if err != nil {
            return c.JSON(http.StatusInternalServerError, map[string]string{
                "error": "Failed to parse user data: " + err.Error(),
            })
        }

        // Fetch interests for the current user
        userData, err := services.FetchUserData(user.ID)
        if err != nil {
            // Log the error and skip adding interests
            // Alternatively, you can return an error response
            // depending on your application's requirements
            continue
        }

        users = append(users, userData)
    }

    if err = rows.Err(); err != nil {
        return c.JSON(http.StatusInternalServerError, map[string]string{
            "error": "Error iterating over users: " + err.Error(),
        })
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

	query := `
		SELECT id, username, email, password_hash, created_at, age
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
        SET username = $1, email = $2, password_hash = $3, age = $4
        WHERE id = $6
        RETURNING id, created_at
    `
    err = db.DB.QueryRow(query,
        user.Username,
        user.Email,
        user.PasswordHash,
        user.Age,
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


