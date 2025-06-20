// Copyright (c) 2025 Bytedance Ltd. and/or its affiliates
// SPDX-License-Identifier: MIT

package image2text

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"
	"sync"

	"github.com/cloudwego/eino-ext/components/model/ark"
	"github.com/cloudwego/eino/schema"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/kitex/pkg/klog"
)

var (
	// visionModelAPIKey is the API key used to authenticate requests to the vision model service.
	// It's a unique identifier that allows the application to access the vision model.
	visionModelAPIKey string

	// visionModelEP specifies the name of the Vision model to be used.
	// This constant determines which version or type of the Ark model the application will interact with.
	visionModelEP string

	once sync.Once
)

// isValidBase64 checks if the input string is a valid Base64 - encoded string.
// If the string starts with "data:image", it extracts the Base64 part and validates it.
//
// Parameters:
//   - s: The input string to be validated.
//
// Returns:
//   - bool: `true` if the string is a valid Base64 string, `false` otherwise.
func isValidBase64(s string) bool {
	// Check if the string contains the "data:image" prefix
	if strings.HasPrefix(s, "data:image") {
		// Extract the Base64 part
		parts := strings.Split(s, ",")
		if len(parts) != 2 {
			return false
		}
		s = parts[1]
	}

	// Validate the Base64 format
	_, err := base64.StdEncoding.DecodeString(s)
	return err == nil
}

// Image2Text is an HTTP handler function that handles image - to - text conversion requests.
// It expects a JSON request body containing a Base64 - encoded image.
// After validating the input, it calls the LLM to get a text description of the image.
//
// Parameters:
//   - ctx: The context for the request.
//   - c: A pointer to the Hertz `app.RequestContext`, used to handle HTTP requests and responses.
func Image2Text(ctx context.Context, c *app.RequestContext) {
	var requestBody struct {
		Base64 string `json:"base64"`
	}

	// Bind the JSON request body
	if err := c.BindJSON(&requestBody); err != nil {
		klog.CtxErrorf(ctx, "Bind request body failed: %v", err)
		return
	}

	// Validate the Base64 data
	if requestBody.Base64 == "" {
		klog.CtxErrorf(ctx, "Empty Base64")
		return
	}

	// Validate if it's a standard Base64 string
	if !isValidBase64(requestBody.Base64) {
		c.JSON(http.StatusBadRequest, map[string]interface{}{
			"code":    400,
			"message": "Invalid base64 string",
		})
		return
	}

	// Process the business logic
	resp, err := callLLM(ctx, requestBody.Base64)
	if err != nil {
		klog.CtxErrorf(ctx, "Image2text processing failed: %v", err)
		c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"code":    500,
			"message": "Failed to process image",
			"error":   err.Error(),
		})
		return
	}

	// Return the success response
	c.JSON(http.StatusOK, map[string]interface{}{
		"code":    200,
		"message": "success",
		"data":    resp,
	})
}

// callLLM calls the large language model (LLM) to generate a text description of the given image.
// It initializes an ARK chat model, constructs a message with the image URL,
// and sends the message to the LLM to get a response.
//
// Parameters:
//   - ctx: The context for the request.
//   - base64: The Base64 - encoded image string.
//
// Returns:
//   - string: The text description of the image generated by the LLM.
//   - error: An error object if an unexpected error occurs during the process.
func callLLM(ctx context.Context, base64 string) (text string, err error) {
	arkModel, err := ark.NewChatModel(ctx, &ark.ChatModelConfig{
		APIKey: visionModelAPIKey,
		Model:  visionModelEP,
	})
	if err != nil {
		err = fmt.Errorf("arkModel init failed, err=%v", err)
		fmt.Println(err)
		return "", err
	}

	message := &schema.Message{
		Role: schema.User,
		MultiContent: []schema.ChatMessagePart{
			{
				Type: "text",
				Text: "Please use a short English sentence to describe the current image, within 30 words.",
			},
			{
				Type: "image_url",
				ImageURL: &schema.ChatMessageImageURL{
					URL: base64, // This should be replaced with the actual image URL
				},
			},
		},
	}

	chatMsg, err := arkModel.Generate(ctx, []*schema.Message{message})
	if err != nil {
		err = fmt.Errorf("arkModel create Stream failed, err=%v", err)
		klog.CtxErrorf(ctx, err.Error())
		return "", err
	}

	return chatMsg.Content, nil
}

func InitModelCfg(apiKey, model string) error {
	if apiKey == "" || model == "" {
		return fmt.Errorf("VisionModel.APIKey=%s, VisionModel.Model=%s, "+
			"please check your config file", apiKey, model)
	}
	once.Do(func() {
		visionModelAPIKey = apiKey
		visionModelEP = model
	})
	return nil
}
