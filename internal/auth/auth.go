package auth

import (
	"chat/internal/database"
	JWT "chat/internal/jwt"
	"database/sql"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func AuthHandler(c *gin.Context) {
	username, password := c.PostForm("username"), c.PostForm("password")
	var hashedPassword []byte
	var userid int

	err := database.DB.QueryRow("SELECT id, password FROM users WHERE username = $1", username).Scan(&userid, &hashedPassword)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.String(http.StatusUnauthorized, "invalid username or password")
			return
		}
		log.Println("Error querying user from database:", err)
		c.String(http.StatusInternalServerError, "login failed: error querying user")
		return
	}
	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		c.String(http.StatusUnauthorized, "invalid username or password")
		return
	}

	accessToken, err := JWT.GenerateAcessToken(userid)
	if err != nil {
		log.Println("Error generating access token:", err)
		c.String(http.StatusInternalServerError, "login failed: error generating access token")
		return
	}

	c.SetCookie(
		"access_token",
		accessToken,
		int(time.Minute*15),
		"/",
		"",
		false, // Secure (false для разработки, true в production - только HTTPS)
		true,
	)

	c.String(http.StatusOK, "Login successful")
}
