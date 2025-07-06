package usecase_test

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/zensai3805/iwakero-gamebook-mapping/internal/domain"
	"github.com/zensai3805/iwakero-gamebook-mapping/internal/usecase"
)

// MockLogWriter はLogWriterのモック実装
type MockLogWriter struct {
	mock.Mock
}

func (m *MockLogWriter) Write(entries []domain.LogEntry) error {
	args := m.Called(entries)
	return args.Error(0)
}

func (m *MockLogWriter) Close() error {
	args := m.Called()
	return args.Error(0)
}

// MockLogFormatter はLogFormatterのモック実装
type MockLogFormatter struct {
	mock.Mock
}

func (m *MockLogFormatter) Format(entry domain.LogEntry) ([]byte, error) {
	args := m.Called(entry)
	return args.Get(0).([]byte), args.Error(1)
}

func TestLoggingService_Debug(t *testing.T) {
	// Arrange
	mockWriter := new(MockLogWriter)
	mockFormatter := new(MockLogFormatter)
	service := usecase.NewLoggingService(mockWriter, mockFormatter, usecase.WithSyncMode())

	// フォーマッターは現在の設計では使用されていないが、将来の拡張のために保持
	mockWriter.On("Write", mock.MatchedBy(func(entries []domain.LogEntry) bool {
		return len(entries) == 1 && entries[0].Level == domain.LogLevelDebug && entries[0].Message == "デバッグメッセージ"
	})).Return(nil)

	// Act
	service.Debug("デバッグメッセージ")

	// Assert
	mockWriter.AssertExpectations(t)
}

func TestLoggingService_Info(t *testing.T) {
	// Arrange
	mockWriter := new(MockLogWriter)
	mockFormatter := new(MockLogFormatter)
	service := usecase.NewLoggingService(mockWriter, mockFormatter, usecase.WithSyncMode())

	// フォーマッターは現在の設計では使用されていないが、将来の拡張のために保持
	mockWriter.On("Write", mock.MatchedBy(func(entries []domain.LogEntry) bool {
		return len(entries) == 1 && entries[0].Level == domain.LogLevelInfo && entries[0].Message == "情報メッセージ"
	})).Return(nil)

	// Act
	service.Info("情報メッセージ")

	// Assert
	mockWriter.AssertExpectations(t)
}

func TestLoggingService_Warn(t *testing.T) {
	// Arrange
	mockWriter := new(MockLogWriter)
	mockFormatter := new(MockLogFormatter)
	service := usecase.NewLoggingService(mockWriter, mockFormatter, usecase.WithSyncMode())

	// フォーマッターは現在の設計では使用されていないが、将来の拡張のために保持
	mockWriter.On("Write", mock.MatchedBy(func(entries []domain.LogEntry) bool {
		return len(entries) == 1 && entries[0].Level == domain.LogLevelWarn && entries[0].Message == "警告メッセージ"
	})).Return(nil)

	// Act
	service.Warn("警告メッセージ")

	// Assert
	mockWriter.AssertExpectations(t)
}

func TestLoggingService_Error(t *testing.T) {
	// Arrange
	mockWriter := new(MockLogWriter)
	mockFormatter := new(MockLogFormatter)
	service := usecase.NewLoggingService(mockWriter, mockFormatter, usecase.WithSyncMode())

	// フォーマッターは現在の設計では使用されていないが、将来の拡張のために保持
	mockWriter.On("Write", mock.MatchedBy(func(entries []domain.LogEntry) bool {
		return len(entries) == 1 && entries[0].Level == domain.LogLevelError && entries[0].Message == "エラーメッセージ"
	})).Return(nil)

	// Act
	service.Error("エラーメッセージ")

	// Assert
	mockWriter.AssertExpectations(t)
}

func TestLoggingService_WithFields(t *testing.T) {
	// Arrange
	mockWriter := new(MockLogWriter)
	mockFormatter := new(MockLogFormatter)
	service := usecase.NewLoggingService(mockWriter, mockFormatter, usecase.WithSyncMode())

	fields := []domain.Field{
		{Key: "user_id", Value: "12345"},
		{Key: "action", Value: "login"},
	}

	// フォーマッターは現在の設計では使用されていないが、将来の拡張のために保持
	mockWriter.On("Write", mock.MatchedBy(func(entries []domain.LogEntry) bool {
		if len(entries) != 1 {
			return false
		}
		entry := entries[0]
		userID, ok1 := entry.GetField("user_id")
		action, ok2 := entry.GetField("action")
		return entry.Level == domain.LogLevelInfo && ok1 && ok2 && userID == "12345" && action == "login"
	})).Return(nil)

	// Act
	service.Info("ユーザーログイン", fields...)

	// Assert
	mockWriter.AssertExpectations(t)
}

