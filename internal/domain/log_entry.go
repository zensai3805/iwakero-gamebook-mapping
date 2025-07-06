package domain

import (
	"encoding/json"
	"time"
)

// Field はログエントリーのフィールドを表す
type Field struct {
	Key   string
	Value interface{}
}

// LogEntry はログエントリーを表すエンティティ
type LogEntry struct {
	Timestamp time.Time
	Level     LogLevel
	Message   string
	Fields    []Field
}

// NewLogEntry は新しいログエントリーを生成する
func NewLogEntry(timestamp time.Time, level LogLevel, message string, fields []Field) LogEntry {
	// フィールドのコピーを作成して不変性を保証
	fieldsCopy := make([]Field, len(fields))
	copy(fieldsCopy, fields)

	return LogEntry{
		Timestamp: timestamp,
		Level:     level,
		Message:   message,
		Fields:    fieldsCopy,
	}
}

// ToJSON はログエントリーをJSON形式に変換する
func (e LogEntry) ToJSON() (string, error) {
	data := map[string]interface{}{
		"timestamp": e.Timestamp.Format(time.RFC3339Nano),
		"level":     e.Level.String(),
		"message":   e.Message,
	}

	// フィールドをマップに追加
	for _, field := range e.Fields {
		data[field.Key] = field.Value
	}

	jsonBytes, err := json.Marshal(data)
	if err != nil {
		return "", err
	}

	return string(jsonBytes), nil
}

// WithField は新しいフィールドを追加した新しいログエントリーを返す（不変性を保持）
func (e LogEntry) WithField(key string, value interface{}) LogEntry {
	newFields := make([]Field, len(e.Fields)+1)
	copy(newFields, e.Fields)
	newFields[len(e.Fields)] = Field{Key: key, Value: value}

	return LogEntry{
		Timestamp: e.Timestamp,
		Level:     e.Level,
		Message:   e.Message,
		Fields:    newFields,
	}
}

// WithFields は複数のフィールドを追加した新しいログエントリーを返す（不変性を保持）
func (e LogEntry) WithFields(fields []Field) LogEntry {
	newFields := make([]Field, len(e.Fields)+len(fields))
	copy(newFields, e.Fields)
	copy(newFields[len(e.Fields):], fields)

	return LogEntry{
		Timestamp: e.Timestamp,
		Level:     e.Level,
		Message:   e.Message,
		Fields:    newFields,
	}
}

// GetField は指定されたキーのフィールド値を取得する
func (e LogEntry) GetField(key string) (interface{}, bool) {
	for _, field := range e.Fields {
		if field.Key == key {
			return field.Value, true
		}
	}
	return nil, false
}
