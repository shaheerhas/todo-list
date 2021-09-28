package tasks

import (
	"github.com/gin-gonic/gin"
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
}
type TaskApp struct {
	Db *gorm.DB
}

var cachey *cache.Cache

type TaskCache struct {
	gin.ResponseWriter
	cache         *cache.Cache
	requestString string
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

func (svc *TaskApp) InitTaskDb() {
	err := svc.Db.AutoMigrate(&Task{})
	if err != nil {
		log.Println(err)
		return
	}
}
