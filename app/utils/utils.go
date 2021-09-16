package utils

import (
	"log"
	"strconv"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func ConvertInterfaceToUint(i interface{}) uint {

	var uId uint64
	switch v := i.(type) {
	case string:
		uId, err := strconv.ParseUint(v, 10, 32)
		if err != nil {
			log.Println("err, uId", err, uId)
			return 0
		}
	case float64:
		{
			uId = uint64(v)
		}
	case int:
		{
			uId = uint64(v)
		}

	}

	return uint(uId)
}
