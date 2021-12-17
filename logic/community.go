package logic

import (
	"03.threeGames/dao/mysql"
	"03.threeGames/models"
	"03.threeGames/pkg/snowflake"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"strconv"
	"strings"
	"time"
)

func CommunityCreatePage(page *models.CreatePage) (pageID int64, err error) {
	pageID = snowflake.GenID()
	createTime := time.Now().UnixMilli()
	pageSql := &models.CreatePageSql{
		UserId:  page.UserId,
		PageId:  pageID,
		Title:   page.Title,
		Content: page.Content,
		CreateTime: createTime,
	}
	err = mysql.CommunityCreatePageComment(pageID)
	if err != nil {
		return
	}
	err = mysql.CommunityCreatePage(pageSql)
	return
}

func CommunityGetPageList(pageList *[]models.GetPage, pageStart, pageSize int) (err error) {
	var pageListSql []models.GetPageSql
	err = mysql.CommunityGetPageList(&pageListSql, pageStart-1, pageSize)
	if err != nil {
		return
	}
	var commentsNum, replyNum int
	var pageID string
	user := new(models.User)
	for i := 0; i < len(pageListSql); i++ {
		commentsNum, err = mysql.CommunityGetCommentsNum(pageListSql[i].PageID)
		if err != nil {
			return
		}
		user, err = mysql.AuthGetUserByID(pageListSql[i].UserID)
		if err != nil {
			return
		}
		replyNum, err = mysql.CommunityGetReplyListNum(pageListSql[i].PageID)
		if err != nil {
			return
		}
		pageID = strconv.FormatInt(pageListSql[i].PageID, 10)
		*pageList = append(*pageList, models.GetPage{
			Views:       pageListSql[i].Views,
			CommentsNum: commentsNum + replyNum,
			PageID:      pageID,
			Username:    user.Username,
			UserAvatar:  user.Avatar,
			Title:   pageListSql[i].Title,
			PreContent: pageListSql[i].PreContent,
			Approve:     pageListSql[i].Approve,
			CreateTime: pageListSql[i].CreateTime,
		})
	}
	return
}

func CommunityGetPageContent(pageID int64) (detail *models.GetPageDetail, err error) {
	var detailSql *models.GetPageDetailSql
	var user *models.User
	var commentsNum, replyNum int
	detailSql, err = mysql.CommunityGetPageDetail(pageID)
	if err != nil {
		return
	}
	user, err = mysql.AuthGetUserByID(detailSql.UserID)
	if err != nil {
		return
	}
	commentsNum, err = mysql.CommunityGetCommentsNum(pageID)
	if err != nil {
		return
	}
	replyNum, err = mysql.CommunityGetReplyListNum(pageID)
	if err != nil {
		return
	}
	detail = &models.GetPageDetail{
		Views: detailSql.Views,
		PageID: strconv.FormatInt(detailSql.PageID, 10),
		Username: user.Username,
		UserAvatar: user.Avatar,
		Title: detailSql.Title,
		Content:    detailSql.Content,
		Approve:    detailSql.Approve,
		CommentsNum: commentsNum + replyNum,
		CreateTime: detailSql.CreateTime,
		UpdateTime: detailSql.UpdateTime,
	}
	err = mysql.CommunityPageViewsAdd(pageID)
	return
}

func CommunityModifyPage(pageID, userID int64, content string) (err error) {
	var exist bool
	exist, err = mysql.CommunityCheckPageExist(pageID, userID)
	if err != nil {
		return
	}
	if exist == false {
		err = errors.New("pageID or userID is wrong")
		zap.L().Error("community modify page failed", zap.Error(err))
		return
	}
	err = mysql.CommunityModifyPage(pageID, content)
	return
}

func CommunityPageAddApprove(pageID, userID int64) (err error) {
	var user *models.User
	user, err = mysql.AuthGetUserByID(userID)
	if err != nil {
		return
	}
	err = mysql.CommunityPageApproveAdd(pageID, user.Username)
	return
}

func CommunityPageRemoveApprove(pageID, userID int64) (err error) {
	var user *models.User
	user, err = mysql.AuthGetUserByID(userID)
	if err != nil {
		return
	}
	var approves string
	approves, err = mysql.CommunityPageGetApprove(pageID)
	index := strings.Index(approves, user.Username)
	length := len(user.Username)
	if index != -1 {
		approves = approves[:index] + approves[index+length+1:]
		err = mysql.CommunityPageModifyApprove(pageID, approves)
	} else {
		err = errors.New("no target username in approves")
		zap.L().Error("community page remove approve failed", zap.Error(err))
	}
	return
}

func CommunityGetPageComments(comments *[]models.GetComment, pageID int64, pageStart, pageSize int, reverse bool) (err error) {
	var commentsSql []models.GetCommentSql
	err = mysql.CommunityGetCommentList(&commentsSql, pageID, pageStart-1, pageSize, reverse)
	if err != nil {
		return
	}
	user := new(models.User)
	for i := 0; i < len(commentsSql); i++ {
		user, err = mysql.AuthGetUserByID(commentsSql[i].UserID)
		if err != nil {
			return
		}
		*comments = append(*comments, models.GetComment{
			CreateTime: commentsSql[i].CreateTime,
			ID: commentsSql[i].ID,
			Content:    commentsSql[i].Content,
			Approve:    commentsSql[i].Approve,
			Reply:      commentsSql[i].Reply,
			Username:   user.Username,
			UserAvatar: user.Avatar,
		})
	}
	return
}

func CommunityPageInsertComment(pageID, userID int64, content string) (lastInsertID int, err error) {
	var id int64
	id, err = mysql.CommunityPageInsertComment(pageID, userID, content)
	if err != nil {
		return
	}
	lastInsertID = int(id)
	return
}

func CommunityCommentsAddReply(query *models.Reply) (reply string, err error) {
	var user *models.User
	user, err = mysql.AuthGetUserByID(query.UserID)
	if err != nil {
		return
	}
	err = mysql.CommunityCommentInsertReply(query.PageID, query.CommentID, user.Username, user.Avatar, query.TargetName, query.Content)
	if err != nil {
		return
	}
	reply, err = mysql.CommunityCommentGetReply(query.PageID, query.CommentID)
	return
}

func CommunityCommentAddApprove(pageID, commentID, userID int64) (err error) {
	var user *models.User
	user, err = mysql.AuthGetUserByID(userID)
	if err != nil {
		return
	}
	err = mysql.CommunityCommentApproveAdd(pageID, commentID, user.Username)
	return
}

func CommunityCommentRemoveApprove(pageID, commentID, userID int64) (err error) {
	var user *models.User
	user, err = mysql.AuthGetUserByID(userID)
	if err != nil {
		return
	}
	var approves string
	approves, err = mysql.CommunityCommentGetApprove(pageID, commentID)
	index := strings.Index(approves, user.Username)
	length := len(user.Username)
	if index != -1 {
		approves = approves[:index] + approves[index+length+1:]
		err = mysql.CommunityCommentModifyApprove(pageID, commentID, approves)
	} else {
		err = errors.New("no target username in approves")
		zap.L().Error("community comment remove approve failed", zap.Error(err))
	}
	return
}