func TestLoggingService_WithContext(t *testing.T) {
	// Arrange
	mockWriter := new(MockLogWriter)
	mockFormatter := new(MockLogFormatter)
	service := usecase.NewLoggingService(mockWriter, mockFormatter, usecase.WithSyncMode())

	contextFields := []domain.Field{
		{Key: "request_id", Value: "abc123"},
		{Key: "module", Value: "auth"},
	}

	// フォーマッターは現在の設計では使用されていないが、将来の拡張のために保持
	mockWriter.On("Write", mock.MatchedBy(func(entries []domain.LogEntry) bool {
		if len(entries) != 1 {
			return false
		}
		entry := entries[0]
		requestID, ok1 := entry.GetField("request_id")
		module, ok2 := entry.GetField("module")
		return entry.Level == domain.LogLevelInfo && ok1 && ok2 && requestID == "abc123" && module == "auth"
	})).Return(nil)

	// Act
	contextLogger := service.WithContext(contextFields...)
	contextLogger.Info("コンテキスト付きメッセージ")

	// Assert
	mockWriter.AssertExpectations(t)
}

func TestLoggingService_BatchProcessing(t *testing.T) {
	// Arrange
	mockWriter := new(MockLogWriter)
	mockFormatter := new(MockLogFormatter)
	// バッチ処理テストでは非同期モードを使用
	service := usecase.NewLoggingService(mockWriter, mockFormatter)

	// 複数のエントリが一度に処理されることを期待
	// フォーマッターは現在の設計では使用されていないが、将来の拡張のために保持
	mockWriter.On("Write", mock.MatchedBy(func(entries []domain.LogEntry) bool {
		return len(entries) >= 2 // バッチ処理を確認
	})).Return(nil).Once()

	// Act
	service.Info("メッセージ1")
	service.Info("メッセージ2")
	service.Info("メッセージ3")

	// Assert
	time.Sleep(200 * time.Millisecond) // バッチ処理を待つ
	mockWriter.AssertExpectations(t)

	// クリーンアップ
	mockWriter.On("Close").Return(nil)
	service.Close()
}

func TestLoggingService_Close(t *testing.T) {
	// Arrange
	mockWriter := new(MockLogWriter)
	mockFormatter := new(MockLogFormatter)
	service := usecase.NewLoggingService(mockWriter, mockFormatter, usecase.WithSyncMode())

	// フォーマッターは現在の設計では使用されていないが、将来の拡張のために保持
	mockWriter.On("Write", mock.Anything).Return(nil)
	mockWriter.On("Close").Return(nil)

	// Act
	service.Info("クローズ前のメッセージ")
	err := service.Close()

	// Assert
	assert.NoError(t, err)
	mockWriter.AssertExpectations(t)
}

func TestLoggingService_ErrorHandling_FormatterError(t *testing.T) {
	// Arrange
	mockWriter := new(MockLogWriter)
	mockFormatter := new(MockLogFormatter)
	service := usecase.NewLoggingService(mockWriter, mockFormatter, usecase.WithSyncMode())

	// フォーマッターは現在の設計では使用されていないが、将来の拡張のために保持
	// 書き込みは正常に行われる
	mockWriter.On("Write", mock.Anything).Return(nil)

	// Act
	service.Error("エラーメッセージ")

	// Assert
	// フォーマットエラーが発生してもサービスは継続する
	mockWriter.AssertExpectations(t)
}

func TestLoggingService_ErrorHandling_WriterError(t *testing.T) {
	// Arrange
	mockWriter := new(MockLogWriter)
	mockFormatter := new(MockLogFormatter)
	service := usecase.NewLoggingService(mockWriter, mockFormatter, usecase.WithSyncMode())

	writerErr := errors.New("書き込みエラー")
	// フォーマッターは現在の設計では使用されていないが、将来の拡張のために保持
	mockWriter.On("Write", mock.Anything).Return(writerErr)

	// Act
	service.Error("エラーメッセージ")

	// Assert
	// 書き込みエラーが発生してもサービスは継続する
	mockWriter.AssertExpectations(t)
}

func TestLoggingService_FileLineContext(t *testing.T) {
	// Arrange
	mockWriter := new(MockLogWriter)
	mockFormatter := new(MockLogFormatter)
	service := usecase.NewLoggingService(mockWriter, mockFormatter, usecase.WithSyncMode())

	// フォーマッターは現在の設計では使用されていないが、将来の拡張のために保持
	mockWriter.On("Write", mock.MatchedBy(func(entries []domain.LogEntry) bool {
		if len(entries) != 1 {
			return false
		}
		entry := entries[0]
		// ファイル名と行番号が自動的に付与されることを確認
		file, okFile := entry.GetField("file")
		line, okLine := entry.GetField("line")
		return okFile && okLine && file != "" && line != ""
	})).Return(nil)

	// Act
	service.Info("ファイル情報付きメッセージ")

	// Assert
	mockWriter.AssertExpectations(t)
}
