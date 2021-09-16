package tasks

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/shaheerhas/todo-list/app/utils"
)

const ATTACHMENTFOLDER = "downloads"

func (svc TaskApp) getTasksList(c *gin.Context) {
	// should add functionality which gets id from context
	userId, exists := getId(c)
	if exists != nil {
		log.Println(exists)
	}
	tasks, err := allTasks(svc, userId)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, err)
		log.Println(err)
		return
	}

	c.IndentedJSON(http.StatusOK, tasks)
}

func getId(c *gin.Context) (uint, error) {
	id, exists := c.Get("userId")

	if exists {
		uid, ok := id.(float64)
		if ok {
			return utils.ConvertInterfaceToUint(uid), nil
		}
	}
	return 0, fmt.Errorf("couldn't parse id")
}

func (svc TaskApp) patchTask(c *gin.Context) {
	var reqBody map[string]interface{}
	if err := c.BindJSON(&reqBody); err != nil {
		log.Println(err)
		c.IndentedJSON(http.StatusBadRequest, "json format not correct")
		return
	}

	var e error
	reqBody["user_id"], e = getId(c)
	if e != nil {
		log.Println(e)
	}

	if err := checkUserTask(svc, utils.ConvertInterfaceToUint(reqBody["user_id"]), utils.ConvertInterfaceToUint(reqBody["id"])); err != nil {
		log.Println(err)
		c.IndentedJSON(http.StatusBadRequest, "this user doesn't have this task")
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
		c.IndentedJSON(http.StatusUnprocessableEntity, "json format not correct")
		return
	}

	uId, err := getId(c)
	if err != nil {
		log.Println(err)
	}
	task.UserID = utils.ConvertInterfaceToUint(uId)

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
	// User's ID here from context
	userId, err := getId(c)
	if err != nil {
		log.Println(err)
	}
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

	userId, err := getId(c)
	if err != nil {
		log.Println(err)
	}
	if err := checkUserTask(svc, userId, uint(taskId)); err != nil {
		log.Println(err)
		c.IndentedJSON(http.StatusBadRequest, "this user doesn't have this task")
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
