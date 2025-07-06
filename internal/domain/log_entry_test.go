package domain

import (
	"encoding/json"
	"testing"
	"time"
)

func TestNewLogEntry_ValidInput_ReturnsLogEntry(t *testing.T) {
	// Arrange
	timestamp := time.Now()
	level := LogLevelInfo
	message := "テストメッセージ"
	fields := []Field{
		{Key: "user", Value: "test-user"},
		{Key: "action", Value: "login"},
	}

	// Act
	entry := NewLogEntry(timestamp, level, message, fields)

	// Assert
	if entry.Timestamp != timestamp {
		t.Errorf("expected timestamp %v, got %v", timestamp, entry.Timestamp)
	}
	if entry.Level != level {
		t.Errorf("expected level %v, got %v", level, entry.Level)
	}
	if entry.Message != message {
		t.Errorf("expected message %s, got %s", message, entry.Message)
	}
	if len(entry.Fields) != len(fields) {
		t.Errorf("expected %d fields, got %d", len(fields), len(entry.Fields))
	}
}

func TestNewLogEntry_EmptyFields_ReturnsLogEntryWithEmptyFields(t *testing.T) {
	// Arrange
	timestamp := time.Now()
	level := LogLevelWarn
	message := "警告メッセージ"
	var fields []Field

	// Act
	entry := NewLogEntry(timestamp, level, message, fields)

	// Assert
	if len(entry.Fields) != 0 {
		t.Errorf("expected 0 fields, got %d", len(entry.Fields))
	}
}

func TestLogEntry_ToJSON_ValidEntry_ReturnsJSONString(t *testing.T) {
	// Arrange
	timestamp := time.Date(2025, 1, 6, 10, 30, 0, 0, time.UTC)
	entry := NewLogEntry(timestamp, LogLevelInfo, "テストメッセージ", []Field{
		{Key: "user", Value: "test-user"},
		{Key: "count", Value: 42},
	})

	// Act
	jsonStr, err := entry.ToJSON()

	// Assert
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal([]byte(jsonStr), &result); err != nil {
		t.Fatalf("failed to unmarshal JSON: %v", err)
	}

	if result["timestamp"] != timestamp.Format(time.RFC3339Nano) {
		t.Errorf("expected timestamp %s, got %v", timestamp.Format(time.RFC3339Nano), result["timestamp"])
	}
	if result["level"] != "INFO" {
		t.Errorf("expected level INFO, got %v", result["level"])
	}
	if result["message"] != "テストメッセージ" {
		t.Errorf("expected message 'テストメッセージ', got %v", result["message"])
	}
	if result["user"] != "test-user" {
		t.Errorf("expected user 'test-user', got %v", result["user"])
	}
	if result["count"] != float64(42) {
		t.Errorf("expected count 42, got %v", result["count"])
	}
}

func TestLogEntry_ToJSON_EmptyFields_ReturnsJSONWithoutFields(t *testing.T) {
	// Arrange
	timestamp := time.Now()
	entry := NewLogEntry(timestamp, LogLevelError, "エラーメッセージ", nil)

	// Act
	jsonStr, err := entry.ToJSON()

	// Assert
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal([]byte(jsonStr), &result); err != nil {
		t.Fatalf("failed to unmarshal JSON: %v", err)
	}

	if result["level"] != "ERROR" {
		t.Errorf("expected level ERROR, got %v", result["level"])
	}
	if result["message"] != "エラーメッセージ" {
		t.Errorf("expected message 'エラーメッセージ', got %v", result["message"])
	}
}

func TestLogEntry_WithField_AddsNewField_ReturnsNewEntry(t *testing.T) {
	// Arrange
	timestamp := time.Now()
	entry := NewLogEntry(timestamp, LogLevelDebug, "デバッグメッセージ", []Field{
		{Key: "original", Value: "value"},
	})

	// Act
	newEntry := entry.WithField("new", "field")

	// Assert
	if len(entry.Fields) != 1 {
		t.Errorf("original entry should have 1 field, got %d", len(entry.Fields))
	}
	if len(newEntry.Fields) != 2 {
		t.Errorf("new entry should have 2 fields, got %d", len(newEntry.Fields))
	}
	if newEntry.Fields[1].Key != "new" || newEntry.Fields[1].Value != "field" {
		t.Errorf("expected new field with key 'new' and value 'field', got %+v", newEntry.Fields[1])
	}
}

func TestLogEntry_WithField_OriginalUnchanged_ReturnsNewEntry(t *testing.T) {
	// Arrange
	timestamp := time.Now()
	original := NewLogEntry(timestamp, LogLevelInfo, "メッセージ", nil)
	originalFieldCount := len(original.Fields)

	// Act
	modified := original.WithField("key", "value")

	// Assert
	if len(original.Fields) != originalFieldCount {
		t.Error("original entry was modified")
	}
	if len(modified.Fields) != 1 {
		t.Errorf("expected 1 field in modified entry, got %d", len(modified.Fields))
	}
}

func TestLogEntry_WithFields_AddsMultipleFields_ReturnsNewEntry(t *testing.T) {
	// Arrange
	timestamp := time.Now()
	entry := NewLogEntry(timestamp, LogLevelWarn, "警告", nil)
	newFields := []Field{
		{Key: "field1", Value: "value1"},
		{Key: "field2", Value: 123},
		{Key: "field3", Value: true},
	}

	// Act
	newEntry := entry.WithFields(newFields)

	// Assert
	if len(entry.Fields) != 0 {
		t.Errorf("original entry should have 0 fields, got %d", len(entry.Fields))
	}
	if len(newEntry.Fields) != 3 {
		t.Errorf("new entry should have 3 fields, got %d", len(newEntry.Fields))
	}
}

func TestLogEntry_GetField_ExistingField_ReturnsValueAndTrue(t *testing.T) {
	// Arrange
	entry := NewLogEntry(time.Now(), LogLevelInfo, "メッセージ", []Field{
		{Key: "user", Value: "test-user"},
		{Key: "id", Value: 12345},
	})

	// Act
	value, exists := entry.GetField("user")

	// Assert
	if !exists {
		t.Error("expected field to exist")
	}
	if value != "test-user" {
		t.Errorf("expected value 'test-user', got %v", value)
	}
}

func TestLogEntry_GetField_NonExistingField_ReturnsNilAndFalse(t *testing.T) {
	// Arrange
	entry := NewLogEntry(time.Now(), LogLevelInfo, "メッセージ", []Field{
		{Key: "user", Value: "test-user"},
	})

	// Act
	value, exists := entry.GetField("nonexistent")

	// Assert
	if exists {
		t.Error("expected field not to exist")
	}
	if value != nil {
		t.Errorf("expected nil value, got %v", value)
	}
}
