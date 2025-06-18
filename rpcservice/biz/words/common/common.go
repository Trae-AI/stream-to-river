// Copyright (c) 2025 Bytedance Ltd. and/or its affiliates
// SPDX-License-Identifier: MIT

package common

import "github.com/Trae-AI/stream-to-river/rpcservice/kitex_gen/base"

// BuildBaseResp constructs a new instance of the base.BaseResp struct.
// This function is used to generate a base response with a specified status code and message.
//
// Parameters:
//   - code: An int32 representing the status code of the response.
//   - msg: A string containing the status message of the response.
//
// Returns:
//   - *base.BaseResp: A pointer to the newly created base.BaseResp struct.
func BuildBaseResp(code int32, msg string) *base.BaseResp {
	return &base.BaseResp{
		StatusCode:    code,
		StatusMessage: msg,
	}
}

// BuildSuccBaseResp constructs a new instance of the base.BaseResp struct for successful operations.
// This function sets the status code to 0 and the status message to "success" by default.
//
// Returns:
//   - *base.BaseResp: A pointer to the newly created base.BaseResp struct with success status.
func BuildSuccBaseResp() *base.BaseResp {
	return &base.BaseResp{
		StatusCode:    0,
		StatusMessage: "success",
	}
}

// BuildBaseRespWithExtra constructs a new instance of the base.BaseResp struct with additional extra information.
// This function allows users to attach a map of extra data to the base response.
//
// Parameters:
//   - code: An int32 representing the status code of the response.
//   - msg: A string containing the status message of the response.
//   - extra: A map[string]string that holds additional information to be included in the response.
//
// Returns:
//   - *base.BaseResp: A pointer to the newly created base.BaseResp struct with extra data.
func BuildBaseRespWithExtra(code int32, msg string, extra map[string]string) *base.BaseResp {
	return &base.BaseResp{
		StatusCode:    code,
		StatusMessage: msg,
		Extra:         extra,
	}
}
