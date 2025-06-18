// Copyright (c) 2025 Bytedance Ltd. and/or its affiliates
// SPDX-License-Identifier: MIT

package model

// WordsRisiteRecordTableName is the name of the database table corresponding to the WordsRisiteRecord model.
const WordsRisiteRecordTableName = "words_risite_record"

// WordsRisiteRecord represents a record in the `words_risite_record` database table.
type WordsRisiteRecord struct {
	Id             int64 `gorm:"column:id;primaryKey;autoIncrement"` // Primary key, auto - incremented.
	WordId         int   `gorm:"column:word_id"`                     // Unique identifier for the word.
	Level          int   `gorm:"column:level"`                       // Level of the word review.
	NextReviewTime int64 `gorm:"column:next_review_time"`            // Timestamp for the next review.
	DowngradeStep  int   `gorm:"column:downgrade_step"`              // Downgrade step for the word review.
	TotalCorrect   int   `gorm:"column:total_correct"`               // Total number of correct answers.
	TotalWrong     int   `gorm:"column:total_wrong"`                 // Total number of wrong answers.
	Score          int   `gorm:"column:score"`                       // Score of the word review.
	UserId         int64 `gorm:"column:user_id"`                     // Unique identifier for the user.
}

// TableName returns the name of the database table for the WordsRisiteRecord model.
func (w WordsRisiteRecord) TableName() string {
	return WordsRisiteRecordTableName
}
