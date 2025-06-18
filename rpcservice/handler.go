// Copyright (c) 2025 Bytedance Ltd. and/or its affiliates
// SPDX-License-Identifier: MIT

package main

import (
	"context"

	"github.com/cloudwego/kitex/pkg/klog"

	"github.com/Trae-AI/stream-to-river/rpcservice/biz/chat"
	biz_words "github.com/Trae-AI/stream-to-river/rpcservice/biz/words"
	"github.com/Trae-AI/stream-to-river/rpcservice/kitex_gen/words"
)

// WordServiceImpl implements the last service interface defined in the IDL.
// It serves as the implementation of the RPC service for word - related operations.
type WordServiceImpl struct{}

// GetWordList implements the WordServiceImpl interface.
// It calls the GetWordListHandler function in the biz_words package to handle requests for getting a word list.
//
// Parameters:
//   - ctx: The context for the request, used for cancellation, deadlines, and passing request - scoped values.
//   - req: A pointer to the WordListReq struct containing the request information.
//
// Returns:
//   - *words.WordListResp: A pointer to the WordListResp struct containing the response information.
//   - error: An error object if an unexpected error occurs during the process.
func (s *WordServiceImpl) GetWordList(ctx context.Context, req *words.WordListReq) (resp *words.WordListResp, err error) {
	return biz_words.GetWordListHandler(ctx, req)
}

// Chat implements the WordServiceImpl interface.
// It calls the ChatHandler function in the chat package to handle chat - related requests with streaming.
//
// Parameters:
//   - ctx: The context for the request, used for cancellation, deadlines, and passing request - scoped values.
//   - req: A pointer to the ChatReq struct containing the request information.
//   - stream: A WordService_ChatServer interface for handling streaming responses.
//
// Returns:
//   - error: An error object if an unexpected error occurs during the process.
func (s *WordServiceImpl) Chat(ctx context.Context, req *words.ChatReq, stream words.WordService_ChatServer) (err error) {
	return chat.ChatHandler(ctx, req, stream)
}

// AddNewWord implements the WordServiceImpl interface.
// It calls the AddNewWord function in the biz_words package to handle requests for adding a new word.
//
// Parameters:
//   - ctx: The context for the request, used for cancellation, deadlines, and passing request - scoped values.
//   - req: A pointer to the AddWordReq struct containing the request information.
//
// Returns:
//   - *words.WordResp: A pointer to the WordResp struct containing the response information.
//   - error: An error object if an unexpected error occurs during the process.
func (s *WordServiceImpl) AddNewWord(ctx context.Context, req *words.AddWordReq) (resp *words.WordResp, err error) {
	return biz_words.AddNewWord(ctx, req)
}

// GetReviewWordList implements the WordServiceImpl interface.
// It retrieves the user ID from the request, logs it, and then calls the GetReviewWordList function
// in the biz_words package to handle requests for getting the review word list.
//
// Parameters:
//   - ctx: The context for the request, used for cancellation, deadlines, and passing request - scoped values.
//   - req: A pointer to the ReviewListReq struct containing the request information.
//
// Returns:
//   - *words.ReviewListResp: A pointer to the ReviewListResp struct containing the response information.
//   - error: An error object if an unexpected error occurs during the process.
func (s *WordServiceImpl) GetReviewWordList(ctx context.Context, req *words.ReviewListReq) (resp *words.ReviewListResp, err error) {
	user_id := req.UserId
	klog.Infof("first user_id: %v", user_id)
	return biz_words.GetReviewWordList(ctx, user_id)
}

