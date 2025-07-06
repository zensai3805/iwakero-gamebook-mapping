package logger

import (
	"fmt"
	"strings"
	"time"

	"github.com/zensai3805/iwakero-gamebook-mapping/internal/domain"
)

// TextFormatter はテキスト形式でログエントリをフォーマットする実装
type TextFormatter struct{}

// NewTextFormatter は新しいTextFormatterを生成する
func NewTextFormatter() *TextFormatter {
	return &TextFormatter{}
}

// Format はログエントリをテキスト形式のバイト列に変換する
func (f *TextFormatter) Format(entry domain.LogEntry) ([]byte, error) {
	var builder strings.Builder

	// タイムスタンプ、レベル、メッセージの基本情報を出力
	builder.WriteString(fmt.Sprintf("%s %s %s",
		entry.Timestamp.Format(time.RFC3339),
		entry.Level.String(),
		entry.Message,
	))

	// フィールドがある場合は追加
	for _, field := range entry.Fields {
		// 特殊文字のエスケープ処理
		valueStr := fmt.Sprintf("%v", field.Value)
		valueStr = strings.ReplaceAll(valueStr, "\n", "\\n")
		valueStr = strings.ReplaceAll(valueStr, "\t", "\\t")

		builder.WriteString(fmt.Sprintf(" %s=%s", field.Key, valueStr))
	}

	return []byte(builder.String()), nil
}
