// Copyright (c) 2025 Bytedance Ltd. and/or its affiliates
// SPDX-License-Identifier: MIT

package model

// WordTagTableName is the name of the database table corresponding to the WordTag model.
const WordTagTableName = "word_tags"

// WordTag represents a record in the `word_tag` database table.
type WordTag struct {
	TagId         int64  `gorm:"column:id;primaryKey;autoIncrement"`        // Primary key, auto - incremented.
	TagName       string `gorm:"column:tag_name;not null;unique"`           // Name of the word tag, cannot be null and must be unique.
	QuestionTypes int    `gorm:"column:question_types;not null;default:15"` // Bit - operation stored question type combinations, default value is 15.
	MaxScore      int    `gorm:"column:max_score;not null;default:0"`       // Maximum score, cannot be null and default value is 0.
}

// TableName returns the name of the database table for the WordTag model.
func (w WordTag) TableName() string {
	return WordTagTableName
}
