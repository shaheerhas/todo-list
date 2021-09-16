package tasks

import (
	"github.com/gin-gonic/gin"
	"github.com/shaheerhas/todo-list/app/auth"
)

func Route(router *gin.Engine, svc TaskApp) {

	authorized := router.Group("/").Use(auth.AuthMiddleware(&gin.Context{}))

	authorized.GET("/tasks", svc.getTasksList)
	authorized.PATCH("/tasks", svc.patchTask)
	authorized.POST("/tasks", svc.postTask)
	authorized.DELETE("/tasks/:taskid", svc.deleteTask)
	authorized.POST("/attachment/:taskid", svc.attachFile)
	authorized.GET("/attachment/:taskid", svc.downloadFile)
	authorized.DELETE("attachment/:taskid", svc.deleteFile)

}
