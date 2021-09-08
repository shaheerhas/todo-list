package tasks

import (
	"time"

	"gorm.io/gorm"
)

type Task struct {
	gorm.Model
	Title          string
	Details        string
	CompletionTime time.Time
	DueTime        time.Time
	Status         bool
	Files          []string `gorm:"type:text[]"`
	UserID         uint     `gorm:"notnull"`
}
type TaskApp struct {
	Db *gorm.DB
}

func (t *TaskApp) InitTaskDb() {
	t.Db.AutoMigrate(&Task{})
}
