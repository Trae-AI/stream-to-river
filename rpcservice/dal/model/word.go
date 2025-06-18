// Copyright (c) 2025 Bytedance Ltd. and/or its affiliates
// SPDX-License-Identifier: MIT

package model

// WordTableName is the name of the database table corresponding to the Word model.
const WordTableName = "words"

// Word represents a record in the `word` database table.
type Word struct {
	WordId      int64  `gorm:"column:word_id;primaryKey;autoIncrement"` // Primary key, auto - incremented.
	UserId      int64  `gorm:"column:user_id"`                          // Unique identifier for the user.
	WordName    string `gorm:"column:word_name"`                        // Name of the word.
	Description string `gorm:"column:description"`                      // Description of the word.
	Explains    string `gorm:"column:explains"`                         // Explanations of the word.
	PronounceUs string `gorm:"column:pronounce_us"`                     // American pronunciation of the word.
	PronounceUk string `gorm:"column:pronounce_uk"`                     // British pronunciation of the word.
	PhoneticUs  string `gorm:"column:phonetic_us"`                      // American phonetic symbol of the word.
	PhoneticUk  string `gorm:"column:phonetic_uk"`                      // British phonetic symbol of the word.
	YoudaoUrl   string `gorm:"column:youdao_url"`                       // Youdao dictionary URL for the word.
	TagId       int32  `gorm:"column:tag_id;default:null"`              // Associated tag ID, can be null.
}

// TableName returns the name of the database table for the Word model.
func (w Word) TableName() string {
	return WordTableName
}
