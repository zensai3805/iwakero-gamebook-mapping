package main

import (
	"bytes"
	"os"
	"strings"
	"testing"

	"github.com/zensai3805/iwakero-gamebook-mapping/internal/domain"
)

func TestIntegration_LoggingSystem_EndToEnd(t *testing.T) {
	// Arrange
	// テスト用の設定
	originalLevel := os.Getenv("LOG_LEVEL")
	originalOutput := os.Getenv("LOG_OUTPUT")
	originalFormat := os.Getenv("LOG_FORMAT")
	defer func() {
		os.Setenv("LOG_LEVEL", originalLevel)
		os.Setenv("LOG_OUTPUT", originalOutput)
		os.Setenv("LOG_FORMAT", originalFormat)
	}()

	os.Setenv("LOG_LEVEL", "DEBUG")
	os.Setenv("LOG_OUTPUT", "console")
	os.Setenv("LOG_FORMAT", "text")

	// テスト用のバッファ
	logBuffer := &bytes.Buffer{}

	// Act - ログシステムを初期化
	config := LoggingConfig{
		Level:      domain.LogLevelDebug,
		Format:     "text",
		OutputType: "console",
		Writer:     logBuffer, // テスト用ライター
		SyncMode:   true,      // テスト用同期モード
	}

	controller, err := NewLoggingController(config)
	if err != nil {
		t.Fatalf("LoggingController作成に失敗: %v", err)
	}
	defer controller.Close()

	logger := controller.GetLogger()
	SetGlobalLogger(logger)
	defer SetGlobalLogger(nil)

	// Act - 各種操作を実行してログが出力されるか確認
	logger.Debug("デバッグメッセージ", domain.Field{Key: "test", Value: "debug"})
	logger.Info("情報メッセージ", domain.Field{Key: "test", Value: "info"})
	logger.Warn("警告メッセージ", domain.Field{Key: "test", Value: "warn"})
	logger.Error("エラーメッセージ", domain.Field{Key: "test", Value: "error"})

	// コンテキスト付きロガーのテスト
	contextLogger := logger.WithContext(domain.Field{Key: "component", Value: "test"})
	contextLogger.Info("コンテキスト付きメッセージ")

	// Assert - ログが出力されているか確認
	logOutput := logBuffer.String()

	expectedMessages := []string{
		"デバッグメッセージ",
		"情報メッセージ",
		"警告メッセージ",
		"エラーメッセージ",
		"コンテキスト付きメッセージ",
	}

	for _, expectedMsg := range expectedMessages {
		if !strings.Contains(logOutput, expectedMsg) {
			t.Errorf("期待されるメッセージが見つかりません: %s\nActual output: %s", expectedMsg, logOutput)
		}
	}

	// レベルフィルタリングのテスト
	if !strings.Contains(logOutput, "DEBUG") {
		t.Error("DEBUGレベルのログが出力されていません")
	}
	if !strings.Contains(logOutput, "INFO") {
		t.Error("INFOレベルのログが出力されていません")
	}
	if !strings.Contains(logOutput, "WARN") {
		t.Error("WARNレベルのログが出力されていません")
	}
	if !strings.Contains(logOutput, "ERROR") {
		t.Error("ERRORレベルのログが出力されていません")
	}
}

