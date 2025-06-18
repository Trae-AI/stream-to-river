// Copyright (c) 2025 Bytedance Ltd. and/or its affiliates
// SPDX-License-Identifier: MIT

package words

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/cloudwego/kitex/pkg/kerrors"
	"github.com/cloudwego/kitex/pkg/klog"

	"github.com/Trae-AI/stream-to-river/rpcservice/biz/words/common"
	"github.com/Trae-AI/stream-to-river/rpcservice/biz/words/vocapi"
	"github.com/Trae-AI/stream-to-river/rpcservice/dal/dao"
	"github.com/Trae-AI/stream-to-river/rpcservice/dal/model"
	"github.com/Trae-AI/stream-to-river/rpcservice/global"
	"github.com/Trae-AI/stream-to-river/rpcservice/kitex_gen/words"
)

// GetWordDetail retrieves the detailed information of a word based on its name.
// Currently, it does not use the database to store the detailed information of the given word.
// (Using the database can reduce the pressure on external API requests, which can be an optimization method for the future.)
//
// Parameters:
//   - ctx: The context used to control the lifetime of the request.
//   - wordName: The name of the word to retrieve details for.
//
// Returns:
//   - *words.WordDetail: A pointer to the WordDetail struct containing the word details.
//   - error: An error object if an error occurs during the process.
func GetWordDetail(ctx context.Context, wordName string) (*words.WordDetail, error) {
	// Clean up the input word name and call the vocapi package to get word details
	wordExplains, err := vocapi.ProcessWord(cleanupWord(wordName))
	if err != nil {
		return nil, kerrors.NewBizStatusError(-1, err.Error())
	}

	// Convert the vocapi.wordDetail to words.WordDetail
	wordDetail := &words.WordDetail{
		NewWordName_: wordExplains.NewWordName,
		Description:  wordExplains.ExplainsOxford,
		Explains:     wordExplains.ExplainsYoudao,
		PronounceUs:  wordExplains.PronounceUS,
		PronounceUk:  wordExplains.PronounceUK,
		Sentences:    wordExplains.Sentences,
	}
	return wordDetail, nil
}

