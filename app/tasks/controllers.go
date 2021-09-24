package tasks

import (
	"github.com/gin-gonic/gin"
	"github.com/shaheerhas/todo-list/app/auth"
)

func Route(router *gin.Engine, svc TaskApp, authApp auth.AuthApp) {

	authorized := router.Group("/").Use(authApp.AuthMiddleware(&gin.Context{}))

	authorized.GET("/tasks", svc.getTasksList)
	authorized.POST("/tasks", svc.postTask)
	authorized.PATCH("/tasks/:taskid", svc.patchTask)
	authorized.DELETE("/tasks/:taskid", svc.deleteTask)

	authorized.GET("/tasks/getTaskCounts", svc.getTaskCounts)
	authorized.GET("/tasks/getTaskAverages", svc.getTaskAverages)
	authorized.GET("/tasks/getOverDueTask", svc.getOverDueTask)
	authorized.GET("/tasks/getMaxTaskCompletedDay", svc.getMaxTaskCompletedDay)
	authorized.GET("/tasks/getOpenedTasksPerDay", svc.getOpenedTasksPerDay)

	authorized.GET("/tasks/similar", svc.similarTasks)

	authorized.POST("/attachment/:taskid", svc.attachFile)
	authorized.GET("/attachment/:taskid", svc.downloadFile)
	authorized.DELETE("/attachment/:taskid", svc.deleteFile)

}
