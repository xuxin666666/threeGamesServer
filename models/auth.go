package models

type UserLogin struct {
	Username  string `json:"username" binding:"required"`
	Password  string `json:"password" binding:"required"`
	AutoLogin *bool   `json:"auto_login" binding:"required"`
}

type UserLoginAutoResponse struct {
	Username string `json:"username" db:"username"`
	Avatar   string `json:"avatar" db:"avatar"`
}

type UserRegister struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Avatar   string `json:"avatar" binding:"required"`
}

type User struct {
	UserID   int64  `json:"user_id" db:"user_id"`
	Username string `json:"username" db:"username"`
	Password string `json:"password" db:"password"`
	Avatar   string `json:"avatar" db:"avatar"`
}

type UserLoginResponse struct {
	Avatar string `json:"avatar"`
	Token  string `json:"token"`
}

type UserAvatar struct {
	Avatar string `json:"avatar" binding:"required"`
}
