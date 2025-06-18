// Copyright (c) 2025 Bytedance Ltd. and/or its affiliates
// SPDX-License-Identifier: MIT

package words

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/cloudwego/kitex/pkg/kerrors"
	"github.com/cloudwego/kitex/pkg/klog"

	"github.com/Trae-AI/stream-to-river/rpcservice/dal/dao"
	"github.com/Trae-AI/stream-to-river/rpcservice/dal/model"
	"github.com/Trae-AI/stream-to-river/rpcservice/global"
	"github.com/Trae-AI/stream-to-river/rpcservice/kitex_gen/words"
)

// QuestionOption represents an option for a review question.
// It contains the answer ID and a description of the option.
type QuestionOption struct {
	AnswerId    int64  `json:"answer_id"`
	Description string `json:"description"`
}

// ReviewQuestion represents a review question.
// It includes information such as the word ID, question type, question text, options, correct answer index, and additional show information.
type ReviewQuestion struct {
	WordId        int64            `json:"word_id"`
	QuestionType  int64            `json:"question_type"`
	Question      string           `json:"question"`
	Options       []QuestionOption `json:"options"`
	CorrectAnswer int              `json:"correct_answer"`      // Index of the correct answer (0 - 3)
	ShowInfo      []string         `json:"show_info,omitempty"` // Additional field for fill - in - the - blank questions
}

// ReviewWordList represents the response structure for the review word list.
// It contains a slice of review questions.
type ReviewWordList struct {
	Questions []ReviewQuestion `json:"questions"`
}

// GetReviewWordList retrieves the list of review words and generates corresponding review questions for a user.
// It first fetches the review records, then retrieves the word details, generates questions for each word,
// shuffles the questions, and finally constructs the response.
//
// Parameters:
//   - ctx: The context used to control the lifetime of the request.
//   - userId: The unique identifier of the user.
//
// Returns:
//   - *words.ReviewListResp: A pointer to the response structure containing the review questions and total number of questions.
//   - error: An error object if an unexpected error occurs during the process.
func GetReviewWordList(ctx context.Context, userId int64) (resp *words.ReviewListResp, err error) {
	resp = words.NewReviewListResp()

	// Log the start of the operation
	klog.CtxInfof(ctx, "[GetReviewWordList] Starting for userId: %d", userId)

	// Get the timestamp of 23:59:59 today in the local timezone
	now := time.Now()
	today2359 := time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 0, now.Location())
	currentTime := today2359.Unix()
	klog.CtxInfof(ctx, "[GetReviewWordList] Querying review records with currentTime: %d (today 23:59:59 in local timezone)", currentTime)

	// Retrieve review records from the database
	reviewRecords, err := dao.GetReviewRecords(userId, currentTime)
	if err != nil {
		klog.CtxErrorf(ctx, "[GetReviewWordList] Failed to get review records: %v", err)
		return nil, kerrors.NewBizStatusError(-1, fmt.Sprintf("query review words from db failed, err=%v", err))
	}
	klog.CtxInfof(ctx, "[GetReviewWordList] Got %d review records", len(reviewRecords))

	if len(reviewRecords) == 0 {
		klog.CtxInfof(ctx, "[GetReviewWordList] No review records found, returning empty response")
		return resp, nil
	}

	// Extract word IDs from review records
	wordIds := make([]int64, len(reviewRecords))
	for i, record := range reviewRecords {
		wordIds[i] = int64(record.WordId)
	}
	klog.Debugf("get wordIds: %v", wordIds)

	// Retrieve word details from the database
	klog.CtxInfof(ctx, "[GetReviewWordList] Querying words info for %d wordIds", len(wordIds))
	wordsInfo, err := dao.GetWordsByIds(wordIds)
	if err != nil {
		klog.CtxErrorf(ctx, "[GetReviewWordList] Failed to get words info: %v", err)
		return nil, kerrors.NewBizStatusError(-1, fmt.Sprintf("query words from db failed , err=%v", err))
	}
	klog.CtxInfof(ctx, "[GetReviewWordList] Got %d words info", len(wordsInfo))

	// Create a mapping from word ID to word details
	wordMap := make(map[int64]*model.Word)
	for _, word := range wordsInfo {
		wordMap[word.WordId] = word
		klog.CtxInfof(ctx, "word_id: %v word_name:%v explains:%v tag_id:%v ",
			word.WordId, word.WordName, word.Explains, word.TagId)
	}

	// Generate questions for each review record
	var allQuestions []*ReviewQuestion
	klog.CtxInfof(ctx, "[GetReviewWordList] Starting to generate questions for %d review records", len(reviewRecords))

	for i, record := range reviewRecords {
		klog.CtxInfof(ctx, "[GetReviewWordList] Processing record %d/%d, wordId: %d", i+1, len(reviewRecords), record.WordId)

		word, exists := wordMap[int64(record.WordId)]
		if !exists {
			klog.CtxWarnf(ctx, "[GetReviewWordList] Word not found in wordMap for wordId: %d", record.WordId)
			continue
		}

		// Generate four types of questions for the word
		questions, err := generateQuestionsForWord(ctx, userId, word)
		if err != nil {
			klog.CtxErrorf(ctx, "Failed to generate questions for word_id=%v: %v", word.WordId, err)
			continue
		}
		klog.CtxInfof(ctx, "[GetReviewWordList] Generated %d questions for wordId: %d", len(questions), word.WordId)

		// Convert slice elements to pointers and append to allQuestions
		for i := range questions {
			allQuestions = append(allQuestions, &questions[i])
		}
	}

	klog.CtxInfof(ctx, "[GetReviewWordList] Total questions generated: %d", len(allQuestions))

	// Shuffle the questions using the Fisher - Yates algorithm
	shuffleQuestionPointers(allQuestions)

	// Set the total number of questions in the response
	resp.SetTotalNum(fmt.Sprintf("%d", len(allQuestions)))

	// Convert ReviewQuestion to the words package's ReviewQuestion type
	questionPtrs := make([]*words.ReviewQuestion, len(allQuestions))
	for i, q := range allQuestions {
		questionPtrs[i] = &words.ReviewQuestion{
			Question:     q.Question,                      // Use the Question field
			WordId:       q.WordId,                        // Add the WordId field
			QuestionType: q.QuestionType,                  // Add the QuestionType field
			Options:      convertToOptionItems(q.Options), // Convert to OptionItem type
		}

		if q.QuestionType == global.FILL_IN_BLANK {
			questionPtrs[i].ShowInfo = q.ShowInfo // Add the ShowInfo field
		}
	}
	resp.SetQuestions(questionPtrs)

	klog.CtxInfof(ctx, "[GetReviewWordList] Successfully completed for userId: %d, returning %d questions", userId, len(allQuestions))
	return resp, nil
}

