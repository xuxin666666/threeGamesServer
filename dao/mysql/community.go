package mysql

import (
	"03.threeGames/models"
	"database/sql"
	"go.uber.org/zap"
	"strconv"
	"strings"
	"time"
)

// CommunityCreatePage 创建帖子
func CommunityCreatePage(page *models.CreatePageSql) (err error) {
	sqlStr := `insert into community(user_id, page_id, title, content, pre_content, create_time) values(?,?,?,?,?,?)`
	var ret sql.Result
	var preContent string
	contentRune := []rune(page.Content)
	if len(contentRune) > 120 {
		preContent = string(contentRune[0:120]) + "..."
	} else {
		preContent = page.Content
	}
	ret, err = db.Exec(sqlStr, page.UserId, page.PageId, page.Title, page.Content, preContent, page.CreateTime)
	if err != nil {
		zap.L().Error("community create page failed", zap.Error(err))
		return
	}
	_, err = ret.LastInsertId()
	if err != nil {
		zap.L().Error("community get insert page id failed", zap.Error(err))
	}
	return
}

// CommunityGetPageList 得到一列帖子，内容为预览内容
func CommunityGetPageList(pageList *[]models.GetPageSql, start int, len int) (err error) {
	sqlStr := `select views, page_id, user_id, title, pre_content, approve, create_time from community order by id desc limit ?,?`
	err = db.Select(pageList, sqlStr, strconv.Itoa(start), strconv.Itoa(start+len))
	if err == sql.ErrNoRows {
		zap.L().Warn("no pages in community", zap.Error(err))
		err = nil
	}
	if err != nil {
		zap.L().Error("select GetPageSql from community failed", zap.Error(err))
	}
	return
}

// CommunityGetPageDetail 得到某个帖子的详情
func CommunityGetPageDetail(pageID int64) (detail *models.GetPageDetailSql, err error) {
	detail = new(models.GetPageDetailSql)
	sqlStr := `select page_id, user_id, views, title, content, approve, create_time, update_time from community where page_id=?`
	err = db.Get(detail, sqlStr, pageID)
	if err != nil {
		zap.L().Error("community get page content by pageID failed", zap.Error(err))
	}
	return
}

// CommunityModifyPage 修改某个帖子
func CommunityModifyPage(pageID int64, content string) (err error) {
	sqlStr := `update community set content=?, update_time=? where page_id=?`
	var ret sql.Result
	updateTime := time.Now().UnixMilli()
	ret, err = db.Exec(sqlStr, content, updateTime, pageID)
	if err != nil {
		zap.L().Error("community modify page failed", zap.Error(err))
		return
	}
	_, err = ret.RowsAffected()
	if err != nil {
		zap.L().Error("community query page modified failed", zap.Error(err))
	}
	return
}

// CommunityCheckPageExist 根据pageID和userID检查帖子是否存在
func CommunityCheckPageExist(pageID, userID int64) (exist bool, err error) {
	sqlStr := `select count(id) from community where page_id=? and user_id=?`
	var count int
	err = db.Get(&count, sqlStr, pageID, userID)
	if err != nil {
		zap.L().Error("community get page by pageID and userID failed", zap.Error(err))
		return
	}
	exist = false
	if count > 0 {
		exist = true
	}
	return
}

// CommunityPageViewsAdd 浏览数+1
func CommunityPageViewsAdd(pageID int64) (err error) {
	sqlStr := `update community set views=views+1 where page_id=?`
	var ret sql.Result
	ret, err = db.Exec(sqlStr, pageID)
	if err != nil {
		zap.L().Error("community modify page views failed", zap.Any("pageID", pageID), zap.Error(err))
		return
	}
	_, err = ret.RowsAffected()
	if err != nil {
		zap.L().Error("community query page views modified failed", zap.Any("pageID", pageID), zap.Error(err))
	}
	return
}

// CommunityPageApproveAdd 增加点赞者的id
func CommunityPageApproveAdd(pageID int64, username string) (err error) {
	username = username + ","
	sqlStr := `update community set approve=concat(approve,?) where page_id=?`
	var ret sql.Result
	ret, err = db.Exec(sqlStr, username, pageID)
	if err != nil {
		zap.L().Error("community page add approve failed", zap.Any("pageID", pageID), zap.Error(err))
		return
	}
	_, err = ret.RowsAffected()
	if err != nil {
		zap.L().Error("community query page add approve modified failed", zap.Any("pageID", pageID), zap.Error(err))
	}
	return
}

