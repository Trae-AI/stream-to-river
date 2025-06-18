// Copyright (c) 2025 Bytedance Ltd. and/or its affiliates
// SPDX-License-Identifier: MIT

package words

import (
	"context"
	"testing"

	"github.com/Trae-AI/stream-to-river/internal/test"
	"github.com/Trae-AI/stream-to-river/rpcservice/dal/model"
	"github.com/Trae-AI/stream-to-river/rpcservice/kitex_gen/words"
)

// TestGetReviewWordList tests the GetReviewWordList function.
func TestGetReviewWordList(t *testing.T) {
	recreateMockWordTable()
	userId := int64(100)

	// 1. mock table not exist err
	mockDB.Migrator().DropTable(&model.WordsRisiteRecord{})
	_, err := GetReviewWordList(context.Background(), userId)
	test.Assert(t, err != nil)

	// 2. mock no record
	recreateMockWordReviewRecordTable()
	resp, err := GetReviewWordList(context.Background(), userId)
	test.Assert(t, err == nil)
	test.Assert(t, resp.TotalNum == "")

	// 3. mock has record
	wordName := "Challenge"
	tagId := int32(1)
	awr := &words.AddWordReq{UserId: userId, WordName: wordName, TagId: tagId}
	_, err = AddNewWord(context.Background(), awr)
	test.Assert(t, err == nil, err)
	resp, err = GetReviewWordList(context.Background(), userId)
	test.Assert(t, err == nil)
	test.Assert(t, resp.TotalNum != "")
	test.Assert(t, len(resp.Questions) == 4)

}
