package users

import (
	"github.com/gin-gonic/gin"
	"github.com/shaheerhas/todo-list/app/auth"
)

func Route(router *gin.Engine, svc UserModelApp, authApp auth.AuthApp) {
	router.POST("/signup", svc.signup)
	router.POST("/login", svc.login)

	router.GET("/fblogin", svc.fbLogin)
	router.GET("/fbcallback", svc.fbCallBack)

	router.POST("/confirm/:emailToken", svc.confirmUser)

	router.POST("/forgotpassword", svc.forgotPassword)
	router.POST("/resetpassword/:emailToken", svc.resetPassword)

	authorized := router.Group("/").Use(authApp.AuthMiddleware(&gin.Context{}))
	authorized.POST("/logout", svc.logout)

}
