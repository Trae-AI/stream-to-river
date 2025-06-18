// Copyright (c) 2025 Bytedance Ltd. and/or its affiliates
// SPDX-License-Identifier: MIT

package common

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestBuildBaseResp tests the BuildBaseResp function.
func TestBuildBaseResp(t *testing.T) {
	code := int32(1)
	msg := "test message"
	resp := BuildBaseResp(code, msg)

	assert.Equal(t, code, resp.StatusCode)
	assert.Equal(t, msg, resp.StatusMessage)
}

// TestBuildSuccBaseResp tests the BuildSuccBaseResp function.
func TestBuildSuccBaseResp(t *testing.T) {
	resp := BuildSuccBaseResp()

	assert.Equal(t, int32(0), resp.StatusCode)
	assert.Equal(t, "success", resp.StatusMessage)
}

// TestBuildBaseRespWithExtra tests the BuildBaseRespWithExtra function.
func TestBuildBaseRespWithExtra(t *testing.T) {
	code := int32(2)
	msg := "test message with extra"
	extra := map[string]string{"key": "value"}
	resp := BuildBaseRespWithExtra(code, msg, extra)

	assert.Equal(t, code, resp.StatusCode)
	assert.Equal(t, msg, resp.StatusMessage)
	assert.Equal(t, extra, resp.Extra)
}
