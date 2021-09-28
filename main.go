package main

import (
	"fmt"
	"io"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-co-op/gocron"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/shaheerhas/todo-list/app/auth"
	"github.com/shaheerhas/todo-list/app/tasks"
	"github.com/shaheerhas/todo-list/app/users"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func init() {
	err := godotenv.Load(".env")

	if err != nil {
		log.Println(err)
	}

}

func setupDb() (*gorm.DB, error) {
	var dbName = os.Getenv("DB_NAME")
	var password = os.Getenv("DB_PASSWORD")
	var port = os.Getenv("DB_PORT")
	dsn := fmt.Sprintf("host=localhost user=postgres password=%s dbname=%s port=%v sslmode=disable TimeZone=Asia/Shanghai", password, dbName, port)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return db, nil
}

var scheduler *gocron.Scheduler

func scheduleEmail(userApp users.UserModelApp) {
	scheduler = gocron.NewScheduler(time.Now().Location())
	err, _ := scheduler.Every(1).Day().At("00:00").Do(userApp.SendReminderEmails)
	if err != nil {
		log.Println(err)
	}
	scheduler.StartAsync()
}

func start() {
	db, err := setupDb()
	if err != nil {
		panic("couldn't connect to db")
	}

	router := gin.Default()

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
	gin.SetMode(gin.ReleaseMode)
	myFile, _ := os.Create("requestLogs.log")
	gin.DefaultWriter = io.MultiWriter(os.Stdout, myFile)
	start()
}
