package config

import (
	"log"
	"path"
	"sync"

	"github.com/spf13/viper"
)

const (
	filePath = "conf/"
	fileName = "stream2river"
)

var once sync.Once

func LoadConfig(relativePath string) {
	fp := filePath
	if relativePath != "" {
		fp = path.Join(relativePath, fp)
	}
	once.Do(func() {
		// Load the configuration file
		viper.SetConfigName(fileName)
		viper.SetConfigType("yaml")
		viper.AddConfigPath(fp)

		if err := viper.ReadInConfig(); err != nil {
			log.Fatalf("Error reading config file: %s", err)
		}
	})
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
