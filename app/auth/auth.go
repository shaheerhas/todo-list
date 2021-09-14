package auth

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

var secretKey = os.Getenv("SECRET_JWT_KEY")

func GenerateJWTToken(id uint, email string) (string, error) {
	authClaims := jwt.MapClaims{}
	authClaims["email"] = email
	authClaims["id"] = id
	authClaims["expiry"] = time.Now().Add(time.Minute * 30).Unix()
	auth := jwt.NewWithClaims(jwt.SigningMethodHS256, authClaims)
	token, err := auth.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}
	return token, nil
}

func ValidateJWTToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return secretKey, nil
	})
	if err != nil {
		return nil, fmt.Errorf("unauthorized")
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		return claims, nil
	} else {
		return nil, err
	}

}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		const BEARER_SCHEMA = "Bearer"
		authHeader := c.GetHeader("Authorization")
		tokenString := authHeader[len(BEARER_SCHEMA):]
		claims, err := ValidateJWTToken(tokenString)
		if err != nil {
			c.IndentedJSON(http.StatusUnauthorized, "token not valid")
			return
		} else {
			c.Set("userId", claims["id"])
		}
	}
}
