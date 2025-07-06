package logger

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zensai3805/iwakero-gamebook-mapping/internal/domain"
)

func TestFileWriter_NewFileWriter(t *testing.T) {
	// Arrange
	tempDir := t.TempDir()
	filePath := filepath.Join(tempDir, "test.log")
	formatter := NewTextFormatter()

	// Act
	writer, err := NewFileWriter(filePath, formatter)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, writer)

	// Clean up
	err = writer.Close()
	require.NoError(t, err)
}

func TestFileWriter_Write_SingleEntry(t *testing.T) {
	// Arrange
	tempDir := t.TempDir()
	filePath := filepath.Join(tempDir, "test.log")
	formatter := NewTextFormatter()

	writer, err := NewFileWriter(filePath, formatter)
	require.NoError(t, err)
	defer writer.Close()

	timestamp := time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC)
	entry := domain.NewLogEntry(timestamp, domain.LogLevelInfo, "テストメッセージ", nil)

	// Act
	err = writer.Write([]domain.LogEntry{entry})

	// Assert
	require.NoError(t, err)

	// ファイルの内容を確認
	content, err := os.ReadFile(filePath)
	require.NoError(t, err)

	contentStr := string(content)
	assert.Contains(t, contentStr, "2023-01-01T12:00:00Z")
	assert.Contains(t, contentStr, "INFO")
	assert.Contains(t, contentStr, "テストメッセージ")
}

func TestFileWriter_Write_MultipleEntries(t *testing.T) {
	// Arrange
	tempDir := t.TempDir()
	filePath := filepath.Join(tempDir, "test.log")
	formatter := NewTextFormatter()

	writer, err := NewFileWriter(filePath, formatter)
	require.NoError(t, err)
	defer writer.Close()

	timestamp := time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC)
	entries := []domain.LogEntry{
		domain.NewLogEntry(timestamp, domain.LogLevelInfo, "メッセージ1", nil),
		domain.NewLogEntry(timestamp, domain.LogLevelError, "メッセージ2", nil),
		domain.NewLogEntry(timestamp, domain.LogLevelWarn, "メッセージ3", nil),
	}

	// Act
	err = writer.Write(entries)

	// Assert
	require.NoError(t, err)

	// ファイルの内容を確認
	content, err := os.ReadFile(filePath)
	require.NoError(t, err)

	contentStr := string(content)
	lines := strings.Split(strings.TrimSpace(contentStr), "\n")
	assert.Len(t, lines, 3)

	assert.Contains(t, lines[0], "メッセージ1")
	assert.Contains(t, lines[1], "メッセージ2")
	assert.Contains(t, lines[2], "メッセージ3")
}

func TestFileWriter_Write_AppendMode(t *testing.T) {
	// Arrange
	tempDir := t.TempDir()
	filePath := filepath.Join(tempDir, "test.log")
	formatter := NewTextFormatter()

	// 最初の書き込み
	writer1, err := NewFileWriter(filePath, formatter)
	require.NoError(t, err)

	timestamp := time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC)
	entry1 := domain.NewLogEntry(timestamp, domain.LogLevelInfo, "メッセージ1", nil)

	err = writer1.Write([]domain.LogEntry{entry1})
	require.NoError(t, err)

	err = writer1.Close()
	require.NoError(t, err)

	// 2回目の書き込み
	writer2, err := NewFileWriter(filePath, formatter)
	require.NoError(t, err)
	defer writer2.Close()

	entry2 := domain.NewLogEntry(timestamp, domain.LogLevelError, "メッセージ2", nil)

	// Act
	err = writer2.Write([]domain.LogEntry{entry2})

	// Assert
	require.NoError(t, err)

	// ファイルの内容を確認
	content, err := os.ReadFile(filePath)
	require.NoError(t, err)

	contentStr := string(content)
	lines := strings.Split(strings.TrimSpace(contentStr), "\n")
	assert.Len(t, lines, 2)

	assert.Contains(t, lines[0], "メッセージ1")
	assert.Contains(t, lines[1], "メッセージ2")
}

func TestFileWriter_Write_EmptyEntries(t *testing.T) {
	// Arrange
	tempDir := t.TempDir()
	filePath := filepath.Join(tempDir, "test.log")
	formatter := NewTextFormatter()

	writer, err := NewFileWriter(filePath, formatter)
	require.NoError(t, err)
	defer writer.Close()

	// Act
	err = writer.Write([]domain.LogEntry{})

	// Assert
	require.NoError(t, err)

	// ファイルが存在し、空であることを確認
	content, err := os.ReadFile(filePath)
	require.NoError(t, err)
	assert.Empty(t, string(content))
}

func TestFileWriter_Close(t *testing.T) {
	// Arrange
	tempDir := t.TempDir()
	filePath := filepath.Join(tempDir, "test.log")
	formatter := NewTextFormatter()

	writer, err := NewFileWriter(filePath, formatter)
	require.NoError(t, err)

	// Act
	err = writer.Close()

	// Assert
	require.NoError(t, err)

	// ファイルがクローズされた後の書き込みは失敗することを確認
	timestamp := time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC)
	entry := domain.NewLogEntry(timestamp, domain.LogLevelInfo, "テストメッセージ", nil)

	err = writer.Write([]domain.LogEntry{entry})
	assert.Error(t, err)
}

func TestFileWriter_NewFileWriter_InvalidPath(t *testing.T) {
	// Arrange
	invalidPath := "/invalid/path/that/does/not/exist/test.log"
	formatter := NewTextFormatter()

	// Act
	writer, err := NewFileWriter(invalidPath, formatter)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, writer)
}

func TestFileWriter_Write_WithJSONFormatter(t *testing.T) {
	// Arrange
	tempDir := t.TempDir()
	filePath := filepath.Join(tempDir, "test.log")
	formatter := NewJSONFormatter()

	writer, err := NewFileWriter(filePath, formatter)
	require.NoError(t, err)
	defer writer.Close()

	timestamp := time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC)
	fields := []domain.Field{
		{Key: "user_id", Value: "123"},
		{Key: "action", Value: "login"},
	}
	entry := domain.NewLogEntry(timestamp, domain.LogLevelInfo, "ユーザーログイン", fields)

	// Act
	err = writer.Write([]domain.LogEntry{entry})

	// Assert
	require.NoError(t, err)

	// ファイルの内容を確認
	content, err := os.ReadFile(filePath)
	require.NoError(t, err)

	contentStr := string(content)
	assert.Contains(t, contentStr, `"timestamp":"2023-01-01T12:00:00Z"`)
	assert.Contains(t, contentStr, `"level":"INFO"`)
	assert.Contains(t, contentStr, `"message":"ユーザーログイン"`)
	assert.Contains(t, contentStr, `"user_id":"123"`)
	assert.Contains(t, contentStr, `"action":"login"`)
}
