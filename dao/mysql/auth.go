package mysql

import (
	"03.threeGames/models"
	"database/sql"
	"go.uber.org/zap"
)

func AuthGetUserByUsername(user *models.User) (err error) {
	sqlStr := `select user_id, username, password, avatar from users where username=?`
	if err = db.Get(user, sqlStr, user.Username); err != nil {
		zap.L().Error("auth get user failed", zap.Error(err), zap.Any("user", user))
	}
	return
}

func AuthCheckExist(username string) (exist bool, err error) {
	var count int
	sqlStr := `select count(user_id) from users where username=?`
	exist = false
	err = db.Get(&count, sqlStr, username)
	if err != nil {
		zap.L().Error("auth check exist failed", zap.Error(err), zap.Any("username", username))
		return
	}
	if count > 0 {
		zap.L().Warn("username existed", zap.Any("username", username))
		exist = true
	}
	return
}

func AuthRegister(user *models.User) (err error) {
	sqlStr := `insert into users(user_id, username, password, avatar) values(?,?,?,?)`
	_, err = db.Exec(sqlStr, user.UserID, user.Username, user.Password, user.Avatar)
	if err != nil {
		zap.L().Error("auth register failed", zap.Error(err), zap.Any("user", user))
	}
	return
}

func AuthChangeAvatar(uid int64, avatar string) (err error) {
	sqlStr := `update users set avatar=? where user_id=?`
	var ret sql.Result
	ret, err = db.Exec(sqlStr, avatar, uid)
	if err != nil {
		zap.L().Error("changeAvatar failed", zap.Error(err))
		return
	}
	_, err = ret.RowsAffected()
	if err != nil {
		zap.L().Error("query rows failed", zap.Error(err))
	}
	return
}

func AuthGetUserByID(userID int64) (user *models.User, err error) {
	user = new(models.User)
	sqlStr := `select user_id, username, password, avatar from users where user_id=?`
	if err = db.Get(user, sqlStr, userID); err != nil {
		zap.L().Error("auth get user by user_id failed", zap.Error(err), zap.Any("user_id", userID))
	}
	return
}