func TestIntegration_CommandExecution_WithLogging(t *testing.T) {
	// Arrange
	// テスト用ディレクトリ設定
	tempDir, err := os.MkdirTemp("", "integration_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	originalDataDir := dataDir
	defer func() { dataDir = originalDataDir }()
	dataDir = tempDir

	// テスト用ログバッファ
	logBuffer := &bytes.Buffer{}
	config := LoggingConfig{
		Level:      domain.LogLevelInfo,
		Format:     "text",
		OutputType: "console",
		Writer:     logBuffer,
		SyncMode:   true,
	}

	controller, err := NewLoggingController(config)
	if err != nil {
		t.Fatalf("LoggingController作成に失敗: %v", err)
	}
	defer controller.Close()

	logger := controller.GetLogger()
	SetGlobalLogger(logger)
	defer SetGlobalLogger(nil)

	// Act - CLIExecutorでコマンドを実行
	executor := NewCLIExecutor()

	// ゲームブック作成
	err = executor.ExecuteNewCommand("テストゲーム")
	if err != nil {
		t.Fatalf("ExecuteNewCommand failed: %v", err)
	}

	// パラグラフ追加
	err = executor.ExecuteAddCommand(1, "物語の始まり")
	if err != nil {
		t.Fatalf("ExecuteAddCommand failed: %v", err)
	}

	// 選択肢追加
	err = executor.ExecuteChoiceCommand(1, "森に進む", 2)
	if err != nil {
		t.Fatalf("ExecuteChoiceCommand failed: %v", err)
	}

	// Assert - ログにコマンド実行が記録されているか確認
	logOutput := logBuffer.String()

	expectedLogEntries := []string{
		"新しいゲームブック作成",
		"テストゲーム",
		"パラグラフ更新", // プレースホルダー更新になる
		"物語の始まり",
		"選択肢追加",
		"森に進む",
	}

	for _, expectedEntry := range expectedLogEntries {
		if !strings.Contains(logOutput, expectedEntry) {
			t.Errorf("期待されるログエントリが見つかりません: %s\nActual output: %s", expectedEntry, logOutput)
		}
	}
}

func TestIntegration_AIMode_LogLevelEscalation(t *testing.T) {
	// Arrange
	originalAIMode := os.Getenv("GAMEBOOK_AI_DEV")
	originalLogLevel := os.Getenv("LOG_LEVEL")
	defer func() {
		os.Setenv("GAMEBOOK_AI_DEV", originalAIMode)
		os.Setenv("LOG_LEVEL", originalLogLevel)
	}()

	// AI開発モードを有効化
	os.Setenv("GAMEBOOK_AI_DEV", "true")
	os.Setenv("LOG_LEVEL", "INFO") // 通常レベル

	// Act
	logger, err := initializeApplicationLogger()
	if err != nil {
		t.Fatalf("Logger initialization failed: %v", err)
	}
	defer cleanupApplicationLogger()

	// Assert - AI開発モード時はDEBUGレベルに自動エスカレートされているか確認
	isAIMode := IsAIDevelopmentMode()
	if !isAIMode {
		t.Error("AI開発モードが有効になっていません")
	}

	// 環境変数がDEBUGに変更されているか確認
	currentLogLevel := os.Getenv("LOG_LEVEL")
	if currentLogLevel != "DEBUG" {
		t.Errorf("AI開発モード時のログレベルエスカレートが機能していません。期待値: DEBUG, 実際: %s", currentLogLevel)
	}

	if logger == nil {
		t.Fatal("Logger is nil")
	}
}

func TestIntegration_DynamicLogLevel_Change(t *testing.T) {
	// Arrange
	config := LoggingConfig{
		Level:      domain.LogLevelInfo,
		Format:     "text",
		OutputType: "console",
		SyncMode:   true,
	}

	controller, err := NewLoggingController(config)
	if err != nil {
		t.Fatalf("LoggingController作成に失敗: %v", err)
	}
	defer controller.Close()

	// 初期レベルを確認
	initialLevel := controller.GetLevel()
	if initialLevel != domain.LogLevelInfo {
		t.Errorf("初期レベルが正しくありません。期待値: %v, 実際: %v", domain.LogLevelInfo, initialLevel)
	}

	// Act - 動的にレベルを変更
	err = SetDynamicLogLevel(controller, domain.LogLevelDebug)
	if err != nil {
		t.Fatalf("SetDynamicLogLevel failed: %v", err)
	}

	// Assert - レベルが変更されているか確認
	newLevel := controller.GetLevel()
	if newLevel != domain.LogLevelDebug {
		t.Errorf("ログレベルが変更されていません。期待値: %v, 実際: %v", domain.LogLevelDebug, newLevel)
	}
}
