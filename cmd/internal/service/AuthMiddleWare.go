package service

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Токен отсутствует"})
			c.Abort()
			return
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(jwtSecret), nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Неверный токен"})
			c.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Ошибка в токене"})
			c.Abort()
			return
		}

		userID, ok := claims["user_id"].(float64) // jwt хранит числа как float64
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Некорректный user_id"})
			c.Abort()
			return
		}

		c.Set("user_id", int(userID))
		c.Set("name", claims["name"])

		c.Next()
	}
}
