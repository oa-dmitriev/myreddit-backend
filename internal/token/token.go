package token

import (
	"fmt"
	"log"
	"myreddit/pkg/user"
	"os"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

var (
	privateKey = "ljaflsdkfjaldsfj"
)

func CreateToken(u *user.User) (string, error) {
	userid := u.Id
	os.Setenv("ACCESS_SECRET", privateKey)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"ID":  userid,
		"exp": time.Now().Add(time.Hour * 15).Unix(),
		"user": map[string]interface{}{
			"login": u.Login,
			"id":    u.Id,
		},
	})
	tokenString, err := token.SignedString([]byte(os.Getenv("ACCESS_SECRET")))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func ExtractToken(c *gin.Context) (*jwt.Token, error) {
	token, err := verifyToken(c)
	if err != nil {
		log.Println("verifyToken is failed: ", err)
		return nil, err
	}
	return token, nil
}

func verifyToken(c *gin.Context) (*jwt.Token, error) {
	tokenString := extractTokenString(c)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("ACCESS_SECRET")), nil
	})
	if err != nil {
		return nil, err
	}
	return token, nil
}

func extractTokenString(c *gin.Context) string {
	bearToken := c.GetHeader("Authorization")
	strArr := strings.Split(bearToken, " ")
	if len(strArr) == 2 {
		return strArr[1]
	}
	return ""
}
