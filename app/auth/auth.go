package auth

import (
	"fmt"
	"log"
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
	authClaims["expiry"] = time.Now().Add(time.Hour * 24).Unix()
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
		return []byte(secretKey), nil
	})
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(jwt.MapClaims)

	expiryTime_, exists := claims["expiry"].(float64)
	expiryTime := time.Unix(int64(expiryTime_), 0)

	if exists && expiryTime.Before(time.Now()) {
		return nil, fmt.Errorf("token expired")
	}
	if ok && token.Valid {
		return claims, nil
	} else {
		return nil, fmt.Errorf("token not valid")
	}

}

func AuthMiddleware(c *gin.Context) gin.HandlerFunc {
	return func(c *gin.Context) {
		const BEARER_SCHEMA = "Bearer "
		authHeader := c.GetHeader("Authorization")

		tokenString := authHeader[len(BEARER_SCHEMA):]
		fmt.Println(tokenString)
		claims, err := ValidateJWTToken(tokenString)
		if err != nil {
			log.Println("Authmiddleware error:", err)
			c.IndentedJSON(http.StatusUnauthorized, "token not valid")
			c.Abort()
			return
		} else {
			c.Set("userId", claims["id"])
		}
		c.Next()
	}
}
