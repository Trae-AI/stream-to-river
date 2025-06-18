// Copyright (c) 2025 Bytedance Ltd. and/or its affiliates
// SPDX-License-Identifier: MIT

package words

import (
	"context"
	"fmt"

	"github.com/cloudwego/kitex/pkg/klog"

	"github.com/Trae-AI/stream-to-river/rpcservice/biz/words/common"
	"github.com/Trae-AI/stream-to-river/rpcservice/dal/dao"
	"github.com/Trae-AI/stream-to-river/rpcservice/kitex_gen/base"
	"github.com/Trae-AI/stream-to-river/rpcservice/kitex_gen/words"
)

// GetSupportedTags retrieves all supported word tag types from the database.
// It queries the database for all word tags and converts them into the response format.
// If the query fails, it logs the error and returns a response indicating the failure.
//
// Parameters:
//   - ctx: The context used to control the lifetime of the request.
//   - req: A pointer to the GetTagsReq struct containing the request information.
//
// Returns:
//   - *words.GetTagsResp: A pointer to the GetTagsResp struct containing the list of tags and the base response.
//   - error: An error object if an unexpected error occurs during the process.
func GetSupportedTags(ctx context.Context, req *words.GetTagsReq) (*words.GetTagsResp, error) {
	// Retrieve all word tags from the database
	tags, err := dao.GetAllWordTags()
	if err != nil {
		// Log the error if the retrieval fails
		klog.CtxErrorf(ctx, "GetAllWordTags failed: %v", err)
		return &words.GetTagsResp{
			BaseResp: &base.BaseResp{
				StatusCode:    -1,
				StatusMessage: "获取标签失败",
			},
		}, nil
	}

	// Convert the retrieved tags into the response format
	var tagInfos []*words.TagInfo
	for _, tag := range tags {
		tagInfo := &words.TagInfo{
			TagId:         tag.TagId,
			TagName:       tag.TagName,
			QuestionTypes: int32(tag.QuestionTypes),
			MaxScore:      int32(tag.MaxScore),
		}
		tagInfos = append(tagInfos, tagInfo)
	}

	// Log the success message with the number of retrieved tags
	klog.Debugf("GetSupportedTags success, count: %d", len(tagInfos))
	return &words.GetTagsResp{
		Tags:     tagInfos,
		BaseResp: common.BuildSuccBaseResp(),
	}, nil
}

// UpdateWordTagID updates the tag ID of a word in the database.
// It first attempts to update the tag ID for the given word and user.
// If the update is successful, it queries the word to verify the change.
// If any step fails, it logs the error and returns a response indicating the failure.
//
// Parameters:
//   - ctx: The context used to control the lifetime of the request.
//   - req: A pointer to the UpdateWordTag struct containing the word ID, user ID, and new tag ID.
//
// Returns:
//   - *words.WordResp: A pointer to the WordResp struct containing the base response and the updated word information.
//   - error: An error object if an unexpected error occurs during the process.
func UpdateWordTagID(ctx context.Context, req *words.UpdateWordTag) (resp *words.WordResp, err error) {
	// Update the tag ID for the given word and user
	err = dao.UpdateWordTagID(req.WordId, req.UserId, req.TagId)
	if err != nil {
		// Format and log the error message if the update fails
		errMsg := fmt.Sprintf("update tagID with wordID=%d and userID=%d failed, err=%s", req.WordId, req.UserId, err.Error())
		klog.CtxErrorf(ctx, "%s", errMsg)
		return &words.WordResp{BaseResp: common.BuildBaseResp(-1, errMsg)}, nil
	}

	// Query the word after the update
	wordModel, err := dao.QueryWord(req.WordId)
	if err != nil {
		// Format and log the error message if the query fails
		errMsg := fmt.Sprintf("query word with wordID=%d failed, err=%v", req.WordId, err)
		klog.CtxErrorf(ctx, "%s", errMsg)
		return &words.WordResp{BaseResp: common.BuildBaseResp(-1, errMsg)}, nil
	}

	// Check if the queried word is nil
	if wordModel == nil {
		// Format and log the error message if the word is not found
		errMsg := fmt.Sprintf("query word with wordID=%d failed, word not found", req.WordId)
		klog.CtxErrorf(ctx, "%s", errMsg)
		return &words.WordResp{BaseResp: common.BuildBaseResp(-1, errMsg)}, nil
	}

	// Return a successful response with the updated word information
	return &words.WordResp{
		BaseResp: common.BuildSuccBaseResp(),
		Word:     model2Word(wordModel),
	}, nil
}
