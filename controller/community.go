package controller

import (
	"03.threeGames/logic"
	"03.threeGames/models"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"strconv"
)

func CommunityCreatePage(c *gin.Context)  {
	userID, err := GetCurrentUser(c)
	if err != nil {
		ResponseError(c, CodeNotLogin)
		return
	}
	p := new(models.CreatePage)
	p.UserId = userID
	err = c.ShouldBindJSON(p)
	if err != nil {
		zap.L().Error("community create page with invalid params", zap.Error(err))
		ResponseError(c, CodeInvalidParams)
		return
	}
	var pageID int64
	pageID, err = logic.CommunityCreatePage(p)
	if err != nil {
		ResponseError(c, CodeServerBusy)
		return
	}
	var pageContent *models.GetPageDetail
	pageContent, err = logic.CommunityGetPageContent(pageID)
	if err != nil {
		ResponseError(c, CodeServerBusy)
		return
	}
	ResponseSuccess(c, pageContent)
}

func CommunityModifyPage(c *gin.Context) {
	pageIDStr := c.Param("page_id")
	pageID, err := strconv.ParseInt(pageIDStr, 10, 64)
	if err != nil {
		zap.L().Error("community get page detail by pageID failed, pageIDStr parse into pageID failed", zap.Error(err))
		ResponseError(c, CodeInvalidParams)
		return
	}
	query := new(models.QueryPageModify)
	if err = c.ShouldBindJSON(query); err != nil {
		zap.L().Error("community modify page with invalid params", zap.Error(err))
		ResponseError(c, CodeInvalidParams)
		return
	}
	var userID int64
	userID, err = GetCurrentUser(c)
	if err != nil {
		ResponseError(c, CodeNotLogin)
		return
	}
	if err != nil {
		zap.L().Error("community get page detail by pageID failed, pageIDStr parse into pageID failed", zap.Error(err))
		ResponseError(c, CodeInvalidParams)
		return
	}
	err = logic.CommunityModifyPage(pageID, userID, query.Content)
	if err != nil {
		ResponseError(c, CodeServerBusy)
		return
	}
	var pageDetail *models.GetPageDetail
	pageDetail, err = logic.CommunityGetPageContent(pageID)
	if err != nil {
		ResponseError(c, CodeServerBusy)
		return
	}
	ResponseSuccess(c, pageDetail)
}

func CommunityGetPageList(c *gin.Context) {
	query := new(models.QueryPage)
	err := c.ShouldBindJSON(query)
	if err != nil {
		zap.L().Error("community get page list with invalid params", zap.Error(err))
		ResponseError(c, CodeInvalidParams)
		return
	}
	if query.Length > 20 {
		zap.L().Error("community get page list with too long length", zap.Any("query_length", query.Length))
		ResponseError(c, CodeInvalidParams)
		return
	}
	pageList := new([]models.GetPage)
	err = logic.CommunityGetPageList(pageList, query.Start, query.Length)
	if err != nil {
		ResponseError(c, CodeServerBusy)
		return
	}
	ResponseSuccess(c, pageList)
}

func CommunityGetPageDetail(c *gin.Context) {
	pageIDStr := c.Param("page_id")
	pageID, err := strconv.ParseInt(pageIDStr, 10, 64)
	if err != nil {
		zap.L().Error("community get page detail by pageID failed, pageIDStr parse into pageID failed", zap.Error(err))
		ResponseError(c, CodeInvalidParams)
		return
	}
	var pageDetail *models.GetPageDetail
	pageDetail, err = logic.CommunityGetPageContent(pageID)
	if err != nil {
		ResponseError(c, CodeServerBusy)
		return
	}
	ResponseSuccess(c, pageDetail)
}

func CommunityPageApprove(c *gin.Context) {
	pageIDStr := c.Param("page_id")
	pageID, err := strconv.ParseInt(pageIDStr, 10, 64)
	if err != nil {
		zap.L().Error("community modify page approve by pageID failed, pageIDStr parse into pageID failed", zap.Error(err))
		ResponseError(c, CodeInvalidParams)
		return
	}
	query := new(models.QueryPageApprove)
	if err = c.ShouldBindJSON(query); err != nil {
		zap.L().Error("community modify page approve with invalid params", zap.Error(err))
		ResponseError(c, CodeInvalidParams)
	}
	var userID int64
	userID, err = GetCurrentUser(c)
	if err != nil {
		ResponseError(c, CodeNotLogin)
		return
	}
	if *query.Add {
		err = logic.CommunityPageAddApprove(pageID, userID)
	} else {
		err = logic.CommunityPageRemoveApprove(pageID, userID)
	}
	if err != nil {
		ResponseError(c, CodeServerBusy)
		return
	}
	ResponseSuccess(c, nil)
}

