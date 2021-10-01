package users

import (
	"github.com/joho/godotenv"
	"github.com/shaheerhas/todo-list/app/utils"
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
)

func init() {
	err := godotenv.Load("/home/shaheer/workspace/golang/todo-list/.env")
	if err != nil {
		log.Println(err)
	}
}

func TestShouldCreateUser(t *testing.T) {

	var (
		fname    = "Test1"
		lname    = "Test2"
		email    = "test@gmail.com"
		password = "12345"
	)
	db, err := utils.SetupDb()
	taskApp := UserModelApp{Db: db}
	user := UserModel{FirstName: fname, LastName: lname, Email: email, Password: password}
	res, err := createUser(taskApp, user)
	if err != nil {
		return
	}

	assert.NoError(t, err)
	assert.Equal(t, user.FirstName, res.FirstName)
	assert.Equal(t, user.LastName, res.LastName)
	assert.Equal(t, user.Email, res.Email)
	assert.Equal(t, user.Password, res.Password)

}

func TestShouldGetUser(t *testing.T) {

	db, err := utils.SetupDb()
	taskApp := UserModelApp{Db: db}

	email := "test@gmail.com"
	res, err := getUser(taskApp, email)
	if err != nil {
		return
	}

	assert.NoError(t, err)
	assert.Equal(t, res.Email, email)
}