// convertToOptionItems converts a slice of QuestionOption to a slice of words.OptionItem.
//
// Parameters:
//   - options: A slice of QuestionOption to be converted.
//
// Returns:
//   - []*words.OptionItem: A slice of pointers to words.OptionItem.
func convertToOptionItems(options []QuestionOption) []*words.OptionItem {
	result := make([]*words.OptionItem, len(options))
	for i, opt := range options {
		result[i] = &words.OptionItem{
			Description:  opt.Description,
			AnswerListId: opt.AnswerId, // Map AnswerId to AnswerListId
		}
	}
	return result
}

// shuffleQuestionPointers shuffles a slice of pointers to ReviewQuestion using the Fisher - Yates algorithm.
// This function modifies the input slice in - place.
//
// Parameters:
//   - questions: A slice of pointers to ReviewQuestion to be shuffled.
func shuffleQuestionPointers(questions []*ReviewQuestion) {
	rand.Seed(time.Now().UnixNano())
	for i := len(questions) - 1; i > 0; i-- {
		j := rand.Intn(i + 1)
		// Swap pointers to avoid copying the entire struct
		questions[i], questions[j] = questions[j], questions[i]
	}
}

// generateQuestionsForWord generates four types of review questions for a single word.
// It determines which question types to generate based on the word's tag and the user's review score.
//
// Parameters:
//   - ctx: The context used to control the lifetime of the request.
//   - userId: The unique identifier of the user.
//   - word: A pointer to the model.Word struct representing the word.
//
// Returns:
//   - []ReviewQuestion: A slice of ReviewQuestion generated for the word.
//   - error: An error object if an unexpected error occurs during the process.
func generateQuestionsForWord(ctx context.Context, userId int64, word *model.Word) ([]ReviewQuestion, error) {
	klog.CtxInfof(ctx, "[generateQuestionsForWord] Starting for userId: %d, wordId: %d, wordName: %s", userId, word.WordId, word.WordName)

	var questions []ReviewQuestion

	// Get the user's review record for the word
	record, err := dao.GetWordsRisiteRecord(userId, int64(word.WordId))
	var score int
	if err != nil {
		// If no record exists, assume the word is new and set the score to 0
		klog.CtxInfof(ctx, "[generateQuestionsForWord] No revisit record found for wordId: %d, setting score to 0", word.WordId)
		score = 0
	} else {
		score = record.Score
		klog.CtxInfof(ctx, "[generateQuestionsForWord] Found revisit record for wordId: %d, score: %d", word.WordId, score)
	}

	questionTypes := global.MAX_SCORE
	if word.TagId != 0 {
		tagInfo, err := dao.GetWordTagById(word.TagId)
		if err == nil {
			questionTypes = tagInfo.QuestionTypes
		} else {
			klog.CtxWarnf(ctx, "[generateQuestionsForWord] Failed to get tag info for tagId: %d, using default questionTypes", word.TagId)
		}
	}
	klog.CtxInfof(ctx, "[generateQuestionsForWord] tag_id:%d questionTypes: %v, current score: %d",
		word.TagId, questionTypes, score)

	// Generate a question to choose the correct Chinese meaning
	if (questionTypes&(1<<(global.CHOOSE_CN-1))) != 0 && (score&(1<<(global.CHOOSE_CN-1))) == 0 {
		klog.CtxInfof(ctx, "[generateQuestionsForWord] Generating CHOOSE_CN question for wordId: %d", word.WordId)
		q1, err := generateChooseQuestion(ctx, userId, word, global.CHOOSE_CN)
		if err != nil {
			klog.CtxErrorf(ctx, "[generateQuestionsForWord] Failed to generate CHOOSE_CN question for wordId: %d, err: %v", word.WordId, err)
			return nil, fmt.Errorf("failed to generate CHOOSE_CN question: %v", err)
		}
		questions = append(questions, q1)
		klog.CtxInfof(ctx, "[generateQuestionsForWord] Successfully generated CHOOSE_CN question for wordId: %d", word.WordId)
	}

	// Generate a question to choose the correct English meaning
	if (questionTypes&(1<<(global.CHOOSE_EN-1))) != 0 && (score&(1<<(global.CHOOSE_EN-1))) == 0 {
		klog.CtxInfof(ctx, "[generateQuestionsForWord] Generating CHOOSE_EN question for wordId: %d", word.WordId)
		q2, err := generateChooseQuestion(ctx, userId, word, global.CHOOSE_EN)
		if err != nil {
			klog.CtxErrorf(ctx, "[generateQuestionsForWord] Failed to generate CHOOSE_EN question for wordId: %d, err: %v", word.WordId, err)
			return nil, fmt.Errorf("failed to generate CHOOSE_EN question: %v", err)
		}
		questions = append(questions, q2)
		klog.CtxInfof(ctx, "[generateQuestionsForWord] Successfully generated CHOOSE_EN question for wordId: %d", word.WordId)
	}

	// Generate a question to choose the correct Chinese meaning based on pronunciation
	if (questionTypes&(1<<(global.PRONOUNCE_CHOOSE-1))) != 0 && (score&(1<<(global.PRONOUNCE_CHOOSE-1))) == 0 {
		klog.CtxInfof(ctx, "[generateQuestionsForWord] Generating PRONOUNCE_CHOOSE question for wordId: %d", word.WordId)
		q3, err := generateChooseQuestion(ctx, userId, word, global.PRONOUNCE_CHOOSE)
		if err != nil {
			klog.CtxErrorf(ctx, "[generateQuestionsForWord] Failed to generate PRONOUNCE_CHOOSE question for wordId: %d, err: %v", word.WordId, err)
			return nil, fmt.Errorf("failed to generate PRONOUNCE_CHOOSE question: %v", err)
		}
		questions = append(questions, q3)
		klog.CtxInfof(ctx, "[generateQuestionsForWord] Successfully generated PRONOUNCE_CHOOSE question for wordId: %d", word.WordId)
	}

	// Generate a fill - in - the - blank question
	if (questionTypes&(1<<(global.FILL_IN_BLANK-1))) != 0 && (score&(1<<(global.FILL_IN_BLANK-1))) == 0 {
		klog.CtxInfof(ctx, "[generateQuestionsForWord] Generating FILL_IN_BLANK question for wordId: %d", word.WordId)
		q4, err := generateFillInBlankQuestion(ctx, word)
		if err != nil {
			klog.CtxErrorf(ctx, "[generateQuestionsForWord] Failed to generate FILL_IN_BLANK question for wordId: %d, err: %v", word.WordId, err)
			return nil, fmt.Errorf("failed to generate FILL_IN_BLANK question: %v", err)
		}
		klog.CtxInfof(ctx, "my check q4 ShowInfo: %v", q4.ShowInfo)
		questions = append(questions, q4)
		klog.CtxInfof(ctx, "[generateQuestionsForWord] Successfully generated FILL_IN_BLANK question for wordId: %d", word.WordId)
	}

	klog.CtxInfof(ctx, "[generateQuestionsForWord] Completed for wordId: %d, generated %d questions", word.WordId, len(questions))
	return questions, nil
}

