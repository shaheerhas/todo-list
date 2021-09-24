package users

import (
	"github.com/shaheerhas/todo-list/app/tasks"
	"gorm.io/gorm"
	"log"
)

type UserModel struct {
	ID         uint   `gorm:"autoIncrement; primaryKey"`
	FirstName  string `gorm:"not null"`
	LastName   string `gorm:"not null"`
	Email      string `gorm:"not null;unique"`
	Password   string
	IsVerified bool `gorm:"not null; default:false"`
	FbUser     bool
	Tasks      []tasks.Task `gorm:"foreignKey:UserID; onDelete:cascade"`
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
