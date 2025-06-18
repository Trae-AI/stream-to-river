// Copyright (c) 2025 Bytedance Ltd. and/or its affiliates
// SPDX-License-Identifier: MIT

package dao

import (
	"errors"

	"gorm.io/gorm"

	"github.com/Trae-AI/stream-to-river/rpcservice/dal/model"
	"github.com/Trae-AI/stream-to-river/rpcservice/dal/mysql"
)

var (
	// ErrNoRecord represents the error when no record is found in the database.
	// It's initialized with gorm's ErrRecordNotFound.
	ErrNoRecord = gorm.ErrRecordNotFound
	// ErrNoUpdate represents the error when no records are updated in the database operation.
	ErrNoUpdate = errors.New("no record updated")
)

// AddWord inserts a new word record into the `word` table.
//
// Parameters:
//   - word: A pointer to the `model.Word` struct representing the word record to be added.
//
// Returns:
//   - error: An error object if an unexpected error occurs during the database operation.
func AddWord(word *model.Word) error {
	ret := mysql.GetDB().Table(model.WordTableName).Create(word)
	if ret.Error != nil {
		return ret.Error
	}
	return nil
}

// QueryWordByUserIdAndName queries the `word` table for a word record based on `user_id` and `word_name`.
//
// Parameters:
//   - userId: The unique identifier of the user.
//   - wordName: The name of the word.
//
// Returns:
//   - *model.Word: A pointer to the retrieved `Word` record. Returns `nil` if no record is found.
//   - error: An error object if an unexpected error occurs during the process.
//     Returns `ErrNoRecord` if no record matches the query.
func QueryWordByUserIdAndName(userId int64, wordName string) (*model.Word, error) {
	var word *model.Word
	ret := mysql.GetDB().Table(model.WordTableName).Where("user_id = ? AND word_name = ?", userId, wordName).First(&word)
	if ret.Error != nil {
		if ret.Error == gorm.ErrRecordNotFound {
			return nil, ErrNoRecord
		}
		return nil, ret.Error
	}
	if ret.RowsAffected == 0 {
		// actually, if RowsAffected == 0, then will get ErrRecordNotFound
		return nil, ErrNoRecord
	}
	return word, nil
}

// GetWordsByUserIdWithPagination queries the `word` table for words belonging to a user with pagination.
//
// Parameters:
//   - userId: The unique identifier of the user.
//   - offset: The number of records to skip.
//   - limit: The maximum number of records to return.
//
// Returns:
//   - []*model.Word: A slice of pointers to the retrieved `Word` records.
//   - error: An error object if an unexpected error occurs during the process.
func GetWordsByUserIdWithPagination(userId int64, offset int64, limit int64) ([]*model.Word, error) {
	var words []*model.Word

	ret := mysql.GetDB().Table(model.WordTableName).
		Where("user_id = ?", userId).
		Order("word_id DESC"). // 按word_id倒序，等同于按添加时间倒序
		Offset(int(offset)).
		Limit(int(limit)).
		Find(&words)

	if ret.Error != nil {
		return nil, ret.Error
	}
	return words, nil
}

// GetTotalWordsCount retrieves the total number of words for a user from the `word` table.
//
// Parameters:
//   - userId: The unique identifier of the user.
//
// Returns:
//   - int32: The total number of words.
//   - error: An error object if an unexpected error occurs during the process.
func GetTotalWordsCount(userId int64) (int32, error) {
	var count int64
	ret := mysql.GetDB().Table(model.WordTableName).
		Where("user_id = ?", userId).
		Count(&count)

	if ret.Error != nil {
		return 0, ret.Error
	}
	return int32(count), nil
}

// UpdateWordTagID updates the `tag_id` for a word record.
//
// Parameters:
//   - wordID: The unique identifier of the word.
//   - userID: The unique identifier of the user.
//   - tagID: The new `tag_id`.
//
// Returns:
//   - error: An error object if an unexpected error occurs during the process.
//     Returns `ErrNoUpdate` if no records are updated.
func UpdateWordTagID(wordID, userID int64, tagID int32) error {
	ret := mysql.GetDB().Table(model.WordTableName).Where(map[string]interface{}{
		"user_id": userID,
		"word_id": wordID,
	}).Updates(map[string]interface{}{"tag_id": tagID})
	if ret.Error != nil {
		return ret.Error
	}
	if ret.RowsAffected == 0 {
		return ErrNoUpdate
	}
	return nil
}

