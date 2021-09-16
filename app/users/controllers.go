package users

import "github.com/gin-gonic/gin"

func Route(router *gin.Engine, svc UserModelApp) {
	router.GET("/users", svc.getUsers)
	router.POST("/users", svc.signup)
	router.POST("/login", svc.login)
	router.POST("/confirm/:emailToken", svc.confirmUser)
}
