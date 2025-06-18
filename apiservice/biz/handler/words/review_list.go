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

// GetReviewWordListHandler is an HTTP handler function designed to handle GET requests for retrieving a review word list.
// It is registered at the `/api/review-list` endpoint.
// This function first extracts the user ID from the request context, then makes an RPC call to fetch the review word list.
// Based on the results of these operations, it returns appropriate HTTP responses to the client.
//
// @router /api/review-list [GET]
func GetReviewWordListHandler(ctx context.Context, c *app.RequestContext) {
	// Attempt to extract the user ID from the request context.
	userId, err := user.GetUserIDFromContext(c)
	if err != nil {
		// If the extraction fails, log the error and return a 401 Unauthorized response.
		hlog.CtxErrorf(ctx, "GetReviewWordListHandler failed, err: %v", err)
		c.JSON(http.StatusUnauthorized, utils.HertzErrorResp("failed to get userId", err.Error()))
		return
	}

	// Log the incoming request with the user ID.
	hlog.CtxInfof(ctx, "GetReviewWordList Request: userId:%d", userId)
	// Make an RPC call to the WordsRPCCli to get the review word list.
	resp, err := rpcclient.WordsRPCCli.GetReviewWordList(ctx, &words.ReviewListReq{UserId: userId})
	if err != nil {
		// If the RPC call fails, log the error and return a 502 Bad Gateway response.
		hlog.CtxErrorf(ctx, "rpc call failed - GetReviewWordList, err: %v", err)
		c.JSON(http.StatusBadGateway, utils.HertzErrorResp("GetReviewWordList failed", err.Error()))
		return
	}
	// Log the response received from the RPC call.
	hlog.CtxDebugf(ctx, "GetReviewWordList Response: %v", resp)

	// If everything succeeds, return a 200 OK response with the review word list.
	c.JSON(http.StatusOK, resp)
}
