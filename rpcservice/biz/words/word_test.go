// Copyright (c) 2025 Bytedance Ltd. and/or its affiliates
// SPDX-License-Identifier: MIT

package words

import (
	"context"
	"testing"

	"github.com/Trae-AI/stream-to-river/internal/test"
	"github.com/Trae-AI/stream-to-river/rpcservice/biz/words/vocapi"
	"github.com/Trae-AI/stream-to-river/rpcservice/dal/dao"
	"github.com/Trae-AI/stream-to-river/rpcservice/dal/model"
	"github.com/Trae-AI/stream-to-river/rpcservice/dal/mysql"
	"github.com/Trae-AI/stream-to-river/rpcservice/dal/redis"
	"github.com/Trae-AI/stream-to-river/rpcservice/kitex_gen/words"
)

func setupTest() {
	vocapi.InitLingoConfig(&vocapi.LingoConfig{
		URL: "https://sstr.trae.com.cn/api/word-detail?word=",
	})
	// 初始化cache
	redis.InitCache()
}

// TestGetWordDetail tests the GetWordDetail function.
func TestGetWordDetail(t *testing.T) {
	setupTest()

	wordName := "pragmatic"
	expectedDetail := &words.WordDetail{
		NewWordName_: "pragmatic",
		Explains:     "adj. 讲求实际的，务实的；实用主义的；（语言学）语用的",
	}

	ctx := context.Background()
	detail, err := GetWordDetail(ctx, wordName)
	test.Assert(t, err == nil, "GetWordDetail should not return error")
	test.Assert(t, detail.NewWordName_ == expectedDetail.NewWordName_, "Word name should match")
	test.Assert(t, detail.Explains == expectedDetail.Explains, "Word explains should match")
}

// TestAddNewWord tests the AddNewWord function.
func TestAddNewWord(t *testing.T) {
	setupTest()

	db, err := mysql.InitMockDB()
	test.Assert(t, err == nil, "Failed to init mock DB:", err)

	err = db.Migrator().DropTable(&model.Word{}, &model.AnswerList{}, &model.WordsRisiteRecord{})
	test.Assert(t, err == nil, "Failed to drop tables:", err)

	err = db.Migrator().CreateTable(&model.Word{}, &model.AnswerList{}, &model.WordsRisiteRecord{})
	test.Assert(t, err == nil, "Failed to create tables:", err)

	wordName := "pragmatic"
	userId := int64(1)
	tagId := int32(1)

	ctx := context.Background()
	wordDetail, err := GetWordDetail(ctx, wordName)
	test.Assert(t, err == nil, "GetWordDetail should not return error")

	req := &words.AddWordReq{
		UserId:   userId,
		WordName: wordName,
		TagId:    tagId,
	}

	resp, err := AddNewWord(ctx, req)
	test.Assert(t, err == nil, "AddNewWord should not return error")
	test.Assert(t, resp.BaseResp.StatusCode == 0, "Response code should be 0")
	test.Assert(t, resp.Word.WordName == wordDetail.NewWordName_, "Word name should match")
}

// TestGetWordByID tests the GetWordByID function.
func TestGetWordByID(t *testing.T) {
	db, err := mysql.InitMockDB()
	test.Assert(t, err == nil, "Failed to init mock DB:", err)

	err = db.Migrator().DropTable(&model.Word{})
	test.Assert(t, err == nil, "Failed to drop table:", err)

	err = db.Migrator().CreateTable(&model.Word{})
	test.Assert(t, err == nil, "Failed to create table:", err)

	word := &model.Word{
		WordName:    "test",
		Description: "test description",
		Explains:    "test explains",
		PronounceUs: "us_url",
		PronounceUk: "uk_url",
		UserId:      1,
		TagId:       1,
	}
	err = dao.AddWord(word)
	test.Assert(t, err == nil, "Failed to add word:", err)

	ctx := context.Background()
	resp, err := GetWordByID(ctx, word.WordId)
	test.Assert(t, err == nil, "GetWordByID should not return error")
	test.Assert(t, resp.BaseResp.StatusCode == 0, "Response code should be 0")
	test.Assert(t, resp.Word.WordName == word.WordName, "Word name should match")
}

