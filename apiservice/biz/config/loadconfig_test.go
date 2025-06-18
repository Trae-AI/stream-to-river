package config

import (
	"testing"

	"github.com/Trae-AI/stream-to-river/internal/config"
	"github.com/Trae-AI/stream-to-river/internal/test"
)

func TestInitWithConfig(t *testing.T) {
	err := config.LoadConfig("../../")
	test.Assert(t, err == nil)
	test.Assert(t, len(config.GetStringMap("LLM")) != 0)
	test.Assert(t, len(config.GetStringMap("LLM.AsrModel")) != 0)
	test.Assert(t, len(config.GetStringMap("LLM.VisionModel")) != 0)

	err = config.CreateTmpConfFile(config.APIServiceMockContent)
	test.Assert(t, err == nil)
	// load tmp conf file
	config.ResetConfig()
	err = config.LoadConfig("")
	test.Assert(t, err == nil)
	test.Assert(t, len(config.GetStringMap("LLM")) != 0)
	test.Assert(t, len(config.GetStringMap("LLM.AsrModel")) != 0)
	test.Assert(t, len(config.GetStringMap("LLM.VisionModel")) != 0)

	InitWithConfig("")
	err = config.RemoveTmpFile()
	test.Assert(t, err == nil)
}
