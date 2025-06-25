// Copyright (c) 2025 Bytedance Ltd. and/or its affiliates
// SPDX-License-Identifier: MIT

package dao

import (
	"github.com/cloudwego/kitex/pkg/klog"

	"github.com/Trae-AI/stream-to-river/rpcservice/dal/model"
	"github.com/Trae-AI/stream-to-river/rpcservice/dal/mysql"
)

// AddWordsReciteRecord inserts a new review record into the `words_recite_record` table.
//
// Parameters:
//   - record: A pointer to the `model.WordsReciteRecord` struct that represents the review record to be added.
//
// Returns:
//   - error: An error object if an unexpected error occurs during the database operation.
func AddWordsReciteRecord(record *model.WordsReciteRecord) error {
	ret := mysql.GetDB().Table(model.WordsReciteRecordTableName).Create(record)
	if ret.Error != nil {
		return ret.Error
	}
	klog.Infof("Insert words_recite_record=%v into table=%s", record, model.WordsReciteRecordTableName)
	return nil
}

// GetWordsReciteRecord queries the `words_recite_record` table for a review record using `user_id` and `word_id`.
//
// Parameters:
//   - userId: The unique identifier of the user.
//   - wordId: The unique identifier of the word.
//
// Returns:
//   - *model.WordsReciteRecord: A pointer to the retrieved `WordsReciteRecord` if found.
//   - error: An error object if an unexpected error occurs during the database operation.
func GetWordsReciteRecord(userId int64, wordId int64) (*model.WordsReciteRecord, error) {
	var record model.WordsReciteRecord

	ret := mysql.GetDB().Table(model.WordsReciteRecordTableName).
		Where("user_id = ? AND word_id = ?", userId, wordId).
		First(&record)

	if ret.Error != nil {
		return nil, ret.Error
	}

	klog.Infof("Found words_recite_record for user_id=%v, word_id=%v", userId, wordId)
	return &record, nil
}

// GetWordsReciteRecordsByUserAndWordIds performs a batch query on the `words_recite_record` table.
// It fetches review records based on `user_id` and a list of `word_ids`.
//
// Parameters:
//   - userId: The unique identifier of the user.
//   - wordIds: A slice of `int64` representing the `word_ids`.
//
// Returns:
//   - map[int64]*model.WordsReciteRecord: A map where the key is `word_id` and the value is a pointer to the retrieved `WordsReciteRecord`.
//   - error: An error object if an unexpected error occurs during the database operation.
func GetWordsReciteRecordsByUserAndWordIds(userId int64, wordIds []int64) (map[int64]*model.WordsReciteRecord, error) {
	var records []*model.WordsReciteRecord

	ret := mysql.GetDB().Table(model.WordsReciteRecordTableName).
		Where("user_id = ? AND word_id IN ?", userId, wordIds).
		Find(&records)

	if ret.Error != nil {
		return nil, ret.Error
	}

	// Convert the slice of records to a map with word_id as the key
	recordMap := make(map[int64]*model.WordsReciteRecord)
	for _, record := range records {
		recordMap[int64(record.WordId)] = record
	}

	klog.Infof("Found %d words_recite_records for user_id=%v", len(records), userId)
	return recordMap, nil
}

// UpdateWordsReciteRecord updates an existing review record in the `words_recite_record` table.
// It updates specific fields (`Level`, `NextReviewTime`, `TotalWrong`, `TotalCorrect`, `Score`) of the record.
//
// Parameters:
//   - record: A pointer to the `model.WordsReciteRecord` struct that represents the updated review record.
//
// Returns:
//   - error: An error object if an unexpected error occurs during the database operation.
func UpdateWordsReciteRecord(record *model.WordsReciteRecord) error {
	ret := mysql.GetDB().Table(model.WordsReciteRecordTableName).
		Where("user_id = ? AND word_id = ?", record.UserId, record.WordId).
		Select("Level", "NextReviewTime", "TotalWrong", "TotalCorrect", "Score").
		Updates(record)

	if ret.Error != nil {
		return ret.Error
	}

	klog.Infof("Updated words_recite_record: %+v", record)
	return nil
}

// GetReviewRecords queries the `words_recite_record` table for review records that need to be reviewed.
// It filters records where the `next_review_time` is less than or equal to the `currentTime`.
//
// Parameters:
//   - userId: The unique identifier of the user.
//   - currentTime: The current timestamp used for filtering.
//
// Returns:
//   - []*model.WordsReciteRecord: A slice of pointers to the retrieved `WordsReciteRecord` that meet the criteria.
//   - error: An error object if an unexpected error occurs during the database operation.
func GetReviewRecords(userId int64, currentTime int64) ([]*model.WordsReciteRecord, error) {
	var records []*model.WordsReciteRecord

	ret := mysql.GetDB().Table(model.WordsReciteRecordTableName).
		Where("user_id = ? AND next_review_time <= ?", userId, currentTime).
		Find(&records)

	if ret.Error != nil {
		return nil, ret.Error
	}

	klog.Infof("Found %d review records for user_id=%v", len(records), userId)
	return records, nil
}

// GetCompletedWordsCountFromRecord retrieves the number of completed review words for a user.
// A word is considered completed if its `level` is greater than or equal to 8.
//
// Parameters:
//   - userId: The unique identifier of the user.
//
// Returns:
//   - int32: The number of completed review words.
//   - error: An error object if an unexpected error occurs during the database operation.
func GetCompletedWordsCountFromRecord(userId int64) (int32, error) {
	var count int64
	ret := mysql.GetDB().Table(model.WordsReciteRecordTableName).
		Where("user_id = ? AND level >= ?", userId, 8).
		Count(&count)

	if ret.Error != nil {
		return 0, ret.Error
	}

	klog.Infof("Completed words count from record for user_id=%v: %d", userId, count)
	return int32(count), nil
}

// DelWordsReciteRecordByUserID deletes all review records for a specific user from the `words_recite_record` table.
//
// Parameters:
//   - userId: The unique identifier of the user.
//
// Returns:
//   - error: An error object if an unexpected error occurs during the database operation.
func DelWordsReciteRecordByUserID(userId int64) error {
	ret := mysql.GetDB().Table(model.WordsReciteRecordTableName).
		Where("user_id = ?", userId).
		Delete(&model.WordsReciteRecord{})

	if ret.Error != nil {
		return ret.Error
	}

	klog.Infof("Deleted words_recite_record for user_id=%v", userId)
	return nil
}
