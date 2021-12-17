package mysql

import (
	"database/sql"
	"go.uber.org/zap"
)

func MineSweepGetScores(uid int64) (scores string, err error) {
	sqlStr := `select scores from mineSweep where user_id=?`
	err = db.Get(&scores, sqlStr, uid)
	if err != nil {
		zap.L().Error("get mineSweep scores failed", zap.Error(err))
	}
	return
}

func MineSweepUpdateScores(uid int64, scores string) (err error) {
	sqlStr := `update mineSweep set scores=? where user_id=?`
	var ret sql.Result
	ret, err = db.Exec(sqlStr, scores, uid)
	if err != nil {
		zap.L().Error("update mineSweep scores failed", zap.Error(err))
		return
	}
	_, err = ret.RowsAffected()
	if err != nil {
		zap.L().Error("mineSweep query rows failed", zap.Error(err))
	}
	return
}

func MineSweepInsertScores(uid int64, scores string) (err error) {
	sqlStr := `insert into mineSweep(user_id, scores) values(?,?)`
	var ret sql.Result
	ret, err = db.Exec(sqlStr, uid, scores)
	if err != nil {
		zap.L().Error("mineSweep insert scores failed", zap.Error(err))
		return
	}
	_, err = ret.LastInsertId()
	if err != nil {
		zap.L().Error("mineSweep get insert scores id failed", zap.Error(err))
	}
	return
}

func MineSweepJudgeExist(uid int64) (exist bool, err error) {
	sqlStr := `select count(scores) from mineSweep where user_id=?`
	var count int
	err = db.Get(&count, sqlStr, uid)
	if err != nil {
		zap.L().Error("mineSweep get scores by uid failed", zap.Error(err))
		return
	}
	exist = false
	if count > 0 {
		exist = true
	}
	return
}
