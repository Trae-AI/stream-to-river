package config

import (
	"log"
	"sync"

	"github.com/Trae-AI/stream-to-river/apiservice/biz/asr"
	"github.com/Trae-AI/stream-to-river/apiservice/biz/image2text"
	"github.com/Trae-AI/stream-to-river/apiservice/biz/user"
	"github.com/Trae-AI/stream-to-river/internal/config"
)

// LocalConfigKey is the config key name
type LocalConfigKey = string

// The key of LocalConfigKey corresponds to the key in the stream2river.yml configuration.
const (
	JWT LocalConfigKey = "JWT_SECRET"

	VM_APIKEY LocalConfigKey = "LLM.VisionModel.APIKey"
	VM_MODEL  LocalConfigKey = "LLM.VisionModel.Model"

	AM_APPID   LocalConfigKey = "LLM.AsrModel.AppID"
	AM_TOKEN   LocalConfigKey = "LLM.AsrModel.Token"
	AM_CLUSTER LocalConfigKey = "LLM.AsrModel.Cluster"
)

var once sync.Once

// InitWithConfig initializes the application configuration.
// It uses sync.Once to ensure that the configuration is loaded only once.
// If an error occurs during configuration loading, the application will terminate.
func InitWithConfig(relativePath string) {
	once.Do(func() {
		// Load local config file
		if err := config.LoadConfig(relativePath); err != nil {
			log.Fatalf("loadConfig failed: err=%v", err)
		}

		// Init JWT config
		user.InitJWTSecret(config.GetString(JWT))

		// Init vision model config
		if err := image2text.InitModelCfg(config.GetString(VM_APIKEY), config.GetString(VM_MODEL)); err != nil {
			log.Fatalf("load LLM.VisionModel config failed, err=%v", err)
		}

		// Init ASR model config
		if err := asr.InitModelCfg(config.GetString(AM_APPID), config.GetString(AM_TOKEN), config.GetString(AM_CLUSTER)); err != nil {
			log.Fatalf("load LLM.AsrModel config failed, err=%v", err)
		}
	})
}
