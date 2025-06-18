// Copyright (c) 2025 Bytedance Ltd. and/or its affiliates
// SPDX-License-Identifier: MIT

package asr

import (
	"context"
	"fmt"
	"sync"

	"github.com/cloudwego/hertz/pkg/app"
)

var (
	appID   string
	token   string
	cluster string

	once sync.Once
)

// RecognizeAudioHandler handles HTTP requests for audio recognition.
// It extracts the audio format and data from the request, validates the audio data,
// then calls the audio recognition service. Finally, it returns the recognition result or an error.
//
// Parameters:
//   - ctx: The context for the request, used for cancellation, deadlines, and passing request-scoped values.
//   - c: A pointer to the app.RequestContext provided by Hertz, containing request and response information.
func RecognizeAudioHandler(ctx context.Context, c *app.RequestContext) {
	// Get the audio format parameter from the query string.
	// If not provided, default to "wav".
	format := c.Query("format")
	if format == "" {
		format = "wav" // Default format
	}

	// Get the audio data from the request body.
	audioData := c.Request.Body()
	// Validate the audio data. If empty, return a 400 Bad Request response.
	if len(audioData) == 0 {
		c.JSON(400, map[string]interface{}{
			"code":    -1,
			"message": "Empty audio data",
		})
		return
	}

	// Call the RecognizeAudio method of the AsrService to perform audio recognition.
	response, err := NewAsrService().RecognizeAudio(ctx, audioData, format)
	// If an error occurs during audio recognition, return a 500 Internal Server Error response.
	if err != nil {
		c.JSON(500, map[string]interface{}{
			"code":    -1,
			"message": err.Error(),
		})
		return
	}

	// Return the audio recognition result with a 200 OK response.
	c.JSON(200, map[string]interface{}{
		"code":    0,
		"message": "success",
		"data":    response,
	})
}

func InitModelCfg(appID, token, cluster string) error {
	if appID == "" || token == "" || cluster == "" {
		return fmt.Errorf("AsrModel.AppID=%s, AsrModel.Token=%s, AsrModel.Cluster=%s,"+
			"please check your config file", appID, token, cluster)
	}
	once.Do(func() {
		appID = appID
		token = token
		cluster = cluster
	})
	return nil
}