// AddNewWord adds a new word to the system based on the provided request.
// It first checks if the word already exists for the user. If not, it adds the word and related records.
// After that, it updates the user's review progress.
//
// Parameters:
//   - ctx: The context used to control the lifetime of the request.
//   - req: A pointer to the AddWordReq struct containing the request information.
//
// Returns:
//   - *words.WordResp: A pointer to the WordResp struct containing the response information.
//   - error: An error object if an error occurs during the process.
func AddNewWord(ctx context.Context, req *words.AddWordReq) (*words.WordResp, error) {
	user_id := req.UserId
	tag_id := req.TagId

	klog.CtxInfof(ctx, "add word user_id: %v wordName before cleanup: %v", user_id, req.WordName)
	req.SetWordName(cleanupWord(req.GetWordName()))
	// Call ProcessWord to get the processed word information
	WordDetails, err := vocapi.ProcessWord(req.WordName)
	// According to the current ProcessWord logic, err will never be non-nil here
	if err != nil {
		klog.CtxErrorf(ctx, "vocapi ProcessWord error: err=%v", err)
		return nil, kerrors.NewBizStatusError(-1, err.Error())
	}
	// If the word does not exist, return an error directly
	if WordDetails.ErrorNo != 0 || WordDetails.NewWordName == "" {
		klog.CtxErrorf(ctx, "Word not found for word_name=%v", req.WordName)
		return &words.WordResp{BaseResp: common.BuildBaseResp(-1, "Word not found")}, nil
	}

	// Check if the word already exists using the NewWordName returned by ProcessWord
	existingWord, err := dao.QueryWordByUserIdAndName(user_id, WordDetails.NewWordName)
	if err != nil && err != dao.ErrNoRecord {
		klog.CtxErrorf(ctx, "CheckWordExists error: err=%v", err)
		// Return a database query error
		return &words.WordResp{BaseResp: common.BuildBaseResp(-1, err.Error())}, nil
	}
	// Now, either existingWord is nil, or existingWord is not nil and err is dao.ErrNoRecord

	// If the word already exists, return the existing record
	if existingWord != nil {
		klog.CtxInfof(ctx, "Word already exists for user_id=%v, req_word_name:%s word_name=%v",
			user_id, req.WordName, WordDetails.NewWordName)
		return &words.WordResp{BaseResp: common.BuildBaseResp(1, "Word already exists"),
			Word: &words.Word{
				WordId:      existingWord.WordId,
				WordName:    existingWord.WordName,
				Description: existingWord.Description,
				Explains:    existingWord.Explains,
				// Use the new PronounceInfo structure
				PronounceUs: &words.PronounceInfo{
					Phonetic: existingWord.PhoneticUs, // Read phonetic symbols from the database
					Url:      existingWord.PronounceUs,
				},
				PronounceUk: &words.PronounceInfo{
					Phonetic: existingWord.PhoneticUk, // Read phonetic symbols from the database
					Url:      existingWord.PronounceUk,
				},
				TagId: existingWord.TagId,
			},
		}, nil
	}

	wordResp := &words.Word{
		WordId:      1,
		WordName:    WordDetails.NewWordName,
		Description: "",
		Explains:    WordDetails.ExplainsYoudao,
		PronounceUs: WordDetails.PronounceUS,
		PronounceUk: WordDetails.PronounceUK,
		Sentences:   WordDetails.Sentences,
		TagId:       tag_id,
	}

	modelWord := model.Word{
		WordName:    WordDetails.NewWordName,
		Description: WordDetails.ExplainsOxford,
		Explains:    WordDetails.ExplainsYoudao,
		PronounceUs: WordDetails.PronounceUS.Url,
		PronounceUk: WordDetails.PronounceUK.Url,
		// New phonetic symbol fields
		PhoneticUs: WordDetails.PronounceUS.Phonetic,
		PhoneticUk: WordDetails.PronounceUK.Phonetic,
		UserId:     user_id,
		TagId:      tag_id,
		YoudaoUrl:  fmt.Sprintf(global.YoudaoUrl, WordDetails.NewWordName),
	}

	// Prepare relevant record data
	currentTime := time.Now().Unix()

	// Prepare answer_list data
	answerList := &model.AnswerList{
		UserId:      user_id,
		WordName:    WordDetails.NewWordName,
		Description: WordDetails.ExplainsYoudao, // Use the explains field as the description
	}

	// Prepare review record data
	reviewRecord := &model.WordsRisiteRecord{
		Level:          1,
		NextReviewTime: currentTime,
		DowngradeStep:  1,
		TotalCorrect:   0,
		TotalWrong:     0,
		Score:          0,
		UserId:         user_id,
	}

	// Call the DAO layer atomic operation method to add the word and related records
	queriedWord, err := dao.AddWordWithRelatedRecords(&modelWord, answerList, reviewRecord)
	if err != nil {
		klog.CtxErrorf(ctx, "AddWordWithRelatedRecords error: err=%v", err)
		return &words.WordResp{BaseResp: common.BuildBaseResp(-1, err.Error())}, nil
	}

	klog.CtxInfof(ctx, "Successfully added word and related records: word_id=%v, user_id=%v, word_name=%v",
		queriedWord.WordId, user_id, WordDetails.NewWordName)

	// Perform non - critical operations outside the transaction: update the review progress
	// Call getTodayReviewProgressWithInitFlag to ensure the existence of today's review progress record and get the initialization flag
	progressReq := &words.ReviewProgressReq{UserId: user_id}
	_, wasInitialized, err := getTodayReviewProgressWithInitFlag(progressReq)
	if err != nil {
		klog.CtxErrorf(ctx, "Failed to get/init today review progress for user_id=%v: %v", user_id, err)
		// Do not return an error here because the word has been successfully added, only the review progress update failed
	} else {
		if !wasInitialized {
			// If no initialization was performed (i.e., the review progress record already exists), manually increase the number of words to be reviewed
			err = dao.IncrementPendingReviewCount(user_id)
			if err != nil {
				klog.CtxErrorf(ctx, "Failed to increment pending review count for user_id=%v: %v", user_id, err)
			} else {
				klog.CtxInfof(ctx, "Successfully updated review progress for user_id=%v after adding new word", user_id)
			}
		} else {
			// If initialization was performed, it means a new user or a new day, and the newly added word is already included in the statistics
			klog.CtxInfof(ctx, "Review progress was initialized for user_id=%v, new word already included in count", user_id)
		}
	}

	// Update the WordId in the returned wordResp
	wordResp.WordId = queriedWord.WordId

	return &words.WordResp{
		BaseResp: common.BuildSuccBaseResp(),
		Word:     wordResp,
	}, nil
}

