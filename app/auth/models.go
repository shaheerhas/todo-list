package auth

import (
	"gorm.io/gorm"
	"log"
)

type BlackListToken struct {
	TokenVal string `gorm:"primaryKey"`
}
type AuthApp struct {
	Db *gorm.DB
}

var Db *gorm.DB

func (authApp *AuthApp) InitBlackListModel() {
	err := authApp.Db.AutoMigrate(&BlackListToken{})
	Db = authApp.Db
	if err != nil {
		log.Println(err)
		return
	}
}
