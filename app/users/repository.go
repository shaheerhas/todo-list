package users

import "fmt"

type NoUsers struct{}

func (m *NoUsers) Error() string {
	return "no users in db"
}

func allUsers(svc UserModelApp, users *[]UserModel) error {
	if err := svc.Db.Find(&users).Error; err != nil {
		return err
	}
	if len(*users) == 0 {
		return &NoUsers{}
	}
	return nil
}

func updateStatus(svc UserModelApp, user UserModel, status bool) error {
	err := svc.Db.Model(&user).Where("email = ?", user.Email).Update("is_verified", status).Error
	if err != nil {
		return err
	}
	return nil
}

func getUser(svc UserModelApp, email string) (UserModel, error) {
	var loginUser UserModel
	result := svc.Db.Where("email = ?", email).Find(&loginUser)
	if result.Error != nil {
		return loginUser, result.Error
	}
	if result.RowsAffected < 1 {
		return loginUser, fmt.Errorf("no user found with this email")
	}
	return loginUser, nil

}

func createUser(svc UserModelApp, user UserModel) (UserModel, error) {
	err := svc.Db.Create(&user)
	if err.Error != nil {
		return user, err.Error
	}
	return user, nil
}
