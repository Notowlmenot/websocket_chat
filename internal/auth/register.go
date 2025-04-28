package auth

import (
	"chat/internal/database"
	"errors"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

func RegisterHandler(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")
	if username == "" || password == "" {
		c.String(http.StatusBadRequest, "missing username or password")
		return
	}
	if strings.ContainsAny(username, "* || - || / || &") {
		c.String(http.StatusBadRequest, "username must not contain special characters")
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		if err.Error() == "bcrypt: password length exceeds 72 bytes" {
			c.String(http.StatusBadRequest, "The password must be less than 72 characters long")
			return
		}
		log.Println("Error hashing password:", err)
		c.String(http.StatusInternalServerError, "registration failed: error hashing password")
		return
	}

	_, err = database.DB.Exec("INSERT INTO users (username, password) VALUES ($1, $2)", username, hashedPassword)
	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == "23505" {
			c.String(http.StatusBadRequest, "Пользователь с таким именем уже существует")
			return
		}
		log.Println("Error inserting user into database:", err)
		c.String(http.StatusInternalServerError, "registration failed: database error")
		return
	}
	c.String(http.StatusOK, "User successfully created")
}
