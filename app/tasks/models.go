package tasks

import (
	"time"

	"gorm.io/gorm"
)

type Task struct {
	gorm.Model
	Title          string `gorm:"not null"`
	Details        string
	CompletionTime time.Time
	DueTime        time.Time
	Status         bool     `gorm:"not null"`
	Files          []string `gorm:"type:text[]"`
	UserID         uint     `gorm:"not null"`
}
type TaskApp struct {
	Db *gorm.DB
}

func (t *TaskApp) InitTaskDb() {
	t.Db.AutoMigrate(&Task{})
}
