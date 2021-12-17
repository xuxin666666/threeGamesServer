package models

type MineSweep struct {
	UserID string `json:"user_id" db:"user_id"`
	Scores string `json:"scores" db:"scores"`
}

type MineSweepScores struct {
	Scores *int `json:"scores" binding:"required"`
}

type MineSweepScoreList struct {
	Scores *string `json:"scores" binding:"required"`
}
