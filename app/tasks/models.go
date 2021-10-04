package tasks

import (
	"github.com/patrickmn/go-cache"
	"gorm.io/gorm"
	"log"
	"time"
)

type Task struct {
	ID             uint   `gorm:"autoIncrement;primaryKey"`
	Title          string `gorm:"not null"`
	Details        string
	CreatedAt      time.Time
	DueTime        time.Time
	CompletionTime time.Time
	Status         bool `gorm:"not null; default:false"`
	File           string
	UserID         uint `gorm:"not null"`
	TaskRepo
}
type TaskApp struct {
	Db *gorm.DB
}

var cachey *cache.Cache
var userId uint

type UserCache struct {
	//gin.ResponseWriter
	//cache            *cache.Cache
	//requestStrings   []string
	requestStringMap map[string]interface{}
}

type TaskCount struct {
	Total     int64
	Completed int64
	Remaining int64
}

type TaskCompleted struct {
	Avg float64
	Day time.Time
}

type TaskRepo interface {
	CreateTask(svc TaskApp, task Task) (Task, error)
	UpdateTask(svc TaskApp, updatedTask map[string]interface{}) error
}

func (svc *TaskApp) InitTaskDb() {
	err := svc.Db.AutoMigrate(&Task{})
	if err != nil {
		log.Println(err)
		return
	}
}
