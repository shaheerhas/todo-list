package main

import (
	"github.com/shaheerhas/todo-list/app/tasks"
	"github.com/shaheerhas/todo-list/app/user"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func setupDb() (*gorm.DB, error) {
	dsn := "host=localhost user=postgres password=tiger123 dbname=dvdrental port=5432 sslmode=disable TimeZone=Asia/Shanghai"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return db, nil
}
func start() {
	db, err := setupDb()
	if err != nil {
		panic("couldn't connect to db")
	}

	taskApp := tasks.TaskApp{db}
	tasks.Route(taskApp)
	userApp := user.UserApp{db}

	//userApp := user.UserApp{db}

}
func main() {

}
