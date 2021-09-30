package users

import (
	"fmt"
	"gorm.io/gorm"
)

func AllUsers(db *gorm.DB) ([]UserModel, error) {
	var users []UserModel
	if err := db.Find(&users).Error; err != nil {
		return users, err
	}
	if len(users) == 0 {
		return users, gorm.ErrRecordNotFound
	}
	return users, nil
}

func updateStatus(svc UserModelApp, userId uint, status bool) error {
	err := svc.Db.Model(&UserModel{}).Where("id = ?", userId).Update("is_verified", status).Error
	return err
}

func updatePassword(svc UserModelApp, userId uint, password string) error {
	err := svc.Db.Model(&UserModel{}).Where("id = ?", userId).Update("password", password).Error
	return err
}

func getUser(svc UserModelApp, email string) (UserModel, error) {
	var user UserModel
	result := svc.Db.Where("email = ?", email).Find(&user)
	if result.Error != nil {
		return user, result.Error
	}
	if result.RowsAffected < 1 {
		return user, fmt.Errorf("no user found with this email")
	}
	return user, nil

}

func createUser(svc UserModelApp, user UserModel) (UserModel, error) {
	err := svc.Db.Create(&user)
	if err.Error != nil {
		return user, err.Error
	}
	return user, nil
}
