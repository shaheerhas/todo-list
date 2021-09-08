package user

import (
	"github.com/shaheerhas/todo-list/app/tasks"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username string `gorm:"not null"`
	Email    string `gorm:"not null"`
	Password string `gorm:"not null"`
	Tasks    []tasks.Task
}

type UserApp struct {
	Db *gorm.DB
}

func (u *UserApp) InitUserDB() {
	u.Db.AutoMigrate(&User{})
}
