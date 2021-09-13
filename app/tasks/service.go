package tasks

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
)

const ATTACHMENTFOLDER = "downloads"

func (svc TaskApp) getTasksList(c *gin.Context) {
	// should add functionality which gets id from context
	userId := 1
	tasks, err := allTasks(svc, userId)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, err)
		log.Println(err)
		return
	}
	c.IndentedJSON(http.StatusOK, tasks)
}

func (svc TaskApp) patchTask(c *gin.Context) {
	//var task Task
	var reqBody map[string]interface{}
	if err := c.BindJSON(&reqBody); err != nil {
		log.Println(err)
		c.IndentedJSON(http.StatusBadRequest, "json format not correct")
		return
	}
	if err := updateTask(svc, reqBody); err != nil {
		log.Println(err)
		c.IndentedJSON(http.StatusInternalServerError, err)
		return
	}
	c.IndentedJSON(http.StatusOK, "task successfully updated")
}

func (svc TaskApp) postTask(c *gin.Context) {
	var task Task
	if err := c.BindJSON(&task); err != nil {
		log.Println(err)
		c.IndentedJSON(http.StatusBadRequest, "json format not correct")
		return
	}
	if err := createTask(svc, task); err != nil {
		log.Println(err)
		c.IndentedJSON(http.StatusBadRequest, "couldn't create record in db")
		return
	}
	c.IndentedJSON(http.StatusCreated, "record created successfully")
}

func (svc TaskApp) getTaskById(c *gin.Context) {

}

func (svc TaskApp) attachFile(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		log.Println(err)
		c.IndentedJSON(http.StatusBadRequest, "error in file format")
		return
	}
	log.Println(file.Filename)

	taskId, err := strconv.Atoi(c.Param("taskid"))
	if err != nil {
		log.Println(err)
		c.IndentedJSON(http.StatusBadRequest, "taskid should be numeric")
		return
	}

	if _, err := getTaskById(svc, taskId); err != nil {
		log.Println(err)
		c.IndentedJSON(http.StatusNotFound, "task with this id not found")
		return
	}
	// add User's ID here from context
	userId := 1
	fileName := fmt.Sprintf("%s/%v_%s", ATTACHMENTFOLDER, userId, file.Filename)
	err = c.SaveUploadedFile(file, fileName)
	if err != nil {
		log.Println(err)
		c.IndentedJSON(http.StatusInternalServerError, "error saving file")
		return
	}
	//save file path to db
	err = addFilePath(svc, fileName, taskId)
	if err != nil {
		log.Println(err)
		c.IndentedJSON(http.StatusInternalServerError, "some issue with adding filename to db")
		return
	}
	c.IndentedJSON(http.StatusOK, fmt.Sprintf("'%s' uploaded!", file.Filename))
}

func (svc TaskApp) downloadFile(c *gin.Context) {
	taskId, err := strconv.Atoi(c.Param("taskid"))
	if err != nil {
		log.Println(err)
		c.IndentedJSON(http.StatusBadRequest, "taskid should be numeric")
		return
	}
	filePath, err := getFilePath(svc, taskId)
	if err != nil {
		log.Println(err)
		c.IndentedJSON(http.StatusNotFound, "file not found in db")
		return
	}
	file, err := os.Open(filePath) //Create a file
	if err != nil {
		log.Println(err)
		c.IndentedJSON(http.StatusInternalServerError, "file not found in server")
		return
	}
	defer file.Close()

	c.Writer.Header().Add("Content-type", "application/octet-stream")
	_, err = io.Copy(c.Writer, file)
	if err != nil {
		log.Println(err)
		c.IndentedJSON(http.StatusNotFound, "file cannot be copied")
		return
	}

}

func (svc TaskApp) deleteTask(c *gin.Context) {

	taskId, err := strconv.Atoi(c.Param("taskid"))
	if err != nil {
		log.Println(err)
		c.IndentedJSON(http.StatusBadRequest, "taskid should be numeric")
		return
	}
	if err := deleteTask(svc, taskId); err != nil {
		log.Println(err)
		c.IndentedJSON(http.StatusNotFound, "task with this id not found in db")
		return
	}
	c.IndentedJSON(http.StatusOK, "record deleted successfully")

}

func (svc TaskApp) deleteFile(c *gin.Context) {
	taskId, err := strconv.Atoi(c.Param("taskid"))
	if err != nil {
		log.Println(err)
		c.IndentedJSON(http.StatusBadRequest, "taskid should be numeric")
		return
	}
	filePath, err := getFilePath(svc, taskId)
	if err != nil {
		log.Println(err)
		c.IndentedJSON(http.StatusNotFound, "file not found in db")
		return
	}
	err = os.Remove(filePath)
	if err != nil {
		log.Println(err)
		c.IndentedJSON(http.StatusInternalServerError, "couldn't delete file from file system")
		return
	}
	err = deleteFilePath(svc, taskId)
	if err != nil {
		log.Println(err)
		c.IndentedJSON(http.StatusInternalServerError, "couldn't delete file path from db")
		return
	}
	c.IndentedJSON(http.StatusOK, "file successfully deleted")

}