// SubmitAnswer implements the WordServiceImpl interface.
// It calls the SubmitAnswer function in the biz_words package to handle requests for submitting answers.
//
// Parameters:
//   - ctx: The context for the request, used for cancellation, deadlines, and passing request - scoped values.
//   - req: A pointer to the SubmitAnswerReq struct containing the request information.
//
// Returns:
//   - *words.SubmitAnswerResp: A pointer to the SubmitAnswerResp struct containing the response information.
//   - error: An error object if an unexpected error occurs during the process.
func (s *WordServiceImpl) SubmitAnswer(ctx context.Context, req *words.SubmitAnswerReq) (resp *words.SubmitAnswerResp, err error) {
	return biz_words.SubmitAnswer(ctx, req.UserId, req.WordId, req.AnswerId, req.QuestionType, req.FilledName)
}

// GetTodayReviewProgress implements the WordServiceImpl interface.
// It calls the GetTodayReviewProgress function in the biz_words package to handle requests
// for getting today's review progress.
//
// Parameters:
//   - ctx: The context for the request, used for cancellation, deadlines, and passing request - scoped values.
//   - req: A pointer to the ReviewProgressReq struct containing the request information.
//
// Returns:
//   - *words.ReviewProgressResp: A pointer to the ReviewProgressResp struct containing the response information.
//   - error: An error object if an unexpected error occurs during the process.
func (s *WordServiceImpl) GetTodayReviewProgress(ctx context.Context, req *words.ReviewProgressReq) (resp *words.ReviewProgressResp, err error) {
	return biz_words.GetTodayReviewProgress(ctx, req)
}

// GetSupportedTags implements the WordServiceImpl interface.
// It calls the GetSupportedTags function in the biz_words package to handle requests
// for getting supported tags.
//
// Parameters:
//   - ctx: The context for the request, used for cancellation, deadlines, and passing request - scoped values.
//   - req: A pointer to the GetTagsReq struct containing the request information.
//
// Returns:
//   - *words.GetTagsResp: A pointer to the GetTagsResp struct containing the response information.
//   - error: An error object if an unexpected error occurs during the process.
func (s *WordServiceImpl) GetSupportedTags(ctx context.Context, req *words.GetTagsReq) (resp *words.GetTagsResp, err error) {
	return biz_words.GetSupportedTags(ctx, req)
}

// GetWordDetail implements the WordServiceImpl interface.
// It calls the GetWordDetail function in the biz_words package to handle requests for getting word details.
//
// Parameters:
//   - ctx: The context for the request, used for cancellation, deadlines, and passing request - scoped values.
//   - wordName: The name of the word for which details are requested.
//
// Returns:
//   - *words.WordDetail: A pointer to the WordDetail struct containing the word details.
//   - error: An error object if an unexpected error occurs during the process.
func (s *WordServiceImpl) GetWordDetail(ctx context.Context, wordName string) (resp *words.WordDetail, err error) {
	return biz_words.GetWordDetail(ctx, wordName)
}

// GetWordByID implements the WordServiceImpl interface.
// It calls the GetWordByID function in the biz_words package to handle requests for getting a word by its ID.
//
// Parameters:
//   - ctx: The context for the request, used for cancellation, deadlines, and passing request - scoped values.
//   - wordID: The ID of the word to be retrieved.
//
// Returns:
//   - *words.WordResp: A pointer to the WordResp struct containing the response information.
//   - error: An error object if an unexpected error occurs during the process.
func (s *WordServiceImpl) GetWordByID(ctx context.Context, wordID int64) (resp *words.WordResp, err error) {
	return biz_words.GetWordByID(ctx, wordID)
}

// UpdateWordTagID implements the WordServiceImpl interface.
// It calls the UpdateWordTagID function in the biz_words package to handle requests for updating a word's tag ID.
//
// Parameters:
//   - ctx: The context for the request, used for cancellation, deadlines, and passing request - scoped values.
//   - req: A pointer to the UpdateWordTag struct containing the request information.
//
// Returns:
//   - *words.WordResp: A pointer to the WordResp struct containing the response information.
//   - error: An error object if an unexpected error occurs during the process.
func (s *WordServiceImpl) UpdateWordTagID(ctx context.Context, req *words.UpdateWordTag) (resp *words.WordResp, err error) {
	return biz_words.UpdateWordTagID(ctx, req)
}
