// Copyright (c) 2025 Bytedance Ltd. and/or its affiliates
// SPDX-License-Identifier: MIT

package chat

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/bytedance/sonic"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/cloudwego/hertz/pkg/protocol/sse"
	"github.com/google/uuid"

	"github.com/Trae-AI/stream-to-river/apiservice/biz/rpcclient"
	"github.com/Trae-AI/stream-to-river/apiservice/biz/utils"
	"github.com/Trae-AI/stream-to-river/rpcservice/kitex_gen/words"
)

// ChatHandler handles HTTP GET requests for the chat service. It uses Server - Sent Events (SSE)
// to stream chat responses back to the client. The handler fetches query parameters from the request,
// initiates a gRPC streaming call to the chat service, and streams the responses to the client.
//
// @router /api/chat?q=xxx&conversation_id=y [GET]
// eg: http://localhost:8889/api/chat?q=英文介绍一下火山引擎
func ChatHandler(ctx context.Context, c *app.RequestContext) {
	newCtx, cancel := context.WithCancel(ctx)

	// Initiate a gRPC streaming call to the chat service
	st, sErr := rpcclient.WordsRPCCli.Chat(newCtx, &words.ChatReq{
		QueryMsg:       c.Query("q"),
		ConversationId: getConversationId(c),
	})

	// Handle errors during the gRPC stream creation
	if sErr != nil {
		cancel()
		errMsg := fmt.Sprintf("kitex create Stream failed, err: %v", sErr)
		hlog.CtxErrorf(newCtx, errMsg)
		c.JSON(http.StatusInternalServerError, utils.HertzErrorResp("chat service failed", errMsg))
		return
	}
	defer cancel()

	// Initialize SSE writer and set response content type
	w := sse.NewWriter(c)
	c.Response.Header.SetContentType("text/event-stream; charset=utf-8")

	var id int
	for {
		id++
		// Receive a response from the gRPC stream
		resp, rErr := st.Recv(ctx)
		if rErr != nil {
			if rErr == io.EOF {
				return
			}
			errMsg := fmt.Sprintf("kitex streaming recv err: %v", rErr)
			hlog.CtxErrorf(newCtx, errMsg)
			if wErr := w.WriteEvent(strconv.Itoa(id), "error", []byte(errMsg)); wErr != nil {
				hlog.CtxErrorf(newCtx, "write event error: %v", wErr)
			}
			return
		}
		// Marshal the response to JSON and send it as an SSE event
		respBytes, _ := sonic.Marshal(resp)
		if wErr := w.WriteEvent(strconv.Itoa(id), "message", respBytes); wErr != nil {
			hlog.CtxErrorf(newCtx, "write event error: %v", wErr)
			return
		}
	}
}

// getConversationId retrieves the conversation ID from the request query parameters.
// If the conversation ID is not provided, it generates a new UUID as the conversation ID.
//
// Parameters:
//   - c: A pointer to the app.RequestContext provided by Hertz, containing request information.
//
// Returns:
//   - string: The conversation ID.
func getConversationId(c *app.RequestContext) string {
	convId := c.Query("conversation_id")
	if len(convId) > 0 {
		return convId
	}
	return uuid.New().String()
}
