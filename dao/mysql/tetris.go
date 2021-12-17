package mysql

import (
	"database/sql"
	"go.uber.org/zap"
)

func TetrisGetScores(uid int64) (scores string, err error) {
	sqlStr := `select scores from tetris where user_id=?`
	err = db.Get(&scores, sqlStr, uid)
	if err != nil {
		zap.L().Error("get tetris scores failed", zap.Error(err))
	}
	return
}

func TetrisUpdateScores(uid int64, scores string) (err error) {
	sqlStr := `update tetris set scores=? where user_id=?`
	var ret sql.Result
	ret, err = db.Exec(sqlStr, scores, uid)
	if err != nil {
		zap.L().Error("update tetris scores failed", zap.Error(err))
		return
	}
	_, err = ret.RowsAffected()
	if err != nil {
		zap.L().Error("tetris query rows failed", zap.Error(err))
	}
	return
}

func TetrisInsertScores(uid int64, scores string) (err error) {
	sqlStr := `insert into tetris(user_id, scores) values(?,?)`
	var ret sql.Result
	ret, err = db.Exec(sqlStr, uid, scores)
	if err != nil {
		zap.L().Error("tetris insert scores failed", zap.Error(err))
		return
	}
	_, err = ret.LastInsertId()
	if err != nil {
		zap.L().Error("tetris get insert scores id failed", zap.Error(err))
	}
	return
}

func TetrisJudgeExist(uid int64) (exist bool, err error) {
	sqlStr := `select count(scores) from tetris where user_id=?`
	var count int
	err = db.Get(&count, sqlStr, uid)
	if err != nil {
		zap.L().Error("tetris get scores by uid failed", zap.Error(err))
		return
	}
	exist = false
	if count > 0 {
		exist = true
	}
	return
}
