package auth

import (
	"chat/internal/database"
	JWT "chat/internal/jwt"
	"errors"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

func RegisterHandler(c *gin.Context) {
	var userid int
	username, password := c.PostForm("username"), c.PostForm("password")
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
	err = database.DB.QueryRow("INSERT INTO users (username, password) VALUES ($1, $2) returning id", username, hashedPassword).Scan(&userid)
	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == "23505" {
			c.String(http.StatusBadRequest, "Пользователь с таким именем уже существует")
			return
		}
		log.Println("Error inserting user into database:", err)
		c.String(http.StatusInternalServerError, "registration failed: error inserting user")
		return
	}

	// Сгенерировать refresh token
	_, err = JWT.GenerateRefreshToken(userid)
	if err != nil {
		log.Println("Error generating refresh token:", err)
		c.String(http.StatusInternalServerError, "registration failed: error generating refresh token")
		return
	}

	// Сгенерировать access token
	accessToken, err := JWT.GenerateAcessToken(userid)
	if err != nil {
		log.Println("Error generating access token:", err)
		c.String(http.StatusInternalServerError, "registration failed: error generating access token")
		return
	}

	// Установить access token в HTTP-only cookie
	c.SetCookie(
		"access_token",
		accessToken,
		int(time.Minute*15),
		"/",
		"",
		false, // Secure (false для разработки, true в production - только HTTPS)
		true,
	)

	c.String(http.StatusOK, "User registered successfully")
}
