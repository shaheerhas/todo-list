package user

import "github.com/gin-gonic/gin"

func Route(svc UserApp) {
	router := gin.Default()
	router.GET("/users")
}
