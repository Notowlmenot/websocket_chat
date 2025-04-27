package auth

import (
	"chat/internal/database"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// func UserExistCheck(username) {
// 	err := database.Exec("SELECT from users (username) with VALUES ($1)", username)
// 	if(){
// 		return false
// 	}else{return true}

// }

func RegisterHandler(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")

	if username == "" || password == "" {
		c.String(http.StatusBadRequest, "registration failed: missing username or password")
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Println("Error hashing password:", err)
		c.String(http.StatusInternalServerError, "registration failed: error hashing password")
		return
	}

	_, err = database.DB.Exec("INSERT INTO users (username, password) VALUES ($1, $2)", username, hashedPassword)
	if err != nil {
		log.Println("Error inserting user into database:", err)
		c.String(http.StatusInternalServerError, "registration failed: database error")
		return
	}

	c.String(http.StatusOK, "User successfully created")
}
