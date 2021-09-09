package user

type NoUsers struct{}

func (m *NoUsers) Error() string {
	return "no users in db"
}

func allUsers(svc UserApp, users *[]User) error {
	if err := svc.Db.Find(&users).Error; err != nil {
		return err
	}
	if len(*users) == 0 {
		return &NoUsers{}
	}
	return nil
}

func getUserById(svc UserApp, id int) {}
func createUser(svc UserApp, user User) (User, error) {

	err := svc.Db.Create(&user)
	if err.Error != nil {
		return user, err.Error
	}
	return user, nil
}
