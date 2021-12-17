package controller

import (
	"03.threeGames/logic"
	"03.threeGames/models"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

func AuthLogin(c *gin.Context) {
	p := new(models.UserLogin)
	if err :=c.ShouldBindJSON(p); err != nil {
		zap.L().Error("invalid login params", zap.Error(err))
		ResponseErrorWithMsg(c, CodeInvalidParams, "参数有误")
		return
	}
	if len(p.Username) < 3 || len(p.Username) > 12 || len(p.Password) < 6 || len(p.Password) > 18 {
		err := errors.New("invalid length of username or password")
		zap.L().Error("invalid login params", zap.Error(err))
		ResponseErrorWithMsg(c, CodeInvalidParams, "参数有误")
		return
	}
	avatar, token, err := logic.AuthLogin(p)
	if err != nil {
		ResponseErrorWithMsg(c, CodeActiveFailed, "登录失败")
		return
	}
	res := &models.UserLoginResponse{
		Avatar: avatar,
		Token:  token,
	}
	ResponseSuccess(c, res)
}

func AuthRegister(c *gin.Context) {
	p := new(models.UserRegister)
	if err := c.ShouldBindJSON(p); err != nil {
		zap.L().Error("register with invalid params", zap.Error(err))
		ResponseErrorWithMsg(c, CodeInvalidParams, "参数有误")
		return
	}
	if len(p.Username) < 3 || len(p.Username) > 12 || len(p.Password) < 6 || len(p.Password) > 18 {
		err := errors.New("invalid length of username or password")
		zap.L().Error("invalid login params", zap.Error(err))
		ResponseErrorWithMsg(c, CodeInvalidParams, "参数有误")
		return
	}
	if err := logic.AuthRegister(p); err != nil {
		ResponseErrorWithMsg(c, CodeServerBusy, "注册失败")
		return
	}
	ResponseSuccess(c, nil)
}

func AuthChangeAvatar(c *gin.Context)  {
	p := new(models.UserAvatar)
	if err := c.ShouldBindJSON(p); err != nil {
		zap.L().Error("change avatar with invalid params", zap.Error(err))
		ResponseErrorWithMsg(c, CodeInvalidParams, "参数有误")
		return
	}
	userID, err :=GetCurrentUser(c)
	if err != nil {
		ResponseError(c, CodeNotLogin)
		return
	}
	err = logic.AuthChangeAvatar(userID, p.Avatar)
	if err != nil {
		ResponseError(c, CodeActiveFailed)
	}
	ResponseSuccess(c, nil)
}

func AuthLoginAuto(c *gin.Context) {
	uid, err := GetCurrentUser(c)
	if err != nil {
		ResponseError(c, CodeNotLogin)
		return
	}
	userR := new(models.UserLoginAutoResponse)
	if err = logic.AuthLoginAuto(uid, userR); err != nil {
		ResponseError(c, CodeServerBusy)
		return
	}
	ResponseSuccess(c, userR)
}
