// Copyright (c) 2025 Bytedance Ltd. and/or its affiliates
// SPDX-License-Identifier: MIT

package config

import (
	"testing"

	"github.com/Trae-AI/stream-to-river/internal/config"
	"github.com/Trae-AI/stream-to-river/internal/test"
)

const (
	filePath = "conf"
	fileName = "stream2river.yml"
)

func TestLoadConfig(t *testing.T) {
	_, _, err := LoadConfig("../../")
	test.Assert(t, err != nil)
	test.Assert(t, len(config.GetStringMap("Lingo")) != 0)
	test.Assert(t, len(config.GetStringMap("LLM")) != 0)
	test.Assert(t, len(config.GetStringMap("LLM.ChatModel")) != 0)
	test.Assert(t, len(config.GetStringMap("Coze")) != 0)

	err = config.CreateTmpConfFile(config.RPCServiceMockContent)
	test.Assert(t, err == nil)
	// load tmp conf file
	config.ResetConfig()
	_, _, err = LoadConfig("")
	test.Assert(t, err == nil, err)
	test.Assert(t, len(config.GetStringMap("Lingo")) != 0)
	test.Assert(t, len(config.GetStringMap("LLM")) != 0)
	test.Assert(t, len(config.GetStringMap("LLM.ChatModel")) != 0)
	test.Assert(t, config.GetString(CM_APIKEY) != "")
	test.Assert(t, config.GetString(CM_MODEL) != "")
	test.Assert(t, len(config.GetStringMap("Coze")) != 0)

	err = config.RemoveTmpFile()
	test.Assert(t, err == nil, err)
}
