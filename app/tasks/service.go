package tasks

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/patrickmn/go-cache"
	"github.com/shaheerhas/todo-list/app/utils"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

func init() {
	cachey = cache.New(15*time.Minute, 20*time.Minute)
}

func isSimilar(task1, task2 string) bool {
	task2Splitted := strings.Split(task2, " ")
	for i, s1 := range task2Splitted {
		if !strings.Contains(task1, s1) {
			return false
		} else {
			// to cater duplicate words
			end := len(task2Splitted) - 1
			if end >= 0 && i < len(task2Splitted) {
				task2Splitted[end], task2Splitted[i] = task2Splitted[i], task2Splitted[end]
				task2Splitted = task2Splitted[:end]
			}
		}
	}
	return true
}

func findSimilarTasks(docs []Task) [][]Task {
	var similar [][]Task
	for i := 0; i < len(docs); i++ {
		task1 := docs[i].Title + " " + docs[i].Details
		for j := i + 1; j < len(docs); j++ {
			if i == j {
				continue
			}
			task2 := docs[j].Title + " " + docs[j].Details
			if len(task1) < len(task2) {
				if isSimilar(task2, task1) {
					var taskSlice []Task
					t1 := docs[i]
					t2 := docs[j]
					taskSlice = append(taskSlice, t1)
					taskSlice = append(taskSlice, t2)
					similar = append(similar, taskSlice)
				}
			}
		}
	}
	return similar
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

func (svc TaskApp) getTasksList(c *gin.Context) {
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
	if err := taskExists(svc, utils.ConvertInterfaceToUint(reqBody["user_id"]), taskId); err != nil {
		log.Println(err)
		msg := "this user doesn't have this task"
		c.JSON(utils.Response(http.StatusNotFound, msg))
		return
	}

	if err := UpdateTask(svc, reqBody); err != nil {
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
	userId = uId
	task.UserID = uId
	if _, err := CreateTask(svc, task); err != nil {
		log.Println(err)
		msg := "couldn't create record in db"
		c.JSON(utils.Response(http.StatusBadRequest, msg))
		return
	}
	msg := "record created successfully"
	c.JSON(utils.Response(http.StatusCreated, msg))
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
	// User's ID here from context
	userId, err := getId(c)
	if err != nil {
		log.Println(err)
	}
	if _, err := getTaskById(svc, taskId, int(userId)); err != nil {
		log.Println(err)
		msg := "task with this id not found"
		c.JSON(utils.Response(http.StatusNotFound, msg))
		return
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

	c.Writer.Header().Add("Content-type", "application/zip")
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
	if err := taskExists(svc, userId, uint(taskId)); err != nil {
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

func checkCache(c *gin.Context) bool {
	userID := strconv.Itoa(int(userId))
	val, ok := cachey.Get(userID)
	if ok && val != nil {
		userCache, ok1 := val.(UserCache)
		if ok1 && userCache.requestStringMap != nil {
			byteData, ok2 := userCache.requestStringMap[c.Request.RequestURI]
			if ok2 && byteData != nil {
				log.Println("cached exists")
				c.JSON(http.StatusOK, userCache.requestStringMap[c.Request.RequestURI])
				return true
			} else {
				return false
			}
		}
	}
	return false
}

func insertToCache(c *gin.Context, userID string, data interface{}) {
	cached, exists := cachey.Get(userID)
	if exists {
		cachedCon, ok := cached.(UserCache)
		if ok && cachedCon.requestStringMap != nil {
			cachedCon.requestStringMap[c.Request.RequestURI] = data
		}
	} else {
		userCache := UserCache{
			requestStringMap: map[string]interface{}{c.Request.RequestURI: data},
		}
		cachey.Set(userID, userCache, cache.DefaultExpiration)
	}
}

func (svc TaskApp) getTaskCounts(c *gin.Context) {
	userID := strconv.Itoa(int(userId))
	if checkCache(c) {
		return
	}
	counts, err := getTasksCount(svc, userId)
	if err != nil {
		log.Println(err)
		msg := "couldn't generate report"
		c.JSON(utils.Response(http.StatusInternalServerError, msg))
	}

	c.JSON(http.StatusOK, counts)
	insertToCache(c, userID, counts)

}

func (svc TaskApp) getTaskAverages(c *gin.Context) {
	userID := strconv.Itoa(int(userId))
	if checkCache(c) {
		return
	}
	userId, err := getId(c)
	if err != nil {
		log.Println(err)
	}
	avgTasks, err := getTasksAverage(svc, userId)
	if err != nil {
		log.Println(err)
		msg := "couldn't generate report"
		c.JSON(utils.Response(http.StatusInternalServerError, msg))
		return
	}

	c.JSON(http.StatusOK, avgTasks)
	insertToCache(c, userID, avgTasks)
}

func (svc TaskApp) getOverDueTask(c *gin.Context) {
	userID := strconv.Itoa(int(userId))
	if checkCache(c) {
		return
	}
	userId, err := getId(c)
	if err != nil {
		log.Println(err)
	}
	countOverDueTasks, err := getOverDueTasks(svc, userId)
	if err != nil {
		log.Println(err)
		msg := "couldn't generate report"
		c.JSON(utils.Response(http.StatusInternalServerError, msg))
		return
	}

	c.JSON(http.StatusOK, map[string]int64{"Over due Tasks:": countOverDueTasks})
	insertToCache(c, userID, map[string]int64{"Over due Tasks:": countOverDueTasks})
}

func (svc TaskApp) getMaxTaskCompletedDay(c *gin.Context) {
	userID := strconv.Itoa(int(userId))
	if checkCache(c) {
		return
	}
	userId, err := getId(c)
	if err != nil {
		log.Println(err)
	}
	mostCompletedTasksDay, err := getMaxTasksCompletedDay(svc, userId)
	if err != nil {
		log.Println(err)
		msg := "couldn't generate report"
		c.JSON(utils.Response(http.StatusInternalServerError, msg))
		return
	}
	res := map[string]string{"Max Tasks Completed Date:": mostCompletedTasksDay}
	c.JSON(http.StatusOK, res)
	insertToCache(c, userID, res)
}

func (svc TaskApp) getOpenedTasksPerDay(c *gin.Context) {
	userID := strconv.Itoa(int(userId))
	if checkCache(c) {
		return
	}
	userId, err := getId(c)
	if err != nil {
		log.Println(err)
	}
	openedTasksPerDay, err := getOpenedTaskPerDay(svc, userId)
	if err != nil {
		log.Println(err)
		msg := "couldn't generate report"
		c.JSON(utils.Response(http.StatusInternalServerError, msg))
		return
	}

	c.JSON(http.StatusOK, openedTasksPerDay)
	insertToCache(c, userID, openedTasksPerDay)
}

func (svc TaskApp) similarTasks(c *gin.Context) {
	userID := strconv.Itoa(int(userId))
	if checkCache(c) {
		return
	}
	userId, _ := getId(c)
	tasks, _ := allTasks(svc, userId)
	similar := findSimilarTasks(tasks)
	if len(similar) == 0 {
		msg := "no similar tasks found"
		c.JSON(utils.Response(http.StatusOK, msg))
		return
	}
	c.JSON(http.StatusOK, similar)
	insertToCache(c, userID, similar)
}
