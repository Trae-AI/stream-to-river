// Copyright (c) 2025 Bytedance Ltd. and/or its affiliates
// SPDX-License-Identifier: MIT

package dao

import (
	"testing"

	"github.com/Trae-AI/stream-to-river/internal/test"
	"github.com/Trae-AI/stream-to-river/rpcservice/dal/model"
	"github.com/Trae-AI/stream-to-river/rpcservice/dal/mysql"
)

func TestWordDBOperation(t *testing.T) {
	db, _ := mysql.InitMockDB()
	// create the required table
	_ = db.AutoMigrate(&model.Word{})
	var wordID int64 = 1

	err := DelWord(wordID)
	test.Assert(t, err == nil)

	w, err := QueryWord(wordID)
	test.Assert(t, err == ErrNoRecord, err)

	newWord := &model.Word{
		WordId:      wordID,
		WordName:    "world",
		Description: "test",
		Explains:    "test",
	}
	err = AddWord(newWord)
	test.Assert(t, err == nil)

	w, err = QueryWord(wordID)
	test.Assert(t, err == nil)
	test.Assert(t, w.WordId == wordID)
}

func TestWordOperationWithMockDB(t *testing.T) {
	db, err := mysql.InitMockDB()
	test.Assert(t, err == nil)

	err = db.Migrator().DropTable(model.Word{})
	test.Assert(t, err == nil)

	exist := db.Migrator().HasTable(model.Word{})
	test.Assert(t, !exist)

	err = db.Migrator().CreateTable(model.Word{})
	test.Assert(t, err == nil)

	exist = db.Migrator().HasTable(model.Word{})
	test.Assert(t, exist)

	var wordID int64 = 1
	w, err := QueryWord(wordID)
	test.Assert(t, err == ErrNoRecord, err)

	newWord := &model.Word{
		WordId:      wordID,
		WordName:    "world",
		Description: "test",
		Explains:    "test",
	}
	err = AddWord(newWord)
	test.Assert(t, err == nil, err)

	w, err = QueryWord(wordID)
	test.Assert(t, err == nil)
	test.Assert(t, w.WordId == wordID)

	err = DelWord(wordID)
	test.Assert(t, err == nil)

	w, err = QueryWord(wordID)
	test.Assert(t, err == ErrNoRecord, err)
}

func TestUpdateWordTagID(t *testing.T) {
	db, err := mysql.InitMockDB()
	test.Assert(t, err == nil)
	err = db.Migrator().DropTable(model.Word{})
	test.Assert(t, err == nil)
	exist := db.Migrator().HasTable(model.Word{})
	test.Assert(t, !exist)
	err = db.Migrator().CreateTable(model.Word{})
	test.Assert(t, err == nil)

	var wordID, userID int64 = 1, 2
	var oldTagID int32 = 3
	var newTagID int32 = 4

	err = UpdateWordTagID(wordID, userID, newTagID)
	test.Assert(t, err == ErrNoUpdate)

	_, err = QueryWord(wordID)
	test.Assert(t, err == ErrNoRecord, err)
	newWord := &model.Word{
		WordId: wordID,
		UserId: userID,
		TagId:  oldTagID,
	}
	err = AddWord(newWord)
	test.Assert(t, err == nil, err)

	err = UpdateWordTagID(wordID, userID, newTagID)
	test.Assert(t, err == nil, err)

	w, err := QueryWord(wordID)
	test.Assert(t, err == nil, err)
	test.Assert(t, w.TagId == newTagID)
}
