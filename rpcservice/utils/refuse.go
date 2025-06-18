// Copyright (c) 2025 Bytedance Ltd. and/or its affiliates
// SPDX-License-Identifier: MIT

package utils

import (
	"strings"
)

var (
	refuseStringArr = []string{
		"as an ai",
		"can't answer",
		"can't help",
		"can't reply",
		"cannot anwser",
		"cannot help",
		"cannot reply",
		"do not have the ability to",
		"no related",
		"no relevant",
		"sorry",
		"unable to reply",
		"一个 ai",
		"不会对",
		"不知道",
		"不能为你",
		"不能为您",
		"不能回答",
		"不能对",
		"不能提供",
		"对不起",
		"我不会",
		"我无法",
		"我没有",
		"抱歉",
		"无法为你",
		"无法为您",
		"无法回答",
		"无法建议",
		"无法提供",
		"无法给出",
		"没有找到",
		"没有能力"}
)

// IsInRefusedString 是否包含rs拒绝回复的字符串
func IsInRefusedString(resp string) bool {
	for i := 0; i < len(refuseStringArr); i++ {
		if strings.Contains(resp, refuseStringArr[i]) {
			return true
		}
	}
	return false
}
