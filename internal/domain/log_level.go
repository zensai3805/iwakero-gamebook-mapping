package domain

import (
	"fmt"
	"strings"
)

// LogLevel はログレベルを表す値オブジェクト
type LogLevel int

// ログレベル定数
const (
	LogLevelDebug LogLevel = iota
	LogLevelInfo
	LogLevelWarn
	LogLevelError
	LogLevelFatal
)

// ログレベル文字列定数
const (
	LogLevelDebugString = "DEBUG"
	LogLevelInfoString  = "INFO"
	LogLevelWarnString  = "WARN"
	LogLevelErrorString = "ERROR"
	LogLevelFatalString = "FATAL"
)

// String はログレベルを文字列に変換する
func (l LogLevel) String() string {
	switch l {
	case LogLevelDebug:
		return LogLevelDebugString
	case LogLevelInfo:
		return LogLevelInfoString
	case LogLevelWarn:
		return LogLevelWarnString
	case LogLevelError:
		return LogLevelErrorString
	case LogLevelFatal:
		return LogLevelFatalString
	default:
		return "UNKNOWN"
	}
}

// IsHigherThan は現在のログレベルが引数のログレベルより高いかを判定する
func (l LogLevel) IsHigherThan(other LogLevel) bool {
	return l > other
}

// ParseLogLevel は文字列からログレベルを生成する
func ParseLogLevel(s string) (LogLevel, error) {
	switch strings.ToUpper(s) {
	case LogLevelDebugString:
		return LogLevelDebug, nil
	case LogLevelInfoString:
		return LogLevelInfo, nil
	case LogLevelWarnString:
		return LogLevelWarn, nil
	case LogLevelErrorString:
		return LogLevelError, nil
	case LogLevelFatalString:
		return LogLevelFatal, nil
	default:
		return LogLevelDebug, fmt.Errorf("不正なログレベル: %s", s)
	}
}
