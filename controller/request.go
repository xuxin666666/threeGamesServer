package controller

import (
	"errors"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

var ErrorUserNotLogin = errors.New("用户未登录")
const ContextUserIDKey = "userID"

func GetCurrentUser(c *gin.Context) (userID int64, err error) {
	uid, ok := c.Get(ContextUserIDKey)
	if !ok {
		err = ErrorUserNotLogin
		zap.L().Error("get userID from context failed", zap.Error(err))
		return
	}
	userID, ok = uid.(int64)
	if !ok {
		err = ErrorUserNotLogin
		zap.L().Error("userID turn into int64 from string failed", zap.Error(err))
		return
	}
	return
}
