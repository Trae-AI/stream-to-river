// Copyright (c) 2025 Bytedance Ltd. and/or its affiliates
// SPDX-License-Identifier: MIT

package words

import (
	"context"
	"fmt"
	"time"

	"github.com/cloudwego/kitex/pkg/klog"

	"github.com/Trae-AI/stream-to-river/rpcservice/biz/words/common"
	"github.com/Trae-AI/stream-to-river/rpcservice/dal/dao"
	"github.com/Trae-AI/stream-to-river/rpcservice/global"
	"github.com/Trae-AI/stream-to-river/rpcservice/kitex_gen/words"
)

// SubmitAnswer submits the answer for a word review.
// It validates the answer based on the question type, updates the review record,
// and increments the completed word count if the review level reaches a certain threshold.
//
// Parameters:
//   - ctx: The context used to control the lifetime of the request.
//   - userId: The unique identifier of the user submitting the answer.
//   - wordId: The unique identifier of the word being reviewed.
//   - answerId: The identifier of the submitted answer.
//   - questionType: The type of the question (e.g., fill - in - the - blank, multiple - choice).
//   - filledName: A pointer to the filled - in answer for fill - in - the - blank questions.
//
// Returns:
//   - *words.SubmitAnswerResp: A pointer to the response struct containing the result of the answer submission.
//   - error: An error object if an unexpected error occurs during the process.
func SubmitAnswer(ctx context.Context, userId int64, wordId int64, answerId int64, questionType int64, filledName *string) (resp *words.SubmitAnswerResp, err error) {
	// Initialize the response object and set a successful base response
	resp = words.NewSubmitAnswerResp()
	resp.BaseResp = common.BuildSuccBaseResp()

	// Log the start of the answer submission process
	klog.CtxInfof(ctx, "Hello SubmitAnswerHandler, userId=%v, wordId=%v, answerId=%v, questionType=%v, filledName=%v",
		userId, wordId, answerId, questionType, filledName)

	// Flag to indicate whether the answer is correct
	var isCorrect bool

	// Query the word from the database
	word, err := dao.QueryWord(wordId)
	if err != nil {
		// Log the error and return a failure response if the word query fails
		klog.CtxErrorf(ctx, "Failed to get word: word_id=%v, error=%v", wordId, err)
		resp.BaseResp = common.BuildBaseResp(-1, "单词不存在")
		return resp, nil
	}

	// Initialize the correct answer ID
	correctAnswerId := int64(0)

	// Determine the correctness of the answer based on the question type
	if questionType == global.FILL_IN_BLANK {
		// Log the initial check for fill - in - the - blank questions
		klog.CtxInfof(ctx, "My Check Fill in blank check: filled=%v, actual=%v", *filledName, wordId)
		// Validate that the filled - in answer is not nil
		if filledName == nil {
			resp.BaseResp = common.BuildBaseResp(-1, "填空题答案不能为空")
			return resp, nil
		}

		// Compare the filled - in answer with the actual word name
		isCorrect = *filledName == word.WordName
		klog.CtxInfof(ctx, "Fill in blank check: filled=%v, actual=%v, correct=%v", *filledName, word.WordName, isCorrect)
	} else {
		// Retrieve the answer list for multiple - choice questions
		answerList, err := dao.GetAnswerListByAnswerId(userId, answerId)
		if err != nil {
			if err == dao.ErrNoRecord {
				// If the answer is invalid, mark it as incorrect
				isCorrect = false
			} else {
				// Log the error and return a failure response if the answer list query fails
				klog.CtxErrorf(ctx, "Failed to get answer_list: user_id=%v, answer_id=%v, error=%v", userId, answerId, err)
				resp.BaseResp = common.BuildBaseResp(-1, fmt.Sprintf("查询答案失败，err=%s", err.Error()))
				return resp, nil
			}
		} else {
			// Determine the correctness of the answer based on the word ID
			isCorrect = answerList.WordId == wordId
		}
	}

	// Retrieve the correct answer list
	correctAnswerList, err := dao.GetAnswerListByWordId(userId, wordId)
	if err != nil {
		// Log the error and return a failure response if the correct answer list query fails
		klog.CtxErrorf(ctx, "Failed to get correct answer: user_id=%v, word_id=%v, error=%v", userId, wordId, err)
		resp.BaseResp = common.BuildBaseResp(-1, "正确答案记录不存在")
		return resp, nil
	}
	correctAnswerId = correctAnswerList.AnswerId

	// Query the word review record
	risiteRecord, err := dao.GetWordsRisiteRecord(userId, wordId)
	if err != nil {
		// Log the error and return a failure response if the review record query fails
		klog.CtxErrorf(ctx, "Failed to get words_risite_record: user_id=%v, word_id=%v, error=%v", userId, wordId, err)
		resp.BaseResp = common.BuildBaseResp(-1, "复习记录不存在")
		return resp, nil
	}

	// Set the correctness flag and correct answer ID in the response
	resp.IsCorrect = isCorrect
	resp.CorrectAnswerId = correctAnswerId

	// Flag to indicate whether the review level has increased
	completedLevel := false

	// Update the review record based on the answer correctness
	if isCorrect {
		// Log the old review record information
		klog.CtxInfof(ctx, "correct old_info:%v my_tag_id:%v", risiteRecord, word.TagId)
		// Initialize the maximum score
		maxScore := global.MAX_SCORE

		// Update the maximum score if the word has a tag
		if word.TagId != 0 {
			tagInfo, err := dao.GetWordTagById(word.TagId)
			if err == nil {
				maxScore = tagInfo.MaxScore
				klog.CtxInfof(ctx, "reset maxScore:%v", maxScore)
			}
		}
		// Increment the total correct count
		risiteRecord.TotalCorrect++
		// Update the score based on the question type
		risiteRecord.Score |= (1 << (questionType - 1))
		resp.Message = "答案正确"

		// Check if the review level needs to be increased
		if (risiteRecord.Score & maxScore) == maxScore {
			klog.Info("update level")
			// Reset the score
			risiteRecord.Score = 0
			// Calculate the next review interval
			interval := global.GetReviewInterval(risiteRecord.Level)

			// Calculate the next review time
			now := time.Now()
			nextReviewTime := now.Add(time.Duration(interval) * time.Second)
			risiteRecord.NextReviewTime = nextReviewTime.Unix()
			klog.CtxInfof(ctx, "Setting next review time: current=%v, interval=%d seconds, next=%v",
				now.Format("2006-01-02 15:04:05"), interval, nextReviewTime.Format("2006-01-02 15:04:05"))

			// Increase the review level
			risiteRecord.Level++
			completedLevel = true // Mark that the level has increased
		}
		// Log the new review record information
		klog.CtxInfof(ctx, "correct new_info:%v", risiteRecord)
	} else {
		// Log the old review record information
		klog.CtxInfof(ctx, "wrong old_info:%v", risiteRecord)
		// Increment the total wrong count
		risiteRecord.TotalWrong++
		// Store the previous score
		prevScore := risiteRecord.Score
		// Update the score based on the question type
		risiteRecord.Score &= ^(1 << (questionType - 1))
		resp.Message = "答案错误"

		// Check if the review level needs to be decreased
		if risiteRecord.Score == 0 && prevScore > 0 {
			if risiteRecord.Level > 0 {
				risiteRecord.Level--
			}
		}
		// Set the next review time to the current time
		now := time.Now()
		risiteRecord.NextReviewTime = now.Unix()
		klog.CtxInfof(ctx, "Setting immediate review time: %v", now.Format("2006-01-02 15:04:05"))
		// Log the new review record information
		klog.CtxInfof(ctx, "wrong new_info:%v", risiteRecord)
	}

	// Update the review record in the database
	err = dao.UpdateWordsRisiteRecord(risiteRecord)
	if err != nil {
		// Log the error and return a failure response if the update fails
		klog.CtxErrorf(ctx, "Failed to update words_risite_record: %v", err)
		resp.BaseResp = common.BuildBaseResp(-1, "更新记录失败")
		return resp, nil
	}

	// If the review level has increased and reached a certain threshold, update the completed word count
	if completedLevel {
		klog.CtxInfof(ctx, "Level reached 8 for user_id=%v, word_id=%v, incrementing all_completed_count", userId, wordId)

		// Ensure that today's review progress record exists
		progressReq := &words.ReviewProgressReq{UserId: userId}
		_, err = GetTodayReviewProgress(ctx, progressReq)
		if err != nil {
			klog.CtxErrorf(ctx, "Failed to get/init today review progress for user_id=%v: %v", userId, err)
		} else {
			// Increment the number of completed review words for today
			err = dao.UpdateCompletedReviewCount(userId, 1)
			if err != nil {
				klog.CtxErrorf(ctx, "Failed to increment completed review count for user_id=%v: %v", userId, err)
			}

			// Increment the total number of completely completed words
			if risiteRecord.Level >= global.WORDS_LEVEL_TOTALLY_GRASK {
				err = dao.IncrementAllCompletedCount(userId)
				if err != nil {
					klog.CtxErrorf(ctx, "Failed to increment all_completed_count for user_id=%v: %v", userId, err)
				} else {
					klog.CtxInfof(ctx, "Successfully incremented all_completed_count for user_id=%v due to level reaching 7", userId)
				}
			}
		}
	}

	// Log the completion of the answer submission process
	klog.CtxInfof(ctx, "Submit answer completed: user_id=%v, word_id=%v, is_correct=%v, ", userId, wordId, isCorrect)
	return resp, nil
}
