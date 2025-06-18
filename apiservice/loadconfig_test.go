package main

import (
	"testing"

	"github.com/Trae-AI/stream-to-river/internal/config"
	"github.com/Trae-AI/stream-to-river/internal/test"
)

func TestInitWithConfig(t *testing.T) {
	initWithConfig("")
	test.Assert(t, len(config.GetStringMap("LLM")) != 0)
	test.Assert(t, len(config.GetStringMap("LLM.AsrModel")) != 0)
	test.Assert(t, len(config.GetStringMap("LLM.VisionModel")) != 0)
}