// generateFillInBlankQuestion generates a fill - in - the - blank review question for a word.
//
// Parameters:
//   - ctx: The context used to control the lifetime of the request.
//   - word: A pointer to the model.Word struct representing the word.
//
// Returns:
//   - ReviewQuestion: A ReviewQuestion of the fill - in - the - blank type.
//   - error: An error object if an unexpected error occurs during the process.
func generateFillInBlankQuestion(ctx context.Context, word *model.Word) (ReviewQuestion, error) {
	// Generate the show information for the fill - in - the - blank question
	showInfo := generateBlankShowInfo(word.WordName)
	klog.CtxInfof(ctx, "my check showInfo: %v", showInfo)

	options := make([]QuestionOption, 1)

	return ReviewQuestion{
		WordId:       word.WordId,
		QuestionType: int64(global.FILL_IN_BLANK),
		Question:     word.Explains, // The question is the Chinese explanation
		Options:      options,
		ShowInfo:     showInfo, // Add the ShowInfo field
	}, nil
}

// generateBlankShowInfo generates the show information for a fill - in - the - blank question.
// It randomly selects some positions in the word to reveal and hides the rest.
//
// Parameters:
//   - wordName: The name of the word.
//
// Returns:
//   - []string: A slice of strings representing the show information.
func generateBlankShowInfo(wordName string) []string {
	if len(wordName) == 0 {
		return []string{}
	}

	// Calculate the number of letters to keep (at least 1, 1/5 of the word length)
	keepCount := len(wordName) / 5
	if keepCount == 0 {
		keepCount = 1
	}

	// Randomly select positions to keep
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	keepPositions := make(map[int]bool)
	for len(keepPositions) < keepCount {
		pos := r.Intn(len(wordName))
		keepPositions[pos] = true
	}

	// Generate the show information
	showInfo := make([]string, len(wordName))
	for i, char := range wordName {
		if keepPositions[i] {
			showInfo[i] = string(char)
		} else {
			showInfo[i] = "_"
		}
	}

	return showInfo
}

