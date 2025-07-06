package main

import (
	"fmt"
	"runtime"
	"strings"
	"time"

	"github.com/zensai3805/iwakero-gamebook-mapping/internal/domain"
)

// StructuredLogEntry は構造化ログエントリを表す
type StructuredLogEntry struct {
	Timestamp   string                 `json:"timestamp"`
	Level       string                 `json:"level"`
	Category    string                 `json:"category"`    // operation, validation, error, ui, system
	Action      string                 `json:"action"`      // new_gamebook, add_paragraph, etc.
	Component   string                 `json:"component"`   // frameworks, commands, interactive, etc.
	Message     string                 `json:"message"`
	Context     map[string]interface{} `json:"context,omitempty"`
	Performance *PerformanceInfo       `json:"performance,omitempty"`
	Location    *SourceLocation        `json:"location,omitempty"`
	Session     *SessionInfo           `json:"session,omitempty"`
}

// PerformanceInfo はパフォーマンス情報（軽量版）
type PerformanceInfo struct {
	DurationMs float64 `json:"duration_ms,omitempty"`
	Success    bool    `json:"success"`
}

// SourceLocation はソースコード位置情報
type SourceLocation struct {
	File     string `json:"file"`
	Line     int    `json:"line"`
	Function string `json:"function"`
}

// SessionInfo はセッション情報
type SessionInfo struct {
	HasCurrentGame   bool   `json:"has_current_game"`
	TotalParagraphs  int    `json:"total_paragraphs,omitempty"`
	CurrentOperation string `json:"current_operation,omitempty"`
}

// getSourceLocation は呼び出し元の位置情報を取得する
func getSourceLocation(skip int) *SourceLocation {
	pc, file, line, ok := runtime.Caller(skip + 1)
	if !ok {
		return nil
	}

	// ファイル名を短縮
	if idx := strings.LastIndex(file, "/"); idx != -1 {
		file = file[idx+1:]
	}

	// 関数名を取得
	function := "unknown"
	if fn := runtime.FuncForPC(pc); fn != nil {
		fullName := fn.Name()
		if idx := strings.LastIndex(fullName, "."); idx != -1 {
			function = fullName[idx+1:]
		} else {
			function = fullName
		}
	}

	return &SourceLocation{
		File:     file,
		Line:     line,
		Function: function,
	}
}

// getCurrentSessionInfo は現在のセッション情報を取得する
func getCurrentSessionInfo() *SessionInfo {
	info := &SessionInfo{
		HasCurrentGame: currentGame != nil,
	}

	if currentGame != nil {
		info.TotalParagraphs = len(currentGame.Paragraphs)
	}

	return info
}

// LogStructuredOperation は構造化されたオペレーションログを出力する
func LogStructuredOperation(category, action, message string, context map[string]interface{}, duration time.Duration) {
	if logger := GetGlobalLogger(); logger != nil {
		entry := StructuredLogEntry{
			Timestamp: time.Now().Format(time.RFC3339Nano),
			Level:     "INFO",
			Category:  category,
			Action:    action,
			Component: "frameworks",
			Message:   message,
			Context:   context,
			Location:  getSourceLocation(1),
			Session:   getCurrentSessionInfo(),
		}

		if duration > 0 {
			entry.Performance = &PerformanceInfo{
				DurationMs: float64(duration.Nanoseconds()) / 1000000,
				Success:    true,
			}
		}

		// AI開発モード時は全詳細を出力
		if IsAIDevelopmentMode() {
			logger.Info(formatStructuredLog(entry), convertToFields(entry)...)
		} else {
			// 通常モードは簡潔に
			logger.Info(message,
				domain.Field{Key: "category", Value: category},
				domain.Field{Key: "action", Value: action},
				domain.Field{Key: "component", Value: "frameworks"})
		}
	}
}

// LogStructuredError は構造化されたエラーログを出力する
func LogStructuredError(category, action, message string, err error, context map[string]interface{}) {
	if logger := GetGlobalLogger(); logger != nil {
		entry := StructuredLogEntry{
			Timestamp: time.Now().Format(time.RFC3339Nano),
			Level:     "ERROR",
			Category:  category,
			Action:    action,
			Component: "frameworks",
			Message:   message,
			Context:   context,
			Location:  getSourceLocation(1),
			Session:   getCurrentSessionInfo(),
		}

		if err != nil {
			entry.Context["error"] = err.Error()
			entry.Context["error_type"] = fmt.Sprintf("%T", err)
		}

		entry.Performance = &PerformanceInfo{
			Success: false,
		}

		logger.Error(formatStructuredLog(entry), convertToFields(entry)...)
	}
}

