package main

import (
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/shaheerhas/todo-list/app/tasks"
	"github.com/shaheerhas/todo-list/app/users"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var dbName = os.Getenv("DB_NAME")
var password = os.Getenv("DB_PASSWORD")

func setupDb() (*gorm.DB, error) {
	//	dsn := "host=localhost user=postgres password=tiger123 dbname=todo-list port=5432 sslmode=disable TimeZone=Asia/Shanghai"
	dsn := fmt.Sprintf("host=localhost user=postgres password=%s dbname=%s port=%v sslmode=disable TimeZone=Asia/Shanghai", password, dbName, 5432)
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

	router := gin.Default()

	taskApp := tasks.TaskApp{Db: db}
	tasks.Route(router, taskApp)
	taskApp.InitTaskDb()

	userApp := users.UserModelApp{Db: db}
	users.Route(router, userApp)
	userApp.InitUserModelDB()

	router.Run()
}
func main() {
	start()
}
