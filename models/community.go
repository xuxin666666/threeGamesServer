package models

type CreatePageSql struct {
	UserId     int64  `json:"user_id" db:"user_id"`
	PageId     int64  `json:"page_id" db:"page_id"`
	CreateTime int64  `json:"create_time" db:"create_time"`
	Title      string `json:"title" db:"title"`
	Content    string `json:"content" db:"content"`
}

type CreatePage struct {
	UserId  int64  `json:"user_id"`
	Title   string `json:"title" binding:"required"`
	Content string `json:"content" binding:"required"`
}

type ModifyPageSql struct {
	UserID     int64  `json:"user_id" db:"user_id"`
	PageID     int64  `json:"page_id" db:"page_id"`
	UpdateTime int64  `json:"update_time" db:"update_time"`
	Content    string `json:"content" db:"content"`
}

type ModifyPageContent struct {
	UserID  int64  `json:"user_id" db:"user_id"`
	PageID  int64  `json:"page_id" db:"page_id"`
	Content string `json:"content" db:"content"`
}

type ModifyPageViews struct {
	UserID int64 `json:"user_id" db:"user_id"`
	PageID int64 `json:"page_id" db:"page_id"`
	Views  int   `json:"views" db:"views"`
}

type GetPageSql struct {
	Views      int    `json:"views" db:"views"`
	PageID     int64  `json:"page_id" db:"page_id"`
	UserID     int64  `json:"user_id" db:"user_id"`
	CreateTime int64  `json:"create_time" db:"create_time"`
	Title      string `json:"title" db:"title"`
	PreContent string `json:"pre_content" db:"pre_content"`
	Approve    string `json:"approve" db:"approve"`
}

type GetPage struct {
	Views       int    `json:"views"`
	CommentsNum int    `json:"comments_num"`
	CreateTime  int64  `json:"create_time"`
	PageID      string `json:"page_id"`
	Username    string `json:"username"`
	UserAvatar  string `json:"user_avatar"`
	Title       string `json:"title"`
	PreContent  string `json:"pre_content"`
	Approve     string `json:"approve"`
}

type QueryPage struct {
	Start  int `json:"start" binding:"required"`
	Length int `json:"length" binding:"required"`
}

type QueryPageModify struct {
	Content string `json:"content" binding:"required"`
}

type QueryPageApprove struct {
	Add *bool `json:"add" binding:"required"`
}

type GetPageDetailSql struct {
	Views      int    `json:"views" db:"views"`
	PageID     int64  `json:"page_id" db:"page_id"`
	UserID     int64  `json:"user_id" db:"user_id"`
	CreateTime int64  `json:"create_time" db:"create_time"`
	UpdateTime int64  `json:"update_time" db:"update_time"`
	Title      string `json:"title" db:"title"`
	Content    string `json:"content" db:"content"`
	Approve    string `json:"approve" db:"approve"`
}

type GetPageDetail struct {
	Views       int    `json:"views"`
	CommentsNum int    `json:"comments_num"`
	CreateTime  int64  `json:"create_time" db:"create_time"`
	UpdateTime  int64  `json:"update_time"`
	PageID      string `json:"page_id"`
	Username    string `json:"username"`
	UserAvatar  string `json:"user_avatar"`
	Title       string `json:"title"`
	Content     string `json:"content"`
	Approve     string `json:"approve"`
}

type CreateComment struct {
	PageID  int64  `json:"page_id"`
	UserID  int64  `json:"user_id"`
	Content string `json:"content"`
}

type QueryCommentAdd struct {
	Content string `json:"content" binding:"required"`
}

type QueryComment struct {
	Start   int  `json:"start" binding:"required"`
	Length  int  `json:"length" binding:"required"`
	Reverse *bool `json:"reverse" binding:"required"`
}

type QueryCommentApprove struct {
	ID  int `json:"id" binding:"required"`
	Add *bool   `json:"add" binding:"required"`
}

type GetCommentSql struct {
	ID         int    `json:"id" db:"id"`
	CreateTime int    `json:"create_time" db:"create_time"`
	UserID     int64  `json:"user_id" db:"user_id"`
	Content    string `json:"content" db:"content"`
	Approve    string `json:"approve" db:"approve"`
	Reply      string `json:"reply" db:"reply"`
}

type GetComment struct {
	CreateTime int    `json:"create_time"`
	ID         int `json:"id"`
	Content    string `json:"content"`
	Approve    string `json:"approve"`
	Reply      string `json:"reply"`
	Username   string `json:"username"`
	UserAvatar string `json:"user_avatar"`
}

type Reply struct {
	PageID     int64  `json:"page_id"`
	UserID     int64  `json:"user_id"`
	CommentID  int64  `json:"comment_id"`
	TargetName string `json:"target_name"`
	Content    string `json:"content"`
}

type QueryReply struct {
	CommentID  string `json:"comment_id" binding:"required"`
	TargetName *string `json:"target_name" binding:"required"`
	Content    string `json:"content" binding:"required"`
}
