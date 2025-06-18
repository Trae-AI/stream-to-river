// Copyright (c) 2025 Bytedance Ltd. and/or its affiliates
// SPDX-License-Identifier: MIT

package hertz

import (
	"os"
	"path"
)

// Predefined env variables and default configurations.
const (
	EnvConfDir  = "HERTZ_CONF_DIR"
	EnvConfFile = "HERTZ_CONF_FILE"
	EnvLogDir   = "HERTZ_LOG_DIR"

	DefaultConfDir  = "conf"
	DefaultConfFile = "hertz.yml"
	DefaultLogDir   = "log"
)

// GetConfDir gets dir of config file.
func GetConfDir() string {
	if confDir := os.Getenv(EnvConfDir); confDir != "" {
		return confDir
	}
	return DefaultConfDir
}

// GetConfFile gets config file path.
func GetConfFile() string {
	file := DefaultConfFile
	if confFile := os.Getenv(EnvConfFile); confFile != "" {
		file = confFile
	}
	return path.Join(GetConfDir(), file)
}

// GetEnvLogDir is to get log dir from env.
func GetEnvLogDir() string {
	return os.Getenv(EnvLogDir)
}

// GetLogDir gets dir of log file.
// Deprecated: it is suggested to use GetEnvLogDir instead of GetLogDir, and GetEnvLogDir won't return default log dir.
func GetLogDir() string {
	if logDir := os.Getenv(EnvLogDir); logDir != "" {
		return logDir
	}
	return DefaultLogDir
}
