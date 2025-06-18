package config

import (
	"log"

	"github.com/Trae-AI/stream-to-river/internal/config"
	"github.com/Trae-AI/stream-to-river/rpcservice/biz/chat/coze"
	"github.com/Trae-AI/stream-to-river/rpcservice/biz/llm"
	"github.com/Trae-AI/stream-to-river/rpcservice/biz/words/vocapi"
	"github.com/Trae-AI/stream-to-river/rpcservice/dal/mysql"
)
// LocalConfigKey is the config key name
type LocalConfigKey = string

// The key of LocalConfigKey corresponds to the key in the stream2river.yml configuration.
const (
	DB_TYPE   LocalConfigKey = "DatabaseType"
	MYSQL_DSN LocalConfigKey = "MySQL.DSN"
	// SQLLITE_DBPATH is only used for test
	SQLLITE_DBPATH LocalConfigKey = "Sqlite.DBPath"

	CM_APIKEY LocalConfigKey = "LLM.ChatModel.APIKey"
	CM_MODEL  LocalConfigKey = "LLM.ChatModel.Model"

	COZE LocalConfigKey = "Coze"

	LINGO_URL LocalConfigKey = "Lingo.URL"
)

// LoadConfig loads the application configuration from a YAML file.
// It sets the configuration file name and type, adds the configuration path,
// and reads the configuration file. Then it extracts the database and Lingo configurations.
//
// Returns:
//   - *mysql.DataBaseConfig: A pointer to the database configuration.
//   - *vocapi.LingoConfig: A pointer to the Lingo configuration.
//   - error: An error object if an unexpected error occurs during the configuration loading process.
func LoadConfig() (*mysql.DataBaseConfig, *vocapi.LingoConfig, error) {
	config.LoadConfig("")
	if err := coze.InitCozeConfig(config.GetStringMapString(COZE)); err != nil {
		log.Fatalf("init coze error: err=%v", err)
	}

	if err := llm.InitModelCfg(config.GetString(CM_APIKEY), config.GetString(CM_MODEL)); err != nil {
		log.Fatalf("load LLM.ChatModel config failed, err=%v", err)
	}

	// Load dbConfig
	return &mysql.DataBaseConfig{
			DBType:       config.GetString(DB_TYPE),
			MySQLDSN:     config.GetString(MYSQL_DSN),
			SqliteDBPath: config.GetString(SQLLITE_DBPATH),
		}, &vocapi.LingoConfig{
			URL: config.GetString(LINGO_URL),
		}, nil
}
