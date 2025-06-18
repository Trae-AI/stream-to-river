// Copyright (c) 2025 Bytedance Ltd. and/or its affiliates
// SPDX-License-Identifier: MIT

package words

import (
	"context"
	"net/http"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/hlog"

	"github.com/Trae-AI/stream-to-river/apiservice/biz/rpcclient"
	"github.com/Trae-AI/stream-to-river/apiservice/biz/user"
	"github.com/Trae-AI/stream-to-river/apiservice/biz/utils"
	"github.com/Trae-AI/stream-to-river/rpcservice/kitex_gen/words"
)

// GetTodayReviewProgressHandler handles HTTP GET requests to retrieve today's review progress.
// It is registered at the `/api/review-progress` endpoint.
// The function first extracts the user ID from the request context.
// If the extraction is successful, it makes an RPC call to get the review progress.
// Based on the outcome of these operations, it returns an appropriate HTTP response.
//
// @router /api/review-progress [GET]
func GetTodayReviewProgressHandler(ctx context.Context, c *app.RequestContext) {
	// Extract the user ID from the request context.
	userId, err := user.GetUserIDFromContext(c)
	if err != nil {
		// Log the error if user ID extraction fails and return a 401 Unauthorized response.
		hlog.CtxErrorf(ctx, "GetTodayReviewProgressHandler failed, err: %v", err)
		c.JSON(http.StatusUnauthorized, utils.HertzErrorResp("failed to get userId", err.Error()))
		return
	}

	// Make an RPC call to get today's review progress.
	resp, err := rpcclient.WordsRPCCli.GetTodayReviewProgress(ctx, &words.ReviewProgressReq{UserId: userId})
	if err != nil {
		// Log the error if the RPC call fails and return a 500 Internal Server Error response.
		c.JSON(http.StatusInternalServerError, utils.HertzErrorResp("GetTodayReviewProgress failed", err.Error()))
		return
	}

	// If everything succeeds, return a 200 OK response with the review progress data.
	c.JSON(http.StatusOK, resp)
}
