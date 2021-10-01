package tasks

import (
	"github.com/joho/godotenv"
	"github.com/shaheerhas/todo-list/app/utils"
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
	"time"
)

func init() {
	err := godotenv.Load("/home/shaheer/workspace/golang/todo-list/.env")
	if err != nil {
		log.Println(err)
	}
}

func TestShouldCreateTask(t *testing.T) {

	var (
		title   = "task1"
		details = "detail1"
		dueTime = time.Now().Add(time.Minute * 120)
		status  = false
		userID  = 1
	)
	db, err := utils.SetupDb()
	taskApp := TaskApp{Db: db}
	task := Task{Title: title, Details: details, DueTime: dueTime, Status: status, UserID: uint(userID)}
	res, err := createTask(taskApp, task)
	if err != nil {
		return
	}

	assert.NoError(t, err)
	assert.Equal(t, task.UserID, res.UserID)
	assert.Equal(t, task.Title, res.Title)
	assert.Equal(t, task.Details, res.Details)
	assert.Equal(t, task.Status, res.Status)

}

func TestShouldGetTask(t *testing.T) {

	db, err := utils.SetupDb()
	taskApp := TaskApp{Db: db}

	id := 2
	userID := 1
	res, err := getTaskById(taskApp, id, userID)
	if err != nil {
		return
	}

	assert.NoError(t, err)
	assert.Equal(t, res.UserID, uint(userID))
	assert.Equal(t, res.ID, uint(id))
}
