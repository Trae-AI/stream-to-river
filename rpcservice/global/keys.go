// Copyright (c) 2025 Bytedance Ltd. and/or its affiliates
// SPDX-License-Identifier: MIT

package global

// ==================== API-related constants ====================
// YoudaoUrl is the URL template for Youdao dictionary query results, taking a word as a parameter.
const YoudaoUrl = "https://www.youdao.com/result?word=%s&lang=en"

// ==================== Review interval-related constants and functions ====================
// Review interval constants, in seconds.
const (
	LEVEL_0_INTERVAL = 5 * 60             // 5 minutes
	LEVEL_1_INTERVAL = 24 * 60 * 60       // 24 hours
	LEVEL_2_INTERVAL = 24 * 60 * 60       // 24 hours
	LEVEL_3_INTERVAL = 48 * 60 * 60       // 48 hours
	LEVEL_4_INTERVAL = 48 * 60 * 60       // 48 hours
	LEVEL_5_INTERVAL = 72 * 60 * 60       // 72 hours
	LEVEL_6_INTERVAL = 192 * 60 * 60      // 192 hours
	LEVEL_7_INTERVAL = 365 * 24 * 60 * 60 // 365 days
)

// GetReviewInterval retrieves the review interval for a given level.
// Parameters:
//   - level: The review level.
//
// Returns:
//   - int64: The review interval in seconds corresponding to the given level.
func GetReviewInterval(level int) int64 {
	switch level {
	case 0:
		return LEVEL_0_INTERVAL
	case 1:
		return LEVEL_1_INTERVAL
	case 2:
		return LEVEL_2_INTERVAL
	case 3:
		return LEVEL_3_INTERVAL
	case 4:
		return LEVEL_4_INTERVAL
	case 5:
		return LEVEL_5_INTERVAL
	case 6:
		return LEVEL_6_INTERVAL
	default: // level 7 and above
		return LEVEL_7_INTERVAL
	}
}

// ==================== Question type-related constants ====================
// QuestionType defines the question type using the int64 type.
type QuestionType = int64

// Question type constants.
const (
	CHOOSE_CN        QuestionType = 1 // Select the correct Chinese meaning.
	CHOOSE_EN        QuestionType = 2 // Select the correct English meaning.
	PRONOUNCE_CHOOSE QuestionType = 3 // Select the correct Chinese meaning based on pronunciation.
	FILL_IN_BLANK    QuestionType = 4 // Fill-in-the-blank question.
)

// ==================== Word-related constants ====================
// WORDS_NUM_PER_PAGE_DEFAULT is the default number of words displayed per page in the word list.
const WORDS_NUM_PER_PAGE_DEFAULT = 20

// WORDS_LEVEL_TOTALLY_GRASK is the level at which a word is considered fully mastered.
const WORDS_LEVEL_TOTALLY_GRASK = 8

// ==================== Fake user-related constants ====================
// FAKE_USER_ID_FOR_DEFAULT is the ID of the default fake user.
const FAKE_USER_ID_FOR_DEFAULT = 1

// FAKE_USER_ORDER_ID_BASE_OFFSET is the base offset for the fake user's order ID.
const FAKE_USER_ORDER_ID_BASE_OFFSET = 1000000

// ==================== Review list module-related constants ====================
// Review list module-related constants.
const (
	MAX_SCORE         = 15 // Binary 1111, indicating all four question types are answered correctly.
	REVIEW_OPTION_NUM = 4  // The total number of answer options for each question.
)

// ==================== Review level-related constants ====================
// MAX_RISITE_LEVEL defines the maximum review level.
const MAX_RISITE_LEVEL = 7
