package tasks

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (svc TaskApp) getTasksList(c *gin.Context) {
	// should add functionality which gets id from context
	tasks, err := allTasks(svc)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, err)
		log.Println(err)
		return
	}
	c.IndentedJSON(http.StatusOK, tasks)
}

func (svc TaskApp) patchTask(c *gin.Context) {

}

func (svc TaskApp) postTask(c *gin.Context) {
	var task Task
	if err := c.BindJSON(&task); err != nil {
		c.IndentedJSON(http.StatusBadRequest, "json format not correct")
		return
	}
	if err := createTask(svc, task); err != nil {
		c.IndentedJSON(http.StatusBadRequest, " couldn't create record in db")
		return
	}
	c.IndentedJSON(http.StatusCreated, " record created successfully")
}

func (svc TaskApp) getTaskById(c *gin.Context) {

}

func (svc TaskApp) attachFile(c *gin.Context) {

}

func (svc TaskApp) downloadFile(c *gin.Context) {

}

func (svc TaskApp) deleteTask(c *gin.Context) {

}
