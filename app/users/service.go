package users

import (
	"encoding/base64"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/shaheerhas/todo-list/app/auth"
	"github.com/shaheerhas/todo-list/app/utils"
	"gopkg.in/gomail.v2"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

func (svc UserModelApp) getUsers(c *gin.Context) {
	var users []UserModel
	err := allUsers(svc, &users)
	if err != nil {
		log.Println(err)
		msg := "record not found"
		c.JSON(utils.Response(http.StatusInternalServerError, msg))
		return
	}
	c.JSON(http.StatusOK, users)

}

func (svc UserModelApp) login(c *gin.Context) {
	var user UserModel
	if err := c.BindJSON(&user); err != nil {
		log.Println(err)
		msg := "invalid or malformed payload"
		c.JSON(utils.Response(http.StatusBadRequest, msg))
		return
	}
	loginUser, err := getUser(svc, user.Email)
	if err != nil {
		log.Println(err)
		msg := "user with this email not found"
		c.JSON(utils.Response(http.StatusNotFound, msg))
		return
	}
	verified := utils.CheckPasswordHash(user.Password, loginUser.Password)
	if !verified {
		msg := "incorrect password"
		c.JSON(utils.Response(http.StatusForbidden, msg))
		return
	}
	token, err := auth.GenerateJWTToken(loginUser.ID, loginUser.Email)
	if err != nil {
		log.Println(err)
		msg := "token creation error"
		c.JSON(utils.Response(http.StatusInternalServerError, msg))
		return
	}
	c.SetCookie("accessToken", token, 0, "", "", false, false)
	msg := "login success"
	c.JSON(utils.Response(http.StatusOK, msg))

}

func (svc UserModelApp) logout(c *gin.Context) {
	//userId := getId(c)
	tokenString, err := c.Cookie("accessToken")
	if err != nil {
		log.Println(err)
	}
	token := auth.BlackListToken{TokenVal: tokenString}
	err = auth.CreateToken(token)
	if err != nil {
		log.Println(err)
		return
	}
	c.SetCookie("accessToken", "", -1, "/", "", false, false)
	msg := "user logged out successfully"
	c.JSON(utils.Response(http.StatusOK, msg))
}

func (svc UserModelApp) signup(c *gin.Context) {
	var user UserModel
	if err := c.BindJSON(&user); err != nil {
		log.Println(err)
		msg := "invalid or malformed payload"
		c.JSON(utils.Response(http.StatusBadRequest, msg))
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
		if strings.Contains(err.Error(), "(SQLSTATE 23505)") {
			msg := "user with this email already registered"
			c.JSON(utils.Response(http.StatusConflict, msg))
			return
		}
		log.Println(err)
		msg := "couldn't create record in db"
		c.JSON(utils.Response(http.StatusInternalServerError, msg))
		return
	}
	var confirmationUrl = os.Getenv("CONFIRMATION_URL")
	body := "Hi, click this link below to confirm your Todo-list Account\n"
	body += confirmationUrl + "/" + encode(user.Email, user.ID)
	err = sendEmail(user, body, "Todo-list Account Confirmation")
	if err != nil {
		log.Println(err)
		msg := "couldn't send email"
		c.JSON(utils.Response(http.StatusInternalServerError, msg))
		return
	}
	msg := "user successfully created, check your email for confirmation email"
	c.JSON(utils.Response(http.StatusCreated, msg))
}

func encode(email string, id uint) string {
	email += "?" + strconv.Itoa(int(id))
	encoded := base64.URLEncoding.EncodeToString([]byte(email))
	return encoded
}

func decode(encodedUrl string) (string, error) {
	decoded, err := base64.URLEncoding.DecodeString(encodedUrl)
	return string(decoded), err
}

func sendEmail(user UserModel, body, subject string) error {
	var senderEmail = os.Getenv("SENDER_EMAIL")
	var senderPassword = os.Getenv("SENDER_PASSWORD")
	msg := gomail.NewMessage()
	msg.SetHeader("From", senderEmail)
	fmt.Println(senderEmail)
	msg.SetHeader("To", user.Email)
	msg.SetHeader("Subject", subject)

	msg.SetBody("text/html", body)
	d := gomail.NewDialer("smtp.gmail.com", 587, senderEmail, senderPassword)

	if err := d.DialAndSend(msg); err != nil {
		return err
	}
	return nil
}

func (svc UserModelApp) forgotPassword(c *gin.Context) {
	var userReq map[string]string
	if err := c.BindJSON(&userReq); err != nil {
		log.Println(err)
		msg := "invalid or malformed payload"
		c.JSON(utils.Response(http.StatusBadRequest, msg))
		return
	}

	email := userReq["email"]
	user, err := getUser(svc, email)
	if err != nil {
		log.Println(err)
		msg := "user with this email not found"
		c.JSON(utils.Response(http.StatusNotFound, msg))
		return
	}

	resetPasswordLink := os.Getenv("RESET_PASSWORD")
	subject := "Password Reset Request for Todo-List App"
	body := "Hi, " + user.FirstName + ", we got a request from your account to reset your password"
	body += "\nClick the link  below to reset your account, or just ignore this email if you didn't request resetting your password.\n"
	body += "\n" + resetPasswordLink + "/" + encode(user.Email, user.ID)
	err = sendEmail(user, body, subject)
	if err != nil {
		log.Println(err)
		msg := "couldn't send email"
		c.JSON(utils.Response(http.StatusInternalServerError, msg))
		return
	}
	msg := "email reset link sent, check your email"
	c.JSON(utils.Response(http.StatusOK, msg))

}

func (svc UserModelApp) resetPassword(c *gin.Context) {
	encodedString := c.Param("emailToken")
	if auth.IsBlackListed(encodedString) {
		msg := "user already has reset their password"
		c.JSON(utils.Response(http.StatusConflict, msg))
		return
	}

	decodedEmailID, err := decode(encodedString)
	if err != nil {
		log.Println(err)
	}
	decoded := strings.Split(decodedEmailID, "?")
	decodedEmail := decoded[0]
	decodedID := decoded[1]
	user, err := getUser(svc, decodedEmail)
	if err != nil {
		log.Println(err)
		msg := "user with email doesn't exist"
		c.JSON(utils.Response(http.StatusNotFound, msg))
		return
	}
	if strconv.Itoa(int(user.ID)) != decodedID {
		log.Println("decoded ID not equal to user's id")
		msg := "couldn't reset password"
		c.JSON(utils.Response(http.StatusConflict, msg))
		return
	}
	var reqBody map[string]string
	if err := c.BindJSON(&reqBody); err != nil {
		log.Println(err)
		msg := "incorrect or malformed payload"
		c.JSON(utils.Response(http.StatusBadRequest, msg))
		return
	}

	hashedPassword, err := utils.HashPassword(reqBody["password"])
	if err != nil {
		log.Println("couldn't hash password", err)
	}
	err = updatePassword(svc, user.ID, hashedPassword)
	if err != nil {
		log.Println(err)
		msg := "couldn't update password"
		c.JSON(utils.Response(http.StatusInternalServerError, msg))
		return
	}
	err = auth.CreateToken(auth.BlackListToken{TokenVal: encodedString})
	if err != nil {
		log.Println(err)
		return
	}
	msg := "password successfully updated"
	c.JSON(utils.Response(http.StatusOK, msg))
}

func (svc UserModelApp) confirmUser(c *gin.Context) {
	encodedString := c.Param("emailToken")
	decodedEmailID, err := decode(encodedString)
	if err != nil {
		log.Println(err)
	}
	decoded := strings.Split(decodedEmailID, "?")
	decodedEmail := decoded[0]
	decodedID := decoded[1]

	user, err := getUser(svc, decodedEmail)

	if err != nil {
		log.Println(err)
		msg := "user with email doesn't exist"
		c.JSON(utils.Response(http.StatusNotFound, msg))
		return
	}
	if user.IsVerified {
		msg := "user already confirmed!"
		c.JSON(utils.Response(http.StatusOK, msg))
		return
	}
	if strconv.Itoa(int(user.ID)) != decodedID {
		msg := "couldn't confirm user"
		c.JSON(utils.Response(http.StatusConflict, msg))
		return
	}
	err = updateStatus(svc, user.ID, true)
	if err != nil {
		log.Println(err)
		msg := "couldn't confirm user"
		c.JSON(utils.Response(http.StatusInternalServerError, msg))
		return
	}
	msg := "user confirmed!"
	c.JSON(utils.Response(http.StatusOK, msg))

}
