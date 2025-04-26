package auth

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func AuthHandler(c *gin.Context) {

	/* Принцип авторизации:
	Отправка get-запроса с входными данными login/pasword
	Поиск в базе данных пользователя
	сравнение хешированного пароля (через bcrypt?) с введеным.
	В случае удачной авторизации - выдача пользователю JWT (сохранение его в куки)
	*/

	//Максимально тупая авторизация без хеширования
	password := "idk"
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal(err)
	}

	err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
	if err != nil {
		fmt.Println("Пароль не верен")
	} else {
		fmt.Println("Пароль верен")
	}
	//Использовать JWT токен для авторизации
	c.JSON(http.StatusOK, gin.H{"token": "JWT_TOKEN"})
}
