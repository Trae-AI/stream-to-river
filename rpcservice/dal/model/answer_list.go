// Copyright (c) 2025 Bytedance Ltd. and/or its affiliates
// SPDX-License-Identifier: MIT

package model

// AnswerListTableName is the name of the database table corresponding to the AnswerList model.
const AnswerListTableName = "answer_list"

// AnswerList represents a record in the `answer_list` database table.
type AnswerList struct {
	AnswerId    int64  `gorm:"column:answer_id;primaryKey;autoIncrement"` // Primary key, auto - incremented.
	WordId      int64  `gorm:"column:word_id"`                            // Unique identifier for the associated word.
	UserId      int64  `gorm:"column:user_id"`                            // Unique identifier for the associated user.
	WordName    string `gorm:"column:word_name"`                          // Name of the word.
	Description string `gorm:"column:description"`                        // Description of the answer.
	OrderId     int    `gorm:"column:order_id"`                           // Order ID of the answer.
}

// TableName returns the name of the database table for the AnswerList model.
func (a AnswerList) TableName() string {
	return AnswerListTableName
}
