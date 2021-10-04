package users

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-co-op/gocron"
	"github.com/shaheerhas/todo-list/app/auth"
	"github.com/shaheerhas/todo-list/app/tasks"
	"github.com/shaheerhas/todo-list/app/utils"
	"golang.org/x/oauth2"
	"gorm.io/gorm"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

var oauthConf *oauth2.Config
var oauthStateString string

func isUserLoggedIn(c *gin.Context) bool {
	accessToken, _ := c.Cookie("accessToken")
	//check if user is already logged in //dos
	_, validErr := auth.ValidateJWTToken(accessToken)
	if accessToken != "" && validErr == nil && !auth.IsBlackListed(accessToken) {
		return true
	}
	return false
}

func initOAuth() {
	oauthStateString = os.Getenv("FB_STATE_STRING")
	oauthConf = &oauth2.Config{
		ClientID:     os.Getenv("FB_APP_ID"),
		ClientSecret: os.Getenv("FB_CLIENT_SECRET"),
		RedirectURL:  os.Getenv("FB_CALLBACK_URL"),
		Scopes:       []string{"public_profile", "email"},
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://www.facebook.com/v12.0/dialog/oauth",
			TokenURL: "https://graph.facebook.com/v12.0/oauth/access_token",
		},
	}
}

var scheduler *gocron.Scheduler

func ScheduleEmail(userApp UserModelApp) {
	scheduler = gocron.NewScheduler(time.Now().Location())
	//err, _ := scheduler.Every(1).Day().At(time.Now().Add(time.Second * 1)).Do(userApp.SendReminderEmails)
	err, _ := scheduler.Every(1).Day().At("00:00").Do(userApp.SendReminderEmails)
	if err != nil {
		log.Println(err)
	}
	scheduler.StartAsync()
}

