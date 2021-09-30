package auth

import (
	"github.com/gin-gonic/gin"
	"github.com/shaheerhas/todo-list/app/utils"
	"log"
	"net/http"
)

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
