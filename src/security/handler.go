package security

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func ValidateJWT(c *gin.Context) {
	tokenString := c.GetHeader("Authorization")
	jwtKey := []byte(os.Getenv("SECRET"))
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return jwtKey, nil
	})
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		c.Set("user_id", claims["user_id"])
	} else {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
	}
}

func IdFromJWT(c *gin.Context) (uint, error) {
	tokenString := c.GetHeader("Authorization")
	jwtKey := []byte(os.Getenv("SECRET"))
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return jwtKey, nil
	})
	var id uint
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		c.Set("user_id", claims["user_id"])
		id = uint(claims["user_id"].(float64))
	} else {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
	}
	return id, err
}
