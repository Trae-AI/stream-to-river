// Copyright (c) 2025 Bytedance Ltd. and/or its affiliates
// SPDX-License-Identifier: MIT

package model

// ReviewProgressTableName is the name of the database table corresponding to the ReviewProgress model.
const ReviewProgressTableName = "review_progress"

// ReviewProgress represents a record in the `review_progress` database table.
type ReviewProgress struct {
	Id                   int64 `gorm:"column:id;primaryKey;autoIncrement"` // Primary key, auto - incremented.
	UserId               int64 `gorm:"column:user_id;uniqueIndex"`         // Unique identifier for the user, with a unique index.
	PendingReviewCount   int   `gorm:"column:pending_review_count"`        // Number of pending reviews.
	CompletedReviewCount int   `gorm:"column:completed_review_count"`      // Number of completed reviews.
	AllCompletedCount    int   `gorm:"column:all_completed_count"`         // Total number of completed words.
	LastUpdateTime       int64 `gorm:"column:last_update_time"`            // Timestamp of the last update.
}

// TableName returns the name of the database table for the ReviewProgress model.
func (r ReviewProgress) TableName() string {
	return ReviewProgressTableName
}
