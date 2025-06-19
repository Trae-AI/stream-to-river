// Copyright (c) 2025 Bytedance Ltd. and/or its affiliates
// SPDX-License-Identifier: MIT

package vocapi

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/Trae-AI/stream-to-river/rpcservice/dal/redis"
)

func TestQueryWord(t *testing.T) {
	wordList := []string{"ambiguous", "bureaucracy", "cynical", "dilemma",
		"eloquent", "grim", "Hypocrisy", "meticulous", "nostalgia", "pragmatic"}

	InitLingoConfig(&LingoConfig{
		URL: "https://sstr.trae.com.cn/api/word-detail?word=",
	})
	// 初始化cache
	redis.InitCache()

	for _, word := range wordList {
		wordExplains, err := ProcessWord(word)
		if err != nil {
			if wordExplains != nil && wordExplains.ErrorNo == http.StatusTooManyRequests {
				continue
			}
			t.Errorf("ProcessWord失败: %v", err.Error())
		}
		if wordExplains.ErrorNo != 0 || wordExplains.NewWordName == "" {
			t.Errorf("Word not found for word_name=%v wordExplains.ErrorNo=%v wordExplains.NewWordName=%v",
				word, wordExplains.ErrorNo, wordExplains.NewWordName)
		} else {
			println(fmt.Sprintf("word=%v, explains=%v", word, wordExplains))
		}
	}
}