// TestModel2Word tests the model2Word function.
func TestModel2Word(t *testing.T) {
	modelWord := &model.Word{
		WordName:    "test",
		Description: "test description",
		Explains:    "test explains",
		PronounceUs: "us_url",
		PronounceUk: "uk_url",
		PhoneticUs:  "us_phonetic",
		PhoneticUk:  "uk_phonetic",
		TagId:       1,
	}

	word := model2Word(modelWord)
	test.Assert(t, word.WordName == modelWord.WordName, "Word name should match")
	test.Assert(t, word.PronounceUs.Url == modelWord.PronounceUs, "US pronunciation URL should match")
	test.Assert(t, word.PronounceUs.Phonetic == modelWord.PhoneticUs, "US phonetic should match")
}

// TestCleanupWord tests the cleanupWord function.
func TestCleanupWord(t *testing.T) {
	dirtyWord := " \t\n test \r\n "
	cleanWord := cleanupWord(dirtyWord)
	test.Assert(t, cleanWord == "test", "Cleaned word should match")
}

func TestAddOrQueryWord(t *testing.T) {
	db, err := mysql.InitMockDB()
	test.Assert(t, err == nil, err)

	err = db.Migrator().DropTable(model.Word{})
	test.Assert(t, err == nil)

	err = db.Migrator().CreateTable(model.Word{})
	test.Assert(t, err == nil)

	var wordID int64 = 1
	_, err = dao.QueryWord(wordID)
	test.Assert(t, err == dao.ErrNoRecord)

	newWord := &words.Word{
		WordId:      wordID,
		WordName:    "world",
		Description: "test",
		Explains:    "test",
	}
	err = dao.AddWord(word2Model(newWord))
	test.Assert(t, err == nil)

	mw, err := dao.QueryWord(wordID)
	test.Assert(t, err == nil)

	w := model2Word(mw)
	test.Assert(t, newWord.WordName == w.WordName)
}

func word2Model(w *words.Word) *model.Word {
	// 处理PronounceUs
	var pronounceUsUrl string
	if w.PronounceUs != nil {
		pronounceUsUrl = w.PronounceUs.Url
	}

	// 处理PronounceUk
	var pronounceUkUrl string
	if w.PronounceUk != nil {
		pronounceUkUrl = w.PronounceUk.Url
	}

	return &model.Word{
		WordId:      w.WordId,
		WordName:    w.WordName,
		Description: w.Description,
		Explains:    w.Explains,
		// 修改：从PronounceInfo结构提取URL存储到数据库
		PronounceUs: pronounceUsUrl,
		PronounceUk: pronounceUkUrl,
	}
}

