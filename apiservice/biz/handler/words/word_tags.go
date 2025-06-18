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

// GetSupportedTagsHandler handles HTTP GET requests to retrieve the supported tag types.
// It is registered at the `/api/tags` endpoint.
// This function creates an RPC request, calls the RPC service to get the supported tags,
// and returns the RPC response to the client. If an error occurs during the RPC call,
// it returns an internal server error response.
//
// @router /api/tags [GET]
func GetSupportedTagsHandler(ctx context.Context, c *app.RequestContext) {
	// Create an RPC request
	req := &words.GetTagsReq{}

	// Call the RPC service
	resp, err := rpcclient.WordsRPCCli.GetSupportedTags(ctx, req)
	if err != nil {
		// Log the error and return an internal server error response
		c.JSON(http.StatusInternalServerError, utils.HertzErrorResp("GetSupportedTags failed", err.Error()))
		return
	}

	// Return the RPC response
	c.JSON(http.StatusOK, resp)
}

// UpdateWordTagID handles HTTP POST requests to update the tag of a word.
// It is registered at the `/api/word-tag` endpoint.
// The function first retrieves the user ID from the request context.
// If the retrieval is successful, it binds and validates the request data,
// sets the user ID in the request, calls the RPC service to update the word tag,
// and returns the RPC response. If any step fails, it returns an appropriate error response.
//
// @router /api/word-tag [POST]
func UpdateWordTagID(ctx context.Context, c *app.RequestContext) {
	// Retrieve the user ID from the request context
	userId, err := user.GetUserIDFromContext(c)
	if err != nil {
		// Log the error and return an unauthorized response
		hlog.CtxErrorf(ctx, "UpdateWordTagID failed, err: %v", err)
		c.JSON(http.StatusUnauthorized, utils.HertzErrorResp("failed to get userId", err.Error()))
		return
	}

	// Bind and validate the request data
	var req words.UpdateWordTag
	err = c.BindAndValidate(&req)
	if err != nil {
		// Return a bad request response
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	// Set the user ID in the request
	req.UserId = userId

	// Call the RPC service to update the word tag
	resp, err := rpcclient.WordsRPCCli.UpdateWordTagID(ctx, &req)
	if err != nil {
		// Log the error and return an internal server error response
		c.JSON(http.StatusInternalServerError, utils.HertzErrorResp("update word tag failed", err.Error()))
		return
	}
	// Return the RPC response
	c.JSON(http.StatusOK, resp)
}
