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

// SubmitAnswerHandler handles HTTP POST requests to submit word answers.
// It is registered at the `/api/answer` endpoint. The function first extracts the user ID
// from the request context, then binds and validates the request data. After that, it makes
// an RPC call to the backend service to submit the answer. Finally, it returns an appropriate
// HTTP response based on the result of these operations.
//
// @router /api/answer [POST]
func SubmitAnswerHandler(ctx context.Context, c *app.RequestContext) {
	// 1. Extract the user ID from the request context.
	userId, err := user.GetUserIDFromContext(c)
	if err != nil {
		// Log the error if user ID extraction fails and return a 401 Unauthorized response.
		hlog.CtxErrorf(ctx, "SubmitAnswerHandler failed, err: %v", err)
		c.JSON(http.StatusUnauthorized, utils.HertzErrorResp("failed to get userId", err.Error()))
		return
	}

	var req words.SubmitAnswerReq
	// 2. Bind and validate the request data.
	err = c.BindAndValidate(&req)
	if err != nil {
		// If binding or validation fails, return a 400 Bad Request response.
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	req.UserId = userId

	// Log the incoming request.
	hlog.CtxInfof(ctx, "SubmitAnswer Request: %v", req)
	// 3. Make an RPC call to submit the answer.
	resp, err := rpcclient.WordsRPCCli.SubmitAnswer(ctx, &req)
	if err != nil {
		// Log the error if the RPC call fails and return a 500 Internal Server Error response.
		hlog.CtxInfof(ctx, "rpc call failed - SubmitAnswer, err: %v", err)
		c.JSON(http.StatusInternalServerError, utils.HertzErrorResp("SubmitAnswer failed", err.Error()))
		return
	}
	// Log the response from the RPC call.
	hlog.CtxDebugf(ctx, "SubmitAnswer Response: %v", resp)

	// 4. If everything succeeds, return a 200 OK response with the answer submission result.
	c.JSON(http.StatusOK, resp)
}
