package main

import (
	"os"
	"testing"
)

func TestMain_LoggerInitialization(t *testing.T) {
	// Arrange
	originalArgs := os.Args
	defer func() { os.Args = originalArgs }()

	// テスト用引数設定（対話モード避け）
	os.Args = []string{"gamebook", "list"}

	// 環境変数設定
	originalLevel := os.Getenv("LOG_LEVEL")
	originalOutput := os.Getenv("LOG_OUTPUT")
	defer func() {
		os.Setenv("LOG_LEVEL", originalLevel)
		os.Setenv("LOG_OUTPUT", originalOutput)
	}()
	os.Setenv("LOG_LEVEL", "DEBUG")
	os.Setenv("LOG_OUTPUT", "console")

	// Act & Assert
	// main関数の内部処理でログシステムが初期化されることをテスト
	// 実際のmain関数実行はテスト環境では困難なので、
	// initializeApplicationLogger関数を直接テストする
	logger, err := initializeApplicationLogger()
	if err != nil {
		t.Fatalf("Logger initialization failed: %v", err)
	}
	if logger == nil {
		t.Fatal("Logger is nil")
	}

	// グローバルロガーが設定されているか確認
	globalLogger := GetGlobalLogger()
	if globalLogger == nil {
		t.Fatal("Global logger was not set")
	}

	// クリーンアップテスト
	err = cleanupApplicationLogger()
	if err != nil {
		t.Fatalf("Logger cleanup failed: %v", err)
	}

	// クリーンアップ後はnilになっているか確認
	globalLogger = GetGlobalLogger()
	if globalLogger != nil {
		t.Fatal("Global logger was not cleaned up")
	}
}

func TestInitializeApplicationLogger_WithConfig(t *testing.T) {
	// Arrange
	originalLevel := os.Getenv("LOG_LEVEL")
	originalOutput := os.Getenv("LOG_OUTPUT")
	originalFormat := os.Getenv("LOG_FORMAT")
	defer func() {
		os.Setenv("LOG_LEVEL", originalLevel)
		os.Setenv("LOG_OUTPUT", originalOutput)
		os.Setenv("LOG_FORMAT", originalFormat)
	}()

	os.Setenv("LOG_LEVEL", "WARN")
	os.Setenv("LOG_OUTPUT", "console")
	os.Setenv("LOG_FORMAT", "json")

	// Act
	logger, err := initializeApplicationLogger()

	// Assert
	if err != nil {
		t.Fatalf("Logger initialization failed: %v", err)
	}
	if logger == nil {
		t.Fatal("Logger is nil")
	}

	// クリーンアップ
	defer func() {
		if err := cleanupApplicationLogger(); err != nil {
			t.Logf("ログクリーンアップ警告: %v", err)
		}
	}()
}

func TestCleanupApplicationLogger_Multiple(t *testing.T) {
	// Arrange
	logger, err := initializeApplicationLogger()
	if err != nil {
		t.Fatalf("Logger initialization failed: %v", err)
	}
	if logger == nil {
		t.Fatal("Logger is nil")
	}

	// Act - 複数回クリーンアップを実行
	err1 := cleanupApplicationLogger()
	err2 := cleanupApplicationLogger()

	// Assert - 複数回実行してもエラーにならない
	if err1 != nil {
		t.Fatalf("First cleanup failed: %v", err1)
	}
	if err2 != nil {
		t.Fatalf("Second cleanup failed: %v", err2)
	}
}
