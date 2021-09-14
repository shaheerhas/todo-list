package tasks

import (
	"github.com/gin-gonic/gin"
	"github.com/shaheerhas/todo-list/app/auth"
)

func Route(router *gin.Engine, svc TaskApp) {

	authorizationReq := router.Group("/")

	authorizationReq.Use(auth.AuthMiddleware())
	{
		router.GET("/tasks", svc.getTasksList)
		router.PATCH("/tasks", svc.patchTask)
		router.POST("/tasks", svc.postTask)
		router.DELETE("/tasks/:taskid", svc.deleteTask)
		router.POST("/attachment/:taskid", svc.attachFile)
		router.GET("/attachment/:taskid", svc.downloadFile)
		router.DELETE("attachment/:taskid", svc.deleteFile)
	}
}
