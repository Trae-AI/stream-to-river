// Copyright (c) 2025 Bytedance Ltd. and/or its affiliates
// SPDX-License-Identifier: MIT

package words

import (
	"context"
	"testing"

	"github.com/Trae-AI/stream-to-river/internal/test"
	"github.com/Trae-AI/stream-to-river/rpcservice/dal/dao"
	"github.com/Trae-AI/stream-to-river/rpcservice/dal/model"
	"github.com/Trae-AI/stream-to-river/rpcservice/kitex_gen/words"
)

// TestReviewProgressCountCorrectness tests the correctness of the review progress count.
// Verifies that the pending_review_count will not be double-counted when a new user adds a word for the first time.
func TestReviewProgressCountCorrectness(t *testing.T) {
	// 1. Clean up and create the tables required for testing.
	recreateMockWordTable()
	recreateMockAnswerListTable()
	recreateMockWordReviewRecordTable()
	recreateMockReviewProgressTable()

	// 2. Test data
	const (
		userId int64 = 12345
		tagId  int32 = 1
	)

	// 3. Simulate a new user adding a word for the first time.
	t.Run("FirstWordAddition", func(t *testing.T) {
		// Since vocapi.ProcessWord requires a real API call, we only test the core logic here.
		// Directly create word, answer_list, and words_risite_record records.
		word1 := &model.Word{
			WordName:    "hello",
			Description: "Hello world",
			Explains:    "打招呼用语",
			UserId:      userId,
			TagId:       tagId,
		}
		err := dao.AddWord(word1)
		test.Assert(t, err == nil, "Failed to add word:", err)

		// Query the added word to get the word_id.
		queriedWord, err := dao.QueryWordByUserIdAndName(userId, "hello")
		test.Assert(t, err == nil, "Failed to query word:", err)

		// Add an answer_list record.
		answerList := &model.AnswerList{
			WordId:      queriedWord.WordId,
			UserId:      userId,
			WordName:    queriedWord.WordName,
			Description: queriedWord.Explains,
		}
		err = dao.AddAnswerList(answerList)
		test.Assert(t, err == nil, "Failed to add answer list:", err)

		// Add a review record.
		record := &model.WordsRisiteRecord{
			WordId:         int(queriedWord.WordId),
			Level:          1,
			NextReviewTime: 1234567890, // Set to the time when a review is required.
			DowngradeStep:  1,
			TotalCorrect:   0,
			TotalWrong:     0,
			Score:          0,
			UserId:         userId,
		}
		err = dao.AddWordsRisiteRecord(record)
		test.Assert(t, err == nil, "Failed to add review record:", err)

		// 4. Call getTodayReviewProgressWithInitFlag to simulate the logic in AddNewWord.
		progressReq := &words.ReviewProgressReq{UserId: userId}
		resp, wasInitialized, err := getTodayReviewProgressWithInitFlag(progressReq)
		test.Assert(t, err == nil, "Failed to get review progress:", err)
		test.Assert(t, wasInitialized == true, "Should have initialized for new user")

		// 5. Verify: When a new user adds a word for the first time, pending_review_count should be equal to total_words_count.
		test.Assert(t, resp.PendingReviewCount == resp.TotalWordsCount,
			"First word: pending_review_count (%d) should equal total_words_count (%d)",
			resp.PendingReviewCount, resp.TotalWordsCount)
		test.Assert(t, resp.TotalWordsCount == 1, "Total words should be 1")
		test.Assert(t, resp.PendingReviewCount == 1, "Pending review count should be 1")

		t.Logf("First word - Total: %d, Pending: %d, Was Initialized: %v",
			resp.TotalWordsCount, resp.PendingReviewCount, wasInitialized)
	})

	// 6. Test subsequent words logic
	t.Run("SubsequentWordAddition", func(t *testing.T) {
		// 添加第二个单词
		word2 := &model.Word{
			WordName:    "world",
			Description: "World peace",
			Explains:    "世界",
			UserId:      userId,
			TagId:       tagId,
		}
		err := dao.AddWord(word2)
		test.Assert(t, err == nil, "Failed to add second word:", err)

		// 查询添加的单词获取word_id
		queriedWord2, err := dao.QueryWordByUserIdAndName(userId, "world")
		test.Assert(t, err == nil, "Failed to query second word:", err)

		// 添加answer_list记录
		answerList2 := &model.AnswerList{
			WordId:      queriedWord2.WordId,
			UserId:      userId,
			WordName:    queriedWord2.WordName,
			Description: queriedWord2.Explains,
		}
		err = dao.AddAnswerList(answerList2)
		test.Assert(t, err == nil, "Failed to add answer list for second word:", err)

		// add review record
		record2 := &model.WordsRisiteRecord{
			WordId:         int(queriedWord2.WordId),
			Level:          1,
			NextReviewTime: 1234567890, // 设置为需要复习的时间
			DowngradeStep:  1,
			TotalCorrect:   0,
			TotalWrong:     0,
			Score:          0,
			UserId:         userId,
		}
		err = dao.AddWordsRisiteRecord(record2)
		test.Assert(t, err == nil, "Failed to add review record for second word:", err)

		// Get the current review progress (should not be re-initialized).
		progressReq := &words.ReviewProgressReq{UserId: userId}
		resp, wasInitialized, err := getTodayReviewProgressWithInitFlag(progressReq)
		test.Assert(t, err == nil, "Failed to get review progress for second word:", err)
		test.Assert(t, wasInitialized == false, "Should not initialize for existing user")
		test.Assert(t, resp.BaseResp.StatusCode == 0, "BaseResp.StatusCode should be 0")

		// Simulate the incremental logic in AddNewWord.
		if !wasInitialized {
			err = dao.IncrementPendingReviewCount(userId)
			test.Assert(t, err == nil, "Failed to increment pending review count:", err)
		}

		// Get the updated review progress.
		resp, _, err = getTodayReviewProgressWithInitFlag(progressReq)
		test.Assert(t, err == nil, "Failed to get updated review progress:", err)

		// 7. Verify: When adding subsequent words, pending_review_count should be incremented correctly.
		test.Assert(t, resp.TotalWordsCount == 2, "Total words should be 2")
		test.Assert(t, resp.PendingReviewCount == 2, "Pending review count should be 2")
		test.Assert(t, resp.PendingReviewCount == resp.TotalWordsCount,
			"Second word: pending_review_count (%d) should equal total_words_count (%d)",
			resp.PendingReviewCount, resp.TotalWordsCount)

		t.Logf("Second word - Total: %d, Pending: %d, Was Initialized: %v",
			resp.TotalWordsCount, resp.PendingReviewCount, wasInitialized)
	})

	t.Log("Review progress count correctness test passed!")
}

// TestGetTodayReviewProgress tests the GetTodayReviewProgress function.
func TestGetTodayReviewProgress(t *testing.T) {
	recreateMockReviewProgressTable()
	recreateMockWordReviewRecordTable()
	recreateMockWordTable()

	userId := int64(12345)

	// no word
	req := &words.ReviewProgressReq{UserId: userId}
	resp, err := GetTodayReviewProgress(context.Background(), req)
	test.Assert(t, err == nil, err)
	test.Assert(t, resp.TotalWordsCount == 0)
	test.Assert(t, resp.PendingReviewCount == 0, resp.PendingReviewCount)

	awr := &words.AddWordReq{
		UserId:   userId,
		WordName: "mock",
		TagId:    1,
	}
	resp2, err := AddNewWord(context.Background(), awr)
	test.Assert(t, err == nil, err)
	test.Assert(t, resp2.BaseResp.StatusCode == 0, resp2.BaseResp.StatusCode)

	// has word
	resp, err = GetTodayReviewProgress(context.Background(), req)
	test.Assert(t, err == nil, err)
	test.Assert(t, resp.TotalWordsCount == 1)
	test.Assert(t, resp.PendingReviewCount == 1)
}
