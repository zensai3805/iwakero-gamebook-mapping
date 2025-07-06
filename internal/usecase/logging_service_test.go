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

	// expectedEntry は使用しないため削除
	expectedOutput := []byte(`{"level":"DEBUG","message":"デバッグメッセージ"}`)

	mockFormatter.On("Format", mock.MatchedBy(func(entry domain.LogEntry) bool {
		return entry.Level == domain.LogLevelDebug && entry.Message == "デバッグメッセージ"
	})).Return(expectedOutput, nil)
	mockWriter.On("Write", mock.MatchedBy(func(entries []domain.LogEntry) bool {
		return len(entries) == 1 && entries[0].Level == domain.LogLevelDebug
	})).Return(nil)

	// Act
	service.Debug("デバッグメッセージ")

	// Assert
	mockFormatter.AssertExpectations(t)
	mockWriter.AssertExpectations(t)
}

func TestLoggingService_Info(t *testing.T) {
	// Arrange
	mockWriter := new(MockLogWriter)
	mockFormatter := new(MockLogFormatter)
	service := usecase.NewLoggingService(mockWriter, mockFormatter, usecase.WithSyncMode())

	expectedOutput := []byte(`{"level":"INFO","message":"情報メッセージ"}`)

	mockFormatter.On("Format", mock.MatchedBy(func(entry domain.LogEntry) bool {
		return entry.Level == domain.LogLevelInfo && entry.Message == "情報メッセージ"
	})).Return(expectedOutput, nil)
	mockWriter.On("Write", mock.MatchedBy(func(entries []domain.LogEntry) bool {
		return len(entries) == 1 && entries[0].Level == domain.LogLevelInfo
	})).Return(nil)

	// Act
	service.Info("情報メッセージ")

	// Assert
	mockFormatter.AssertExpectations(t)
	mockWriter.AssertExpectations(t)
}

func TestLoggingService_Warn(t *testing.T) {
	// Arrange
	mockWriter := new(MockLogWriter)
	mockFormatter := new(MockLogFormatter)
	service := usecase.NewLoggingService(mockWriter, mockFormatter, usecase.WithSyncMode())

	expectedOutput := []byte(`{"level":"WARN","message":"警告メッセージ"}`)

	mockFormatter.On("Format", mock.MatchedBy(func(entry domain.LogEntry) bool {
		return entry.Level == domain.LogLevelWarn && entry.Message == "警告メッセージ"
	})).Return(expectedOutput, nil)
	mockWriter.On("Write", mock.MatchedBy(func(entries []domain.LogEntry) bool {
		return len(entries) == 1 && entries[0].Level == domain.LogLevelWarn
	})).Return(nil)

	// Act
	service.Warn("警告メッセージ")

	// Assert
	mockFormatter.AssertExpectations(t)
	mockWriter.AssertExpectations(t)
}

func TestLoggingService_Error(t *testing.T) {
	// Arrange
	mockWriter := new(MockLogWriter)
	mockFormatter := new(MockLogFormatter)
	service := usecase.NewLoggingService(mockWriter, mockFormatter, usecase.WithSyncMode())

	expectedOutput := []byte(`{"level":"ERROR","message":"エラーメッセージ"}`)

	mockFormatter.On("Format", mock.MatchedBy(func(entry domain.LogEntry) bool {
		return entry.Level == domain.LogLevelError && entry.Message == "エラーメッセージ"
	})).Return(expectedOutput, nil)
	mockWriter.On("Write", mock.MatchedBy(func(entries []domain.LogEntry) bool {
		return len(entries) == 1 && entries[0].Level == domain.LogLevelError
	})).Return(nil)

	// Act
	service.Error("エラーメッセージ")

	// Assert
	mockFormatter.AssertExpectations(t)
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

	mockFormatter.On("Format", mock.MatchedBy(func(entry domain.LogEntry) bool {
		return entry.Level == domain.LogLevelInfo &&
			func() bool {
				userID, ok1 := entry.GetField("user_id")
				action, ok2 := entry.GetField("action")
				return ok1 && ok2 && userID == "12345" && action == "login"
			}()
	})).Return([]byte(`{"level":"INFO","message":"ユーザーログイン","fields":{"user_id":"12345","action":"login"}}`), nil)
	mockWriter.On("Write", mock.Anything).Return(nil)

	// Act
	service.Info("ユーザーログイン", fields...)

	// Assert
	mockFormatter.AssertExpectations(t)
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

	mockFormatter.On("Format", mock.MatchedBy(func(entry domain.LogEntry) bool {
		return entry.Level == domain.LogLevelInfo &&
			func() bool {
				requestID, ok1 := entry.GetField("request_id")
				module, ok2 := entry.GetField("module")
				return ok1 && ok2 && requestID == "abc123" && module == "auth"
			}()
	})).Return([]byte(`{}`), nil)
	mockWriter.On("Write", mock.Anything).Return(nil)

	// Act
	contextLogger := service.WithContext(contextFields...)
	contextLogger.Info("コンテキスト付きメッセージ")

	// Assert
	mockFormatter.AssertExpectations(t)
	mockWriter.AssertExpectations(t)
}

func TestLoggingService_BatchProcessing(t *testing.T) {
	// Arrange
	mockWriter := new(MockLogWriter)
	mockFormatter := new(MockLogFormatter)
	// バッチ処理テストでは非同期モードを使用
	service := usecase.NewLoggingService(mockWriter, mockFormatter)

	// 複数のエントリが一度に処理されることを期待
	mockFormatter.On("Format", mock.Anything).Return([]byte(`{}`), nil)
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

	mockFormatter.On("Format", mock.Anything).Return([]byte(`{}`), nil)
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

	formatterErr := errors.New("フォーマットエラー")
	mockFormatter.On("Format", mock.Anything).Return([]byte{}, formatterErr)
	// フォーマットエラーでも書き込みは行われる
	mockWriter.On("Write", mock.Anything).Return(nil)

	// Act
	service.Error("エラーメッセージ")

	// Assert
	// フォーマットエラーが発生してもサービスは継続する
	mockFormatter.AssertExpectations(t)
	mockWriter.AssertExpectations(t)
}

func TestLoggingService_ErrorHandling_WriterError(t *testing.T) {
	// Arrange
	mockWriter := new(MockLogWriter)
	mockFormatter := new(MockLogFormatter)
	service := usecase.NewLoggingService(mockWriter, mockFormatter, usecase.WithSyncMode())

	writerErr := errors.New("書き込みエラー")
	mockFormatter.On("Format", mock.Anything).Return([]byte(`{}`), nil)
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

	mockFormatter.On("Format", mock.MatchedBy(func(entry domain.LogEntry) bool {
		// ファイル名と行番号が自動的に付与されることを確認
		file, okFile := entry.GetField("file")
		line, okLine := entry.GetField("line")
		return okFile && okLine && file != "" && line != ""
	})).Return([]byte(`{}`), nil)
	mockWriter.On("Write", mock.Anything).Return(nil)

	// Act
	service.Info("ファイル情報付きメッセージ")

	// Assert
	mockFormatter.AssertExpectations(t)
}
