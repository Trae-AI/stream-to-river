// Copyright (c) 2025 Bytedance Ltd. and/or its affiliates
// SPDX-License-Identifier: MIT

package words

import (
	"context"
	"fmt"

	"github.com/cloudwego/kitex/pkg/klog"

	"github.com/Trae-AI/stream-to-river/rpcservice/biz/words/common"
	"github.com/Trae-AI/stream-to-river/rpcservice/dal/dao"
	"github.com/Trae-AI/stream-to-river/rpcservice/dal/model"
	"github.com/Trae-AI/stream-to-river/rpcservice/global"
	"github.com/Trae-AI/stream-to-river/rpcservice/kitex_gen/words"
)

// GetWordListHandler retrieves a paginated list of words for a specific user from the database.
// It also fetches the review records for these words and combines them into the response format.
//
// Parameters:
//   - ctx: The context used to control the lifetime of the request.
//   - req: A pointer to the WordListReq struct containing the user ID, number of words per page, and offset.
//
// Returns:
//   - *words.WordListResp: A pointer to the WordListResp struct containing the list of words and the base response.
//   - error: An error object if an unexpected error occurs during the process.
func GetWordListHandler(ctx context.Context, req *words.WordListReq) (resp *words.WordListResp, err error) {
	// Initialize the response object
	resp = words.NewWordListResp()

	// Set default values for pagination parameters
	num := req.Num
	if num == 0 {
		// Use the global default value if the number of words per page is not provided
		num = global.WORDS_NUM_PER_PAGE_DEFAULT
	}

	offset := req.Offset
	if offset < 0 {
		// Set offset to 0 if it is negative
		offset = 0
	}

	// Retrieve the word list from the database with pagination
	wordModels, err := dao.GetWordsByUserIdWithPagination(req.UserId, offset, num)
	if err != nil {
		// Log the error and return a failure response if the database query fails
		resp.BaseResp = common.BuildBaseResp(-1, fmt.Sprintf("Failed to get words: %v", err))
		return resp, nil
	}

	// Extract all word IDs from the retrieved word list for batch querying review records
	var wordIds []int64
	for _, wordModel := range wordModels {
		wordIds = append(wordIds, wordModel.WordId)
	}

	// Batch query review records for the extracted word IDs
	var recordMap map[int64]*model.WordsRisiteRecord
	if len(wordIds) > 0 {
		recordMap, err = dao.GetWordsRisiteRecordsByUserAndWordIds(req.UserId, wordIds)
		if err != nil {
			// Log the error if the batch query fails, but continue the main process
			klog.CtxErrorf(ctx, "Failed to get risite records: %v", err)
			recordMap = make(map[int64]*model.WordsRisiteRecord)
		}
	} else {
		// Initialize an empty map if there are no word IDs
		recordMap = make(map[int64]*model.WordsRisiteRecord)
	}

	// Convert the retrieved word models and review records into the response format
	var wordsList []*words.Word
	for _, wordModel := range wordModels {
		// Set the default review level to 0
		var level int32 = 0
		if record, exists := recordMap[wordModel.WordId]; exists {
			// Update the review level if the review record exists
			level = int32(record.Level)
		}

		wordsList = append(wordsList, &words.Word{
			WordId:   wordModel.WordId,
			WordName: wordModel.WordName,
			// Description is set to an empty string
			Description: "",
			Explains:    wordModel.Explains,
			// Use the new PronounceInfo structure for US pronunciation
			PronounceUs: &words.PronounceInfo{
				// Read the phonetic symbol from the database
				Phonetic: wordModel.PhoneticUs,
				Url:      wordModel.PronounceUs,
			},
			// Use the new PronounceInfo structure for UK pronunciation
			PronounceUk: &words.PronounceInfo{
				// Read the phonetic symbol from the database
				Phonetic: wordModel.PhoneticUk,
				Url:      wordModel.PronounceUk,
			},
			// Add the tag ID from the model
			TagId: wordModel.TagId,
			// Add the current review level
			Level: level,
			// Add the maximum review level from the global configuration
			MaxLevel: global.MAX_RISITE_LEVEL,
		})
	}

	// Set the word list and a successful base response in the final response
	resp.WordsList = wordsList
	resp.BaseResp = common.BuildSuccBaseResp()
	return resp, nil
}
