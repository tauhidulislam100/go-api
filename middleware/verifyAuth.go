package middleware

import (
	"fmt"
	"gin-api/configs"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

func VerifyAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		reqToken := c.GetHeader("Authorization")
		tokenString := strings.Split(reqToken, " ")[1]
		fmt.Printf("Your token is %v \n", strings.Split(reqToken, " ")[1])

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Don't forget to validate the alg is what you expect:

			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}

			return []byte(configs.EnvSecretKey()), nil
		})

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			fmt.Println("error", claims)
			c.Next()
		} else {
			fmt.Println("error", err)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status": http.StatusUnauthorized, "message": "error", "data": err.Error()})
			return
		}
	}
}
