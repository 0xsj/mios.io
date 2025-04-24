package context

import (
	"errors"

	"github.com/gin-gonic/gin"
)

const (
	UserKey   = "user"
	UserIDKey = "user_id"
	TokenKey  = "token"
)

var (
	ErrUserNotFound = errors.New("user not found in context")
)

func SetUserID(c *gin.Context, userID string) {
	c.Set(UserIDKey, userID)
}

func GetUserID(c *gin.Context) (string, error) {
	userID, exists := c.Get(UserIDKey)
	if !exists {
		return "", ErrUserNotFound
	}
	return userID.(string), nil
}

func SetUser(c *gin.Context, user interface{}) {
	c.Set(UserKey, user)
}

func GetUser(c *gin.Context) (interface{}, error) {
	user, exists := c.Get(UserKey)
	if !exists {
		return nil, ErrUserNotFound
	}
	return user, nil
}

func SetToken(c *gin.Context, token string) {
	c.Set(TokenKey, token)
}

func GetToken(c *gin.Context) (string, error) {
	token, exists := c.Get(TokenKey)
	if !exists {
		return "", errors.New("token not found in context")
	}
	return token.(string), nil
}
