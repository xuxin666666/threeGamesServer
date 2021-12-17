package logic

import (
	"03.threeGames/dao/mysql"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func MineSweepGetScores(uid int64) (scores string, err error) {
	var exist bool
	exist, err = mysql.MineSweepJudgeExist(uid)
	if err != nil {
		return
	}
	if exist {
		return mysql.MineSweepGetScores(uid)
	} else {
		scores = ""
		return
	}
}

func MineSweepUpdateScores(uid int64, scores int) (err error) {
	// 这里的分数是时间，故数值越小则应排的越前
	var exist bool
	var resScores string
	exist, err = mysql.MineSweepJudgeExist(uid) // 判断该用户是否存过得分
	if err != nil {
		return
	}

	newScore := strconv.Itoa(scores)
	dat := strconv.FormatInt(time.Now().UnixMilli(), 10)
	newScore = newScore + "`" + dat + "```"
	if exist {
		resScores, err = mysql.MineSweepGetScores(uid)
		if err != nil {
			return
		}
		scoreArr := strings.Split(resScores, "```") // 分割后一定会多一个空字符串
		scoreArrLen := len(scoreArr)
		var newScoresBuilder strings.Builder
		var temp []string
		var tempInt, flag int
		flag = -1
		if scoreArrLen == 1 { // 如果记录为空
			err = mysql.MineSweepUpdateScores(uid, newScore)
			if err != nil {
				return
			}
		} else if scoreArrLen < 6 { // 如果该用户的得分记录只有不到5个, <=4
			for i := 0; i < scoreArrLen-1; i++ {
				temp = strings.Split(scoreArr[i], "`")
				tempInt, err = strconv.Atoi(temp[0])
				if err != nil {
					zap.L().Error("tetris update scores, parse string into int failed", zap.Error(err))
					return
				}
				if flag == -1 && scores < tempInt {
					flag = 1
					newScoresBuilder.WriteString(newScore)
				}
				newScoresBuilder.WriteString(scoreArr[i])
				newScoresBuilder.WriteString("```")
			}
			if flag == -1 { // 说明这次得分比之前的都低
				newScoresBuilder.WriteString(newScore)
			}
			newScores := newScoresBuilder.String()
			err = mysql.MineSweepUpdateScores(uid, newScores)
			if err != nil {
				return
			}
		} else { // 最多存5条数据，保留得分最高的5个
			for i := 0; i < scoreArrLen-1; i++ {
				temp = strings.Split(scoreArr[i], "`")
				tempInt, err = strconv.Atoi(temp[0])
				if err != nil {
					zap.L().Error("tetris update scores, parse string into int failed", zap.Error(err))
					return
				}
				if flag == -1 && scores < tempInt {
					flag = 1
					newScoresBuilder.WriteString(newScore)
				}
				if flag == -1 || i != scoreArrLen - 2 { // 前面插入了新的得分，那排最后的得分不要了
					newScoresBuilder.WriteString(scoreArr[i])
					newScoresBuilder.WriteString("```")
				}
			}
			newScores := newScoresBuilder.String()
			err = mysql.MineSweepUpdateScores(uid, newScores)
			if err != nil {
				return
			}
		}
	} else {
		err = mysql.MineSweepInsertScores(uid, newScore)
		if err != nil {
			return
		}
	}
	return
}

func MineSweepUpdateScoreList(uid int64, scores string) (resScores string, err error) {
	var exist bool
	exist, err = mysql.MineSweepJudgeExist(uid) // 判断该用户是否存过得分
	if err != nil {
		return
	}

	var scoreArr [][2]int
	scoreArr, err = mHandleScores(scores)
	if err != nil {
		return
	}

	if exist {
		resScores, err = mysql.MineSweepGetScores(uid)
		if err != nil {
			return
		}

		resScoresArr := strings.Split(resScores, "```")
		resScoresArrLen := len(resScoresArr) - 1

		if resScoresArrLen == 0 {
			resScores = scores
		} else {
			var resArr, resultArr [][2]int
			var oneRecord [2]int
			var score, dat int
			var scoreStr []string
			scoreArrLen := len(scoreArr)
			indexPost, indexLocal := 0, 0
			for i := 0; i < resScoresArrLen; i++ {
				scoreStr = strings.Split(resScoresArr[i], "`")
				score, err = strconv.Atoi(scoreStr[0])
				if err != nil {
					zap.L().Error("tetris score parse string into int failed", zap.Error(err))
					return
				}
				dat, err = strconv.Atoi(scoreStr[1])
				if err != nil {
					zap.L().Error("tetris score parse string into int failed", zap.Error(err))
					return
				}
				oneRecord[0] = score
				oneRecord[1] = dat
				resArr = append(resArr, oneRecord)
			}
			for indexPost < scoreArrLen && indexLocal < resScoresArrLen {
				if resArr[indexLocal][1] - scoreArr[indexPost][1] < 10 * 1000 &&
					resArr[indexLocal][1] - scoreArr[indexPost][1] > -10 * 1000 {
					resultArr = append(resultArr, resArr[indexLocal])
					indexPost++
					indexLocal++
				} else if resArr[indexLocal][0] <= scoreArr[indexPost][0] {
					resultArr = append(resultArr, resArr[indexLocal])
					indexLocal++
				} else {
					resultArr = append(resultArr, scoreArr[indexPost])
					indexPost++
				}
			}
			for indexLocal < resScoresArrLen {
				resultArr = append(resultArr, resArr[indexLocal])
				indexLocal++
			}
			for indexPost < scoreArrLen {
				resultArr = append(resultArr, scoreArr[indexPost])
				indexPost++
			}

			length := 5
			if length > len(resultArr) {
				length = len(resultArr)
			}
			var build strings.Builder
			for i := 0; i < length; i++ {
				build.WriteString(strconv.Itoa(resultArr[i][0]))
				build.WriteString("`")
				build.WriteString(strconv.Itoa(resultArr[i][1]))
				build.WriteString("```")
			}
			resScores = build.String()
			err = mysql.MineSweepUpdateScores(uid, resScores)
			if err != nil {
				return
			}
		}
	} else {
		err = mysql.MineSweepInsertScores(uid, scores)
		if err != nil {
			return
		}
	}
	return
}

// 对scores进行格式检验，并根据分数从小到大排序
func mHandleScores(scores string) (scoreArr [][2]int, err error) {
	if scores == "" {
		return
	}
	var match bool
	match, err = regexp.MatchString("[0-9]+`1[0-9]+```", scores)
	if err != nil {
		zap.L().Error("tetris match scores failed", zap.Error(err))
		return
	}
	if match == false {
		err = errors.New("match == false")
		zap.L().Error("tetris match scores failed", zap.Error(err))
		return
	}
	preScoreArr := strings.Split(scores, "```")
	preScoreArrLen := len(preScoreArr) - 1
	var preScoreArrScore [][2]int
	var oneRecord [2]int
	var score, dat int
	var scoreStr []string
	for i := 0; i < preScoreArrLen; i++ {
		scoreStr = strings.Split(preScoreArr[i], "`")
		score, err = strconv.Atoi(scoreStr[0])
		if err != nil {
			zap.L().Error("tetris score parse string into int failed", zap.Error(err))
			return
		}
		dat, err = strconv.Atoi(scoreStr[1])
		if err != nil {
			zap.L().Error("tetris score parse string into int failed", zap.Error(err))
			return
		}
		oneRecord[0] = score
		oneRecord[1] = dat
		preScoreArrScore = append(preScoreArrScore, oneRecord)
	}
	for i := 0; i < preScoreArrLen; i++ {
		for j := i + 1; j < preScoreArrLen; j++ {
			if preScoreArrScore[i][0] > preScoreArrScore[j][0] {
				oneRecord = preScoreArrScore[i]
				preScoreArrScore[i] = preScoreArrScore[j]
				preScoreArrScore[j] = oneRecord
			}
		}
	}
	scoreArr = preScoreArrScore
	return
}