// GetWordByID retrieves a word by its unique ID from the database.
// If the word is found, it constructs and returns a response containing the word information.
// If the query fails, it returns an error response.
//
// Parameters:
//   - ctx: The context used to control the lifetime of the request.
//   - wordID: The unique identifier of the word to retrieve.
//
// Returns:
//   - *words.WordResp: A pointer to the WordResp struct containing the word information or an error message.
//   - error: An error object if an unexpected error occurs during the process.
func GetWordByID(ctx context.Context, wordID int64) (*words.WordResp, error) {
	// Query the word from the database using the provided word ID
	wordModel, err := dao.QueryWord(wordID)
	if err != nil {
		// Format an error message if the query fails
		errMsg := fmt.Sprintf("query word with wordID=%d failed , err=%s", wordID, err.Error())
		// Return an error response with the formatted error message
		return &words.WordResp{BaseResp: common.BuildBaseResp(-1, errMsg)}, nil
	}

	// Convert the model.Word object to a words.Word object
	word := &words.Word{
		WordId:   wordModel.WordId,
		WordName: wordModel.WordName,
		// Description is initially set to an empty string
		Description: "",
		Explains:    wordModel.Explains,
		// Use the new PronounceInfo structure for US pronunciation
		PronounceUs: &words.PronounceInfo{
			// Phonetic information is temporarily set to empty
			Phonetic: "",
			Url:      wordModel.PronounceUs,
		},
		// Use the new PronounceInfo structure for UK pronunciation
		PronounceUk: &words.PronounceInfo{
			// Phonetic information is temporarily set to empty
			Phonetic: "",
			Url:      wordModel.PronounceUk,
		},
		// Add the tag ID from the model
		TagId: wordModel.TagId,
		// Set the current level to 0
		Level: 0,
		// Set the maximum level using the global constant
		MaxLevel: global.MAX_RISITE_LEVEL,
	}

	// Return a successful response with the word information
	return &words.WordResp{BaseResp: common.BuildSuccBaseResp(), Word: word}, nil
}

// model2Word converts a pointer to a model.Word struct to a pointer to a words.Word struct.
// It maps fields from the model layer to the service layer.
//
// Parameters:
//   - w: A pointer to the model.Word struct to be converted.
//
// Returns:
//   - *words.Word: A pointer to the converted words.Word struct.
func model2Word(w *model.Word) *words.Word {
	return &words.Word{
		WordName:    w.WordName,
		Description: w.Description,
		// Note: Explains is incorrectly mapped to Description, which might be a bug
		Explains: w.Description,
		PronounceUs: &words.PronounceInfo{
			Phonetic: w.PhoneticUs,
			Url:      w.PronounceUs,
		},
		PronounceUk: &words.PronounceInfo{
			Phonetic: w.PhoneticUk,
			Url:      w.PronounceUk,
		},
		TagId: w.TagId,
	}
}

// cleanupWord cleans up a word by removing whitespace, newline characters,
// carriage return characters, spaces, and tab characters.
// It is useful for standardizing input words before processing.
//
// Parameters:
//   - word: The input word string to be cleaned.
//
// Returns:
//   - string: The cleaned word string.
func cleanupWord(word string) string {
	// Remove leading and trailing whitespace
	word = strings.TrimSpace(word)
	// Remove newline characters
	word = strings.ReplaceAll(word, "\n", "")
	// Remove carriage return characters
	word = strings.ReplaceAll(word, "\r", "")
	// Remove all spaces
	word = strings.ReplaceAll(word, " ", "")
	// Remove tab characters
	word = strings.ReplaceAll(word, "\t", "")
	return word
}
