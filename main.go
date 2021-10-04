package main

import (
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/shaheerhas/todo-list/app"
	"github.com/shaheerhas/todo-list/app/auth"
	"github.com/shaheerhas/todo-list/app/tasks"
	"github.com/shaheerhas/todo-list/app/users"
	"github.com/shaheerhas/todo-list/app/utils"
	"log"
)

func init() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Println(err)
	}

}

func start() {
	db, err := utils.SetupDb()
	if err != nil {
		panic("couldn't connect to db")
	}

	router := gin.Default()

	router.Use(app.LoggerToFile())

	authApp := auth.AuthApp{Db: db}
	authApp.InitBlackListModel()

	userApp := users.UserModelApp{Db: db}
	users.Route(router, userApp, authApp)
	userApp.InitUserModelDB()
	users.ScheduleEmail(userApp)

	if err != nil {
		log.Println(err)
	}

	taskApp := tasks.TaskApp{Db: db}
	tasks.Route(router, taskApp, authApp)
	taskApp.InitTaskDb()

	err = router.Run()
	if err != nil {
		log.Println(err)
		return
	}

}

func main() {
	start()
}
