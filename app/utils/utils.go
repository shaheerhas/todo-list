package utils

import (
	"fmt"
	"log"
	"strconv"

	"golang.org/x/crypto/bcrypt"
)

func Response(statusCode int, msg string) (int, map[string]string) {
	return statusCode, map[string]string{"msg": msg}
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func ConvertInterfaceToUint(i interface{}) uint {
	var uId int
	switch v := i.(type) {
	case string:
		{
			uId, _ = strconv.Atoi(v)
		}
	case float64:
		{
			uId = int(v)
		}
	case int:
		{
			uId = int(v)
		}
	case int64:
		{
			uId = int(v)
		}
	case uint:
		{
			uId = int(v)
		}
	default:
		{
			log.Println("couldn't parse id")
		}
	}
	fmt.Println("uid func2", uint(uId))
	return uint(uId)
}
