package logger

import (
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zensai3805/iwakero-gamebook-mapping/internal/domain"
)

func TestTextFormatter_Format_BasicEntry(t *testing.T) {
	// Arrange
	formatter := NewTextFormatter()
	timestamp := time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC)
	entry := domain.NewLogEntry(timestamp, domain.LogLevelInfo, "テストメッセージ", nil)

	// Act
	result, err := formatter.Format(entry)

	// Assert
	require.NoError(t, err)

	resultStr := string(result)
	assert.Contains(t, resultStr, "2023-01-01T12:00:00Z")
	assert.Contains(t, resultStr, "INFO")
	assert.Contains(t, resultStr, "テストメッセージ")
}

func TestTextFormatter_Format_WithFields(t *testing.T) {
	// Arrange
	formatter := NewTextFormatter()
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

	resultStr := string(result)
	assert.Contains(t, resultStr, "2023-01-01T12:00:00Z")
	assert.Contains(t, resultStr, "INFO")
	assert.Contains(t, resultStr, "ユーザーログイン")
	assert.Contains(t, resultStr, "user_id=123")
	assert.Contains(t, resultStr, "action=login")
	assert.Contains(t, resultStr, "count=42")
}

func TestTextFormatter_Format_AllLogLevels(t *testing.T) {
	// Arrange
	formatter := NewTextFormatter()
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

			resultStr := string(result)
			assert.Contains(t, resultStr, testCase.expected)
		})
	}
}

func TestTextFormatter_Format_EmptyMessage(t *testing.T) {
	// Arrange
	formatter := NewTextFormatter()
	timestamp := time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC)
	entry := domain.NewLogEntry(timestamp, domain.LogLevelInfo, "", nil)

	// Act
	result, err := formatter.Format(entry)

	// Assert
	require.NoError(t, err)

	resultStr := string(result)
	assert.Contains(t, resultStr, "2023-01-01T12:00:00Z")
	assert.Contains(t, resultStr, "INFO")
	// 空のメッセージでも適切にフォーマットされることを確認
	assert.NotEmpty(t, resultStr)
}

func TestTextFormatter_Format_NilFields(t *testing.T) {
	// Arrange
	formatter := NewTextFormatter()
	timestamp := time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC)
	entry := domain.NewLogEntry(timestamp, domain.LogLevelInfo, "テストメッセージ", nil)

	// Act
	result, err := formatter.Format(entry)

	// Assert
	require.NoError(t, err)

	resultStr := string(result)
	assert.Contains(t, resultStr, "2023-01-01T12:00:00Z")
	assert.Contains(t, resultStr, "INFO")
	assert.Contains(t, resultStr, "テストメッセージ")
	// フィールドがない場合は基本情報のみ出力される
	assert.Equal(t, 3, len(strings.Fields(resultStr)))
}

func TestTextFormatter_Format_SpecialCharacters(t *testing.T) {
	// Arrange
	formatter := NewTextFormatter()
	timestamp := time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC)
	fields := []domain.Field{
		{Key: "path", Value: "/path/with spaces"},
		{Key: "data", Value: "line1\nline2"},
	}
	entry := domain.NewLogEntry(timestamp, domain.LogLevelInfo, "特殊文字テスト", fields)

	// Act
	result, err := formatter.Format(entry)

	// Assert
	require.NoError(t, err)

	resultStr := string(result)
	assert.Contains(t, resultStr, "特殊文字テスト")
	assert.Contains(t, resultStr, "path=/path/with spaces")
	assert.Contains(t, resultStr, "data=line1\\nline2")
}

func TestTextFormatter_Format_FieldsOrder(t *testing.T) {
	// Arrange
	formatter := NewTextFormatter()
	timestamp := time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC)
	fields := []domain.Field{
		{Key: "first", Value: "1"},
		{Key: "second", Value: "2"},
		{Key: "third", Value: "3"},
	}
	entry := domain.NewLogEntry(timestamp, domain.LogLevelInfo, "フィールド順序テスト", fields)

	// Act
	result, err := formatter.Format(entry)

	// Assert
	require.NoError(t, err)

	resultStr := string(result)

	// フィールドが指定された順序で出力されることを確認
	firstPos := strings.Index(resultStr, "first=1")
	secondPos := strings.Index(resultStr, "second=2")
	thirdPos := strings.Index(resultStr, "third=3")

	assert.True(t, firstPos < secondPos)
	assert.True(t, secondPos < thirdPos)
}
