package tasks

import "github.com/gin-gonic/gin"

func Route(svc TaskApp) {
	router := gin.Default()
	router.GET("/tasks", svc.getTasksList)
	router.GET("/tasks/id", svc.getTaskById)
	router.PATCH("/tasks", svc.patchTask)
	router.POST("/tasks", svc.postTask)
	router.DELETE("/tasks/id", svc.deleteTask)

}
