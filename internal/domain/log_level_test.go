package domain

import (
	"testing"
)

func TestLogLevel_String_Debug_ReturnsDebug(t *testing.T) {
	// Arrange
	level := LogLevelDebug

	// Act
	result := level.String()

	// Assert
	if result != "DEBUG" {
		t.Errorf("expected DEBUG, got %s", result)
	}
}

func TestLogLevel_String_Info_ReturnsInfo(t *testing.T) {
	// Arrange
	level := LogLevelInfo

	// Act
	result := level.String()

	// Assert
	if result != "INFO" {
		t.Errorf("expected INFO, got %s", result)
	}
}

func TestLogLevel_String_Warn_ReturnsWarn(t *testing.T) {
	// Arrange
	level := LogLevelWarn

	// Act
	result := level.String()

	// Assert
	if result != "WARN" {
		t.Errorf("expected WARN, got %s", result)
	}
}

func TestLogLevel_String_Error_ReturnsError(t *testing.T) {
	// Arrange
	level := LogLevelError

	// Act
	result := level.String()

	// Assert
	if result != "ERROR" {
		t.Errorf("expected ERROR, got %s", result)
	}
}

func TestLogLevel_String_Fatal_ReturnsFatal(t *testing.T) {
	// Arrange
	level := LogLevelFatal

	// Act
	result := level.String()

	// Assert
	if result != "FATAL" {
		t.Errorf("expected FATAL, got %s", result)
	}
}

func TestLogLevel_String_Unknown_ReturnsUnknown(t *testing.T) {
	// Arrange
	level := LogLevel(999)

	// Act
	result := level.String()

	// Assert
	if result != "UNKNOWN" {
		t.Errorf("expected UNKNOWN, got %s", result)
	}
}

func TestLogLevel_IsHigherThan_DebugLowerThanInfo_ReturnsTrue(t *testing.T) {
	// Arrange
	debug := LogLevelDebug
	info := LogLevelInfo

	// Act
	result := info.IsHigherThan(debug)

	// Assert
	if !result {
		t.Error("expected Info to be higher than Debug")
	}
}

func TestLogLevel_IsHigherThan_InfoLowerThanWarn_ReturnsTrue(t *testing.T) {
	// Arrange
	info := LogLevelInfo
	warn := LogLevelWarn

	// Act
	result := warn.IsHigherThan(info)

	// Assert
	if !result {
		t.Error("expected Warn to be higher than Info")
	}
}

func TestLogLevel_IsHigherThan_WarnLowerThanError_ReturnsTrue(t *testing.T) {
	// Arrange
	warn := LogLevelWarn
	err := LogLevelError

	// Act
	result := err.IsHigherThan(warn)

	// Assert
	if !result {
		t.Error("expected Error to be higher than Warn")
	}
}

func TestLogLevel_IsHigherThan_ErrorLowerThanFatal_ReturnsTrue(t *testing.T) {
	// Arrange
	err := LogLevelError
	fatal := LogLevelFatal

	// Act
	result := fatal.IsHigherThan(err)

	// Assert
	if !result {
		t.Error("expected Fatal to be higher than Error")
	}
}

func TestLogLevel_IsHigherThan_SameLevel_ReturnsFalse(t *testing.T) {
	// Arrange
	info1 := LogLevelInfo
	info2 := LogLevelInfo

	// Act
	result := info1.IsHigherThan(info2)

	// Assert
	if result {
		t.Error("expected same levels to return false")
	}
}

func TestLogLevel_IsHigherThan_LowerLevel_ReturnsFalse(t *testing.T) {
	// Arrange
	info := LogLevelInfo
	debug := LogLevelDebug

	// Act
	result := debug.IsHigherThan(info)

	// Assert
	if result {
		t.Error("expected Debug not to be higher than Info")
	}
}

func TestParseLogLevel_Debug_ReturnsDebugLevel(t *testing.T) {
	// Arrange
	input := "DEBUG"

	// Act
	result, err := ParseLogLevel(input)

	// Assert
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if result != LogLevelDebug {
		t.Errorf("expected LogLevelDebug, got %v", result)
	}
}

func TestParseLogLevel_Info_ReturnsInfoLevel(t *testing.T) {
	// Arrange
	input := "INFO"

	// Act
	result, err := ParseLogLevel(input)

	// Assert
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if result != LogLevelInfo {
		t.Errorf("expected LogLevelInfo, got %v", result)
	}
}

func TestParseLogLevel_Warn_ReturnsWarnLevel(t *testing.T) {
	// Arrange
	input := "WARN"

	// Act
	result, err := ParseLogLevel(input)

	// Assert
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if result != LogLevelWarn {
		t.Errorf("expected LogLevelWarn, got %v", result)
	}
}

func TestParseLogLevel_Error_ReturnsErrorLevel(t *testing.T) {
	// Arrange
	input := "ERROR"

	// Act
	result, err := ParseLogLevel(input)

	// Assert
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if result != LogLevelError {
		t.Errorf("expected LogLevelError, got %v", result)
	}
}

func TestParseLogLevel_Fatal_ReturnsFatalLevel(t *testing.T) {
	// Arrange
	input := "FATAL"

	// Act
	result, err := ParseLogLevel(input)

	// Assert
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if result != LogLevelFatal {
		t.Errorf("expected LogLevelFatal, got %v", result)
	}
}

func TestParseLogLevel_LowerCase_ReturnsLevel(t *testing.T) {
	// Arrange
	input := "info"

	// Act
	result, err := ParseLogLevel(input)

	// Assert
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if result != LogLevelInfo {
		t.Errorf("expected LogLevelInfo, got %v", result)
	}
}

func TestParseLogLevel_Invalid_ReturnsError(t *testing.T) {
	// Arrange
	input := "INVALID"

	// Act
	_, err := ParseLogLevel(input)

	// Assert
	if err == nil {
		t.Error("expected error for invalid log level")
	}
}

func TestParseLogLevel_Empty_ReturnsError(t *testing.T) {
	// Arrange
	input := ""

	// Act
	_, err := ParseLogLevel(input)

	// Assert
	if err == nil {
		t.Error("expected error for empty string")
	}
}
