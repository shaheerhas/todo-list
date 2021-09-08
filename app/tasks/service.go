package tasks

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (svc TaskApp) getTasksList(c *gin.Context) {
	tasks, err := allTasks(svc)
	if err != nil {

		c.IndentedJSON(http.StatusNoContent, err)
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
	}
}

func (svc TaskApp) getTaskById(c *gin.Context) {

}

func (svc TaskApp) attachFile(c *gin.Context) {

}

func (svc TaskApp) downloadFile(c *gin.Context) {

}

func (svc TaskApp) deleteTask(c *gin.Context) {

}
