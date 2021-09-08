package tasks

import (
	"time"

	"gorm.io/gorm"
)

type Task struct {
	gorm.Model
	Title          string `gorm:"unique; notnull"`
	Details        string
	CompletionTime time.Time
	DueTime        time.Time
	Status         bool
	Files          []string
	UserId         int `gorm:"notnull"`
}
type TaskApp struct {
	Db *gorm.DB
}
