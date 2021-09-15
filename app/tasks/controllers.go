package tasks

import (
	"github.com/gin-gonic/gin"
	"github.com/shaheerhas/todo-list/app/auth"
)

type GinCtxFunc interface {
}

func Route(router *gin.Engine, svc TaskApp) {

	//	authorizationReq := router.Group("/")

	//	 authorizationReq.Use(auth.AuthMiddleware())
	//	 {
	router.GET("/tasks", auth.AuthMiddleware(svc.getTasksList))
	router.PATCH("/tasks", auth.AuthMiddleware(svc.patchTask))
	router.POST("/tasks", auth.AuthMiddleware(svc.postTask))
	router.DELETE("/tasks/:taskid", auth.AuthMiddleware(svc.deleteTask))
	router.POST("/attachment/:taskid", auth.AuthMiddleware(svc.attachFile))
	router.GET("/attachment/:taskid", auth.AuthMiddleware(svc.downloadFile))
	router.DELETE("attachment/:taskid", auth.AuthMiddleware(svc.deleteFile))
	//}
}
