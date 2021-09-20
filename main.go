package main

import (
	"fmt"
	"github.com/shaheerhas/todo-list/app/auth"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
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

// create a struct for db
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
func start() {
	db, err := setupDb()
	if err != nil {
		panic("couldn't connect to db")
	}

	router := gin.Default()

	authApp := auth.AuthApp{Db: db}
	authApp.InitBlackListModel()

	taskApp := tasks.TaskApp{Db: db}
	tasks.Route(router, taskApp, authApp)
	taskApp.InitTaskDb()

	userApp := users.UserModelApp{Db: db}
	users.Route(router, userApp, authApp)
	userApp.InitUserModelDB()

	err = router.Run()
	if err != nil {
		log.Println(err)
		return
	}
}
func main() {
	start()
}
