// Copyright (c) 2025 Bytedance Ltd. and/or its affiliates
// SPDX-License-Identifier: MIT

package utils

import "testing"

func TestIsInRefusedString(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"as an ai", true},
		{"hello world", false},
		{"不能回答", true},
	}
	for _, tt := range tests {
		result := IsInRefusedString(tt.input)
		if result != tt.expected {
			t.Errorf("IsInRefusedString(%q) = %v; want %v", tt.input, result, tt.expected)
		}
	}
}
