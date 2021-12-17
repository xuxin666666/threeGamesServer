package controller

import (
	"03.threeGames/logic"
	"03.threeGames/models"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func GetTetrisScores(c *gin.Context)  {
	userID, err := GetCurrentUser(c)
	if err != nil {
		ResponseError(c, CodeNotLogin)
		return
	}
	var scores string
	scores, err = logic.TetrisGetScores(userID)
	if err != nil {
		ResponseError(c, CodeServerBusy)
		return
	}
	ResponseSuccess(c, gin.H{
		"scores": scores,
	})
}

func PostTetrisScores(c *gin.Context) {
	p := new(models.TetrisScores)
	if err := c.ShouldBindJSON(p); err != nil {
		zap.L().Error("post tetris score with invalid params", zap.Error(err))
		ResponseError(c, CodeInvalidParams)
		return
	}
	userID, err := GetCurrentUser(c)
	if err != nil {
		ResponseError(c, CodeNotLogin)
		return
	}
	err = logic.TetrisUpdateScores(userID, *p.Scores)
	if err != nil {
		ResponseError(c, CodeServerBusy)
		return
	}
	ResponseSuccess(c, nil)
}

func PostTetrisScoresList(c *gin.Context) {
	p := new(models.TetrisScoreList)
	if err := c.ShouldBindJSON(p); err != nil {
		zap.L().Error("post tetris scores with invalid params", zap.Error(err))
		ResponseError(c, CodeInvalidParams)
		return
	}
	userID, err := GetCurrentUser(c)
	if err != nil {
		ResponseError(c, CodeNotLogin)
		return
	}

	var response string
	response, err = logic.TetrisUpdateScoreList(userID, *p.Scores)
	if err != nil {
		ResponseError(c, CodeServerBusy)
		return
	}
	ResponseSuccess(c, gin.H{
		"scores": response,
	})
}
