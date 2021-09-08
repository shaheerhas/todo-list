package tasks

import "github.com/shaheerhas/todo-list/app/tasks"

type User struct {
	id       int
	username string
	email    string
	tasks    []tasks.Task
}
