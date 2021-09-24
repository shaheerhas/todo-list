package auth

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/shaheerhas/todo-list/app/utils"
)

var secretKey = os.Getenv("SECRET_JWT_KEY")

func GenerateJWTToken(id uint, email string) (string, error) {
	authClaims := jwt.MapClaims{}
	authClaims["email"] = email
	authClaims["id"] = id
	authClaims["exp"] = time.Now().Add(time.Minute * 15).Unix()

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

	if ok && token.Valid {
		return claims, nil
	} else if ve, ok := err.(*jwt.ValidationError); ok {
		if ve.Errors&jwt.ValidationErrorMalformed != 0 {
			return nil, fmt.Errorf("invalid or malformed token")
		} else if ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
			// Token is either expired or not active yet
			return nil, fmt.Errorf("token expired or not active")
		} else {
			return nil, fmt.Errorf("couldn't handle token")
		}
	} else {
		return nil, fmt.Errorf("couldn't handle token")
	}

}

func (authApp AuthApp) AuthMiddleware(c *gin.Context) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString, err := c.Cookie("accessToken")
		if err != nil {
			log.Println("Middleware error:", err)
		}
		claims, err := ValidateJWTToken(tokenString)
		if err != nil {
			log.Println("Middleware error:", err)
			msg := "token not valid"
			c.JSON(utils.Response(http.StatusUnauthorized, msg))
			c.Abort()
			return
		} else if IsBlackListed(tokenString) {
			log.Println("Middleware error:", "token black listed")
			msg := "token not valid"
			c.JSON(utils.Response(http.StatusUnauthorized, msg))
			c.Abort()
			return
		} else {
			c.Set("userId", claims["id"])
		}
		c.Next()
	}
}