// generateChooseQuestion generates a multiple - choice review question for a word.
// It can generate different types of multiple - choice questions based on the question type.
//
// Parameters:
//   - ctx: The context used to control the lifetime of the request.
//   - userId: The unique identifier of the user.
//   - word: A pointer to the model.Word struct representing the word.
//   - questionType: The type of the multiple - choice question.
//
// Returns:
//   - ReviewQuestion: A ReviewQuestion of the multiple - choice type.
//   - error: An error object if an unexpected error occurs during the process.
func generateChooseQuestion(ctx context.Context, userId int64, word *model.Word, questionType global.QuestionType) (ReviewQuestion, error) {
	klog.CtxInfof(ctx, "[generateChooseQuestion] Starting for userId: %d, wordId: %d, questionType: %d", userId, word.WordId, questionType)

	// Get four options including the correct answer
	answerLists, err := getRandomOptionsWithCorrect(ctx, userId, word.WordId)
	if err != nil {
		klog.CtxErrorf(ctx, "[generateChooseQuestion] Failed to get random options for userId: %d, wordId: %d, err: %v", userId, word.WordId, err)
		return ReviewQuestion{}, err
	}
	klog.CtxInfof(ctx, "[generateChooseQuestion] Got %d answerLists for wordId: %d", len(answerLists), word.WordId)

	options := make([]QuestionOption, global.REVIEW_OPTION_NUM)

	// Set the question content and the function to get option descriptions based on the question type
	var question string
	var getOptionFunc func(*model.AnswerList) string

	switch questionType {
	case global.CHOOSE_CN:
		// Given the English word, choose the Chinese meaning
		question = word.WordName
		getOptionFunc = func(al *model.AnswerList) string {
			return al.Description
		}
	case global.CHOOSE_EN:
		// Given the Chinese meaning, choose the English word
		question = word.Explains
		getOptionFunc = func(al *model.AnswerList) string {
			return al.WordName
		}
	case global.PRONOUNCE_CHOOSE:
		// Given the pronunciation, choose the Chinese meaning
		question = word.PronounceUs
		getOptionFunc = func(al *model.AnswerList) string {
			return al.Description
		}
	case global.FILL_IN_BLANK:
		// Redirect to generate a fill - in - the - blank question
		klog.CtxInfof(ctx, "[generateChooseQuestion] Redirecting to generateFillInBlankQuestion for wordId: %d", word.WordId)
		return generateFillInBlankQuestion(ctx, word)
	default:
		klog.CtxErrorf(ctx, "[generateChooseQuestion] Unknown question type: %d for wordId: %d", questionType, word.WordId)
		return ReviewQuestion{}, fmt.Errorf("unknown question type: %d", questionType)
	}

	// Fill in the options
	for i := 0; i < global.REVIEW_OPTION_NUM; i++ {
		if i < len(answerLists) && answerLists[i] != nil {
			options[i] = QuestionOption{
				AnswerId:    answerLists[i].AnswerId, // Use the AnswerId from AnswerList
				Description: getOptionFunc(answerLists[i]),
			}
		}
	}

	klog.CtxInfof(ctx, "[generateChooseQuestion] Successfully completed for userId: %d, wordId: %d, questionType: %d", userId, word.WordId, questionType)
	return ReviewQuestion{
		WordId:       word.WordId,
		QuestionType: int64(questionType),
		Question:     question,
		Options:      options,
	}, nil
}

