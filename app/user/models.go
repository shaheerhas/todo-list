package user

import (
	"github.com/shaheerhas/todo-list/app/tasks"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username string
	Email    string
	Password string
	Tasks    []tasks.Task
}

type UserApp struct {
	Db *gorm.DB
}

func (u *UserApp) InitUserDB() {
	u.Db.AutoMigrate(&User{})
}