// LogStructuredValidation は構造化されたバリデーションログを出力する
func LogStructuredValidation(field string, value interface{}, reason string, context map[string]interface{}) {
	if logger := GetGlobalLogger(); logger != nil {
		entry := StructuredLogEntry{
			Timestamp: time.Now().Format(time.RFC3339Nano),
			Level:     "WARN",
			Category:  "validation",
			Action:    "field_validation",
			Component: "frameworks",
			Message:   fmt.Sprintf("入力値検証エラー: %s", reason),
			Context: map[string]interface{}{
				"field":  field,
				"value":  fmt.Sprintf("%v", value),
				"reason": reason,
			},
			Location: getSourceLocation(1),
			Session:  getCurrentSessionInfo(),
		}

		// コンテキストをマージ
		for k, v := range context {
			entry.Context[k] = v
		}

		logger.Warn(formatStructuredLog(entry), convertToFields(entry)...)
	}
}

// LogStructuredUI は構造化されたUI操作ログを出力する
func LogStructuredUI(action, message string, context map[string]interface{}, duration time.Duration) {
	if logger := GetGlobalLogger(); logger != nil {
		entry := StructuredLogEntry{
			Timestamp: time.Now().Format(time.RFC3339Nano),
			Level:     "DEBUG",
			Category:  "ui",
			Action:    action,
			Component: "frameworks",
			Message:   message,
			Context:   context,
			Session:   getCurrentSessionInfo(),
		}

		if duration > 0 {
			entry.Performance = &PerformanceInfo{
				DurationMs: float64(duration.Nanoseconds()) / 1000000,
				Success:    true,
			}
		}

		// UI操作は AI開発モード時のみ詳細ログ
		if IsAIDevelopmentMode() {
			logger.Debug(formatStructuredLog(entry), convertToFields(entry)...)
		}
	}
}

// formatStructuredLog は構造化ログエントリを人間が読みやすい形式にフォーマットする
func formatStructuredLog(entry StructuredLogEntry) string {
	var parts []string

	// カテゴリとアクション
	parts = append(parts, fmt.Sprintf("[%s:%s]", entry.Category, entry.Action))

	// メッセージ
	parts = append(parts, entry.Message)

	// セッション情報（重要な場合のみ）
	if entry.Session != nil && entry.Session.HasCurrentGame {
		parts = append(parts, fmt.Sprintf("(ゲーム読み込み済み:%d段落)", entry.Session.TotalParagraphs))
	}

	// パフォーマンス情報（存在する場合）
	if entry.Performance != nil {
		if entry.Performance.DurationMs > 0 {
			parts = append(parts, fmt.Sprintf("%.2fms", entry.Performance.DurationMs))
		}
		if !entry.Performance.Success {
			parts = append(parts, "(失敗)")
		}
	}

	// AI開発モード時は位置情報を追加
	if IsAIDevelopmentMode() && entry.Location != nil {
		parts = append(parts, fmt.Sprintf("@%s:%d", entry.Location.File, entry.Location.Line))
	}

	return strings.Join(parts, " ")
}

// convertToFields は構造化エントリをログフィールドに変換する
func convertToFields(entry StructuredLogEntry) []domain.Field {
	fields := []domain.Field{
		{Key: "category", Value: entry.Category},
		{Key: "action", Value: entry.Action},
		{Key: "component", Value: entry.Component},
	}

	// コンテキストを追加
	for k, v := range entry.Context {
		fields = append(fields, domain.Field{Key: k, Value: fmt.Sprintf("%v", v)})
	}

	// パフォーマンス情報を追加
	if entry.Performance != nil {
		fields = append(fields, domain.Field{Key: "success", Value: entry.Performance.Success})
		if entry.Performance.DurationMs > 0 {
			fields = append(fields, domain.Field{Key: "duration_ms", Value: fmt.Sprintf("%.2f", entry.Performance.DurationMs)})
		}
	}

	// セッション情報を追加
	if entry.Session != nil {
		fields = append(fields, domain.Field{Key: "has_current_game", Value: entry.Session.HasCurrentGame})
		if entry.Session.TotalParagraphs > 0 {
			fields = append(fields, domain.Field{Key: "total_paragraphs", Value: entry.Session.TotalParagraphs})
		}
	}

	// AI開発モード時は位置情報を追加
	if IsAIDevelopmentMode() && entry.Location != nil {
		fields = append(fields, domain.Field{Key: "source_file", Value: entry.Location.File})
		fields = append(fields, domain.Field{Key: "source_line", Value: entry.Location.Line})
		fields = append(fields, domain.Field{Key: "source_function", Value: entry.Location.Function})
	}

	return fields
}