package tasks

import (
	"github.com/golang/mock/gomock"
	"github.com/shaheerhas/todo-list/app/utils"
	"testing"
	"time"
)

func TestShouldCreateTask(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockTaskRepo := NewMockTaskRepo(mockCtrl)
	var (
		title   = "task1"
		details = "detail1"
		dueTime = time.Now().Add(time.Minute * 120)
		status  = false
		userID  = 1
	)
	mockTaskCreator := Task{Title: title, Details: details, DueTime: dueTime, Status: status,
		UserID: uint(userID), TaskRepo: mockTaskRepo}
	db, _ := utils.SetupDb()
	taskApp := TaskApp{db}
	mockTaskRepo.EXPECT().CreateTask(taskApp, mockTaskCreator).Return(mockTaskCreator, nil).Times(1)
	mockTaskCreator.CreateTask(taskApp, mockTaskCreator)

}
