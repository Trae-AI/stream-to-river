// Copyright (c) 2025 Bytedance Ltd. and/or its affiliates
// SPDX-License-Identifier: MIT

package main

import (
	"github.com/Trae-AI/stream-to-river/apiservice/biz/asr"
	"github.com/Trae-AI/stream-to-river/apiservice/biz/handler/chat"
	"github.com/Trae-AI/stream-to-river/apiservice/biz/handler/words"
	"github.com/Trae-AI/stream-to-river/apiservice/biz/image2text"
	"github.com/Trae-AI/stream-to-river/apiservice/biz/repome"
	"github.com/Trae-AI/stream-to-river/apiservice/biz/user"
	"github.com/cloudwego/hertz/pkg/app/server"
)

// RegisterAPI registers API paths and their corresponding handlers to the Hertz server.
// It organizes endpoints into different categories such as user management, chat, word operations, and multimodal services.
//
// Parameters:
//   - hz: A pointer to the Hertz server instance where API routes will be registered.
func RegisterAPI(hz *server.Hertz) {
	// User management related endpoints
	hz.POST("/api/register", user.RegisterHandler)
	hz.POST("/api/login", user.LoginHandler)
	hz.GET("/api/user", user.GetUserHandler)

	// Chat related endpoint
	hz.GET("/api/chat", chat.ChatHandler)

	// Word related endpoints
	hz.GET("/api/word-query", words.QueryWordHandler)
	hz.POST("/api/word-add", words.AddNewWordHandler)
	hz.GET("/api/word-detail", words.WordDetailHandler)
	hz.GET("/api/word-list", words.WordsHandler)
	hz.GET("/api/review-list", words.GetReviewWordListHandler)
	hz.GET("/api/review-progress", words.GetTodayReviewProgressHandler)
	hz.GET("/api/tags", words.GetSupportedTagsHandler)
	hz.POST("/api/answer", words.SubmitAnswerHandler)
	hz.POST("/api/word-tag", words.UpdateWordTagID)

	// Multimodal related endpoints
	hz.POST("/api/asrrecognize", asr.RecognizeAudioHandler)
	hz.POST("/api/image2text", image2text.Image2Text)

	// API and project document services (currently commented out)
	hz.GET("/api/apidoc", repome.RepoMeHandler)
	hz.GET("/api/static/:filename", repome.StaticHandler) // Static file service
}
