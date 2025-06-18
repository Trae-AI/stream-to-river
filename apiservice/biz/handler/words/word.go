// Copyright (c) 2025 Bytedance Ltd. and/or its affiliates
// SPDX-License-Identifier: MIT

package words

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/hlog"

	"github.com/Trae-AI/stream-to-river/apiservice/biz/rpcclient"
	"github.com/Trae-AI/stream-to-river/apiservice/biz/user"
	"github.com/Trae-AI/stream-to-river/apiservice/biz/utils"
	"github.com/Trae-AI/stream-to-river/rpcservice/kitex_gen/words"
)

// AddNewWordHandler handles HTTP POST requests to add a new word.
// It is registered at the `/api/word-add` endpoint.
//
// The function first retrieves the user ID from the request context.
// If the retrieval fails, it logs the error and returns a 401 Unauthorized response.
// Then it binds and validates the request data. If binding fails, it sets a default tag ID.
// After that, it cleans up the word name, checks if it's empty, and makes an RPC call to add the word.
// Based on the result of the RPC call, it returns an appropriate HTTP response.
//
// @router /api/word-add [POST]
func AddNewWordHandler(ctx context.Context, c *app.RequestContext) {
	// Retrieve user ID from the request context
	userId, err := user.GetUserIDFromContext(c)
	if err != nil {
		hlog.CtxErrorf(ctx, "AddNewWordHandler failed, err: %v", err)
		c.JSON(http.StatusUnauthorized, utils.HertzErrorResp("failed to get userId", err.Error()))
		return
	}

	// Bind and validate the request data
	var req words.AddWordReq
	err = c.BindAndValidate(&req)
	if err != nil {
		// Set default tag ID if binding fails
		hlog.CtxErrorf(ctx, "AddNewWordHandler bind failed, err=%v", err)
		req.TagId = 1
	}
	req.UserId = userId

	// Clean up the word name
	wordName := cleanupWord(req.WordName)
	if wordName == "" {
		hlog.CtxErrorf(ctx, "AddNewWord failed, word is empty")
		c.JSON(http.StatusBadRequest, utils.HertzErrorResp("word is empty", ""))
		return
	}

	// Log the request
	hlog.CtxInfof(ctx, "AddNewWord Request: %v", req)
	// Make an RPC call to add the word
	resp, err := rpcclient.WordsRPCCli.AddNewWord(ctx, &req)
	if err != nil {
		hlog.CtxErrorf(ctx, "rpc call failed - AddNewWord, err: %v", err)
		c.JSON(http.StatusInternalServerError, utils.HertzErrorResp(fmt.Sprintf("add new word %s failed", req.WordName), err.Error()))
		return
	}
	// Log the response
	hlog.CtxDebugf(ctx, "AddNewWord Response: %v", resp)

	// Return the RPC response
	c.JSON(http.StatusOK, resp)
}

// WordDetailHandler handles HTTP GET requests to retrieve word details.
// It is registered at the `/api/word-detail` endpoint.
//
// The function first cleans up the word name from the query parameter.
// If the word name is empty, it logs the error and returns a 400 Bad Request response.
// Then it makes an RPC call to get the word details.
// Based on the result of the RPC call, it returns an appropriate HTTP response.
//
// @router /api/word-detail [GET]
func WordDetailHandler(ctx context.Context, c *app.RequestContext) {
	// Clean up the word name from the query parameter
	wordName := cleanupWord(c.Query("word"))
	if wordName == "" {
		hlog.CtxErrorf(ctx, "word is empty")
		c.JSON(http.StatusBadRequest, utils.HertzErrorResp("word is empty", ""))
		return
	}

	// Make an RPC call to get the word details
	resp, err := rpcclient.WordsRPCCli.GetWordDetail(ctx, wordName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.HertzErrorResp("GetWordDetail failed", err.Error()))
		return
	}
	// Return the RPC response
	c.JSON(http.StatusOK, resp)
}

// QueryWordHandler handles HTTP GET requests to query a word by its ID.
// It is registered at the `/api/word-query` endpoint.
//
// The function first checks if the word ID query parameter is empty.
// If it is, it returns a 400 Bad Request response.
// Then it parses the word ID string to an integer.
// If parsing fails, it logs the error and returns a 400 Bad Request response.
// After that, it makes an RPC call to query the word.
// Based on the result of the RPC call, it returns an appropriate HTTP response.
//
// @router /api/word-query [GET]
func QueryWordHandler(ctx context.Context, c *app.RequestContext) {
	// Get the word ID query parameter
	var wordIdStr = c.Query("word_id")
	if wordIdStr == "" {
		c.JSON(http.StatusBadRequest, utils.HertzErrorResp("word_id is empty", ""))
		return
	}

	// Parse the word ID string to an integer
	wordId, err := strconv.ParseInt(wordIdStr, 10, 64)
	if err != nil {
		hlog.CtxErrorf(ctx, "parse offset failed, err: %v", err)
		c.JSON(http.StatusBadRequest, utils.HertzErrorResp("parse offset failed", err.Error()))
		return
	}

	// Make an RPC call to query the word
	resp, err := rpcclient.WordsRPCCli.GetWordByID(ctx, wordId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.HertzErrorResp("queryWord failed", err.Error()))
		return
	}

	// Return the RPC response
	c.JSON(http.StatusOK, resp)
}

// cleanupWord cleans up a word by removing whitespace, newline characters,
// carriage return characters, spaces, and tab characters.
// It is useful for standardizing input words before processing.
//
// Parameters:
//   - word: The input word string to be cleaned.
//
// Returns:
//   - string: The cleaned word string.
func cleanupWord(word string) string {
	word = strings.TrimSpace(word)            // Remove leading and trailing whitespace
	word = strings.ReplaceAll(word, "\n", "") // Remove newline characters
	word = strings.ReplaceAll(word, "\r", "") // Remove carriage return characters
	word = strings.ReplaceAll(word, " ", "")  // Remove spaces
	word = strings.ReplaceAll(word, "\t", "") // Remove tab characters
	return word
}
