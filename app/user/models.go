package user

import (
	"github.com/shaheerhas/todo-list/app/tasks"
	"gorm.io/gorm"
)

type User struct {
	Id       int
	Username string       `gorm:"unique; notnull"`
	Email    string       `gorm:"unique; notnull"`
	Password string       `gorm:"notnull"`
	Tasks    []tasks.Task `gorm:"foreignKey:userId"`
}

type UserApp struct {
	Db *gorm.DB
}