// CommunityPageModifyApprove 修改整个点赞的内容
func CommunityPageModifyApprove(pageID int64, approve string) (err error) {
	sqlStr := `update community set approve=? where page_id=?`
	var ret sql.Result
	ret, err = db.Exec(sqlStr, approve, pageID)
	if err != nil {
		zap.L().Error("community page modify approve failed", zap.Any("pageID", pageID), zap.Error(err))
		return
	}
	_, err = ret.RowsAffected()
	if err != nil {
		zap.L().Error("community query page approve modified failed", zap.Any("pageID", pageID), zap.Error(err))
	}
	return
}

// CommunityPageGetApprove 得到点赞的内容
func CommunityPageGetApprove(pageID int64) (approve string, err error) {
	sqlStr := `select approve from community where page_id=?`
	err = db.Get(&approve, sqlStr, pageID)
	if err != nil {
		zap.L().Error("community page get approve failed", zap.Any("pageID", pageID), zap.Error(err))
		return
	}
	return
}

// CommunityCreatePageComment 创建帖子的评论表格
func CommunityCreatePageComment(pageID int64) (err error) {
	pageIDStr := strconv.FormatInt(pageID, 10)
	sqlStr1 := "drop table if exists `" + pageIDStr + "`"
	_, err = db.Exec(sqlStr1)
	if err != nil {
		zap.L().Error("drop comment table failed", zap.Any("comment table ID", pageID), zap.Error(err))
		return
	}
	sqlStr := "create table `" + pageIDStr + "` (\n    `id` bigint(20) not null auto_increment,\n    `user_id` bigint(20) not null ,\n    `content` varchar(1000) collate utf8mb4_general_ci not null ,\n    `approve` varchar(6000) collate utf8mb4_general_ci default '',\n    `reply` varchar(6000) collate utf8mb4_general_ci default '',\n    `create_time` bigint(20) not null comment '创建时间',\n    `replyNum` int default 0 not null,\n    primary key (`id`)\n) engine = InnoDB default charset = utf8mb4 collate = utf8mb4_general_ci;"
	_, err = db.Exec(sqlStr)
	if err != nil {
		zap.L().Error("create comment table failed", zap.Any("comment table ID", pageID), zap.Error(err))
	}
	return
}

// CommunityGetCommentsNum 得到某个帖子的评论总数
func CommunityGetCommentsNum(pageID int64) (commentNum int, err error) {
	commentNum = 0
	pageIDStr := strconv.FormatInt(pageID, 10)
	sqlStr := "select count(id) from `" + pageIDStr + "`"
	err = db.Get(&commentNum, sqlStr)
	if err != nil {
		zap.L().Error("community get comments length failed", zap.Error(err))
	}
	return
}

// CommunityGetReplyListNum 得到某个帖子的回复总数
func CommunityGetReplyListNum(pageID int64) (replyNum int, err error) {
	pageIDStr := strconv.FormatInt(pageID, 10)
	count := 0
	sqlStr1 := "select count(replyNum) from `" + pageIDStr + "`"
	err = db.Get(&count, sqlStr1)
	if err != nil {
		zap.L().Error("community get replyNum from comments table failed", zap.Error(err))
		return
	}
	if count == 0 {
		replyNum = 0
		return
	}
	sqlStr := "select sum(replyNum) from `" + pageIDStr + "`"
	err = db.Get(&replyNum, sqlStr)
	if err != nil {
		zap.L().Error("community get reply num failed", zap.Error(err))
	}
	return
}

// CommunityGetCommentList 得到某个帖子的评论列表，从start开始长度为len的len条数据
func CommunityGetCommentList(commentList *[]models.GetCommentSql, pageID int64, start, len int, reverse bool) (err error) {
	pageIDStr := strconv.FormatInt(pageID, 10)
	var sqlStr string
	if reverse {
		sqlStr = "select id, user_id, content, approve, reply, create_time from `" + pageIDStr + "` order by id desc limit ?,?"
	} else {
		sqlStr = "select id, user_id, content, approve, reply, create_time from `" + pageIDStr + "` limit ?,?"
	}
	err = db.Select(commentList, sqlStr, start, start+len)
	if err == sql.ErrNoRows {
		zap.L().Error("no comments in page", zap.Any("pageID", pageID), zap.Error(err))
		commentList = new([]models.GetCommentSql)
		err = nil
	}
	if err != nil {
		zap.L().Error("community select comments in page failed", zap.Any("pageID", pageID), zap.Error(err))
	}
	return
}

