package logger

import (
	"encoding/json"
	"time"

	"github.com/zensai3805/iwakero-gamebook-mapping/internal/domain"
)

// JSONFormatter はJSON形式でログエントリをフォーマットする実装
type JSONFormatter struct{}

// NewJSONFormatter は新しいJSONFormatterを生成する
func NewJSONFormatter() *JSONFormatter {
	return &JSONFormatter{}
}

// Format はログエントリをJSON形式のバイト列に変換する
func (f *JSONFormatter) Format(entry domain.LogEntry) ([]byte, error) {
	data := map[string]interface{}{
		"timestamp": entry.Timestamp.Format(time.RFC3339),
		"level":     entry.Level.String(),
		"message":   entry.Message,
	}

	// フィールドをマップに追加
	for _, field := range entry.Fields {
		data[field.Key] = field.Value
	}

	jsonBytes, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	return jsonBytes, nil
}
