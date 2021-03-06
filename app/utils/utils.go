package utils

import (
	"encoding/base64"
	"gopkg.in/gomail.v2"
	"log"
	"os"
	"strconv"
	"time"

	"golang.org/x/crypto/bcrypt"
)

func Response(statusCode int, msg string) (int, map[string]string) {
	return statusCode, map[string]string{"msg": msg}
}

// Encode encodes incoming parameters to a hash, EncodeToString returns the base64 encoding of src.
func Encode(email string, id uint) string {
	email += "?" + strconv.Itoa(int(id)) + "?" + strconv.Itoa(int(time.Now().Unix()))
	encoded := base64.URLEncoding.EncodeToString([]byte(email))
	return encoded
}

// Decode decodes the previously encoded url
func Decode(encodedString string) (string, error) {
	decoded, err := base64.URLEncoding.DecodeString(encodedString)
	return string(decoded), err
}

func SendEmail(userEmail, body, subject string) error {
	var senderEmail = os.Getenv("SENDER_EMAIL")
	var senderPassword = os.Getenv("SENDER_PASSWORD")
	msg := gomail.NewMessage()
	msg.SetHeader("From", senderEmail)
	msg.SetHeader("To", userEmail)
	msg.SetHeader("Subject", subject)

	msg.SetBody("text/html", body)
	d := gomail.NewDialer("smtp.gmail.com", 587, senderEmail, senderPassword)

	if err := d.DialAndSend(msg); err != nil {
		return err
	}
	log.Println("Email sent to", userEmail)
	return nil
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
	return uint(uId)
}
