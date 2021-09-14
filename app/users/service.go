package users

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/shaheerhas/todo-list/app/auth"
	"github.com/shaheerhas/todo-list/app/utils"
)

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
	fmt.Println(loginUser, user)
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

func (svc UserModelApp) postUser(c *gin.Context) {
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
	fmt.Println(user)
	c.IndentedJSON(http.StatusCreated, "record successfully created")
}