// TestReviewProgressCountCorrectness 测试复习进度计数的正确性
// 验证新用户首次添加单词时pending_review_count不会重复计数
func TestReviewProgressCountCorrectness(t *testing.T) {
	// 1. 初始化mock数据库
	db, err := mysql.InitMockDB()
	test.Assert(t, err == nil, "Failed to init mock DB:", err)

	// 2. 清理并创建测试所需的表
	tables := []interface{}{
		&model.Word{},
		&model.AnswerList{},
		&model.WordsRisiteRecord{},
		&model.ReviewProgress{},
	}

	for _, table := range tables {
		err = db.Migrator().DropTable(table)
		test.Assert(t, err == nil, "Failed to drop table:", err)

		err = db.Migrator().CreateTable(table)
		test.Assert(t, err == nil, "Failed to create table:", err)
	}

	// 3. 测试数据
	const (
		userId int64 = 12345
		tagId  int32 = 1
	)

	// 4. 模拟新用户首次添加单词
	t.Run("FirstWordAddition", func(t *testing.T) {
		// 由于vocapi.ProcessWord需要真实的API调用，我们这里只测试核心逻辑
		// 直接创建word、answer_list、words_risite_record记录
		word1 := &model.Word{
			WordName:    "hello",
			Description: "Hello world",
			Explains:    "打招呼用语",
			UserId:      userId,
			TagId:       tagId,
		}
		err = dao.AddWord(word1)
		test.Assert(t, err == nil, "Failed to add word:", err)

		// 查询添加的单词获取word_id
		queriedWord, err := dao.QueryWordByUserIdAndName(userId, "hello")
		test.Assert(t, err == nil, "Failed to query word:", err)

		// 添加answer_list记录
		answerList := &model.AnswerList{
			WordId:      queriedWord.WordId,
			UserId:      userId,
			WordName:    queriedWord.WordName,
			Description: queriedWord.Explains,
		}
		err = dao.AddAnswerList(answerList)
		test.Assert(t, err == nil, "Failed to add answer list:", err)

		// 添加复习记录
		record := &model.WordsRisiteRecord{
			WordId:         int(queriedWord.WordId),
			Level:          1,
			NextReviewTime: 1234567890, // 设置为需要复习的时间
			DowngradeStep:  1,
			TotalCorrect:   0,
			TotalWrong:     0,
			Score:          0,
			UserId:         userId,
		}
		err = dao.AddWordsRisiteRecord(record)
		test.Assert(t, err == nil, "Failed to add review record:", err)

		// 5. 调用getTodayReviewProgressWithInitFlag模拟AddNewWord中的逻辑
		progressReq := &words.ReviewProgressReq{UserId: userId}
		resp, wasInitialized, err := getTodayReviewProgressWithInitFlag(progressReq)
		test.Assert(t, err == nil, "Failed to get review progress:", err)
		test.Assert(t, wasInitialized == true, "Should have initialized for new user")

		// 6. 验证：新用户首次添加单词，pending_review_count应该等于total_words_count
		test.Assert(t, resp.PendingReviewCount == resp.TotalWordsCount,
			"First word: pending_review_count (%d) should equal total_words_count (%d)",
			resp.PendingReviewCount, resp.TotalWordsCount)
		test.Assert(t, resp.TotalWordsCount == 1, "Total words should be 1")
		test.Assert(t, resp.PendingReviewCount == 1, "Pending review count should be 1")

		t.Logf("First word - Total: %d, Pending: %d, Was Initialized: %v",
			resp.TotalWordsCount, resp.PendingReviewCount, wasInitialized)
	})

	// 7. 测试后续添加单词的逻辑
	t.Run("SubsequentWordAddition", func(t *testing.T) {
		// 添加第二个单词
		word2 := &model.Word{
			WordName:    "world",
			Description: "World peace",
			Explains:    "世界",
			UserId:      userId,
			TagId:       tagId,
		}
		err = dao.AddWord(word2)
		test.Assert(t, err == nil, "Failed to add second word:", err)

		// 查询添加的单词获取word_id
		queriedWord2, err := dao.QueryWordByUserIdAndName(userId, "world")
		test.Assert(t, err == nil, "Failed to query second word:", err)

		// 添加answer_list记录
		answerList2 := &model.AnswerList{
			WordId:      queriedWord2.WordId,
			UserId:      userId,
			WordName:    queriedWord2.WordName,
			Description: queriedWord2.Explains,
		}
		err = dao.AddAnswerList(answerList2)
		test.Assert(t, err == nil, "Failed to add answer list for second word:", err)

		// 添加复习记录
		record2 := &model.WordsRisiteRecord{
			WordId:         int(queriedWord2.WordId),
			Level:          1,
			NextReviewTime: 1234567890, // 设置为需要复习的时间
			DowngradeStep:  1,
			TotalCorrect:   0,
			TotalWrong:     0,
			Score:          0,
			UserId:         userId,
		}
		err = dao.AddWordsRisiteRecord(record2)
		test.Assert(t, err == nil, "Failed to add review record for second word:", err)

		// 获取当前复习进度（应该不会重新初始化）
		progressReq := &words.ReviewProgressReq{UserId: userId}
		resp, wasInitialized, err := getTodayReviewProgressWithInitFlag(progressReq)
		test.Assert(t, err == nil, "Failed to get review progress for second word:", err)
		test.Assert(t, wasInitialized == false, "Should not initialize for existing user")
		test.Assert(t, resp.BaseResp.StatusCode == 0, "BaseResp.StatusCode should be 0")

		// 模拟AddNewWord中的增量逻辑
		if !wasInitialized {
			err = dao.IncrementPendingReviewCount(userId)
			test.Assert(t, err == nil, "Failed to increment pending review count:", err)
		}

		// 重新获取复习进度
		resp, _, err = getTodayReviewProgressWithInitFlag(progressReq)
		test.Assert(t, err == nil, "Failed to get updated review progress:", err)

		// 8. 验证：后续添加单词，pending_review_count应该正确递增
		test.Assert(t, resp.TotalWordsCount == 2, "Total words should be 2")
		test.Assert(t, resp.PendingReviewCount == 2, "Pending review count should be 2")
		test.Assert(t, resp.PendingReviewCount == resp.TotalWordsCount,
			"Second word: pending_review_count (%d) should equal total_words_count (%d)",
			resp.PendingReviewCount, resp.TotalWordsCount)

		t.Logf("Second word - Total: %d, Pending: %d, Was Initialized: %v",
			resp.TotalWordsCount, resp.PendingReviewCount, wasInitialized)
	})

	t.Log("Review progress count correctness test passed!")
}
