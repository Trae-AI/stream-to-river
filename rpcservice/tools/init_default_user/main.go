// Copyright (c) 2025 Bytedance Ltd. and/or its affiliates
// SPDX-License-Identifier: MIT

package main

import (
	"fmt"
	"log"

	"gorm.io/gorm"

	"github.com/Trae-AI/stream-to-river/rpcservice/biz/config"
	"github.com/Trae-AI/stream-to-river/rpcservice/biz/words/vocapi"
	"github.com/Trae-AI/stream-to-river/rpcservice/dal/dao"
	"github.com/Trae-AI/stream-to-river/rpcservice/dal/model"
	"github.com/Trae-AI/stream-to-river/rpcservice/dal/mysql"
	"github.com/Trae-AI/stream-to-river/rpcservice/dal/redis"
	"github.com/Trae-AI/stream-to-river/rpcservice/global"
)

// AddDefaultUserWithWordList adds default words to the database for a specified user.
// If the user already has words in the database, the function returns without making any changes.
// For each word in the list, it fetches the word details, checks if the word already exists,
// and inserts the word and its corresponding answer list into the database if it doesn't exist.
// Parameters:
//   - defaultUserID: The ID of the user for whom the words will be added.
//   - defaultWordList: A list of words to be added to the database.
//
// Returns:
//   - error: An error object if an error occurs during the process; otherwise, nil.
func AddDefaultUserWithWordList(defaultUserID int64, defaultWordList []string) error {
	// 1. Check if the specified user_id exists in the words table
	totalWords, err := dao.GetTotalWordsCount(defaultUserID)
	if err != nil && err != gorm.ErrRecordNotFound {
		return err
	}
	// 2. If the user already exists, return directly
	if totalWords > 0 {
		log.Printf("user %d already exists", defaultUserID)
		return nil
	}

	for _, word := range defaultWordList {
		// 3.1 Call the vocapi to get word details
		WordExplains, err := vocapi.ProcessWord(word)
		if err != nil {
			return err
		}
		// 3.2 If the word does not exist, return an error directly
		if WordExplains.ErrorNo != 0 || WordExplains.NewWordName == "" {
			return fmt.Errorf("word not found: %s", word)
		}
		// 3.3 Check if the word already exists using the NewWordName returned by ProcessWord
		existingWord, err := dao.QueryWordByUserIdAndName(defaultUserID, WordExplains.NewWordName)
		if err != nil && err != gorm.ErrRecordNotFound {
			// Return the database query error
			return fmt.Errorf("failed to query word: %w", err)
		}
		if existingWord != nil {
			// Move on to the next word
			continue
		}
		// 3.4 Build data for insertion
		modelWord := model.Word{
			WordName:    WordExplains.NewWordName,
			Description: WordExplains.ExplainsOxford,
			Explains:    WordExplains.ExplainsYoudao,
			PronounceUs: WordExplains.PronounceUS.Url,
			PronounceUk: WordExplains.PronounceUK.Url,
			// New phonetic field
			PhoneticUs: WordExplains.PronounceUS.Phonetic,
			PhoneticUk: WordExplains.PronounceUK.Phonetic,
			UserId:     defaultUserID,
			TagId:      1,
			YoudaoUrl:  fmt.Sprintf(global.YoudaoUrl, WordExplains.NewWordName),
		}

		// 3.5 Add to the word table
		err = dao.AddWord(&modelWord)
		if err != nil {
			return err
		}
		// 3.6 Query the words table by WordName to get the word_id
		queriedWord, err := dao.QueryWordByUserIdAndName(defaultUserID, WordExplains.NewWordName)
		if err != nil {
			return fmt.Errorf("failed to query word: %w", err)
		}

		// 3.7 Create an answer_list record using the word_id queried from the database
		answerList := &model.AnswerList{
			WordId:      queriedWord.WordId,
			UserId:      defaultUserID,
			WordName:    queriedWord.WordName,
			Description: queriedWord.Explains, // Use the explains field as the description
		}

		// 3.8 Add to the answer_list table
		err = dao.AddAnswerList(answerList)
		if err != nil {
			return fmt.Errorf("failed to add answer list: %w", err)
		}
		log.Printf("AddAnswerList success: word_id=%v, word_name=%v", queriedWord.WordId, modelWord.WordName)
	}

	return nil
}

func main() {
	// load config
	dbConfig, lingoConfig, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("loadConfig error: err=%v", err)
	}

	// Initialize the database connection
	if err := mysql.InitDBWithConfig(dbConfig); err != nil {
		panic(err)
	}
	// Initialize the Redis cache
	redis.InitCache()
	// Initialize the Lingo configuration
	vocapi.InitLingoConfig(lingoConfig)

	// Define the list of default words
	word_list := []string{"encourage", "debate", "positioned",
		"ambiguous", "bureaucracy", "cynical",
		"dilemma", "eloquent", "grim", "Hypocrisy", "meticulous",
		"nostalgia", "pragmatic"}

	// Add default words to the database for user ID 1
	if err := AddDefaultUserWithWordList(1, word_list); err != nil {
		panic(err)
	}

	// Log the successful initialization of the default user
	log.Printf("init default user success. total words: %d", len(word_list))
}
