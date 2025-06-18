// Copyright (c) 2025 Bytedance Ltd. and/or its affiliates
// SPDX-License-Identifier: MIT

package coze

import (
	"errors"
)

var (
	CozeConf *CozeConfig

	InvalidConfErr = errors.New("coze configuration is invalid")
)

// The key corresponds to the subkey in the stream2river.yml Coze configuration.
const (
	Token      = "token"
	PublishKey = "publishkey"
	PrivateKey = "privatekey"
	BaseURL    = "baseurl"
	WorkflowID = "workflowid"
	Auth       = "auth"
	ClientID   = "clientid"
)

type CozeConfig struct {
	WorkflowID string
	BaseURL    string
	Token      string
	ClientID   string
	PublishKey string
	PrivateKey string
	Auth       string
}

func InitCozeConfig(cozeCfg map[string]string) error {
	if len(cozeCfg) == 0 {
		return InvalidConfErr
	}
	if cozeCfg[Token] == "" || cozeCfg[ClientID] == "" || cozeCfg[PublishKey] == "" ||
		cozeCfg[PrivateKey] == "" || cozeCfg[WorkflowID] == "" {
		return InvalidConfErr
	}
	CozeConf = &CozeConfig{}
	CozeConf.Token = cozeCfg[Token]
	CozeConf.BaseURL = cozeCfg[BaseURL]
	CozeConf.ClientID = cozeCfg[ClientID]
	CozeConf.PublishKey = cozeCfg[PublishKey]
	CozeConf.PrivateKey = cozeCfg[PrivateKey]
	CozeConf.WorkflowID = cozeCfg[WorkflowID]
	CozeConf.Auth = cozeCfg[Auth]

	return nil
}
