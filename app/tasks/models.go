package tasks

import (
	"time"

	"gorm.io/gorm"
)

type Task struct {
	id      int
	time    time.Time
	title   string
	details string
	files   []string
}

type TaskApp struct {
	Db *gorm.DB
}
