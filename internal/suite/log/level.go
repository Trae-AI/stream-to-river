// Copyright (c) 2025 Bytedance Ltd. and/or its affiliates
// SPDX-License-Identifier: MIT

package log

import (
	"strings"

	"github.com/cloudwego/kitex/pkg/klog"
)

type Level string

// Log levels.
const (
	LevelTrace  Level = "trace"
	LevelDebug  Level = "debug"
	LevelInfo   Level = "info"
	LevelNotice Level = "notice"
	LevelWarn   Level = "warn"
	LevelError  Level = "error"
	LevelFatal  Level = "fatal"
)

// Value return log level constant.
func (level Level) Value() klog.Level {
	l := Level(strings.ToLower(string(level)))
	switch l {
	case LevelTrace:
		return klog.LevelTrace
	case LevelDebug:
		return klog.LevelDebug
	case LevelInfo:
		return klog.LevelInfo
	case LevelNotice:
		return klog.LevelNotice
	case LevelWarn:
		return klog.LevelWarn
	case LevelError:
		return klog.LevelError
	case LevelFatal:
		return klog.LevelFatal
	default:
		return klog.LevelTrace
	}
}
