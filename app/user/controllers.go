package user

import "github.com/gin-gonic/gin"

func Route(router *gin.Engine, svc UserApp) {
	router.GET("/users", svc.getUsers)
	router.POST("/users", svc.postUser)
}
