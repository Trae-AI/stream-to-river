package config

import (
	"fmt"
	"os"
	"path"
	"sync"

	"github.com/spf13/viper"
)

const (
	filePath = "conf/"
	fileName = "stream2river"
)

var lock sync.Mutex

func LoadConfig(relativePath string) error {
	fp := filePath
	if relativePath != "" {
		fp = path.Join(relativePath, fp)
	}
	lock.Lock()
	defer lock.Unlock()
	// Load the configuration file
	viper.SetConfigName(fileName)
	viper.SetConfigType("yaml")
	viper.AddConfigPath(fp)

	if err := viper.ReadInConfig(); err != nil {
		return fmt.Errorf("reading config file failed: %w", err)
	}
	return nil
}

// GetString is to Get global k:v config, the val is string
func GetString(key string) string {
	return viper.GetString(key)
}

// GetStringMap is to Get global k:v config, the val is map[string]interface{}
func GetStringMap(key string) map[string]interface{} {
	return viper.GetStringMap(key)
}

// GetStringMapString is to Get global k:v config, the val is map[string]string
func GetStringMapString(key string) map[string]string {
	return viper.GetStringMapString(key)
}

// ResetConfig is only used for testing
func ResetConfig() {
	viper.Reset()
}

// CreateTmpConfFile is only used for testing
func CreateTmpConfFile(content string) error {
	confDir := "conf"
	if err := os.MkdirAll(confDir, 0755); err != nil {
		return err
	}
	configContent := []byte(content)
	configPath := path.Join(filePath, fileName+".yml")
	if err := os.WriteFile(configPath, configContent, 0644); err != nil {
		return err
	}
	return nil
}

// RemoveTmpFile is only used for testing
func RemoveTmpFile() error {
	if err := os.Remove(path.Join(filePath, fileName+".yml")); err != nil {
		return err
	}
	if err := os.Remove(filePath); err != nil {
		return err
	}
	return nil
}

const RPCServiceMockContent = `# 数据库相关配置
DatabaseType: "sqlite"
MySQL:
  DSN: "root:123456@tcp(127.0.0.1:13306)/stream?charset=utf8mb4&parseTime=True&loc=Local"
SQLite:
  DBPath: "./data/stream.db"

# VOC Lingo 火山引擎词典相关配置
Lingo:
  URL: "https://sstr.trae.com.cn/api/word-detail?word="

LLM:
  ChatModel:
    APIKey: "mock"
    Model: "ep-mock"

Coze:
  BaseURL: "https://api.coze.cn"
  # The following fields are configured with reference to rpcservice/biz/chat/coze/README.md
  WorkflowID: "mock"
  Auth: "mock"
  Token: "mock"
  ClientID: "mock"
  PublishKey: "mock"
  PrivateKey: "mock"
`

const APIServiceMockContent = `LLM:
  AsrModel:
    # You can read the "Sentence Recognition" access document in advance: https://www.volcengine.com/docs/6561/80816, and go to the Volcano Ark platform to access the sentence recognition capability https://console.volcengine.com/speech/service/15, and fill in the following AppID / Token / Cluster provided by the platform
    AppID: "mock"
    Token: "mock"
    Cluster: "mock"
  VisionModel:
    # You need to go to the Volcano Ark platform https://console.volcengine.com/ark/region:ark+cn-beijing/model/detail?Id=doubao-1-5-vision-lite to apply for Doubao's latest Vision lite model and get their latest api_key and model_id
    APIKey: "mock"
    Model: "mock"

# JWT_SECRET is used to sign and verify JWT tokens. It must be a long, random string.
# Recommended to use at least 32 bytes (256 bits) of random data.
# You can generate a secure random string using the following commands:
#   openssl rand -base64 32
#   or in Python: import secrets; print(secrets.token_urlsafe(32))
JWT_SECRET: your_secret_key`