// generateRandomOrderIds generates a list of non - repeating, random order IDs not greater than maxOrderId.
// It ensures that the list contains the existId if provided.
//
// Parameters:
//   - maxOrderId: The maximum order ID.
//   - existId: The existing order ID to include in the list (can be 0 to ignore).
//   - count: The number of order IDs to generate.
//
// Returns:
//   - []int: A slice of integers representing the generated order IDs.
func generateRandomOrderIds(maxOrderId, existId int, count int) []int {
	if count <= 0 || maxOrderId <= 0 {
		return []int{}
	}
	if count == 1 && existId > 0 {
		return []int{existId}
	}

	targetCount := count
	if count > maxOrderId {
		targetCount = maxOrderId
	}

	randOrderSeq := generateRandomSequence(maxOrderId)
	if existId > 0 {
		candidateOrderSeq := randOrderSeq[:targetCount]
		for _, id := range candidateOrderSeq {
			if id == existId {
				return candidateOrderSeq
			}
		}

		r := rand.New(rand.NewSource(time.Now().UnixNano()))
		candidateOrderSeq[r.Intn(targetCount)] = existId
		return candidateOrderSeq
	} else {
		return randOrderSeq[:targetCount]
	}
}

// generateRandomSequence generates a random sequence of integers from 1 to maxId.
//
// Parameters:
//   - maxId: The maximum value of the sequence.
//
// Returns:
//   - []int: A slice of integers representing the random sequence.
func generateRandomSequence(maxId int) []int {
	sequence := make([]int, maxId)
	for i := 0; i < maxId; i++ {
		sequence[i] = i + 1
	}

	rand.Shuffle(maxId, func(i, j int) {
		sequence[i], sequence[j] = sequence[j], sequence[i]
	})

	return sequence
}

