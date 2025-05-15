package controllers

import (
	"database/sql"
	"net/http"
	"user-management/database"
	"user-management/models"

	"github.com/gin-gonic/gin"
)

func CreateUser(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		// Return 400 Bad Request on validation error
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	query := `INSERT INTO users 
    (username, firstname, lastname, email, avatar, phone, dob, country, city, street_name, street_address)
    VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	_, err := database.DB.Exec(query,
		user.Username, user.FirstName, user.LastName, user.Email, user.Avatar, user.Phone,
		user.DOB, user.Country, user.City, user.StreetName, user.StreetAddress)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}

func GetUsers(c *gin.Context) {
	filter := c.Query("filter")
	sort := c.Query("sort")

	baseQuery := "SELECT * FROM users"
	var args []interface{}

	if filter != "" {
		baseQuery += " WHERE username LIKE ? OR firstname LIKE ? OR lastname LIKE ?"
		f := "%" + filter + "%"
		args = append(args, f, f, f)
	}

	if sort != "" && (sort == "username" || sort == "firstname" || sort == "lastname") {
		baseQuery += " ORDER BY " + sort
	}

	rows, err := database.DB.Query(baseQuery, args...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Query failed"})
		return
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var u models.User
		_ = rows.Scan(&u.ID, &u.Username, &u.FirstName, &u.LastName, &u.Email, &u.Avatar, &u.Phone,
			&u.DOB, &u.Country, &u.City, &u.StreetName, &u.StreetAddress)
		users = append(users, u)
	}

	c.JSON(http.StatusOK, users)
}

func GetUser(c *gin.Context) {
	id := c.Param("id")
	var u models.User
	row := database.DB.QueryRow("SELECT * FROM users WHERE id = ?", id)
	err := row.Scan(&u.ID, &u.Username, &u.FirstName, &u.LastName, &u.Email, &u.Avatar, &u.Phone,
		&u.DOB, &u.Country, &u.City, &u.StreetName, &u.StreetAddress)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		}
		return
	}
	c.JSON(http.StatusOK, u)
}

func UpdateUser(c *gin.Context) {
	id := c.Param("id")
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}

	query := `UPDATE users SET 
        username=?, firstname=?, lastname=?, email=?, avatar=?, phone=?, 
        dob=?, country=?, city=?, street_name=?, street_address=? WHERE id=?`
	_, err := database.DB.Exec(query,
		user.Username, user.FirstName, user.LastName, user.Email, user.Avatar, user.Phone,
		user.DOB, user.Country, user.City, user.StreetName, user.StreetAddress, id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Update failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User updated"})
}

func DeleteUser(c *gin.Context) {
	id := c.Param("id")
	_, err := database.DB.Exec("DELETE FROM users WHERE id = ?", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Delete failed"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "User deleted"})
}