// CommunityPageInsertComment 某个帖子添加评论
func CommunityPageInsertComment(pageID, userID int64, content string) (lastInsertID int64, err error) {
	pageIDStr := strconv.FormatInt(pageID, 10)
	tim := time.Now().UnixMilli()
	sqlStr := "insert into `" + pageIDStr +"`(user_id, content, create_time) values (?,?,?)"
	var ret sql.Result
	ret, err = db.Exec(sqlStr, userID, content, tim)
	if err != nil {
		zap.L().Error("community page insert comment failed", zap.Error(err))
		return
	}
	lastInsertID, err = ret.LastInsertId()
	if err != nil {
		zap.L().Error("community page get inserted comment id failed", zap.Error(err))
	}
	return
}

// CommunityCommentInsertReply 某个帖子某个评论添加回复
func CommunityCommentInsertReply(pageID, ID int64, username, userAvatar, targetName, content string) (err error) {
	var build strings.Builder
	build.WriteString("```")
	build.WriteString(username)
	build.WriteString("`")
	build.WriteString(userAvatar)
	build.WriteString("`")
	build.WriteString(targetName)
	build.WriteString("`")
	build.WriteString(content)
	// response的结构是 ```username`userAvatar`targetName`content，每条数据用```分隔，每条数据的每个部分用`分隔
	response := build.String()
	pageIDStr := strconv.FormatInt(pageID, 10)
	sqlStr := "update `" + pageIDStr + "` set reply=concat(reply,?), replyNum=replyNum+1 where id=?"
	var ret sql.Result
	ret, err = db.Exec(sqlStr, response, ID)
	if err != nil {
		zap.L().Error("community comment insert reply failed", zap.Error(err))
		return
	}
	_, err = ret.RowsAffected()
	if err != nil {
		zap.L().Error("community comment query inserted reply failed", zap.Error(err))
	}
	return
}

// CommunityCommentGetReply 得到某个帖子的某条评论的所有回复
func CommunityCommentGetReply(pageID, ID int64) (reply string, err error) {
	pageIDStr := strconv.FormatInt(pageID, 10)
	sqlStr := "select reply from `" + pageIDStr + "` where id=?"
	err = db.Get(&reply, sqlStr, ID)
	if err != nil {
		zap.L().Error("community comment get reply failed", zap.Any("pageID", pageID), zap.Any("ID", ID), zap.Error(err))
	}
	return
}

// CommunityCommentApproveAdd 增加点赞者的id
func CommunityCommentApproveAdd(pageID, ID int64, username string) (err error) {
	username = username + ","
	pageIDStr := strconv.FormatInt(pageID, 10)
	sqlStr := "update `" + pageIDStr + "` set approve=concat(approve,?) where id=?"
	var ret sql.Result
	ret, err = db.Exec(sqlStr, username, ID)
	if err != nil {
		zap.L().Error("community comment add approve failed", zap.Any("pageID", pageID), zap.Any("ID", ID), zap.Error(err))
		return
	}
	_, err = ret.RowsAffected()
	if err != nil {
		zap.L().Error("community query comment add approve modified failed", zap.Any("pageID", pageID), zap.Any("ID", ID), zap.Error(err))
	}
	return
}

// CommunityCommentModifyApprove 修改整个点赞的内容
func CommunityCommentModifyApprove(pageID, ID int64, approve string) (err error) {
	pageIDStr := strconv.FormatInt(pageID, 10)
	sqlStr := "update `" + pageIDStr + "` set approve=? where id=?"
	var ret sql.Result
	ret, err = db.Exec(sqlStr, approve, ID)
	if err != nil {
		zap.L().Error("community comment modify approve failed", zap.Any("pageID", pageID), zap.Any("ID", ID), zap.Error(err))
		return
	}
	_, err = ret.RowsAffected()
	if err != nil {
		zap.L().Error("community query comment approve modified failed", zap.Any("pageID", pageID), zap.Any("ID", ID), zap.Error(err))
	}
	return
}

// CommunityCommentGetApprove 得到点赞的内容
func CommunityCommentGetApprove(pageID, ID int64) (approve string, err error) {
	pageIDStr := strconv.FormatInt(pageID, 10)
	sqlStr := "select approve from `" + pageIDStr + "` where id=?"
	err = db.Get(&approve, sqlStr, ID)
	if err != nil {
		zap.L().Error("community comment get approve failed", zap.Any("pageID", pageID), zap.Any("ID", ID), zap.Error(err))
	}
	return
}

//func int64ToBytes(i int64) []byte {
//	var buf = make([]byte, 8)
//	binary.BigEndian.PutUint64(buf, uint64(i))
//	return buf
//}
