// Copyright (c) 2025 Bytedance Ltd. and/or its affiliates
// SPDX-License-Identifier: MIT

package utils

import "github.com/cloudwego/hertz/pkg/common/utils"

// HertzErrorResp constructs a Hertz response map for error scenarios.
// It takes two string parameters: `msg` and `errMsg`.
// If `errMsg` is empty, it sets both the "message" and "error" fields in the response to `msg`.
// Otherwise, it sets the "message" field to `msg` and the "error" field to `errMsg`.
//
// Parameters:
//   - msg: The main message to be included in the response, typically describing the operation.
//   - errMsg: The error message. If empty, `msg` will be used as the error message.
//
// Returns:
//   - utils.H: A map of type `utils.H` containing the "message" and "error" fields.
func HertzErrorResp(msg string, errMsg string) utils.H {
	if errMsg == "" {
		return utils.H{"message": msg, "error": msg}
	}
	return utils.H{"message": msg, "error": errMsg}
}
