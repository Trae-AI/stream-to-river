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

// TestGetWordListHandler tests the GetWordListHandler function.
func TestGetWordListHandler(t *testing.T) {
	recreateMockWordTable()

	wordName := "Challenge"
	userId := int64(1)
	tagId := int32(1)

	wlr := &words.WordListReq{
		UserId: userId,
		Num:    global.WORDS_NUM_PER_PAGE_DEFAULT,
		Offset: 0,
	}
	// 1. no word
	resp1, err := GetWordListHandler(context.Background(), wlr)
	test.Assert(t, err == nil, err)
	test.Assert(t, resp1.BaseResp.StatusCode == 0, resp1.BaseResp.StatusCode)
	test.Assert(t, len(resp1.WordsList) == 0)

	awr := &words.AddWordReq{
		UserId:   userId,
		WordName: wordName,
		TagId:    tagId,
	}
	resp2, err := AddNewWord(context.Background(), awr)
	test.Assert(t, err == nil, err)
	test.Assert(t, resp2.BaseResp.StatusCode == 0, resp2.BaseResp.StatusCode)

	// 2. has word
	resp1, err = GetWordListHandler(context.Background(), wlr)
	test.Assert(t, err == nil, err)
	test.Assert(t, resp1.BaseResp.StatusCode == 0, resp1.BaseResp.StatusCode)
	test.Assert(t, len(resp1.WordsList) == 1)
	test.Assert(t, len(resp1.WordsList) == 1)
}
