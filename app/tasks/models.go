package tasks

import (
	"time"

	"gorm.io/gorm"
)

type Task struct {
	ID             uint   `gorm:"autoIncrement;primaryKey"`
	Title          string `gorm:"not null"`
	Details        string
	CompletionTime time.Time
	DueTime        time.Time
	Status         bool `gorm:"not null; default:false"`
	File           string
	UserID         uint `gorm:"not null"`
}
type TaskApp struct {
	Db *gorm.DB
}

func (t *TaskApp) InitTaskDb() {
	t.Db.AutoMigrate(&Task{})
}
