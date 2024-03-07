package database

import (
	"log"
	"time"
  "os"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateJWT(userid uint) (string, error) {
	jwtKey := []byte(os.Getenv("SECRET"))
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["authorized"] = true
	claims["user_id"] = userid
	claims["exp"] = time.Now().Add(time.Hour * 3).Unix()
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		log.Println("Error while signing the token")
		return "", err
	}
	return tokenString, nil
}
