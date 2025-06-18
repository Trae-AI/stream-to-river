// Copyright (c) 2025 Bytedance Ltd. and/or its affiliates
// SPDX-License-Identifier: MIT

package dao

import (
	"github.com/cloudwego/kitex/pkg/klog"

	"github.com/Trae-AI/stream-to-river/rpcservice/dal/model"
	"github.com/Trae-AI/stream-to-river/rpcservice/dal/mysql"
)

// AddWordTag adds a new word tag record to the `word_tag` table.
//
// Parameters:
//   - tag: A pointer to the `model.WordTag` struct representing the word tag record to be added.
//
// Returns:
//   - error: An error object if an unexpected error occurs during the database operation.
func AddWordTag(tag *model.WordTag) error {
	ret := mysql.GetDB().Table(model.WordTagTableName).Create(tag)
	if ret.Error != nil {
		return ret.Error
	}
	klog.Infof("Insert word_tag=%v into table=%s", tag, model.WordTagTableName)
	return nil
}

// GetWordTagById queries the `word_tag` table for a word tag record based on the `tag_id`.
//
// Parameters:
//   - tagId: The unique identifier of the word tag.
//
// Returns:
//   - *model.WordTag: A pointer to the retrieved `WordTag` record. Returns `nil` if no record is found.
//   - error: An error object if an unexpected error occurs during the process. Returns `ErrNoRecord` if no record matches the query.
func GetWordTagById(tagId int32) (*model.WordTag, error) {
	var tag *model.WordTag
	ret := mysql.GetDB().Table(model.WordTagTableName).Where("id = ?", tagId).First(&tag)
	if ret.Error != nil {
		return nil, ret.Error
	}
	if ret.RowsAffected == 0 {
		klog.Errorf("Query failed[no record] with tagId=%v in table=%s", tagId, model.WordTagTableName)
		return nil, ErrNoRecord
	}
	klog.Infof("Query word_tag=%v with tagId=%v in table=%s", tag, tagId, model.WordTagTableName)
	return tag, nil
}

// GetWordTagByName queries the `word_tag` table for a word tag record based on the `tag_name`.
//
// Parameters:
//   - tagName: The name of the word tag.
//
// Returns:
//   - *model.WordTag: A pointer to the retrieved `WordTag` record. Returns `nil` if no record is found.
//   - error: An error object if an unexpected error occurs during the process. Returns `ErrNoRecord` if no record matches the query.
func GetWordTagByName(tagName string) (*model.WordTag, error) {
	var tag *model.WordTag
	ret := mysql.GetDB().Table(model.WordTagTableName).Where("tag_name = ?", tagName).First(&tag)
	if ret.Error != nil {
		return nil, ret.Error
	}
	if ret.RowsAffected == 0 {
		klog.Errorf("Query failed[no record] with tagName=%v in table=%s", tagName, model.WordTagTableName)
		return nil, ErrNoRecord
	}
	klog.Infof("Query word_tag=%v with tagName=%v in table=%s", tag, tagName, model.WordTagTableName)
	return tag, nil
}

// GetAllWordTags retrieves all word tag records from the `word_tag` table.
//
// Returns:
//   - []*model.WordTag: A slice of pointers to the retrieved `WordTag` records.
//   - error: An error object if an unexpected error occurs during the process.
func GetAllWordTags() ([]*model.WordTag, error) {
	var tags []*model.WordTag
	ret := mysql.GetDB().Table(model.WordTagTableName).Find(&tags)
	if ret.Error != nil {
		return nil, ret.Error
	}
	klog.Infof("Query all word_tags, count=%d", len(tags))
	return tags, nil
}

// UpdateWordTag updates a word tag record in the `word_tag` table.
//
// Parameters:
//   - tag: A pointer to the `model.WordTag` struct representing the word tag record to be updated.
//
// Returns:
//   - error: An error object if an unexpected error occurs during the database operation.
func UpdateWordTag(tag *model.WordTag) error {
	ret := mysql.GetDB().Table(model.WordTagTableName).Where("tag_id = ?", tag.TagId).Updates(tag)
	if ret.Error != nil {
		return ret.Error
	}
	klog.Infof("Update word_tag=%v in table=%s", tag, model.WordTagTableName)
	return nil
}

// DeleteWordTag deletes a word tag record from the `word_tag` table based on the `tag_id`.
//
// Parameters:
//   - tagId: The unique identifier of the word tag.
//
// Returns:
//   - error: An error object if an unexpected error occurs during the database operation.
func DeleteWordTag(tagId int64) error {
	ret := mysql.GetDB().Table(model.WordTagTableName).Where("tag_id = ?", tagId).Delete(&model.WordTag{})
	if ret.Error != nil {
		return ret.Error
	}
	klog.Infof("Delete word_tag with tagId=%v from table=%s", tagId, model.WordTagTableName)
	return nil
}
