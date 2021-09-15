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
	expiryTime, _ := claims["expiry"].(time.Time)
	fmt.Println(time.Now(), expiryTime, expiryTime.Before(time.Now()))
	if expiryTime.Before(time.Now()) {
		return nil, fmt.Errorf("token expired")
	}
	if ok && token.Valid {
		return claims, nil
	} else {
		return nil, fmt.Errorf("token not valid")
	}

}

//
func AuthMiddleware(obj func(*gin.Context)) gin.HandlerFunc {
	return func(c *gin.Context) {
		const BEARER_SCHEMA = "Bearer "
		authHeader := c.GetHeader("Authorization")
		tokenString := authHeader[len(BEARER_SCHEMA):]
		fmt.Println(tokenString)
		claims, err := ValidateJWTToken(tokenString)
		if err != nil {
			log.Println("Authmiddleware error:", err)
			c.IndentedJSON(http.StatusUnauthorized, "token not valid")
			log.Println("err", err)
			//	c.Abort()
		} else {
			log.Println("no err", claims["id"])
			c.Set("userId", claims["id"])
			obj(c)
		}
		//c.Next()
	}
}
