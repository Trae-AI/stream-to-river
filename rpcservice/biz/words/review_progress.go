// Copyright (c) 2025 Bytedance Ltd. and/or its affiliates
// SPDX-License-Identifier: MIT

package words

import (
	"context"
	"errors"
	"time"

	"github.com/Trae-AI/stream-to-river/rpcservice/biz/words/common"
	"github.com/Trae-AI/stream-to-river/rpcservice/dal/dao"
	"github.com/Trae-AI/stream-to-river/rpcservice/dal/model"
	"github.com/Trae-AI/stream-to-river/rpcservice/kitex_gen/words"
)

// GetTodayReviewProgress retrieves the user's review progress for the current day.
// It's a wrapper function that calls `getTodayReviewProgressWithInitFlag` and discards the initialization flag.
//
// Parameters:
//   - ctx: The context used to control the lifetime of the request.
//   - req: A pointer to the `ReviewProgressReq` struct containing the user ID.
//
// Returns:
//   - *words.ReviewProgressResp: A pointer to the response struct with the review progress information.
//   - error: An error object if an unexpected error occurs during the process.
func GetTodayReviewProgress(ctx context.Context, req *words.ReviewProgressReq) (*words.ReviewProgressResp, error) {
	resp, _, err := getTodayReviewProgressWithInitFlag(req)
	return resp, err
}

// getTodayReviewProgressWithInitFlag retrieves the user's review progress for the current day and indicates whether initialization was performed.
// If the review progress record doesn't exist or is outdated, it initializes the record.
//
// Parameters:
//   - req: A pointer to the `ReviewProgressReq` struct containing the user ID.
//
// Returns:
//   - *words.ReviewProgressResp: A pointer to the response struct with the review progress information.
//   - bool: A flag indicating whether the review progress record was initialized.
//   - error: An error object if an unexpected error occurs during the process.
func getTodayReviewProgressWithInitFlag(req *words.ReviewProgressReq) (*words.ReviewProgressResp, bool, error) {
	now := time.Now()
	todayStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location()).Unix() // Timestamp of 00:00 today (local timezone)

	// Attempt to retrieve the existing review progress record
	progress, err := dao.GetReviewProgress(req.UserId)
	needInit := false

	if err != nil {
		// If the record doesn't exist, initialization is required
		if errors.Is(err, dao.ErrNoRecord) {
			needInit = true
		} else {
			return nil, false, err
		}
	} else {
		// If the last update time is earlier than 00:00 today, re - initialize
		if progress.LastUpdateTime < todayStart {
			needInit = true
		}
	}

	// Re - calculate the number of words pending review when initialization is needed
	if needInit {
		progress, err = initReviewProgress(req.UserId)
		if err != nil {
			return nil, false, err
		}
	}

	// Retrieve word statistics
	totalWords, err := getWordStatistics(req.UserId)
	if err != nil {
		return nil, needInit, err
	}

	// Convert the timestamp to YYYY - MM - DD HH:MM format (using the local timezone)
	lastUpdateTimeStr := time.Unix(progress.LastUpdateTime, 0).In(now.Location()).Format("2006-01-02 15:04")

	// Construct the response
	resp := &words.ReviewProgressResp{
		BaseResp:             common.BuildSuccBaseResp(),
		PendingReviewCount:   int32(progress.PendingReviewCount),
		CompletedReviewCount: int32(progress.CompletedReviewCount),
		LastUpdateTime:       lastUpdateTimeStr,
		TotalWordsCount:      totalWords,                        // Total number of words
		AllCompletedCount:    int32(progress.AllCompletedCount), // Total number of completed words (directly retrieved from the table)
	}

	return resp, needInit, nil
}

// initReviewProgress initializes the user's review progress.
// It calculates the number of words pending review and the total number of completed words,
// then creates or updates the review progress record.
//
// Parameters:
//   - userId: The unique identifier of the user.
//
// Returns:
//   - *model.ReviewProgress: A pointer to the `ReviewProgress` struct with the initialized review progress.
//   - error: An error object if an unexpected error occurs during the process.
func initReviewProgress(userId int64) (*model.ReviewProgress, error) {
	now := time.Now()
	currentTime := now.Unix()
	// Calculate the timestamp of 00:00 tomorrow (using the local timezone)
	tomorrow := now.AddDate(0, 0, 1)
	tomorrowStart := time.Date(tomorrow.Year(), tomorrow.Month(), tomorrow.Day(), 0, 0, 0, 0, now.Location()).Unix()

	// Retrieve the review records that need to be reviewed today (next_review_time < 00:00 tomorrow)
	reviewRecords, err := dao.GetReviewRecords(userId, tomorrowStart)
	if err != nil {
		return nil, err
	}

	// Count the number of words pending review
	pendingCount := len(reviewRecords)

	// Query the total number of completed words (level >= 7)
	allCompletedCount, err := dao.GetCompletedWordsCountFromRecord(userId)
	if err != nil {
		return nil, err
	}

	// Create or update the review progress record
	progress := &model.ReviewProgress{
		UserId:               userId,
		PendingReviewCount:   pendingCount,
		CompletedReviewCount: 0,                      // Reset to 0 for a new day
		AllCompletedCount:    int(allCompletedCount), // Set the total number of completed words
		LastUpdateTime:       currentTime,
	}

	err = dao.CreateOrUpdateReviewProgress(progress)
	if err != nil {
		return nil, err
	}

	return progress, nil
}

// getWordStatistics retrieves word statistics for a user.
// It counts the total number of words belonging to the user.
//
// Parameters:
//   - userId: The unique identifier of the user.
//
// Returns:
//   - int32: The total number of words.
//   - error: An error object if an unexpected error occurs during the process.
func getWordStatistics(userId int64) (totalWords int32, err error) {
	// 1. Count the total number of words: count the user's records in the words table
	totalWords, err = dao.GetTotalWordsCount(userId)
	if err != nil {
		return 0, err
	}

	// 2. Count the number of completed review words: records in the words_risite_record table where the user's level >= 7

	return totalWords, nil
}
