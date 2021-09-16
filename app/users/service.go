package users

import (
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/shaheerhas/todo-list/app/auth"
	"github.com/shaheerhas/todo-list/app/utils"
	"gopkg.in/gomail.v2"
)

var senderEmail = os.Getenv("SENDER_EMAIL")
var senderPassword = os.Getenv("SENDER_PASSWORD")
var confirmationUrl = os.Getenv("CONFIRMATION_URL")

func (svc UserModelApp) getUsers(c *gin.Context) {
	var users []UserModel
	err := allUsers(svc, &users)
	if err != nil {
		log.Println(err)
		c.IndentedJSON(http.StatusBadRequest, err)
		return
	}
	fmt.Println(users)
	c.IndentedJSON(http.StatusOK, users)

}

func (svc UserModelApp) login(c *gin.Context) {
	var user UserModel
	if err := c.BindJSON(&user); err != nil {
		log.Println(err)
		c.IndentedJSON(http.StatusBadRequest, "some issue with your json formatting")
		return
	}
	loginUser, err := getUser(svc, user.Email)
	if err != nil {
		log.Println(err)
		c.IndentedJSON(http.StatusNotFound, "user with this email not found")
		return
	}
	verified := utils.CheckPasswordHash(user.Password, loginUser.Password)
	if !verified {
		c.IndentedJSON(http.StatusUnauthorized, "password not correct")
		return
	}
	token, err := auth.GenerateJWTToken(loginUser.ID, loginUser.Email)
	if err != nil {
		log.Println(err)
		c.IndentedJSON(http.StatusInternalServerError, "token creation error")
		return
	}
	c.SetCookie("accessCookie", token, 0, "", "", false, false)
	c.IndentedJSON(http.StatusOK, "login success")

}
func (svc UserModelApp) patchUser(c *gin.Context) {

}

func (svc UserModelApp) signup(c *gin.Context) {
	var user UserModel
	if err := c.BindJSON(&user); err != nil {
		log.Println(err)
		c.IndentedJSON(http.StatusBadRequest, "some issue with your json formatting")
		return
	}
	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		log.Println("couldn't hash password")
	} else {
		user.Password = hashedPassword
	}
	user, err = createUser(svc, user)
	if err != nil {
		log.Println(err)
		c.IndentedJSON(http.StatusInternalServerError, "couldn't create record in db")
		return
	}

	err = sendEmail(user)
	if err != nil {
		c.IndentedJSON(http)
	}
	c.IndentedJSON(http.StatusCreated, "user successfully created, check your email for confirmation email")
}

func encode(email string, id uint) string {
	email += "?" + strconv.Itoa(int(id))
	encoded := base64.URLEncoding.EncodeToString([]byte(email))
	return encoded
}

func decode(encodedUrl string) (string, error) {
	// for confirm url
	decoded, err := base64.URLEncoding.DecodeString(encodedUrl)
	return string(decoded), err
}

func sendEmail(user UserModel) error {
	msg := gomail.NewMessage()
	msg.SetHeader("From", senderEmail)
	fmt.Println(senderEmail)
	msg.SetHeader("To", user.Email)
	msg.SetHeader("Subject", "Account Confirmation Todo-list Application")
	body := "Hi, click this link below to confirm your Todo-list Account\n"
	body += confirmationUrl + "/" + encode(user.Email, user.ID)
	msg.SetBody("text/html", body)
	d := gomail.NewDialer("smtp.gmail.com", 587, senderEmail, senderPassword)

	if err := d.DialAndSend(msg); err != nil {
		log.Println(err)
		return err
	}
}

func getEmailFromURL(url string) string {
	email := url[9:]
	fmt.Println(len(confirmationUrl), url, email)
	return email
}

func (svc UserModelApp) confirmUser(c *gin.Context) {
	//url := c.Request.URL.String()
	userEmail := c.Param("emailToken")
	decodedEmailID, err := decode(userEmail)
	if err != nil {
		log.Println(err)
	}
	decoded := strings.Split(decodedEmailID, "?")
	decodedEmail := decoded[0]
	decodedID := decoded[1]

	user, err := getUser(svc, decodedEmail)
	fmt.Println(userEmail, decodedEmail, decodedID, user, err)
	if err != nil {
		log.Println(err)
		c.IndentedJSON(http.StatusNotFound, "user with email doesn't exist")
		return
	}
	if user.IsVerified {
		c.IndentedJSON(http.StatusOK, "user already confirmed!")
		return
	}
	if strconv.Itoa(int(user.ID)) != decodedID {
		c.IndentedJSON(http.StatusBadRequest, "user id email mismatch")
		return
	}
	err = updateStatus(svc, user, true)
	if err != nil {
		log.Println(err)
		c.IndentedJSON(http.StatusInternalServerError, "couldn't update status of the user")
		return
	}
	c.IndentedJSON(http.StatusOK, "user confirmed!")
}
