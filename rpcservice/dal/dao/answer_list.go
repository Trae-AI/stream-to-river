// Copyright (c) 2025 Bytedance Ltd. and/or its affiliates
// SPDX-License-Identifier: MIT

package dao

import (
	"github.com/Trae-AI/stream-to-river/rpcservice/dal/model"
	"github.com/Trae-AI/stream-to-river/rpcservice/dal/mysql"
	"github.com/cloudwego/kitex/pkg/klog"
)

// AddAnswerList adds a new record to the `answer_list` table.
// It first fetches the maximum `order_id` for the specified user.
// Then it sets the `order_id` of the new record to the maximum `order_id` plus one.
// If there are no existing records for the user, the `order_id` is set to 1.
// Finally, it inserts the new record into the `answer_list` table.
//
// Parameters:
//   - answerList: A pointer to the `model.AnswerList` struct representing the record to be added.
//
// Returns:
//   - error: An error object if an unexpected error occurs during the process.
func AddAnswerList(answerList *model.AnswerList) error {
	// Query the maximum order_id for the user
	var maxOrderId int
	ret := mysql.GetDB().Table(model.AnswerListTableName).Where("user_id = ?", answerList.UserId).Select("COALESCE(MAX(order_id), 0) as max_order_id").Scan(&maxOrderId)
	if ret.Error != nil {
		return ret.Error
	}

	// Set the new record's order_id to the maximum value plus one. If no records exist, set it to 1.
	answerList.OrderId = maxOrderId + 1

	// Insert the record
	ret = mysql.GetDB().Table(model.AnswerListTableName).Create(answerList)
	if ret.Error != nil {
		return ret.Error
	}
	klog.Infof("Insert answer_list=%v into table=%s with order_id=%d", answerList, model.AnswerListTableName, answerList.OrderId)
	return nil
}

// GetAnswerListByAnswerId queries the `answer_list` table for a record based on the `user_id` and `answer_id`.
//
// Parameters:
//   - userId: The unique identifier of the user.
//   - answerId: The unique identifier of the answer.
//
// Returns:
//   - *model.AnswerList: A pointer to the retrieved `AnswerList` record.
//   - error: An error object if an unexpected error occurs during the process.
func GetAnswerListByAnswerId(userId int64, answerId int64) (*model.AnswerList, error) {
	var answerList model.AnswerList

	ret := mysql.GetDB().Table(model.AnswerListTableName).
		Where("answer_id = ? AND user_id = ?", answerId, userId).
		First(&answerList)

	if ret.Error != nil {
		return nil, ret.Error
	}

	klog.Infof("Found answer_list record for user_id=%v, answer_id=%v", userId, answerId)
	return &answerList, nil
}

// GetAnswerListByOrderIds queries the `answer_list` table for records based on the `user_id` and a list of `order_ids`.
//
// Parameters:
//   - userId: The unique identifier of the user.
//   - orderIds: A slice of integers representing the `order_ids`.
//
// Returns:
//   - []*model.AnswerList: A slice of pointers to the retrieved `AnswerList` records.
//   - error: An error object if an unexpected error occurs during the process.
func GetAnswerListByOrderIds(userId int64, orderIds []int) ([]*model.AnswerList, error) {
	var answerList []*model.AnswerList

	ret := mysql.GetDB().Table(model.AnswerListTableName).
		Where("user_id = ? AND order_id IN ?", userId, orderIds).
		Find(&answerList)

	if ret.Error != nil {
		return nil, ret.Error
	}

	return answerList, nil
}

// GetMaxOrderId retrieves the maximum `order_id` for a user from the `answer_list` table.
//
// Parameters:
//   - userId: The unique identifier of the user.
//
// Returns:
//   - int: The maximum `order_id` for the user.
//   - error: An error object if an unexpected error occurs during the process.
func GetMaxOrderId(userId int64) (int, error) {
	var maxOrderId int
	ret := mysql.GetDB().Table(model.AnswerListTableName).
		Where("user_id = ?", userId).
		Select("COALESCE(MAX(order_id), 0) as max_order_id").
		Scan(&maxOrderId)

	if ret.Error != nil {
		return 0, ret.Error
	}

	return maxOrderId, nil
}

// GetAnswerListByWordId queries the `answer_list` table for a record based on the `user_id` and `word_id`.
//
// Parameters:
//   - userId: The unique identifier of the user.
//   - wordId: The unique identifier of the word.
//
// Returns:
//   - *model.AnswerList: A pointer to the retrieved `AnswerList` record.
//   - error: An error object if an unexpected error occurs during the process.
func GetAnswerListByWordId(userId int64, wordId int64) (*model.AnswerList, error) {
	var answerList model.AnswerList

	ret := mysql.GetDB().Table(model.AnswerListTableName).
		Where("user_id = ? AND word_id = ?", userId, wordId).
		First(&answerList)

	if ret.Error != nil {
		return nil, ret.Error
	}

	klog.Infof("Found answer_list record for user_id=%v, word_id=%v", userId, wordId)
	return &answerList, nil
}

// DelAnswerListByUserID deletes all `answer_list` records for a user.
//
// Parameters:
//   - userId: The unique identifier of the user.
//
// Returns:
//   - error: An error object if an unexpected error occurs during the process.
func DelAnswerListByUserID(userId int64) error {
	ret := mysql.GetDB().Table(model.AnswerListTableName).
		Where("user_id = ?", userId).
		Delete(&model.AnswerList{})

	if ret.Error != nil {
		return ret.Error
	}

	klog.Infof("Deleted answer_list records for user_id=%v", userId)
	return nil
}