func (svc UserModelApp) login(c *gin.Context) {
	if isUserLoggedIn(c) {
		msg := "user already logged in"
		c.JSON(utils.Response(http.StatusOK, msg))
		return
	}
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
	if loginUser.FbUser {
		log.Println("fb user")
		msg := "wrong credentials user registered as an fb user"
		c.JSON(utils.Response(http.StatusBadRequest, msg))
		return
	}
	if !loginUser.IsVerified {
		log.Println("user not verified yet")
		msg := "please verify your email first"
		c.JSON(utils.Response(http.StatusForbidden, msg))
		return
	}

	if verified := utils.CheckPasswordHash(user.Password, loginUser.Password); !verified {
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

func (svc UserModelApp) fbLogin(c *gin.Context) {

	if isUserLoggedIn(c) {
		msg := "user already logged in"
		c.JSON(utils.Response(http.StatusOK, msg))
		return
	}
	initOAuth()
	URL := oauthConf.AuthCodeURL(oauthStateString)
	log.Println(URL)
	c.Redirect(http.StatusTemporaryRedirect, URL)
}

func (svc UserModelApp) fbCallBack(c *gin.Context) {
	state := c.Request.FormValue("state")
	if state != oauthStateString {
		log.Printf("invalid oauth state, expected '%s', got '%s'\n", oauthStateString, state)
		msg := "invalid oauth state"
		c.JSON(utils.Response(http.StatusBadRequest, msg))
		return
	}

	code := c.Request.FormValue("code")
	token, err := oauthConf.Exchange(c, code)

	if err != nil {
		log.Printf("oauthConf.Exchange() failed with '%s'\n", err)
		msg := "invalid or malformed payload"
		c.JSON(utils.Response(http.StatusBadRequest, msg))
		return
	}

	resp, err := http.Get("https://graph.facebook.com/me?fields=id,name,email&access_token=" +
		url.QueryEscape(token.AccessToken))

	if err != nil {
		log.Println(err)
		msg := "invalid or malformed payload"
		c.JSON(utils.Response(http.StatusBadRequest, msg))
		return
	}

	var respJSON map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&respJSON)
	if err != nil {
		log.Println(err)
		msg := "invalid or malformed payload"
		c.JSON(utils.Response(http.StatusBadRequest, msg))
		return
	}

	user, err := svc.fbSignup(respJSON)
	// handle two types of errors, one is that user is already registered
	// second just log the error
	if err != nil {
		if strings.Contains(err.Error(), "already registered") {
			log.Println("user already registered")
			jwtToken, err := auth.GenerateJWTToken(user.ID, user.Email)
			if err != nil {
				log.Println(err)
			}
			c.SetCookie("accessToken", jwtToken, 0, "", "", false, false)
			msg := "login success"
			c.JSON(utils.Response(http.StatusOK, msg))
			return
		}

		log.Println(err)
		msg := "user cannot be registered on server"
		c.JSON(utils.Response(http.StatusInternalServerError, msg))
		return
	}
	msg := "user created successfully"
	log.Println("user email", user.Email, respJSON["email"])
	jwtToken, err := auth.GenerateJWTToken(user.ID, respJSON["email"].(string))
	if err != nil {
		log.Println(err)
	}
	c.SetCookie("accessToken", jwtToken, 0, "", "", false, false)
	msg += "\nlogin success"
	//c.JSON(utils.Response(http.StatusOK, msg))
	c.JSON(utils.Response(http.StatusCreated, msg))

}

func (svc UserModelApp) fbSignup(fbUser map[string]interface{}) (UserModel, error) {
	email := fbUser["email"].(string)
	name := strings.Split(fbUser["name"].(string), " ")

	user := UserModel{
		FirstName:  name[0],
		LastName:   name[1],
		Email:      email,
		Password:   "",
		IsVerified: true,
		FbUser:     true,
		Tasks:      nil,
	}
	dbUser, err := createUser(svc, user)
	if dbUser.ID == 0 {
		msg := "already registered"
		retUser, _ := getUser(svc, email)
		return retUser, fmt.Errorf(msg)
	}
	return dbUser, err
}

func (svc UserModelApp) logout(c *gin.Context) {
	tokenString, err := c.Cookie("accessToken")
	if err != nil {
		log.Println(err)
	}
	token := auth.BlackListToken{TokenVal: tokenString}
	err = auth.CreateBlackListToken(token)
	if err != nil {
		log.Println(err)
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

	if user.Password == "" {
		msg := "password cannot be empty"
		log.Println("empty password error:", msg)
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
	body += confirmationUrl + "/" + utils.Encode(user.Email, user.ID)
	err = utils.SendEmail(user.Email, body, "Todo-list Account Confirmation")
	if err != nil {
		log.Println(err)
		msg := "couldn't send email"
		c.JSON(utils.Response(http.StatusInternalServerError, msg))
		return
	}
	msg := "user successfully created, check your email for confirmation email"
	c.JSON(utils.Response(http.StatusCreated, msg))
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
	if user.FbUser {
		msg := "user is registered as a fb user"
		c.JSON(utils.Response(http.StatusBadRequest, msg))
		return
	}
	if !user.IsVerified {
		msg := "verify your account first"
		c.JSON(utils.Response(http.StatusConflict, msg))
		return
	}

	resetPasswordLink := os.Getenv("RESET_PASSWORD")
	subject := "Password Reset Request for Todo-List App"
	body := "Hi, " + user.FirstName + ", we got a request from your account to reset your password."
	body += "\nClick the link  below to reset your account, or just ignore this email if you didn't request resetting your password.\n"
	body += "\n\r\n" + resetPasswordLink + "/" + utils.Encode(user.Email, user.ID)
	err = utils.SendEmail(user.Email, body, subject)
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
		msg := "password reset link invalid"
		c.JSON(utils.Response(http.StatusConflict, msg))
		return
	}

	decodedEmailID, err := utils.Decode(encodedString)
	if err != nil {
		log.Println(err)
	}
	decoded := strings.Split(decodedEmailID, "?")
	if len(decoded) < 1 {
		msg := "incorrect confirmation token"
		c.JSON(utils.Response(http.StatusBadRequest, msg))
		return
	}
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
	err = auth.CreateBlackListToken(auth.BlackListToken{TokenVal: encodedString})
	if err != nil {
		log.Println(err)
	}
	msg := "password successfully updated"
	c.JSON(utils.Response(http.StatusOK, msg))
}

func (svc UserModelApp) confirmUser(c *gin.Context) {
	encodedString := c.Param("emailToken")
	decodedEmailID, err := utils.Decode(encodedString)
	if err != nil {
		log.Println(err)
	}

	decoded := strings.Split(decodedEmailID, "?")
	if len(decoded) < 2 {
		msg := "incorrect confirmation token"
		c.JSON(utils.Response(http.StatusBadRequest, msg))
		return
	}
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

func sendEmailThread(wg *sync.WaitGroup, user UserModel, todayTasks []tasks.Task) {

	subject := "Task Due Reminder Email"
	body := "Hey" + " " + user.FirstName
	body += "\n This is to remind you that you have the following tasks due for today!\n"
	for _, task := range todayTasks {
		body += task.Title + "\n"
	}
	err := utils.SendEmail(user.Email, body, subject)
	if err != nil {
		log.Println(err)
	}
	wg.Done()
}

//, todayTasks []tasks.Task
func sendEmailAndDBThread(wg *sync.WaitGroup, user UserModel, db *gorm.DB) {
	todayTasks, err := tasks.FindDueTodayTasks(db, user.ID)
	if err != nil {
		log.Println(err)
	}
	if len(todayTasks) > 0 {
		subject := "Task Due Reminder Email"
		body := "Hey" + " " + user.FirstName
		body += "\n This is to remind you that you have the following tasks due for today!\n"
		for _, task := range todayTasks {
			body += task.Title + "\n"
		}
		err = utils.SendEmail(user.Email, body, subject)
		if err != nil {
			log.Println(err)
		}
		time.Sleep(1)
		wg.Done()
	}
}
func (svc UserModelApp) SendReminderEmails() {
	time1 := time.Now().UnixMicro()
	allUsers, err := AllUsers(svc.Db)
	var wg sync.WaitGroup
	if err != nil {
		log.Println(err)
	}
	for _, user := range allUsers {
		wg.Add(1)
		todayTasks, err := tasks.FindDueTodayTasks(svc.Db, user.ID)
		if err != nil {
			log.Println(err)
		}
		if len(todayTasks) > 0 {
			//sendEmailThread(&wg, user, svc.Db)
			sendEmailThread(&wg, user, todayTasks)
		}
	}

	defer fmt.Println("time taken", time.Now().UnixMicro()-time1)
	wg.Wait()
}
