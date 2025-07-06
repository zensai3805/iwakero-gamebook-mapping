package logger

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zensai3805/iwakero-gamebook-mapping/internal/domain"
)

func TestJSONFormatter_Format_BasicEntry(t *testing.T) {
	// Arrange
	formatter := NewJSONFormatter()
	timestamp := time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC)
	entry := domain.NewLogEntry(timestamp, domain.LogLevelInfo, "テストメッセージ", nil)

	// Act
	result, err := formatter.Format(entry)

	// Assert
	require.NoError(t, err)

	var data map[string]interface{}
	err = json.Unmarshal(result, &data)
	require.NoError(t, err)

	assert.Equal(t, "2023-01-01T12:00:00Z", data["timestamp"])
	assert.Equal(t, "INFO", data["level"])
	assert.Equal(t, "テストメッセージ", data["message"])
}

func TestJSONFormatter_Format_WithFields(t *testing.T) {
	// Arrange
	formatter := NewJSONFormatter()
	timestamp := time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC)
	fields := []domain.Field{
		{Key: "user_id", Value: "123"},
		{Key: "action", Value: "login"},
		{Key: "count", Value: 42},
	}
	entry := domain.NewLogEntry(timestamp, domain.LogLevelInfo, "ユーザーログイン", fields)

	// Act
	result, err := formatter.Format(entry)

	// Assert
	require.NoError(t, err)

	var data map[string]interface{}
	err = json.Unmarshal(result, &data)
	require.NoError(t, err)

	assert.Equal(t, "2023-01-01T12:00:00Z", data["timestamp"])
	assert.Equal(t, "INFO", data["level"])
	assert.Equal(t, "ユーザーログイン", data["message"])
	assert.Equal(t, "123", data["user_id"])
	assert.Equal(t, "login", data["action"])
	assert.Equal(t, float64(42), data["count"]) // JSONでは数値はfloat64になる
}

func TestJSONFormatter_Format_AllLogLevels(t *testing.T) {
	// Arrange
	formatter := NewJSONFormatter()
	timestamp := time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC)

	testCases := []struct {
		level    domain.LogLevel
		expected string
	}{
		{domain.LogLevelDebug, "DEBUG"},
		{domain.LogLevelInfo, "INFO"},
		{domain.LogLevelWarn, "WARN"},
		{domain.LogLevelError, "ERROR"},
		{domain.LogLevelFatal, "FATAL"},
	}

	for _, testCase := range testCases {
		t.Run(testCase.expected, func(t *testing.T) {
			// Arrange
			entry := domain.NewLogEntry(timestamp, testCase.level, "テストメッセージ", nil)

			// Act
			result, err := formatter.Format(entry)

			// Assert
			require.NoError(t, err)

			var data map[string]interface{}
			err = json.Unmarshal(result, &data)
			require.NoError(t, err)

			assert.Equal(t, testCase.expected, data["level"])
		})
	}
}

func TestJSONFormatter_Format_EmptyMessage(t *testing.T) {
	// Arrange
	formatter := NewJSONFormatter()
	timestamp := time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC)
	entry := domain.NewLogEntry(timestamp, domain.LogLevelInfo, "", nil)

	// Act
	result, err := formatter.Format(entry)

	// Assert
	require.NoError(t, err)

	var data map[string]interface{}
	err = json.Unmarshal(result, &data)
	require.NoError(t, err)

	assert.Equal(t, "", data["message"])
}

func TestJSONFormatter_Format_NilFields(t *testing.T) {
	// Arrange
	formatter := NewJSONFormatter()
	timestamp := time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC)
	entry := domain.NewLogEntry(timestamp, domain.LogLevelInfo, "テストメッセージ", nil)

	// Act
	result, err := formatter.Format(entry)

	// Assert
	require.NoError(t, err)

	var data map[string]interface{}
	err = json.Unmarshal(result, &data)
	require.NoError(t, err)

	// 基本フィールドのみ存在することを確認
	assert.Contains(t, data, "timestamp")
	assert.Contains(t, data, "level")
	assert.Contains(t, data, "message")
	assert.Len(t, data, 3)
}

func TestJSONFormatter_Format_SpecialCharacters(t *testing.T) {
	// Arrange
	formatter := NewJSONFormatter()
	timestamp := time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC)
	fields := []domain.Field{
		{Key: "path", Value: "/path/with\"quotes"},
		{Key: "data", Value: "line1\nline2"},
	}
	entry := domain.NewLogEntry(timestamp, domain.LogLevelInfo, "特殊文字テスト", fields)

	// Act
	result, err := formatter.Format(entry)

	// Assert
	require.NoError(t, err)

	var data map[string]interface{}
	err = json.Unmarshal(result, &data)
	require.NoError(t, err)

	assert.Equal(t, "/path/with\"quotes", data["path"])
	assert.Equal(t, "line1\nline2", data["data"])
}
