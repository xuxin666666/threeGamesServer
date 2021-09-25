package controller

import (
	"github.com/gin-gonic/gin"
)

type ResCode int

const (
	CodeSuccess       ResCode = 200
	CodeNotLogin      ResCode = 401
	CodeForbidden     ResCode = 403
	CodeNotFound      ResCode = 404
	CodeActiveFailed  ResCode = 430
	CodeRefreshFailed ResCode = 431
	CodeServerBusy    ResCode = 500
	CodeOtherReasons  ResCode = 501
)

var CodeMsgMap = map[ResCode]string{
	CodeSuccess:       "请求成功",
	CodeNotLogin:      "未登录",
	CodeForbidden:     "权限不够",
	CodeNotFound:      "这里空空如也",
	CodeActiveFailed:  "active token失效了",
	CodeRefreshFailed: "refresh token失效了",
	CodeServerBusy:    "服务繁忙",
	CodeOtherReasons:  "未知错误",
}

type ResponseData struct {
	Status ResCode     `json:"status"`
	Msg  interface{} `json:"msg"`
	Data interface{} `json:"data"`
}

func (c ResCode) Msg() string {
	msg, ok := CodeMsgMap[c]
	if !ok {
		msg = CodeMsgMap[CodeServerBusy]
	}
	return msg
}

func ResponseSuccess(c *gin.Context, data interface{})  {
	rd := &ResponseData{
		Status: CodeSuccess,
		Msg:  CodeSuccess.Msg(),
		Data: data,
	}
	c.JSON(int(CodeSuccess), rd)
}

func ResponseError(c *gin.Context, code ResCode)  {
	rd := &ResponseData{
		Status: code,
		Msg:  code.Msg(),
		Data: nil,
	}
	c.JSON(int(code), rd)
}

func ResponseErrorWithMsg(c *gin.Context, code ResCode, msg interface{})  {
	rd := &ResponseData{
		Status: code,
		Msg:  msg,
		Data: nil,
	}
	c.JSON(int(code), rd)
}
