// Copyright (c) 2025 Bytedance Ltd. and/or its affiliates
// SPDX-License-Identifier: MIT

package llm

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/cloudwego/eino-ext/components/model/ark"
)

var (
	// arkModelAPIKey is the API key used to authenticate requests to the Ark model service.
	// It's a unique identifier that allows the application to access the Ark model.
	arkModelAPIKey string

	// arkModelEP specifies the name of the Ark model to be used.
	// This constant determines which version or type of the Ark model the application will interact with.
	arkModelEP string

	once sync.Once

	EmptyConfErr = errors.New("ark model configuration is empty")
)

// GetArkModel initializes and returns a new instance of the Ark chat model.
// It uses the predefined API key and model name to configure the model.
//
// Parameters:
//   - ctx: The context for the operation, used for cancellation, deadlines, and passing request-scoped values.
//
// Returns:
//   - *ark.ChatModel: A pointer to the initialized Ark chat model instance. Returns nil if an error occurs.
//   - error: An error object if an unexpected error occurs during the model initialization process.
func GetArkModel() (*ark.ChatModel, error) {
	if arkModelAPIKey == "" || arkModelEP == "" {
		return nil, EmptyConfErr
	}
	arkModel, err := ark.NewChatModel(context.Background(), &ark.ChatModelConfig{
		APIKey: arkModelAPIKey,
		Model:  arkModelEP,
	})
	if err != nil {
		return nil, err
	}
	return arkModel, nil
}

func InitModelCfg(apiKey, model string) error {
	if apiKey == "" || model == "" {
		return fmt.Errorf("ChatModel.APIKey=%s, ChatModel.Model=%s, "+
			"please check your config file", apiKey, model)
	}
	once.Do(func() {
		arkModelAPIKey = apiKey
		arkModelEP = model
	})
	return nil
}
