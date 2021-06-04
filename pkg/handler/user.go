package handler

import (
	"fmt"
	"log"
	"myreddit/internal/token"
	"myreddit/pkg/user"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	UserRepo *user.UserRepo
}

func (h *UserHandler) Register(c *gin.Context) {
	u := &user.User{}
	err := c.ShouldBindJSON(u)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": err.Error(),
		})
		return
	}
	err = h.UserRepo.Register(u)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"Error": err.Error(),
		})
		return
	}
	token, err := token.CreateToken(u)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"Error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": token,
		"id":    u.Id,
	})
}

func (h *UserHandler) Login(c *gin.Context) {
	u := &user.User{}
	err := c.ShouldBindJSON(u)
	if err != nil {
		log.Println("JSON: ", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": err.Error(),
		})
		return
	}
	err = h.UserRepo.Login(u)
	if err != nil {
		log.Println("repo.Login: ", err)
		c.JSON(http.StatusUnauthorized, gin.H{
			"Error": err.Error(),
		})
		return
	}
	token, err := token.CreateToken(u)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"Error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"token": token,
		"id":    u.Id,
	})
}

func getUser(c *gin.Context) (*user.User, error) {
	if u, ok := c.Get("user"); ok {
		user := u.(*user.User)
		return user, nil
	}
	return nil, fmt.Errorf("no user found")
}
