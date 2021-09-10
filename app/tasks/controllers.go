package tasks

import "github.com/gin-gonic/gin"

func Route(router *gin.Engine, svc TaskApp) {

	router.GET("/tasks", svc.getTasksList)
	router.PATCH("/tasks", svc.patchTask)
	router.POST("/tasks", svc.postTask)
	router.DELETE("/tasks/:taskid", svc.deleteTask)
	router.POST("/attachment/:taskid", svc.attachFile)
	router.GET("/attachment/:taskid", svc.downloadFile)
	router.DELETE("attachment/:taskid", svc.deleteFile)

}
