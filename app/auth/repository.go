package auth

func CreateBlackListToken(token BlackListToken) error {
	result := Db.Create(&token)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func IsBlackListed(token string) bool {
	res := Db.First(&BlackListToken{}, "token_val = ?", token)
	if res.Error != nil {
		return false
	}
	if res.RowsAffected < 1 {
		return false
	}
	return true
}
