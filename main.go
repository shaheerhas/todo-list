package main

import (
	"github.com/shaheerhas/todo-list/app"
	"github.com/shaheerhas/todo-list/app/utils"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-co-op/gocron"
	"github.com/joho/godotenv"
	"github.com/shaheerhas/todo-list/app/auth"
	"github.com/shaheerhas/todo-list/app/tasks"
	"github.com/shaheerhas/todo-list/app/users"
	"log"
)

func init() {
	err := godotenv.Load(".env")

	if err != nil {
		log.Println(err)
	}

}

var scheduler *gocron.Scheduler

func scheduleEmail(userApp users.UserModelApp) {
	scheduler = gocron.NewScheduler(time.Now().Location())
	//err, _ := scheduler.Every(1).Day().At(time.Now().Add(time.Second * 1)).Do(userApp.SendReminderEmails)
	err, _ := scheduler.Every(1).Day().At("00:00").Do(userApp.SendReminderEmails)
	if err != nil {
		log.Println(err)
	}
	scheduler.StartAsync()
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
	scheduleEmail(userApp)

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
	//gin.SetMode(gin.ReleaseMode)
	//myFile, _ := os.Create("requestLogs.log")
	//gin.DefaultWriter = io.MultiWriter(os.Stdout, myFile)
	start()
}