// getRandomOptionsWithCorrect retrieves four options including the correct answer for a word.
// It first gets the correct answer, then generates random order IDs, and retrieves the corresponding options.
// If there are not enough options, it supplements them from the default user.
//
// Parameters:
//   - ctx: The context used to control the lifetime of the request.
//   - userId: The unique identifier of the user.
//   - wordId: The unique identifier of the word.
//
// Returns:
//   - []*model.AnswerList: A slice of pointers to model.AnswerList representing the options.
//   - error: An error object if an unexpected error occurs during the process.
func getRandomOptionsWithCorrect(ctx context.Context, userId int64, wordId int64) ([]*model.AnswerList, error) {
	klog.CtxInfof(ctx, "[getRandomOptionsWithCorrect] Starting for userId: %d, wordId: %d", userId, wordId)

	// Get the correct answer list for the word
	correctAnswerList, err := dao.GetAnswerListByWordId(userId, wordId)
	if err != nil {
		return nil, fmt.Errorf("failed to get correct answer from db, err: %v", err)
	}
	klog.CtxInfof(ctx, "[getRandomOptionsWithCorrect] Found correct answer: orderId=%d, answerId=%d", correctAnswerList.OrderId, correctAnswerList.AnswerId)

	// Get the maximum order ID
	maxOrderId, err := dao.GetMaxOrderId(userId)
	if err != nil {
		return nil, fmt.Errorf("failed to get MaxOrderId from db, err: %v", err)
	}
	klog.CtxInfof(ctx, "[getRandomOptionsWithCorrect] MaxOrderId for userId %d: %d", userId, maxOrderId)

	count := global.REVIEW_OPTION_NUM

	// Generate random order IDs
	randomOrderIds := generateRandomOrderIds(maxOrderId, correctAnswerList.OrderId, count-1)
	klog.CtxInfof(ctx, "[getRandomOptionsWithCorrect] Generated randomOrderIds: %v, correctOrderId=%d", randomOrderIds, correctAnswerList.OrderId)

	// Get the answer lists corresponding to the random order IDs
	answerLists, err := dao.GetAnswerListByOrderIds(userId, randomOrderIds)
	if err != nil {
		return nil, fmt.Errorf("failed to get AnswerList from db, err: %v", err)
	}

	// Supplement options from the default user if there are not enough
	if len(answerLists) < count {
		klog.CtxInfof(ctx, "[getRandomOptionsWithCorrect] answerLists not enough, need %d, got %d", count, len(answerLists))
		neededCount := count - len(answerLists)

		// Get the maximum order ID for the default user
		fakeMaxOrderId, err := dao.GetMaxOrderId(global.FAKE_USER_ID_FOR_DEFAULT)
		if err != nil {
			return nil, err
		}
		klog.CtxInfof(ctx, "[getRandomOptionsWithCorrect] FakeMaxOrderId: %d", fakeMaxOrderId)

		// Generate random order IDs for the default user
		fakeRandomOrderIds := generateRandomOrderIds(fakeMaxOrderId, -1, neededCount)
		klog.CtxInfof(ctx, "fakeRandomOrderIds: %v", fakeRandomOrderIds)

		// Get the answer lists for the default user
		fakeAnswerLists, err := dao.GetAnswerListByOrderIds(global.FAKE_USER_ID_FOR_DEFAULT, fakeRandomOrderIds)
		if err != nil {
			klog.CtxErrorf(ctx, "GetAnswerListByOrderIds failed: %v", err)
			return nil, err
		}

		// Update the order IDs of the answer lists from the default user
		for i, answer := range fakeAnswerLists {
			new_order_id := maxOrderId + global.FAKE_USER_ORDER_ID_BASE_OFFSET + i
			fakeRandomOrderIds[i] = new_order_id
			answer.OrderId = new_order_id
		}
		klog.CtxInfof(ctx, "fakeAnswerLists: %v", fakeAnswerLists)

		answerLists = append(answerLists, fakeAnswerLists...)

		// The order ID list will include results from the default user
		randomOrderIds = append(randomOrderIds, fakeRandomOrderIds...)
	}

	resultAnswerLists := make([]*model.AnswerList, len(randomOrderIds))
	for i, orderId := range randomOrderIds {
		for _, answerList := range answerLists {
			if answerList.OrderId == orderId {
				resultAnswerLists[i] = answerList
			}
		}
	}

	klog.CtxInfof(ctx, "Successfully completed for userId: %d, wordId: %d, returning %d options, resultAnswerLists:%v",
		userId, wordId, len(resultAnswerLists), resultAnswerLists)
	return resultAnswerLists, nil
}
