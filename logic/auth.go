package logic

import (
	"03.threeGames/dao/mysql"
	"03.threeGames/models"
	"03.threeGames/pkg/jwt"
	"03.threeGames/pkg/snowflake"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"go.uber.org/zap"
	"time"
)

// AuthRegister 注册
func AuthRegister(person *models.UserRegister) (err error) {
	var exist bool
	exist, err = mysql.AuthCheckExist(person.Username)
	if err != nil {
		return
	}
	if exist == true {
		err = errors.New("username existed")
		return
	}
	userID := snowflake.GenID()
	user := &models.User{
		UserID:   userID,
		Username: person.Username,
		Password: encryptPassword(person.Password),
		Avatar:   person.Avatar,
	}
	return mysql.AuthRegister(user)
}

// AuthLogin 登录
func AuthLogin(person *models.UserLogin) (avatar string, token string, err error) {
	user := new(models.User)
	user.Username = person.Username
	err = mysql.AuthGetUserByUsername(user)
	if err != nil {
		return
	}
	if user.Password != encryptPassword(person.Password) {
		err = errors.New("something wrong between username and password")
		zap.L().Error("auth login failed", zap.Error(err))
		return
	}
	var duration time.Duration
	if *person.AutoLogin {
		duration = time.Hour * 24 * 30
	} else {
		duration = time.Hour * 12
	}
	avatar = user.Avatar
	token, err = jwt.GenToken(user.UserID, user.Username, duration)
	if err != nil {
		zap.L().Error("create token failed", zap.Error(err))
	}
	return
}

// AuthChangeAvatar 更改头像
func AuthChangeAvatar(uid int64, avatar string) (err error) {
	return mysql.AuthChangeAvatar(uid, avatar)
}

func AuthLoginAuto(uid int64, p *models.UserLoginAutoResponse) (err error) {
	user := new(models.User)
	user, err = mysql.AuthGetUserByID(uid)
	if err != nil {
		return
	}
	p.Username = user.Username
	p.Avatar = user.Avatar
	return
}

func encryptPassword(oPassword string) string {
	h := md5.New()
	h.Write([]byte(oPassword))
	nPassword := h.Sum([]byte("FATAL ERROR"))
	return hex.EncodeToString(nPassword)
}
