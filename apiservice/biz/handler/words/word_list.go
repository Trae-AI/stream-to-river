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

// WordsHandler handles HTTP GET requests to retrieve a list of words.
// It is registered at the `/api/word-list` endpoint.
// This function first extracts the user ID from the request context.
// Then it binds and validates the request parameters, makes an RPC call to fetch the word list,
// and returns an appropriate HTTP response based on the outcome.
//
// @router /api/word-list [GET]
func WordsHandler(ctx context.Context, c *app.RequestContext) {
	// Extract the user ID from the request context.
	userId, err := user.GetUserIDFromContext(c)
	if err != nil {
		// Log the error if user ID extraction fails and return a 401 Unauthorized response.
		hlog.CtxErrorf(ctx, "WordsHandler failed, err: %v", err)
		c.JSON(http.StatusUnauthorized, utils.HertzErrorResp("failed to get userId", err.Error()))
		return
	}

	// Bind and validate the request parameters to a WordListReq struct.
	var req words.WordListReq
	if err = c.BindAndValidate(&req); err != nil {
		// Log the error if binding or validation fails and return a 400 Bad Request response.
		hlog.CtxErrorf(ctx, "parse param for [/api/word-list] failed, err: %v", err)
		c.JSON(http.StatusBadRequest, utils.HertzErrorResp("parse param failed", err.Error()))
		return
	}
	// Set the user ID in the request.
	req.UserId = userId

	// Log the incoming request.
	hlog.Infof("GetWordList Request: userId:%d offset:%d", req.UserId, req.Offset)
	// Make an RPC call to fetch the word list.
	resp, err := rpcclient.WordsRPCCli.GetWordList(ctx, &req)
	if err != nil {
		// Log the error if the RPC call fails and return a 500 Internal Server Error response.
		hlog.CtxErrorf(ctx, "rpc call failed - GetWordList, err: %v", err)
		c.JSON(http.StatusInternalServerError, utils.HertzErrorResp("GetWordList failed", err.Error()))
		return
	}
	// Log the size of the retrieved word list.
	hlog.CtxDebugf(ctx, "GetWordList Response: word list size=%d", len(resp.WordsList))

	// If everything succeeds, return a 200 OK response with the word list.
	c.JSON(http.StatusOK, resp)
}