func CommunityGetComments(c *gin.Context) {
	pageIDStr := c.Param("page_id")
	pageID, err := strconv.ParseInt(pageIDStr, 10, 64)
	if err != nil {
		zap.L().Error("community get page comments by pageID failed, pageIDStr parse into pageID failed", zap.Error(err))
		ResponseError(c, CodeInvalidParams)
		return
	}
	query := new(models.QueryComment)
	err = c.ShouldBindJSON(query)
	if err != nil {
		zap.L().Error("community get page comments list with invalid params", zap.Error(err))
		ResponseError(c, CodeInvalidParams)
		return
	}
	if query.Length > 20 {
		zap.L().Error("community get page comments list with too long length", zap.Any("query_length", query.Length))
		ResponseError(c, CodeInvalidParams)
		return
	}
	commentsList := new([]models.GetComment)
	err = logic.CommunityGetPageComments(commentsList, pageID, query.Start, query.Length, *query.Reverse)
	if err != nil {
		zap.L().Error("111", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	ResponseSuccess(c, commentsList)
}

func CommunityAddComment(c *gin.Context) {
	pageIDStr := c.Param("page_id")
	pageID, err := strconv.ParseInt(pageIDStr, 10, 64)
	if err != nil {
		zap.L().Error("community get pageID failed, pageIDStr parse into pageID failed", zap.Error(err))
		ResponseError(c, CodeInvalidParams)
		return
	}
	var userID int64
	userID, err = GetCurrentUser(c)
	if err != nil {
		ResponseError(c, CodeNotLogin)
		return
	}
	query := new(models.QueryCommentAdd)
	if err = c.ShouldBindJSON(query); err != nil {
		zap.L().Error("community add comment with invalid params", zap.Error(err))
		ResponseError(c, CodeInvalidParams)
		return
	}
	_, err = logic.CommunityPageInsertComment(pageID, userID, query.Content)
	if err != nil {
		ResponseError(c, CodeServerBusy)
		return
	}
	ResponseSuccess(c, nil)
}

func CommunityAddReply(c *gin.Context) {
	pageIDStr := c.Param("page_id")
	pageID, err := strconv.ParseInt(pageIDStr, 10, 64)
	if err != nil {
		zap.L().Error("community add reply failed, pageIDStr parse into pageID failed", zap.Error(err))
		ResponseError(c, CodeInvalidParams)
		return
	}
	var userID int64
	userID, err = GetCurrentUser(c)
	if err != nil {
		ResponseError(c, CodeNotLogin)
		return
	}
	query := new(models.QueryReply)
	if err = c.ShouldBindJSON(query); err != nil {
		zap.L().Error("community add reply with invalid params", zap.Error(err))
		ResponseError(c, CodeInvalidParams)
		return
	}
	var commentID int64
	commentID, err = strconv.ParseInt(query.CommentID, 10, 64)
	if err != nil {
		ResponseError(c, CodeInvalidParams)
		return
	}
	Reply := &models.Reply{
		PageID:     pageID,
		UserID:     userID,
		CommentID: commentID,
		TargetName: *query.TargetName,
		Content:    query.Content,
	}
	var response string
	response, err = logic.CommunityCommentsAddReply(Reply)
	if err != nil {
		ResponseError(c,CodeServerBusy)
		return
	}
	ResponseSuccess(c, gin.H{
		"reply": response,
	})
}

func CommunityCommentApprove(c *gin.Context) {
	pageIDStr := c.Param("page_id")
	pageID, err := strconv.ParseInt(pageIDStr, 10, 64)
	if err != nil {
		zap.L().Error("community modify comment approve by pageID failed, pageIDStr parse into pageID failed", zap.Error(err))
		ResponseError(c, CodeInvalidParams)
		return
	}
	query := new(models.QueryCommentApprove)
	if err = c.ShouldBindJSON(query); err != nil {
		zap.L().Error("community modify comment approve with invalid params", zap.Error(err))
		ResponseError(c, CodeInvalidParams)
	}
	var userID int64
	userID, err = GetCurrentUser(c)
	if err != nil {
		ResponseError(c, CodeNotLogin)
		return
	}
	commentID := int64(query.ID)
	if *query.Add {
		err = logic.CommunityCommentAddApprove(pageID, commentID, userID)
	} else {
		err = logic.CommunityCommentRemoveApprove(pageID, commentID, userID)
	}
	if err != nil {
		ResponseError(c, CodeServerBusy)
		return
	}
	ResponseSuccess(c, nil)
}
