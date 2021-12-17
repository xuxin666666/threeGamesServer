package controller

import (
	"03.threeGames/logic"
	"03.threeGames/models"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func GetMineSweepScores(c *gin.Context)  {
	userID, err := GetCurrentUser(c)
	if err != nil {
		ResponseError(c, CodeNotLogin)
		return
	}
	var scores string
	scores, err = logic.MineSweepGetScores(userID)
	if err != nil {
		ResponseError(c, CodeServerBusy)
		return
	}
	ResponseSuccess(c, scores)
}

func PostMineSweepScores(c *gin.Context) {
	p := new(models.MineSweepScores)
	if err := c.ShouldBindJSON(p); err != nil {
		zap.L().Error("get tetris scores with invalid params", zap.Error(err))
		ResponseError(c, CodeInvalidParams)
		return
	}
	userID, err := GetCurrentUser(c)
	if err != nil {
		ResponseError(c, CodeNotLogin)
		return
	}
	err = logic.MineSweepUpdateScores(userID, *p.Scores)
	if err != nil {
		ResponseError(c, CodeServerBusy)
		return
	}
	ResponseSuccess(c, nil)
}

func PostMineSweepScoresList(c *gin.Context) {
	p := new(models.MineSweepScoreList)
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
	response, err = logic.MineSweepUpdateScoreList(userID, *p.Scores)
	if err != nil {
		ResponseError(c, CodeServerBusy)
		return
	}
	ResponseSuccess(c, gin.H{
		"scores": response,
	})
}
