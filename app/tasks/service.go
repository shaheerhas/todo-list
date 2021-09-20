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

func (svc TaskApp) getTasksList(c *gin.Context) {
	// should add functionality which gets id from context
	userId, exists := getId(c)
	if exists != nil {
		log.Println(exists)
	}
	tasks, err := allTasks(svc, userId)
	if err != nil {
		msg := "no tasks for this user"
		c.JSON(utils.Response(http.StatusNotFound, msg))
		log.Println(err)
		return
	}
	c.JSON(http.StatusOK, tasks)
}

func (svc TaskApp) patchTask(c *gin.Context) {
	var reqBody map[string]interface{}
	if err := c.BindJSON(&reqBody); err != nil {
		log.Println(err)
		msg := "invalid or malformed payload"
		c.JSON(utils.Response(http.StatusBadRequest, msg))
		return
	}
	var err error
	reqBody["user_id"], err = getId(c)
	if err != nil {
		log.Println(err)
	}
	taskId := utils.ConvertInterfaceToUint(c.Param("taskid"))
	reqBody["id"] = taskId
	if err := checkUserTask(svc, utils.ConvertInterfaceToUint(reqBody["user_id"]), taskId); err != nil {
		log.Println(err)
		msg := "this user doesn't have this task"
		c.JSON(utils.Response(http.StatusForbidden, msg))
		return
	}

	if err := updateTask(svc, reqBody); err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	msg := "task successfully updated"
	c.JSON(utils.Response(http.StatusOK, msg))
}

func (svc TaskApp) postTask(c *gin.Context) {
	var task Task
	if err := c.BindJSON(&task); err != nil {
		log.Println(err)
		msg := "invalid or malformed payload"
		c.JSON(utils.Response(http.StatusBadRequest, msg))
		return
	}

	uId, err := getId(c)
	if err != nil {
		log.Println(err)
	}
	task.UserID = uId
	if err := createTask(svc, task); err != nil {
		log.Println(err)
		msg := "couldn't create record in db"
		c.JSON(utils.Response(http.StatusBadRequest, msg))
		return
	}
	msg := "record created successfully"
	c.JSON(utils.Response(http.StatusCreated, msg))
}

func (svc TaskApp) getTaskById(c *gin.Context) {

}

func (svc TaskApp) attachFile(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		log.Println(err)
		msg := "invalid or malformed payload"

		c.JSON(utils.Response(http.StatusBadRequest, msg))
		return
	}

	taskId, err := strconv.Atoi(c.Param("taskid"))
	if err != nil {
		log.Println(err)
		msg := "taskid should be numeric"
		c.JSON(utils.Response(http.StatusBadRequest, msg))
		return
	}

	if _, err := getTaskById(svc, taskId); err != nil {
		log.Println(err)
		msg := "task with this id not found"
		c.JSON(utils.Response(http.StatusNotFound, msg))
		return
	}
	// User's ID here from context
	userId, err := getId(c)
	if err != nil {
		log.Println(err)
	}
	var attachmentFolder = os.Getenv("ATTACHMENT_FOLDER")
	fileName := fmt.Sprintf("%s/%v_%s", attachmentFolder, userId, file.Filename)
	log.Println(fileName, "uploaded")
	err = c.SaveUploadedFile(file, fileName)
	if err != nil {
		log.Println(err)
		msg := "error saving file"
		c.JSON(utils.Response(http.StatusInternalServerError, msg))
		return
	}
	//save file path to db
	err = addFilePath(svc, fileName, taskId)
	if err != nil {
		log.Println(err)
		msg := "some issue with adding filename to db"
		c.JSON(utils.Response(http.StatusInternalServerError, msg))
		return
	}
	msg := fmt.Sprintf("'%s' uploaded!", file.Filename)
	c.JSON(utils.Response(http.StatusOK, msg))
}

func (svc TaskApp) downloadFile(c *gin.Context) {
	taskId, err := strconv.Atoi(c.Param("taskid"))
	if err != nil {
		msg := "taskid parameter should be numeric"
		c.JSON(utils.Response(http.StatusBadRequest, msg))
		return
	}
	filePath, err := getFilePath(svc, taskId)
	if err != nil {
		log.Println(err)
		msg := "file not found in db"
		c.JSON(utils.Response(http.StatusNotFound, msg))
		return
	}
	file, err := os.Open(filePath) //Create a file
	if err != nil {
		log.Println(err)
		msg := "file not found in server"
		c.JSON(utils.Response(http.StatusNotFound, msg))
		return
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Println(err)
		}
	}(file)

	c.Writer.Header().Add("Content-type", "application/octet-stream")
	_, err = io.Copy(c.Writer, file)
	if err != nil {
		log.Println(err)
		msg := "couldn't send file"
		c.JSON(utils.Response(http.StatusInternalServerError, msg))
		return
	}
	msg := "sending file"
	c.JSON(utils.Response(http.StatusOK, msg))
}

func (svc TaskApp) deleteTask(c *gin.Context) {

	taskId, err := strconv.Atoi(c.Param("taskid"))
	if err != nil {
		log.Println(err)
		msg := "taskid parameter should be numeric"
		c.JSON(utils.Response(http.StatusBadRequest, msg))
		return
	}

	userId, err := getId(c)
	if err != nil {
		log.Println(err)
	}
	if err := checkUserTask(svc, userId, uint(taskId)); err != nil {
		log.Println(err)
		msg := "this user doesn't have this task"
		c.JSON(utils.Response(http.StatusBadRequest, msg))
		return
	}

	if err := deleteTask(svc, taskId); err != nil {
		log.Println(err)
		msg := "task with this id not found in db"
		c.JSON(utils.Response(http.StatusNotFound, msg))
		return
	}
	msg := "record delete successfully"
	c.JSON(utils.Response(http.StatusOK, msg))

}

func (svc TaskApp) deleteFile(c *gin.Context) {
	taskId, err := strconv.Atoi(c.Param("taskid"))
	if err != nil {
		log.Println(err)
		msg := "taskid should be numeric"
		c.JSON(utils.Response(http.StatusBadRequest, msg))
		return
	}
	filePath, err := getFilePath(svc, taskId)
	if err != nil {
		log.Println(err)
		msg := "file not found in db"
		c.JSON(utils.Response(http.StatusNotFound, msg))
		return
	}
	err = os.Remove(filePath)
	msg := "couldn't delete file"
	if err != nil {
		log.Println(err)
		c.JSON(utils.Response(http.StatusInternalServerError, msg))
		return
	}
	err = deleteFilePath(svc, taskId)
	if err != nil {
		log.Println(err)
		c.JSON(utils.Response(http.StatusInternalServerError, msg))
		return
	}
	msg = "file successfully deleted"
	c.JSON(utils.Response(http.StatusOK, msg))

}
