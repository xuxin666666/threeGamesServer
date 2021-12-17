package models

type Tetris struct {
	UserID string `json:"user_id" db:"user_id"`
	Scores string `json:"scores" db:"scores"`
}

type TetrisScores struct {
	Scores *int `json:"scores" binding:"required"`
}

type TetrisScoreList struct {
	Scores *string `json:"scores" binding:"required"`
}
