package tasks

import (
	"github.com/gin-gonic/gin"
	"github.com/patrickmn/go-cache"
	"github.com/shaheerhas/todo-list/app/auth"
	"time"
)

func Route(router *gin.Engine, svc TaskApp, authApp auth.AuthApp) {

	authorized := router.Group("/").Use(authApp.AuthMiddleware(&gin.Context{}))

	authorized.GET("/tasks", svc.getTasksList)
	authorized.POST("/tasks", svc.postTask)
	authorized.PATCH("/tasks/:taskid", svc.patchTask)
	authorized.DELETE("/tasks/:taskid", svc.deleteTask)

	cachey = cache.New(15*time.Minute, 30*time.Minute)
	cached := authorized.Use(CacheCheck(&gin.Context{}))

	cached.GET("/tasks/getTaskCounts", svc.getTaskCounts)
	cached.GET("/tasks/getTaskAverages", svc.getTaskAverages)
	cached.GET("/tasks/getOverDueTask", svc.getOverDueTask)
	cached.GET("/tasks/getMaxTaskCompletedDay", svc.getMaxTaskCompletedDay)
	cached.GET("/tasks/getOpenedTasksPerDay", svc.getOpenedTasksPerDay)

	cached.GET("/tasks/similar", svc.similarTasks)

	authorized.POST("/attachment/:taskid", svc.attachFile)
	authorized.GET("/attachment/:taskid", svc.downloadFile)
	authorized.DELETE("/attachment/:taskid", svc.deleteFile)

}
