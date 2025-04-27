package JWT

import (
	"chat/internal/database"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var secretKey = []byte(os.Getenv("JWT_SECRET_KEY"))

func GenerateRefreshToken(userID int) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Hour * 24).Unix(), // Срок действия — 24 часа (UNIX timestamp)
	}

	Stringtoken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(secretKey)
	if err != nil {
		return "", fmt.Errorf("ошибка при генерации refresh token: %w", err)
	}

	// Запись токена в БД (нужно передать контекст)
	_, err = database.DB.Exec("INSERT INTO refresh_tokens (userid, token) VALUES ($1, $2)", userID, Stringtoken)
	if err != nil {
		return "", fmt.Errorf("ошибка при записи refresh token в БД: %w", err)
	}

	return Stringtoken, nil // Возвращаем сгенерированный refresh token
}

func GenerateAcessToken(userID int) (string, error) {

	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Minute * 15), // Срок действия — 15 минут
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secretKey)
}

func ValidateAcessToken(userID int, TokenString string) (bool, error) {
	token, err := jwt.Parse(TokenString, func(token *jwt.Token) (interface{}, error) { //проверка ключом hmacSampleSecret
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return secretKey, nil
	})

	if err != nil {
		fmt.Println("Ошибка при разборе токена:", err)
		return false, nil
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if claims["user_id"] == userID {
			return true, nil
		}
	} else {
		fmt.Println("Токен недействителен:", err)
		return false, nil
	}
	return false, nil
}
