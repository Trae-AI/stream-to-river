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

func TestGetSupportedTags(t *testing.T) {
	wt := &model.WordTag{TagId: 1, TagName: "test_tag", QuestionTypes: 1, MaxScore: 10}
	dao.AddWordTag(wt)

	ctx := context.Background()
	req := &words.GetTagsReq{}
	resp, err := GetSupportedTags(ctx, req)
	test.Assert(t, err == nil, err)
	test.Assert(t, resp.BaseResp.StatusCode == 0)
	test.Assert(t, resp.BaseResp.StatusMessage == "success")
	test.Assert(t, resp.Tags[0].TagName == "test_tag")

	mockDB.Migrator().DropTable(model.WordTag{})
	// mock table not exist err
	resp, err = GetSupportedTags(ctx, req)
	test.Assert(t, err == nil, err)
	test.Assert(t, resp.BaseResp.StatusCode == -1)
	test.Assert(t, resp.BaseResp.StatusMessage == "获取标签失败")
}

func TestUpdateWordTagID(t *testing.T) {
	oldTagID := 1
	newTagID := 2
	var userID int64 = 2
	ctx := context.Background()
	req1 := &words.AddWordReq{UserId: userID, WordName: "pragmatic", TagId: int32(oldTagID)}
	resp, err := AddNewWord(ctx, req1)
	test.Assert(t, err == nil)
	test.Assert(t, resp.Word.TagId == int32(oldTagID))

	req2 := &words.UpdateWordTag{WordId: resp.Word.WordId, UserId: userID, TagId: int32(newTagID)}
	resp, err = UpdateWordTagID(ctx, req2)
	test.Assert(t, err == nil, err)
	test.Assert(t, resp.BaseResp.StatusCode == 0)
	test.Assert(t, resp.BaseResp.StatusMessage == "success")
	test.Assert(t, resp.Word.TagId == int32(newTagID))

	mockDB.Migrator().DropTable(model.Word{})
	// mock table not exist err
	resp, err = UpdateWordTagID(ctx, req2)
	test.Assert(t, err == nil, err)
	test.Assert(t, resp.BaseResp.StatusCode == -1)

}
