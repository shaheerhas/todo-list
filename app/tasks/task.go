package tasks

import "time"

type Task struct {
	time    time.Time
	title   string
	id      int
	details string
	files   []string
}
