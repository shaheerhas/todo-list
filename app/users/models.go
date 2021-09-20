package users

import (
	"github.com/shaheerhas/todo-list/app/tasks"
	"gorm.io/gorm"
	"log"
)

type UserModel struct {
	ID         uint         `gorm:"autoIncrement; primaryKey"`
	FirstName  string       `gorm:"not null"`
	LastName   string       `gorm:"not null"`
	Email      string       `gorm:"not null;unique"`
	Password   string       `gorm:"not null"`
	IsVerified bool         `gorm:"not null; default:false"`
	Tasks      []tasks.Task `gorm:"ForeignKey:UserID; onDelete CASCADE"`
}

type UserModelApp struct {
	Db *gorm.DB
}

func (svc *UserModelApp) InitUserModelDB() {
	err := svc.Db.AutoMigrate(&UserModel{})
	if err != nil {
		log.Println(err)
		return
	}
}
