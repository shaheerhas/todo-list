package users

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

func getUser(svc UserModelApp, email string) (UserModel, error) {
	var loginUser UserModel
	if err := svc.Db.Where("email = ?", email).Find(&loginUser).Error; err != nil {

		return loginUser, err
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
