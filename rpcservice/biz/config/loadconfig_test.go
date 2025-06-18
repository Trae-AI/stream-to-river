package config

import (
	"testing"

	"github.com/Trae-AI/stream-to-river/internal/config"
	"github.com/Trae-AI/stream-to-river/internal/test"
)

func TestLoadConfig(t *testing.T) {
	_, _, err := LoadConfig()
	test.Assert(t, err == nil)
	test.Assert(t, len(config.GetStringMap("Lingo")) != 0)
	test.Assert(t, len(config.GetStringMap("LLM")) != 0)
	test.Assert(t, len(config.GetStringMap("LLM.ChatModel")) != 0)
	test.Assert(t, len(config.GetStringMap("Coze")) != 0)
}
