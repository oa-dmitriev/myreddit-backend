package middleware

import (
	"fmt"
	"myreddit/internal/token"
	"myreddit/pkg/user"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		fmt.Println("WELKRJWELRJWKERJWLERJKW\nwlekrjwlekjrwle")
		user, err := getUser(c)
		if err != nil {
			c.AbortWithStatusJSON(
				http.StatusUnauthorized,
				"unautharized, dude")
		}
		c.Set("user", user)
	}
}

func getUser(c *gin.Context) (*user.User, error) {
	token, err := token.ExtractToken(c)
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("invalid token")
	}
	u := &user.User{}
	err = u.Unpack(claims["user"])
	if err != nil {
		return nil, err
	}
	return u, nil
}
