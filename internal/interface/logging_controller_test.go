package main

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zensai3805/iwakero-gamebook-mapping/internal/domain"
)

func TestLoggingController_NewLoggingController(t *testing.T) {
	// Arrange
	config := LoggingConfig{
		Level:      domain.LogLevelInfo,
		Format:     "text",
		OutputType: "console",
	}

	// Act
	controller, err := NewLoggingController(config)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, controller)

	// Clean up
	err = controller.Close()
	require.NoError(t, err)
}

func TestLoggingController_NewLoggingController_FileOutput(t *testing.T) {
	// Arrange
	tempDir := t.TempDir()
	filePath := filepath.Join(tempDir, "test.log")

	config := LoggingConfig{
		Level:      domain.LogLevelInfo,
		Format:     "json",
		OutputType: "file",
		FilePath:   filePath,
	}

	// Act
	controller, err := NewLoggingController(config)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, controller)

	// Clean up
	err = controller.Close()
	require.NoError(t, err)
}

func TestLoggingController_NewLoggingController_InvalidFormat(t *testing.T) {
	// Arrange
	config := LoggingConfig{
		Level:      domain.LogLevelInfo,
		Format:     "invalid",
		OutputType: "console",
	}

	// Act
	controller, err := NewLoggingController(config)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, controller)
}

func TestLoggingController_NewLoggingController_InvalidOutputType(t *testing.T) {
	// Arrange
	config := LoggingConfig{
		Level:      domain.LogLevelInfo,
		Format:     "text",
		OutputType: "invalid",
	}

	// Act
	controller, err := NewLoggingController(config)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, controller)
}

func TestLoggingController_NewLoggingController_FilePathRequired(t *testing.T) {
	// Arrange
	config := LoggingConfig{
		Level:      domain.LogLevelInfo,
		Format:     "text",
		OutputType: "file",
		FilePath:   "", // 空のファイルパス
	}

	// Act
	controller, err := NewLoggingController(config)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, controller)
}

func TestLoggingController_GetLogger(t *testing.T) {
	// Arrange
	config := LoggingConfig{
		Level:      domain.LogLevelInfo,
		Format:     "text",
		OutputType: "console",
	}

	controller, err := NewLoggingController(config)
	require.NoError(t, err)
	defer controller.Close()

	// Act
	logger := controller.GetLogger()

	// Assert
	assert.NotNil(t, logger)
	assert.Implements(t, (*domain.Logger)(nil), logger)
}

func TestLoggingController_SetLevel(t *testing.T) {
	// Arrange
	config := LoggingConfig{
		Level:      domain.LogLevelInfo,
		Format:     "text",
		OutputType: "console",
	}

	controller, err := NewLoggingController(config)
	require.NoError(t, err)
	defer controller.Close()

	// Act
	controller.SetLevel(domain.LogLevelError)

	// Assert
	assert.Equal(t, domain.LogLevelError, controller.GetLevel())
}

func TestLoggingController_GetLevel(t *testing.T) {
	// Arrange
	config := LoggingConfig{
		Level:      domain.LogLevelWarn,
		Format:     "text",
		OutputType: "console",
	}

	controller, err := NewLoggingController(config)
	require.NoError(t, err)
	defer controller.Close()

	// Act
	level := controller.GetLevel()

	// Assert
	assert.Equal(t, domain.LogLevelWarn, level)
}

func TestLoggingController_Close(t *testing.T) {
	// Arrange
	config := LoggingConfig{
		Level:      domain.LogLevelInfo,
		Format:     "text",
		OutputType: "console",
	}

	controller, err := NewLoggingController(config)
	require.NoError(t, err)

	// Act
	err = controller.Close()

	// Assert
	require.NoError(t, err)

	// 複数回のクローズも問題なし
	err = controller.Close()
	require.NoError(t, err)
}

func TestLoggingController_Integration_ConsoleOutput(t *testing.T) {
	// Arrange
	var buffer bytes.Buffer
	config := LoggingConfig{
		Level:      domain.LogLevelInfo,
		Format:     "text",
		OutputType: "console",
		Writer:     &buffer, // テスト用のバッファ
		SyncMode:   true,    // テスト用の同期モード
	}

	controller, err := NewLoggingController(config)
	require.NoError(t, err)
	defer controller.Close()

	logger := controller.GetLogger()

	// Act
	logger.Info("テストメッセージ", domain.Field{Key: "test", Value: "value"})

	// Assert
	output := buffer.String()
	assert.Contains(t, output, "INFO")
	assert.Contains(t, output, "テストメッセージ")
	assert.Contains(t, output, "test=value")
}

func TestLoggingController_Integration_FileOutput(t *testing.T) {
	// Arrange
	tempDir := t.TempDir()
	filePath := filepath.Join(tempDir, "test.log")

	config := LoggingConfig{
		Level:      domain.LogLevelInfo,
		Format:     "json",
		OutputType: "file",
		FilePath:   filePath,
		SyncMode:   true, // テスト用の同期モード
	}

	controller, err := NewLoggingController(config)
	require.NoError(t, err)
	defer controller.Close()

	logger := controller.GetLogger()

	// Act
	logger.Info("テストメッセージ", domain.Field{Key: "test", Value: "value"})

	// Assert
	// ファイルが作成されていることを確認
	_, statErr := os.Stat(filePath)
	assert.NoError(t, statErr)
}
