// Copyright (c) 2025 Bytedance Ltd. and/or its affiliates
// SPDX-License-Identifier: MIT

package dao

import (
	"time"

	"github.com/cloudwego/kitex/pkg/klog"
	"gorm.io/gorm"

	"github.com/Trae-AI/stream-to-river/rpcservice/dal/model"
	"github.com/Trae-AI/stream-to-river/rpcservice/dal/mysql"
)

// UpdateCompletedReviewCount updates the `completed_review_count` field for a specific user in the `review_progress` table.
// It increments the existing `completed_review_count` by the specified `increment` value.
//
// Parameters:
//   - userId: The unique identifier of the user whose `completed_review_count` needs to be updated.
//   - increment: The amount by which to increase the `completed_review_count`.
//
// Returns:
//   - error: An error object if an unexpected error occurs during the database operation.
func UpdateCompletedReviewCount(userId int64, increment int) error {
	ret := mysql.GetDB().Table(model.ReviewProgressTableName).
		Where("user_id = ?", userId).
		Update("completed_review_count", mysql.GetDB().Raw("completed_review_count + ?", increment))

	if ret.Error != nil {
		return ret.Error
	}

	klog.Infof("Updated completed review count for user_id=%v, increment=%v", userId, increment)
	return nil
}

// IncrementPendingReviewCount increments the `pending_review_count` field for a specific user in the `review_progress` table.
// It also updates the `last_update_time` field to the current Unix timestamp.
//
// Parameters:
//   - userId: The unique identifier of the user whose `pending_review_count` needs to be incremented.
//
// Returns:
//   - error: An error object if an unexpected error occurs during the database operation.
func IncrementPendingReviewCount(userId int64) error {
	currentTime := time.Now().Unix()
	ret := mysql.GetDB().Table(model.ReviewProgressTableName).
		Where("user_id = ?", userId).
		Updates(map[string]interface{}{
			"pending_review_count": mysql.GetDB().Raw("pending_review_count + 1"),
			"last_update_time":     currentTime,
		})

	if ret.Error != nil {
		return ret.Error
	}

	klog.Infof("Incremented pending review count for user_id=%v", userId)
	return nil
}

// IncrementAllCompletedCount increments the `all_completed_count` field for a specific user in the `review_progress` table.
//
// Parameters:
//   - userId: The unique identifier of the user whose `all_completed_count` needs to be incremented.
//
// Returns:
//   - error: An error object if an unexpected error occurs during the database operation.
func IncrementAllCompletedCount(userId int64) error {
	ret := mysql.GetDB().Table(model.ReviewProgressTableName).
		Where("user_id = ?", userId).
		Update("all_completed_count", gorm.Expr("all_completed_count + ?", 1))

	if ret.Error != nil {
		return ret.Error
	}

	klog.Infof("Incremented all_completed_count for user_id=%v", userId)
	return nil
}

// GetReviewProgress retrieves the review progress record for a specific user from the `review_progress` table.
//
// Parameters:
//   - userId: The unique identifier of the user whose review progress needs to be retrieved.
//
// Returns:
//   - *model.ReviewProgress: A pointer to the retrieved `ReviewProgress` record.
//   - error: An error object if an unexpected error occurs during the database operation.
func GetReviewProgress(userId int64) (*model.ReviewProgress, error) {
	var progress model.ReviewProgress
	ret := mysql.GetDB().Table(model.ReviewProgressTableName).
		Where("user_id = ?", userId).
		First(&progress)

	if ret.Error != nil {
		return nil, ret.Error
	}

	klog.Infof("Found review progress for user_id=%v: %+v", userId, progress)
	return &progress, nil
}

// CreateOrUpdateReviewProgress creates a new review progress record or updates an existing one for a specific user.
// It checks if a record exists for the user based on `user_id` and then updates the relevant fields.
// If no record exists, it creates a new one.
//
// Parameters:
//   - progress: A pointer to the `model.ReviewProgress` struct representing the review progress to be created or updated.
//
// Returns:
//   - error: An error object if an unexpected error occurs during the database operation.
func CreateOrUpdateReviewProgress(progress *model.ReviewProgress) error {
	ret := mysql.GetDB().Table(model.ReviewProgressTableName).
		Where("user_id = ?", progress.UserId).
		Assign(map[string]interface{}{
			"pending_review_count":   progress.PendingReviewCount,
			"completed_review_count": progress.CompletedReviewCount,
			"last_update_time":       progress.LastUpdateTime,
		}).
		FirstOrCreate(progress)

	if ret.Error != nil {
		return ret.Error
	}

	klog.Infof("Created/Updated review progress for user_id=%v: %+v", progress.UserId, progress)
	return nil
}

// DelReviewProgressByUserID deletes the review progress record for a specific user from the `review_progress` table.
//
// Parameters:
//   - userId: The unique identifier of the user whose review progress record needs to be deleted.
//
// Returns:
//   - error: An error object if an unexpected error occurs during the database operation.
func DelReviewProgressByUserID(userId int64) error {
	ret := mysql.GetDB().Table(model.ReviewProgressTableName).
		Where("user_id = ?", userId).
		Delete(&model.ReviewProgress{})

	if ret.Error != nil {
		return ret.Error
	}

	klog.Infof("Deleted review progress for user_id=%v", userId)
	return nil
}
