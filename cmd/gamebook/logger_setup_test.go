package main

import (
	"os"
	"testing"

	"github.com/zensai3805/iwakero-gamebook-mapping/internal/domain"
)

func TestInitializeLogger_DefaultConfig(t *testing.T) {
	// Arrange
	// 環境変数をクリア
	originalLevel := os.Getenv("LOG_LEVEL")
	originalOutput := os.Getenv("LOG_OUTPUT")
	defer func() {
		os.Setenv("LOG_LEVEL", originalLevel)
		os.Setenv("LOG_OUTPUT", originalOutput)
	}()
	os.Unsetenv("LOG_LEVEL")
	os.Unsetenv("LOG_OUTPUT")

	// Act
	logger, err := InitializeLogger()

	// Assert
	if err != nil {
		t.Fatalf("InitializeLogger failed: %v", err)
	}
	if logger == nil {
		t.Fatal("InitializeLogger returned nil logger")
	}

	// デフォルトレベルはINFO
	globalLogger := GetGlobalLogger()
	if globalLogger == nil {
		t.Fatal("Global logger was not set")
	}
}

func TestInitializeLogger_WithEnvironmentVariables(t *testing.T) {
	// Arrange
	originalLevel := os.Getenv("LOG_LEVEL")
	originalOutput := os.Getenv("LOG_OUTPUT")
	defer func() {
		os.Setenv("LOG_LEVEL", originalLevel)
		os.Setenv("LOG_OUTPUT", originalOutput)
	}()

	os.Setenv("LOG_LEVEL", "DEBUG")
	os.Setenv("LOG_OUTPUT", "console")

	// Act
	logger, err := InitializeLogger()

	// Assert
	if err != nil {
		t.Fatalf("InitializeLogger failed: %v", err)
	}
	if logger == nil {
		t.Fatal("InitializeLogger returned nil logger")
	}

	// グローバルロガーが設定されているか確認
	globalLogger := GetGlobalLogger()
	if globalLogger == nil {
		t.Fatal("Global logger was not set")
	}
}

func TestSetGlobalLogger(t *testing.T) {
	// Arrange
	testLogger := &simpleMockLogger{}

	// Act
	SetGlobalLogger(testLogger)

	// Assert
	retrievedLogger := GetGlobalLogger()
	if retrievedLogger != testLogger {
		t.Fatal("SetGlobalLogger did not set the expected logger")
	}
}

func TestGetGlobalLogger_WhenNotInitialized(t *testing.T) {
	// Arrange
	// グローバルロガーをリセット
	SetGlobalLogger(nil)

	// Act
	logger := GetGlobalLogger()

	// Assert
	if logger != nil {
		t.Fatal("GetGlobalLogger should return nil when not initialized")
	}
}

func TestCleanupLogger(t *testing.T) {
	// Arrange
	testLogger := &testMockLogger{closed: false}
	SetGlobalLogger(testLogger)

	// Act
	err := CleanupLogger()

	// Assert
	if err != nil {
		t.Fatalf("CleanupLogger failed: %v", err)
	}
	if !testLogger.closed {
		t.Fatal("CleanupLogger did not close the logger")
	}
	if GetGlobalLogger() != nil {
		t.Fatal("CleanupLogger did not reset global logger")
	}
}

func TestCleanupLogger_WhenNotInitialized(t *testing.T) {
	// Arrange
	SetGlobalLogger(nil)

	// Act
	err := CleanupLogger()

	// Assert
	if err != nil {
		t.Fatalf("CleanupLogger should not fail when not initialized: %v", err)
	}
}

// testMockLogger はテスト用のロガーモック（クローズテスト用）
// simpleMockLogger はテスト用のシンプルなモックロガー
type simpleMockLogger struct{}

func (m *simpleMockLogger) Debug(msg string, fields ...domain.Field) {}
func (m *simpleMockLogger) Info(msg string, fields ...domain.Field)  {}
func (m *simpleMockLogger) Warn(msg string, fields ...domain.Field)  {}
func (m *simpleMockLogger) Error(msg string, fields ...domain.Field) {}
func (m *simpleMockLogger) Fatal(msg string, fields ...domain.Field) {}
func (m *simpleMockLogger) WithContext(fields ...domain.Field) domain.Logger {
	return m
}
func (m *simpleMockLogger) Close() error {
	return nil
}

// testMockLogger はクローズテスト用のモックロガー
type testMockLogger struct {
	closed bool
}

func (m *testMockLogger) Debug(msg string, fields ...domain.Field) {}
func (m *testMockLogger) Info(msg string, fields ...domain.Field)  {}
func (m *testMockLogger) Warn(msg string, fields ...domain.Field)  {}
func (m *testMockLogger) Error(msg string, fields ...domain.Field) {}
func (m *testMockLogger) Fatal(msg string, fields ...domain.Field) {}
func (m *testMockLogger) WithContext(fields ...domain.Field) domain.Logger {
	return m
}
func (m *testMockLogger) Close() error {
	m.closed = true
	return nil
}
