package logger

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zensai3805/iwakero-gamebook-mapping/internal/domain"
)

func TestConsoleWriter_NewConsoleWriter(t *testing.T) {
	// Arrange
	var buffer bytes.Buffer
	formatter := NewTextFormatter()

	// Act
	writer := NewConsoleWriter(&buffer, formatter)

	// Assert
	assert.NotNil(t, writer)
}

func TestConsoleWriter_Write_SingleEntry(t *testing.T) {
	// Arrange
	var buffer bytes.Buffer
	formatter := NewTextFormatter()
	writer := NewConsoleWriter(&buffer, formatter)

	timestamp := time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC)
	entry := domain.NewLogEntry(timestamp, domain.LogLevelInfo, "テストメッセージ", nil)

	// Act
	err := writer.Write([]domain.LogEntry{entry})

	// Assert
	require.NoError(t, err)

	output := buffer.String()
	assert.Contains(t, output, "2023-01-01T12:00:00Z")
	assert.Contains(t, output, "INFO")
	assert.Contains(t, output, "テストメッセージ")
}

func TestConsoleWriter_Write_MultipleEntries(t *testing.T) {
	// Arrange
	var buffer bytes.Buffer
	formatter := NewTextFormatter()
	writer := NewConsoleWriter(&buffer, formatter)

	timestamp := time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC)
	entries := []domain.LogEntry{
		domain.NewLogEntry(timestamp, domain.LogLevelInfo, "メッセージ1", nil),
		domain.NewLogEntry(timestamp, domain.LogLevelError, "メッセージ2", nil),
		domain.NewLogEntry(timestamp, domain.LogLevelWarn, "メッセージ3", nil),
	}

	// Act
	err := writer.Write(entries)

	// Assert
	require.NoError(t, err)

	output := buffer.String()
	lines := strings.Split(strings.TrimSpace(output), "\n")
	assert.Len(t, lines, 3)

	assert.Contains(t, lines[0], "メッセージ1")
	assert.Contains(t, lines[1], "メッセージ2")
	assert.Contains(t, lines[2], "メッセージ3")
}

func TestConsoleWriter_Write_EmptyEntries(t *testing.T) {
	// Arrange
	var buffer bytes.Buffer
	formatter := NewTextFormatter()
	writer := NewConsoleWriter(&buffer, formatter)

	// Act
	err := writer.Write([]domain.LogEntry{})

	// Assert
	require.NoError(t, err)

	output := buffer.String()
	assert.Empty(t, output)
}

func TestConsoleWriter_Write_WithJSONFormatter(t *testing.T) {
	// Arrange
	var buffer bytes.Buffer
	formatter := NewJSONFormatter()
	writer := NewConsoleWriter(&buffer, formatter)

	timestamp := time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC)
	fields := []domain.Field{
		{Key: "user_id", Value: "123"},
		{Key: "action", Value: "login"},
	}
	entry := domain.NewLogEntry(timestamp, domain.LogLevelInfo, "ユーザーログイン", fields)

	// Act
	err := writer.Write([]domain.LogEntry{entry})

	// Assert
	require.NoError(t, err)

	output := buffer.String()
	assert.Contains(t, output, `"timestamp":"2023-01-01T12:00:00Z"`)
	assert.Contains(t, output, `"level":"INFO"`)
	assert.Contains(t, output, `"message":"ユーザーログイン"`)
	assert.Contains(t, output, `"user_id":"123"`)
	assert.Contains(t, output, `"action":"login"`)
}

func TestConsoleWriter_Write_WithLevelFilter(t *testing.T) {
	// Arrange
	var buffer bytes.Buffer
	formatter := NewTextFormatter()
	writer := NewConsoleWriterWithLevelFilter(&buffer, formatter, domain.LogLevelWarn)

	timestamp := time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC)
	entries := []domain.LogEntry{
		domain.NewLogEntry(timestamp, domain.LogLevelDebug, "デバッグメッセージ", nil),
		domain.NewLogEntry(timestamp, domain.LogLevelInfo, "情報メッセージ", nil),
		domain.NewLogEntry(timestamp, domain.LogLevelWarn, "警告メッセージ", nil),
		domain.NewLogEntry(timestamp, domain.LogLevelError, "エラーメッセージ", nil),
	}

	// Act
	err := writer.Write(entries)

	// Assert
	require.NoError(t, err)

	output := buffer.String()
	lines := strings.Split(strings.TrimSpace(output), "\n")

	// WARNレベル以上のメッセージのみが出力されることを確認
	assert.Len(t, lines, 2)
	assert.Contains(t, output, "警告メッセージ")
	assert.Contains(t, output, "エラーメッセージ")
	assert.NotContains(t, output, "デバッグメッセージ")
	assert.NotContains(t, output, "情報メッセージ")
}

func TestConsoleWriter_Close(t *testing.T) {
	// Arrange
	var buffer bytes.Buffer
	formatter := NewTextFormatter()
	writer := NewConsoleWriter(&buffer, formatter)

	// Act
	err := writer.Close()

	// Assert
	require.NoError(t, err)

	// コンソールライターのクローズ後も書き込みは可能（標準出力は常に利用可能）
	timestamp := time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC)
	entry := domain.NewLogEntry(timestamp, domain.LogLevelInfo, "テストメッセージ", nil)

	err = writer.Write([]domain.LogEntry{entry})
	require.NoError(t, err)
}