// AddWordWithRelatedRecords atomically adds a word and its related records in a transaction.
//
// Parameters:
//   - word: A pointer to the `model.Word` struct representing the word record.
//   - answerList: A pointer to the `model.AnswerList` struct representing the related answer list record.
//   - reviewRecord: A pointer to the `model.WordsRisiteRecord` struct representing the related review record.
//
// Returns:
//   - *model.Word: A pointer to the added `Word` record.
//   - error: An error object if an unexpected error occurs during the process.
func AddWordWithRelatedRecords(word *model.Word, answerList *model.AnswerList, reviewRecord *model.WordsRisiteRecord) (*model.Word, error) {
	var result *model.Word

	err := WithTransaction(func(tx *gorm.DB) error {
		// 1. 添加单词
		if err := tx.Table(model.WordTableName).Create(word).Error; err != nil {
			return err
		}

		// 2. 查询刚插入的单词获取完整信息（包括auto-generated ID）
		var queriedWord model.Word
		if err := tx.Table(model.WordTableName).
			Where("user_id = ? AND word_name = ?", word.UserId, word.WordName).
			First(&queriedWord).Error; err != nil {
			return err
		}
		result = &queriedWord

		// 3. 使用查询到的word_id更新相关记录
		answerList.WordId = queriedWord.WordId
		reviewRecord.WordId = int(queriedWord.WordId)

		// 4. 计算并设置answer_list的order_id
		var maxOrderId int
		if err := tx.Table(model.AnswerListTableName).
			Where("user_id = ?", answerList.UserId).
			Select("COALESCE(MAX(order_id), 0) as max_order_id").
			Scan(&maxOrderId).Error; err != nil {
			return err
		}
		answerList.OrderId = maxOrderId + 1

		// 5. 添加answer_list记录
		if err := tx.Table(model.AnswerListTableName).Create(answerList).Error; err != nil {
			return err
		}

		// 6. 添加复习记录
		if err := tx.Table(model.WordsRisiteRecordTableName).Create(reviewRecord).Error; err != nil {
			return err
		}

		return nil
	})

	return result, err
}

// QueryWord queries the `word` table for a word record based on `word_id`.
//
// Parameters:
//   - wordId: The unique identifier of the word.
//
// Returns:
//   - *model.Word: A pointer to the retrieved `Word` record.
//   - error: An error object if an unexpected error occurs during the process.
//     Returns `ErrNoRecord` if no record matches the query.
func QueryWord(wordId int64) (*model.Word, error) {
	var word *model.Word
	ret := mysql.GetDB().Table(model.WordTableName).Where("word_id = ?", wordId).Find(&word)
	if ret.Error != nil {
		return nil, ret.Error
	}
	if ret.RowsAffected == 0 {
		return nil, ErrNoRecord
	}
	return word, nil
}

// GetWordsByIds batch - queries the `word` table for word records based on a list of `word_ids`.
//
// Parameters:
//   - wordIds: A slice of integers representing the `word_ids`.
//
// Returns:
//   - []*model.Word: A slice of pointers to the retrieved `Word` records.
//   - error: An error object if an unexpected error occurs during the process.
func GetWordsByIds(wordIds []int64) ([]*model.Word, error) {
	var words []*model.Word

	ret := mysql.GetDB().Table(model.WordTableName).
		Where("word_id IN ?", wordIds).
		Find(&words)

	if ret.Error != nil {
		return nil, ret.Error
	}

	return words, nil
}

// DelWord deletes a word record from the `word` table based on `word_id`.
//
// Parameters:
//   - wordId: The unique identifier of the word.
//
// Returns:
//   - error: An error object if an unexpected error occurs during the process.
func DelWord(wordId int64) error {
	ret := mysql.GetDB().Table(model.WordTableName).Where("word_id = ?", wordId).Delete(&model.Word{})
	if ret.Error != nil {
		return ret.Error
	}
	return nil
}

// DelUserByID deletes all word records for a user from the `word` table.
//
// Parameters:
//   - userID: The unique identifier of the user.
//
// Returns:
//   - error: An error object if an unexpected error occurs during the process.
func DelUserByID(userID int64) error {
	ret := mysql.GetDB().Table(model.WordTableName).Where("user_id = ?", userID).Delete(&model.Word{})
	if ret.Error != nil {
		return ret.Error
	}
	return nil
}

// WithTransaction executes a function within a database transaction.
// If the function returns an error, the transaction is rolled back.
// If the function executes successfully, the transaction is committed.
//
// Parameters:
//   - fn: A function that takes a `*gorm.DB` as a parameter and returns an error.
//
// Returns:
//   - error: An error object if an unexpected error occurs during the transaction.
func WithTransaction(fn func(tx *gorm.DB) error) (err error) {
	tx := mysql.GetDB().Begin()
	if tx.Error != nil {
		return tx.Error
	}

	// Defer rollback in case of panic or error
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			// Convert panic to error instead of re-panicking
			if e, ok := r.(error); ok {
				err = e
			} else {
				err = errors.New("transaction failed due to panic")
			}
		}
	}()

	// Execute the function
	if err := fn(tx); err != nil {
		tx.Rollback()
		return err
	}

	// Commit the transaction
	return tx.Commit().Error
}
