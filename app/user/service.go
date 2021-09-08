package user

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (svc UserApp) getUsers(c *gin.Context) {
	var users []User
	err := allUsers(users, svc)
	if err != nil {
		log.Println(err)
		c.IndentedJSON(http.StatusBadRequest, "couldn't get users")
		return
	}
	c.IndentedJSON(http.StatusOK, users)

}

func (svc UserApp) postUser(c *gin.Context) {
	var user User
	if err := c.BindJSON(&user); err != nil {
		log.Println(err)
		c.IndentedJSON(http.StatusBadRequest, "some issue with your json formatting")
		return
	}
	err := createUser(user, svc)
	if err != nil {
		log.Println(err)
		c.IndentedJSON(http.StatusBadRequest, "couldn't create record in db")
		return
	}
	fmt.Println(user)
	c.IndentedJSON(http.StatusCreated, "record successfully created")
}
