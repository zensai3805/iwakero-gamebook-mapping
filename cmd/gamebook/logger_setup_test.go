package main

import (
	"os"
	"testing"

	"github.com/zensai3805/iwakero-gamebook-mapping/internal/domain"
)

func TestSetupLogger_DefaultConfiguration(t *testing.T) {
	// Arrange
	// 環境変数をクリア
	os.Unsetenv("LOG_LEVEL")
	os.Unsetenv("LOG_OUTPUT")
	os.Unsetenv("LOG_FORMAT")

	// Act
	logger, cleanupFunc, err := SetupLogger()

	// Assert
	if err != nil {
		t.Fatalf("SetupLogger()でエラーが発生: %v", err)
	}

	if logger == nil {
		t.Fatal("loggerがnilです")
	}

	if cleanupFunc == nil {
		t.Fatal("cleanupFuncがnilです")
	}

	// Cleanup
	cleanupFunc()
}

func TestSetupLogger_WithLogLevel(t *testing.T) {
	// Arrange
	os.Setenv("LOG_LEVEL", "DEBUG")
	defer os.Unsetenv("LOG_LEVEL")

	// Act
	logger, cleanupFunc, err := SetupLogger()

	// Assert
	if err != nil {
		t.Fatalf("SetupLogger()でエラーが発生: %v", err)
	}

	if logger == nil {
		t.Fatal("loggerがnilです")
	}

	// Cleanup
	cleanupFunc()
}

func TestSetupLogger_WithInvalidLogLevel(t *testing.T) {
	// Arrange
	os.Setenv("LOG_LEVEL", "INVALID")
	defer os.Unsetenv("LOG_LEVEL")

	// Act
	logger, cleanupFunc, err := SetupLogger()

	// Assert
	if err == nil {
		t.Fatal("無効なログレベルでエラーが発生するべきです")
	}

	if logger != nil {
		t.Fatal("エラー時にloggerはnilであるべきです")
	}

	if cleanupFunc != nil {
		t.Fatal("エラー時にcleanupFuncはnilであるべきです")
	}
}

func TestSetupLogger_WithFileOutput(t *testing.T) {
	// Arrange
	os.Setenv("LOG_OUTPUT", "file")
	os.Setenv("LOG_FILE_PATH", "/tmp/test.log")
	defer os.Unsetenv("LOG_OUTPUT")
	defer os.Unsetenv("LOG_FILE_PATH")

	// Act
	logger, cleanupFunc, err := SetupLogger()

	// Assert
	if err != nil {
		t.Fatalf("SetupLogger()でエラーが発生: %v", err)
	}

	if logger == nil {
		t.Fatal("loggerがnilです")
	}

	// Cleanup
	cleanupFunc()
	os.Remove("/tmp/test.log")
}

func TestSetupLogger_WithJSONFormat(t *testing.T) {
	// Arrange
	os.Setenv("LOG_FORMAT", "json")
	defer os.Unsetenv("LOG_FORMAT")

	// Act
	logger, cleanupFunc, err := SetupLogger()

	// Assert
	if err != nil {
		t.Fatalf("SetupLogger()でエラーが発生: %v", err)
	}

	if logger == nil {
		t.Fatal("loggerがnilです")
	}

	// Cleanup
	cleanupFunc()
}

func TestGetGlobalLogger_BeforeSetup(t *testing.T) {
	// Arrange
	resetGlobalLogger()

	// Act
	logger := GetGlobalLogger()

	// Assert
	if logger == nil {
		t.Fatal("グローバルロガーがnilです")
	}
}

func TestGetGlobalLogger_AfterSetup(t *testing.T) {
	// Arrange
	logger, cleanupFunc, err := SetupLogger()
	if err != nil {
		t.Fatalf("SetupLogger()でエラーが発生: %v", err)
	}
	defer cleanupFunc()

	// Act
	globalLogger := GetGlobalLogger()

	// Assert
	if globalLogger == nil {
		t.Fatal("グローバルロガーがnilです")
	}

	if globalLogger != logger {
		t.Fatal("グローバルロガーがセットアップされたロガーと異なります")
	}
}

func TestParseLogLevel_ValidLevels(t *testing.T) {
	tests := []struct {
		input    string
		expected domain.LogLevel
	}{
		{"DEBUG", domain.LogLevelDebug},
		{"INFO", domain.LogLevelInfo},
		{"WARN", domain.LogLevelWarn},
		{"ERROR", domain.LogLevelError},
		{"FATAL", domain.LogLevelFatal},
		{"debug", domain.LogLevelDebug},
		{"info", domain.LogLevelInfo},
		{"warn", domain.LogLevelWarn},
		{"error", domain.LogLevelError},
		{"fatal", domain.LogLevelFatal},
	}

	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			// Act
			level, err := parseLogLevel(test.input)

			// Assert
			if err != nil {
				t.Fatalf("parseLogLevel(%s)でエラーが発生: %v", test.input, err)
			}

			if level != test.expected {
				t.Fatalf("期待値: %v, 実際の値: %v", test.expected, level)
			}
		})
	}
}

func TestParseLogLevel_InvalidLevel(t *testing.T) {
	// Act
	level, err := parseLogLevel("INVALID")

	// Assert
	if err == nil {
		t.Fatal("無効なログレベルでエラーが発生するべきです")
	}

	if level != domain.LogLevelInfo {
		t.Fatal("無効なログレベルの場合、デフォルトのINFOレベルが返されるべきです")
	}
}

// resetGlobalLogger はテスト用にグローバルロガーをリセットする
func resetGlobalLogger() {
	setGlobalLogger(nil)
}
