package domain

import (
	"testing"
	"time"
)

// MockLogger はテスト用のモックロガー実装
type MockLogger struct {
	entries []LogEntry
	fields  []Field
}

func (m *MockLogger) Debug(msg string, fields ...Field) {
	m.log(LogLevelDebug, msg, fields)
}

func (m *MockLogger) Info(msg string, fields ...Field) {
	m.log(LogLevelInfo, msg, fields)
}

func (m *MockLogger) Warn(msg string, fields ...Field) {
	m.log(LogLevelWarn, msg, fields)
}

func (m *MockLogger) Error(msg string, fields ...Field) {
	m.log(LogLevelError, msg, fields)
}

func (m *MockLogger) Fatal(msg string, fields ...Field) {
	m.log(LogLevelFatal, msg, fields)
}

func (m *MockLogger) WithContext(fields ...Field) Logger {
	newLogger := &MockLogger{
		entries: m.entries,
		fields:  make([]Field, len(m.fields)+len(fields)),
	}
	copy(newLogger.fields, m.fields)
	copy(newLogger.fields[len(m.fields):], fields)
	return newLogger
}

func (m *MockLogger) log(level LogLevel, msg string, fields []Field) {
	allFields := make([]Field, len(m.fields)+len(fields))
	copy(allFields, m.fields)
	copy(allFields[len(m.fields):], fields)

	entry := NewLogEntry(time.Now(), level, msg, allFields)
	m.entries = append(m.entries, entry)
}

func TestLogger_MockImplementation_Debug_LogsCorrectly(t *testing.T) {
	// Arrange
	logger := &MockLogger{}
	msg := "デバッグメッセージ"

	// Act
	logger.Debug(msg, Field{Key: "test", Value: "value"})

	// Assert
	if len(logger.entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(logger.entries))
	}
	entry := logger.entries[0]
	if entry.Level != LogLevelDebug {
		t.Errorf("expected Debug level, got %v", entry.Level)
	}
	if entry.Message != msg {
		t.Errorf("expected message %s, got %s", msg, entry.Message)
	}
}

func TestLogger_MockImplementation_Info_LogsCorrectly(t *testing.T) {
	// Arrange
	logger := &MockLogger{}
	msg := "情報メッセージ"

	// Act
	logger.Info(msg)

	// Assert
	if len(logger.entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(logger.entries))
	}
	if logger.entries[0].Level != LogLevelInfo {
		t.Errorf("expected Info level, got %v", logger.entries[0].Level)
	}
}

func TestLogger_MockImplementation_Warn_LogsCorrectly(t *testing.T) {
	// Arrange
	logger := &MockLogger{}
	msg := "警告メッセージ"

	// Act
	logger.Warn(msg)

	// Assert
	if len(logger.entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(logger.entries))
	}
	if logger.entries[0].Level != LogLevelWarn {
		t.Errorf("expected Warn level, got %v", logger.entries[0].Level)
	}
}

func TestLogger_MockImplementation_Error_LogsCorrectly(t *testing.T) {
	// Arrange
	logger := &MockLogger{}
	msg := "エラーメッセージ"

	// Act
	logger.Error(msg)

	// Assert
	if len(logger.entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(logger.entries))
	}
	if logger.entries[0].Level != LogLevelError {
		t.Errorf("expected Error level, got %v", logger.entries[0].Level)
	}
}

func TestLogger_MockImplementation_Fatal_LogsCorrectly(t *testing.T) {
	// Arrange
	logger := &MockLogger{}
	msg := "致命的エラーメッセージ"

	// Act
	logger.Fatal(msg)

	// Assert
	if len(logger.entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(logger.entries))
	}
	if logger.entries[0].Level != LogLevelFatal {
		t.Errorf("expected Fatal level, got %v", logger.entries[0].Level)
	}
}

func TestLogger_MockImplementation_WithContext_AddsFields(t *testing.T) {
	// Arrange
	logger := &MockLogger{}
	contextFields := []Field{
		{Key: "user", Value: "test-user"},
		{Key: "request_id", Value: "12345"},
	}

	// Act
	contextLogger := logger.WithContext(contextFields...)
	contextLogger.Info("コンテキスト付きメッセージ")

	// Assert
	mockContextLogger := contextLogger.(*MockLogger)
	if len(mockContextLogger.entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(mockContextLogger.entries))
	}

	entry := mockContextLogger.entries[0]
	userValue, hasUser := entry.GetField("user")
	if !hasUser || userValue != "test-user" {
		t.Errorf("expected user field with value 'test-user', got %v", userValue)
	}

	reqIDValue, hasReqID := entry.GetField("request_id")
	if !hasReqID || reqIDValue != "12345" {
		t.Errorf("expected request_id field with value '12345', got %v", reqIDValue)
	}
}

func TestLogger_MockImplementation_WithContext_ChainedContext(t *testing.T) {
	// Arrange
	logger := &MockLogger{}

	// Act
	logger1 := logger.WithContext(Field{Key: "app", Value: "gamebook"})
	logger2 := logger1.WithContext(Field{Key: "version", Value: "1.0.0"})
	logger2.Info("チェーンされたコンテキスト")

	// Assert
	mockLogger := logger2.(*MockLogger)
	if len(mockLogger.entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(mockLogger.entries))
	}

	entry := mockLogger.entries[0]
	appValue, hasApp := entry.GetField("app")
	if !hasApp || appValue != "gamebook" {
		t.Errorf("expected app field with value 'gamebook', got %v", appValue)
	}

	versionValue, hasVersion := entry.GetField("version")
	if !hasVersion || versionValue != "1.0.0" {
		t.Errorf("expected version field with value '1.0.0', got %v", versionValue)
	}
}

func TestLogger_InterfaceContract_AllMethodsRequired(t *testing.T) {
	// Arrange & Act & Assert
	// このテストはコンパイルが通ることで、インターフェースの契約が満たされていることを確認
	var _ Logger = &MockLogger{}
}
