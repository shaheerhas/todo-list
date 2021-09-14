package users

import (
	"github.com/shaheerhas/todo-list/app/tasks"
	"gorm.io/gorm"
)

type UserModel struct {
	ID         uint         `gorm:"autoIncrement; primaryKey"`
	FirstName  string       `gorm:"not null"`
	LastName   string       `gorm:"not null"`
	Email      string       `gorm:"not null;unique"`
	Password   string       `gorm:"not null"`
	IsVerified bool         `gorm:"not null; default:false"`
	Tasks      []tasks.Task `gorm:"ForeignKey:UserID"`
}

type UserModelApp struct {
	Db *gorm.DB
}

func (u *UserModelApp) InitUserModelDB() {
	u.Db.AutoMigrate(&UserModel{})
}
