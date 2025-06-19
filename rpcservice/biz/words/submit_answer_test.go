// Copyright (c) 2025 Bytedance Ltd. and/or its affiliates
// SPDX-License-Identifier: MIT

package words

import (
	"context"
	"testing"

	"github.com/Trae-AI/stream-to-river/internal/test"
	"github.com/Trae-AI/stream-to-river/rpcservice/global"
	"github.com/Trae-AI/stream-to-river/rpcservice/kitex_gen/words"
)

func TestSubmitAnswer(t *testing.T) {
	recreateMockWordTable()

	ctx := context.Background()
	userId := int64(1)
	wordId := int64(2)
	answerId := int64(3)
	questionType := global.FILL_IN_BLANK
	filledName := "test"

	// mock no word record
	resp, err := SubmitAnswer(ctx, userId, wordId, answerId, questionType, &filledName)
	test.Assert(t, err == nil)
	test.Assert(t, resp.BaseResp.StatusCode == -1)

	// mock has word record
	wordName := "Challenge"
	tagId := int32(1)
	awr := &words.AddWordReq{UserId: userId, WordName: wordName, TagId: tagId}
	resp2, err := AddNewWord(context.Background(), awr)
	test.Assert(t, err == nil, err)
	resp, err = SubmitAnswer(ctx, userId, resp2.Word.WordId, answerId, questionType, &filledName)
	test.Assert(t, err == nil)
	test.Assert(t, resp.BaseResp.StatusCode == 0)

	resp, err = SubmitAnswer(ctx, userId, resp2.Word.WordId, answerId, global.CHOOSE_CN, &filledName)
	test.Assert(t, err == nil)
	test.Assert(t, resp.BaseResp.StatusCode == 0)
}
