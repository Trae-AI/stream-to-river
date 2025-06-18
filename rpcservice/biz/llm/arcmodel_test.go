// Copyright (c) 2025 Bytedance Ltd. and/or its affiliates
// SPDX-License-Identifier: MIT

package llm

import (
	"context"
	"io"
	"testing"

	"github.com/Trae-AI/stream-to-river/internal/config"
	"github.com/Trae-AI/stream-to-river/internal/test"
	"github.com/cloudwego/eino/schema"
)

func TestArkModel(t *testing.T) {
	err := config.LoadConfig("../../")
	test.Assert(t, err == nil)

	err = InitModelCfg(config.GetString("LLM.ChatModel.APIKey"), config.GetString("LLM.ChatModel.Model"))
	if err != nil {
		// no ChatModel configuration
		t.Skip()
	}

	arkModel, err := GetArkModel()
	test.Assert(t, err == nil)
	msg, err := arkModel.Stream(context.Background(), []*schema.Message{schema.UserMessage("介绍下字节跳动")})
	test.Assert(t, err == nil, err)
	for {
		chunk, e := msg.Recv()
		if e != nil {
			if e == io.EOF {
				break
			}
			test.Assert(t, err == nil)
		}
		print(chunk.Content)
	}
}
