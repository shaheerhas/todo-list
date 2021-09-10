package user

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/shaheerhas/todo-list/app/utils"
)

func (svc UserApp) getUsers(c *gin.Context) {
	var users []User
	err := allUsers(svc, &users)
	if err != nil {
		log.Println(err)
		c.IndentedJSON(http.StatusBadRequest, err)
		return
	}
	fmt.Println(users)
	c.IndentedJSON(http.StatusOK, users)

}

func (svc UserApp) patchUser(c *gin.Context) {

}

func (svc UserApp) postUser(c *gin.Context) {
	var user User
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
