package main

import (
	"bytes"
	"strings"
	"testing"

	"github.com/zensai3805/iwakero-gamebook-mapping/internal/domain"
)

func TestLogUserOperation(t *testing.T) {
	// Arrange
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

	// Act
	LogUserOperation("test_operation", map[string]interface{}{
		"param1": "value1",
		"param2": 42,
	})

	// Assert
	logOutput := logBuffer.String()
	expectedContents := []string{
		"ユーザー操作: test_operation",
		"category=operation",
		"action=test_operation",
		"component=frameworks",
	}

	for _, expected := range expectedContents {
		if !strings.Contains(logOutput, expected) {
			t.Errorf("期待される内容が見つかりません: %s\nActual output: %s", expected, logOutput)
		}
	}
}

func TestLogValidationError(t *testing.T) {
	// Arrange
	logBuffer := &bytes.Buffer{}
	config := LoggingConfig{
		Level:      domain.LogLevelWarn,
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

	// Act
	LogValidationError("title", "", "空のタイトル", map[string]interface{}{
		"command": "new",
	})

	// Assert
	logOutput := logBuffer.String()
	expectedContents := []string{
		"入力値検証エラー: 空のタイトル",
		"category=validation",
		"action=field_validation",
		"field=title",
		"reason=空のタイトル",
		"command=new",
	}

	for _, expected := range expectedContents {
		if !strings.Contains(logOutput, expected) {
			t.Errorf("期待される内容が見つかりません: %s\nActual output: %s", expected, logOutput)
		}
	}
}

func TestLogCommandResult(t *testing.T) {
	// Arrange
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

	// Act - 成功ケース
	LogCommandResult("test_command", true, map[string]interface{}{
		"result": "success",
	})

	// Act - 失敗ケース
	LogCommandResult("test_command", false, map[string]interface{}{
		"error": "test error",
	})

	// Assert
	logOutput := logBuffer.String()
	
	// 成功ログの確認
	if !strings.Contains(logOutput, "コマンド実行成功") {
		t.Error("成功ログが記録されていません")
	}
	
	// 失敗ログの確認
	if !strings.Contains(logOutput, "コマンド実行失敗") {
		t.Error("失敗ログが記録されていません")
	}
	
	if !strings.Contains(logOutput, "test_command") {
		t.Error("コマンド名が記録されていません")
	}
}

func TestLogUIInteraction_AIDevelopmentMode(t *testing.T) {
	// Arrange
	logBuffer := &bytes.Buffer{}
	config := LoggingConfig{
		Level:      domain.LogLevelDebug,
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

	// AI開発モードのテスト
	originalAIMode := IsAIDevelopmentMode()
	defer func() {
		// テスト後に元の状態に戻す（実際の環境変数は変更しない）
	}()

	// Act
	LogUIInteraction("menu_selection", map[string]interface{}{
		"selected_option": "test_option",
		"options_count": 5,
		"selection_time_ms": 123.45,
	})

	// Assert
	logOutput := logBuffer.String()
	
	// DEBUG レベルなので AI開発モード時のみ出力される
	if IsAIDevelopmentMode() {
		if !strings.Contains(logOutput, "UI操作") {
			t.Error("UI操作ログが記録されていません")
		}
		
		if !strings.Contains(logOutput, "menu_selection") {
			t.Error("操作タイプが記録されていません")
		}
	} else {
		// AI開発モードでない場合はDEBUGログは出力されない
		if strings.Contains(logOutput, "UI操作") {
			t.Error("通常モードでUI操作ログが記録されています")
		}
	}
	
	// AI開発モードでない場合は最小限の情報のみ
	if !originalAIMode && strings.Contains(logOutput, "options_count") {
		t.Error("通常モードで詳細情報が記録されています")
	}
}