// Copyright (c) 2025 Bytedance Ltd. and/or its affiliates
// SPDX-License-Identifier: MIT

package words

import (
	"context"
	"fmt"
	"log"
	"testing"

	"github.com/Trae-AI/stream-to-river/internal/test"
	"github.com/Trae-AI/stream-to-river/rpcservice/biz/words/vocapi"
	"github.com/Trae-AI/stream-to-river/rpcservice/dal/dao"
	"github.com/Trae-AI/stream-to-river/rpcservice/dal/model"
	"github.com/Trae-AI/stream-to-river/rpcservice/dal/mysql"
	"github.com/Trae-AI/stream-to-river/rpcservice/dal/redis"
	"github.com/Trae-AI/stream-to-river/rpcservice/kitex_gen/words"
	"gorm.io/gorm"
)

var mockDB *gorm.DB

func TestMain(m *testing.M) {
	db, err := mysql.InitMockDB()
	if err != nil {
		log.Fatal(err)
	}
	mockDB = db
	// create test tables
	err = createTestTables()
	if err != nil {
		log.Fatal(err)
	}

	// init lingo
	vocapi.InitLingoConfig(&vocapi.LingoConfig{
		URL: "https://sstr.trae.com.cn/api/word-detail?word=",
	})
	// init cache
	redis.InitCache()

	m.Run()
}

// TestGetWordDetail tests the GetWordDetail function.
func TestGetWordDetail(t *testing.T) {
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
	recreateMockWordTable()

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

	// not exist word
	resp, err := AddNewWord(ctx, req)
	test.Assert(t, err == nil, "AddNewWord should not return error")
	test.Assert(t, resp.BaseResp.StatusCode == 0, "Response code should be 0")
	test.Assert(t, resp.Word.WordName == wordDetail.NewWordName_, "Word name should match")

	// exist word
	resp, err = AddNewWord(ctx, req)
	test.Assert(t, err == nil, "AddNewWord should not return error")
	test.Assert(t, resp.BaseResp.StatusCode == 1, "Response code should be 0")
	test.Assert(t, resp.Word.WordName == wordDetail.NewWordName_, "Word name should match")
}

// TestGetWordByID tests the GetWordByID function.
func TestGetWordByID(t *testing.T) {
	recreateMockWordTable()

	word := &model.Word{
		WordName:    "test",
		Description: "test description",
		Explains:    "test explains",
		PronounceUs: "us_url",
		PronounceUk: "uk_url",
		UserId:      1,
		TagId:       1,
	}
	err := dao.AddWord(word)
	test.Assert(t, err == nil, "Failed to add word:", err)

	resp, err := GetWordByID(context.Background(), word.WordId)
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
	recreateMockWordTable()
	var wordID int64 = 1
	_, err := dao.QueryWord(wordID)
	test.Assert(t, err == dao.ErrNoRecord, err)

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

func createTestTables() error {
	if err := recreateMockWordTable(); err != nil {
		return err
	}
	if err := recreateMockTagTable(); err != nil {
		return err
	}
	if err := recreateMockAnswerListTable(); err != nil {
		return err
	}
	if err := recreateMockWordReviewRecordTable(); err != nil {
		return err
	}
	if err := recreateMockReviewProgressTable(); err != nil {
		return err
	}
	return nil
}

func recreateMockWordTable() error {
	if err := mockDB.Migrator().DropTable(&model.Word{}); err != nil {
		return fmt.Errorf("failed to drop Word table: %w", err)
	}
	if err := mockDB.Migrator().CreateTable(&model.Word{}); err != nil {
		return fmt.Errorf("failed to create Word table: %w", err)
	}
	return nil
}

func recreateMockTagTable() error {
	if err := mockDB.Migrator().DropTable(model.WordTag{}); err != nil {
		return fmt.Errorf("failed to drop WordTag table: %w", err)
	}
	if err := mockDB.Migrator().CreateTable(model.WordTag{}); err != nil {
		return fmt.Errorf("failed to create WordTag table: %w", err)
	}
	return nil
}

func recreateMockAnswerListTable() error {
	if err := mockDB.Migrator().DropTable(&model.AnswerList{}); err != nil {
		return fmt.Errorf("failed to drop AnswerList table: %w", err)
	}
	if err := mockDB.Migrator().CreateTable(&model.AnswerList{}); err != nil {
		return fmt.Errorf("failed to create AnswerList table: %w", err)
	}
	return nil
}

func recreateMockWordReviewRecordTable() error {
	if err := mockDB.Migrator().DropTable(&model.WordsRisiteRecord{}); err != nil {
		return fmt.Errorf("failed to drop WordsRisiteRecord table: %w", err)
	}
	if err := mockDB.Migrator().CreateTable(&model.WordsRisiteRecord{}); err != nil {
		return fmt.Errorf("failed to create WordsRisiteRecord table: %w", err)
	}
	return nil
}

func recreateMockReviewProgressTable() error {
	if err := mockDB.Migrator().DropTable(&model.ReviewProgress{}); err != nil {
		return fmt.Errorf("failed to drop ReviewProgress table: %w", err)
	}
	if err := mockDB.Migrator().CreateTable(&model.ReviewProgress{}); err != nil {
		return fmt.Errorf("failed to create ReviewProgress table: %w", err)
	}
	return nil
}
