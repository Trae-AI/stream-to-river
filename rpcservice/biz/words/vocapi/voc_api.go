// Copyright (c) 2025 Bytedance Ltd. and/or its affiliates
// SPDX-License-Identifier: MIT

package vocapi

import (
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/bytedance/sonic"
	"github.com/cloudwego/kitex/pkg/klog"

	"github.com/Trae-AI/stream-to-river/rpcservice/dal/redis"
	"github.com/Trae-AI/stream-to-river/rpcservice/kitex_gen/words"
)

// WordExplainCacheExpireTime is the expiration time for word explanation cache entries.
// By default, it's set to one week.
var WordExplainCacheExpireTime = 24 * 7 * time.Hour

// LingoServiceRequestTimeout is the timeout duration for HTTP requests to the Lingo dictionary service.
// The default timeout is 3 seconds.
var LingoServiceRequestTimeout = 3 * time.Second

// WordNotExistErrorMsg is the error message returned when a word is not found.
var WordNotExistErrorMsg string = "这个词暂时不见了，请联系客服或稍后再试~"

// LingoConfig holds the configuration for the Lingo dictionary service.
// It contains the base URL of the service.
type LingoConfig struct {
	URL string
}

// lingoConfigInstance is a global variable that stores the Lingo service configuration.
// It's initialized using the InitLingoConfig function.
var lingoConfigInstance *LingoConfig

// VocSimpleResponse represents the structure of a simple API response from the Lingo service.
// It includes an error number and vocabulary - related data.
type VocSimpleResponse struct {
	ErrNo int `json:"err_no"`
	Data  struct {
		Voc struct {
			Word string `json:"word"`
		} `json:"voc"`
	} `json:"data"`
}

// Pronounce defines the structure for pronunciation information from the API.
// It includes the pronunciation type, phonetic symbol, and audio details.
type Pronounce struct {
	Type     string `json:"type"`
	Phonetic string `json:"phonetic"`
	Audio    struct {
		AudioVid   string `json:"audio_vid"`
		AudioModel string `json:"audio_model"`
		StartTime  int    `json:"start_time"`
		EndTime    int    `json:"end_time"`
		Duration   int    `json:"duration"`
		DurationMs int    `json:"duration_ms"`
		AudioUrl   string `json:"audio_url"`
		Token      string `json:"token"`
	} `json:"audio"`
}

// SentenceAudio defines the structure for sentence audio information.
// It contains details about the audio such as video ID, model, time range, and URL.
type SentenceAudio struct {
	AudioVid   string `json:"audio_vid"`
	AudioModel string `json:"audio_model"`
	StartTime  int    `json:"start_time"`
	EndTime    int    `json:"end_time"`
	Duration   int    `json:"duration"`
	DurationMs int    `json:"duration_ms"`
	AudioUrl   string `json:"audio_url"`
	Token      string `json:"token"`
}

// PronounceInfo defines the pronunciation information used in the final response.
// It includes the phonetic symbol and the audio URL.
type PronounceInfo struct {
	Phonetic string `json:"phonetic"`
	Url      string `json:"url"`
}

// SentenceInfo defines the structure for sentence information in the final response.
// It contains the sentence text and its audio URL.
type SentenceInfo struct {
	Text     string `json:"text"`
	AudioUrl string `json:"audio_url"`
}

// Sentence represents a sentence with additional metadata from the API.
// It includes highlighted text, favorite status, text, translation, audio, and text indexes.
type Sentence struct {
	HighLightTexts []interface{} `json:"high_light_texts"`
	FavorStatus    int           `json:"favor_status"`
	Text           string        `json:"text"`
	Tran           string        `json:"tran"`
	Audio          SentenceAudio `json:"audio"`
	TextIndexes    []struct {
		Start int `json:"start"`
		End   int `json:"end"`
	} `json:"text_indexes"`
	TranIndexes []interface{} `json:"tran_indexes"`
}

// WordExplains represents the final result structure for word explanations.
// It includes error number, new word name, explanations from different sources,
// pronunciation information, and example sentences.
type WordExplains struct {
	ErrorNo        int                   `json:"error_no"`
	NewWordName    string                `json:"new_word_name"`
	ExplainsYoudao string                `json:"explains"`
	ExplainsOxford string                `json:"description"`
	PronounceUS    *words.PronounceInfo  `json:"pronounce_us"`
	PronounceUK    *words.PronounceInfo  `json:"pronounce_uk"`
	Sentences      []*words.SentenceInfo `json:"sentences"`
}

// InitLingoConfig initializes the global Lingo service configuration.
// It sets the provided LingoConfig instance to the global variable lingoConfigInstance.
//
// Parameters:
//   - lingoConfig: A pointer to the LingoConfig struct containing the service URL.
func InitLingoConfig(lingoConfig *LingoConfig) {
	lingoConfigInstance = lingoConfig
}

// ProcessWord processes a given word by first checking the cache.
// If the word is not cached, it fetches the word details from the Lingo dictionary service,
// caches the result, and then returns it.
//
// Parameters:
//   - originalWordName: The word to be processed.
//
// Returns:
//   - *WordExplains: A pointer to the WordExplains struct containing word details.
//   - error: An error object if an unexpected error occurs during the process.
func ProcessWord(originalWordName string) (*WordExplains, error) {
	klog.Infof("ProcessWord called with originalWordName: %s", originalWordName)

	// Initialize the default response in case of errors
	result := &WordExplains{
		NewWordName:    "",
		ErrorNo:        400,
		ExplainsYoudao: WordNotExistErrorMsg,
		ExplainsOxford: WordNotExistErrorMsg,
		Sentences:      make([]*words.SentenceInfo, 0),
	}

	// Check if the input word is empty
	if originalWordName == "" {
		klog.Warn("ProcessWord: originalWordName is empty")
		return result, nil
	}

	// Step 1: Check the cache
	cacheKey := fmt.Sprintf("voc_%s", originalWordName)
	cacheValue, exist := redis.Cache.Get(cacheKey)
	if exist {
		klog.Infof("ProcessWord hit cache, key: %s, value: %s", cacheKey, cacheValue)
		return cacheValue.(*WordExplains), nil
	}

	// Step 2: Make an HTTP request to the Lingo service
	req, err := http.NewRequest("GET", lingoConfigInstance.URL+originalWordName, nil)
	if err != nil {
		klog.Errorf("Failed to create HTTP request for word detail of '%s': %v", originalWordName, err)
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{
		Timeout: LingoServiceRequestTimeout,
	}
	resp, err := client.Do(req)
	if err != nil {
		klog.Errorf("http.Get failed for word '%s': %v", originalWordName, err)
		return result, nil
	}
	defer resp.Body.Close()

	// Check if the HTTP response is successful
	if resp.StatusCode != http.StatusOK {
		klog.Errorf("http.Get failed for word '%s', status code: %d", originalWordName, resp.StatusCode)
		result.ErrorNo = resp.StatusCode
		return result, nil
	}

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		klog.Errorf("Failed to read response body for word detail of '%s': %v", originalWordName, err)
		return nil, err
	}

	// Unmarshal the response body into the WordExplains struct
	var wordDetail WordExplains
	err = sonic.Unmarshal(body, &wordDetail)
	if err != nil {
		klog.Errorf("Failed to unmarshal response body for word detail of '%s': %v", originalWordName, err)
		return nil, err
	}

	// Step 3: Cache the result
	redis.Cache.Set(cacheKey, &wordDetail, 0)

	return &wordDetail, nil
}
