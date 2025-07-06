package usecase

import (
	"fmt"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/zensai3805/iwakero-gamebook-mapping/internal/domain"
	"github.com/zensai3805/iwakero-gamebook-mapping/internal/usecase/interfaces"
)

const (
	// バッファサイズ
	bufferSize = 100
	// バッチ処理の間隔
	batchInterval = 100 * time.Millisecond
	// バッチフラッシュの閾値
	batchFlushThreshold = 10
)

// LoggingService はロギングサービスの実装
type LoggingService struct {
	writer        interfaces.LogWriter
	formatter     interfaces.LogFormatter
	buffer        chan domain.LogEntry
	contextFields []domain.Field
	wg            sync.WaitGroup
	closeOnce     sync.Once
	closeChan     chan struct{}
	syncMode      bool          // テスト用の同期モード
	flushChan     chan struct{} // テスト用のフラッシュ通知チャネル
}

// LoggingServiceOption はLoggingServiceのオプション
type LoggingServiceOption func(*LoggingService)

// WithSyncMode は同期モードを有効にする（テスト用）
func WithSyncMode() LoggingServiceOption {
	return func(s *LoggingService) {
		s.syncMode = true
	}
}

// WithFlushNotification はフラッシュ通知チャネルを設定する（テスト用）
func WithFlushNotification() LoggingServiceOption {
	return func(s *LoggingService) {
		s.flushChan = make(chan struct{}, 10) // バッファ付きで複数回の通知に対応
	}
}

// NewLoggingService は新しいLoggingServiceを生成する
func NewLoggingService(writer interfaces.LogWriter, formatter interfaces.LogFormatter, opts ...LoggingServiceOption) *LoggingService {
	service := &LoggingService{
		writer:    writer,
		formatter: formatter,
		buffer:    make(chan domain.LogEntry, bufferSize),
		closeChan: make(chan struct{}),
	}

	// オプションを適用
	for _, opt := range opts {
		opt(service)
	}

	// 同期モードでない場合のみバックグラウンドでバッチ処理を開始
	if !service.syncMode {
		service.wg.Add(1)
		go service.processBatch()
	}

	return service
}

// Debug はデバッグレベルのログを出力する
func (s *LoggingService) Debug(message string, fields ...domain.Field) {
	s.log(domain.LogLevelDebug, message, fields...)
}

// Info は情報レベルのログを出力する
func (s *LoggingService) Info(message string, fields ...domain.Field) {
	s.log(domain.LogLevelInfo, message, fields...)
}

// Warn は警告レベルのログを出力する
func (s *LoggingService) Warn(message string, fields ...domain.Field) {
	s.log(domain.LogLevelWarn, message, fields...)
}

// Error はエラーレベルのログを出力する
func (s *LoggingService) Error(message string, fields ...domain.Field) {
	s.log(domain.LogLevelError, message, fields...)
}

// Fatal は致命的レベルのログを出力する
func (s *LoggingService) Fatal(message string, fields ...domain.Field) {
	s.log(domain.LogLevelFatal, message, fields...)
}

// WithContext はコンテキストフィールドを持つ新しいLoggerを返す
func (s *LoggingService) WithContext(fields ...domain.Field) domain.Logger {
	// 新しいコンテキストフィールドをコピー
	newContextFields := make([]domain.Field, len(s.contextFields)+len(fields))
	copy(newContextFields, s.contextFields)
	copy(newContextFields[len(s.contextFields):], fields)

	return &contextLogger{
		service:       s,
		contextFields: newContextFields,
	}
}

// contextLogger はコンテキストフィールドを持つLogger
type contextLogger struct {
	service       *LoggingService
	contextFields []domain.Field
}

func (c *contextLogger) Debug(message string, fields ...domain.Field) {
	allFields := make([]domain.Field, 0, len(c.contextFields)+len(fields))
	allFields = append(allFields, c.contextFields...)
	allFields = append(allFields, fields...)
	c.service.Debug(message, allFields...)
}

func (c *contextLogger) Info(message string, fields ...domain.Field) {
	allFields := make([]domain.Field, 0, len(c.contextFields)+len(fields))
	allFields = append(allFields, c.contextFields...)
	allFields = append(allFields, fields...)
	c.service.Info(message, allFields...)
}

func (c *contextLogger) Warn(message string, fields ...domain.Field) {
	allFields := make([]domain.Field, 0, len(c.contextFields)+len(fields))
	allFields = append(allFields, c.contextFields...)
	allFields = append(allFields, fields...)
	c.service.Warn(message, allFields...)
}

func (c *contextLogger) Error(message string, fields ...domain.Field) {
	allFields := make([]domain.Field, 0, len(c.contextFields)+len(fields))
	allFields = append(allFields, c.contextFields...)
	allFields = append(allFields, fields...)
	c.service.Error(message, allFields...)
}

func (c *contextLogger) Fatal(message string, fields ...domain.Field) {
	allFields := make([]domain.Field, 0, len(c.contextFields)+len(fields))
	allFields = append(allFields, c.contextFields...)
	allFields = append(allFields, fields...)
	c.service.Fatal(message, allFields...)
}

func (c *contextLogger) WithContext(fields ...domain.Field) domain.Logger {
	return c.service.WithContext(append(c.contextFields, fields...)...)
}

// Close はリソースをクリーンアップする
func (s *LoggingService) Close() error {
	var closeErr error
	s.closeOnce.Do(func() {
		// 同期モードでない場合のみバッチ処理を停止
		if !s.syncMode {
			// バッチ処理を停止
			close(s.closeChan)
			// バッファに残っているエントリを処理
			s.wg.Wait()
		}
		// フラッシュ通知チャネルをクローズ
		if s.flushChan != nil {
			close(s.flushChan)
		}
		// Writerをクローズ
		closeErr = s.writer.Close()
	})
	return closeErr
}

// WaitForFlush はバッチ処理のフラッシュを待つ（テスト用）
func (s *LoggingService) WaitForFlush() {
	if s.flushChan != nil {
		<-s.flushChan
	}
}

// log は内部的なログ記録処理
func (s *LoggingService) log(level domain.LogLevel, message string, fields ...domain.Field) {
	// ファイル名と行番号を取得
	_, file, line, _ := runtime.Caller(2)
	// ファイル名を短縮
	if idx := strings.LastIndex(file, "/"); idx != -1 {
		file = file[idx+1:]
	}

	// フィールドを結合
	allFields := make([]domain.Field, 0, len(s.contextFields)+len(fields)+2)
	allFields = append(allFields, s.contextFields...)
	allFields = append(allFields, fields...)
	allFields = append(allFields, domain.Field{Key: "file", Value: file})
	allFields = append(allFields, domain.Field{Key: "line", Value: fmt.Sprintf("%d", line)})

	entry := domain.NewLogEntry(time.Now(), level, message, allFields)

	// 同期モードの場合は直接処理
	if s.syncMode {
		s.writeBatch([]domain.LogEntry{entry})
		return
	}

	// バッファに送信（ノンブロッキング）
	select {
	case s.buffer <- entry:
	default:
		// バッファが満杯の場合はドロップ（本番環境では別の戦略を検討）
	}
}

// processBatch はバッチ処理を実行する
func (s *LoggingService) processBatch() {
	defer s.wg.Done()

	ticker := time.NewTicker(batchInterval)
	defer ticker.Stop()

	batch := make([]domain.LogEntry, 0, bufferSize)

	for {
		select {
		case entry := <-s.buffer:
			batch = append(batch, entry)

			// バッファが一定サイズに達したら即座に処理
			if len(batch) >= batchFlushThreshold {
				s.writeBatch(batch)
				batch = batch[:0]
			}

		case <-ticker.C:
			// 定期的にバッチを処理
			if len(batch) > 0 {
				s.writeBatch(batch)
				batch = batch[:0]
			}

		case <-s.closeChan:
			// 終了時は残りのエントリを全て処理
			for {
				select {
				case entry := <-s.buffer:
					batch = append(batch, entry)
				default:
					if len(batch) > 0 {
						s.writeBatch(batch)
					}
					return
				}
			}
		}
	}
}

// writeBatch はバッチを書き込む
func (s *LoggingService) writeBatch(batch []domain.LogEntry) {
	// 注: 現在の設計では、LogFormatterはLogWriterの実装側で使用されることを想定している
	// そのため、ここではLogEntryをそのままWriterに渡している
	// 将来的には、フォーマット済みのデータを渡すインターフェースへの変更を検討する

	// エラーが発生してもログシステム自体は継続する必要があるため、
	// エラーは標準エラー出力に記録するのみ
	if writeErr := s.writer.Write(batch); writeErr != nil {
		// 書き込みエラーを標準エラー出力に記録
		// ログシステム自体のエラーで無限ループを避けるため、
		// ここでは単純なfmt.Fprintfを使用
		fmt.Fprintf(os.Stderr, "LoggingService: 書き込みエラー: %v\n", writeErr)
	}

	// テスト用のフラッシュ通知
	if s.flushChan != nil {
		select {
		case s.flushChan <- struct{}{}:
		default:
			// バッファが満杯の場合は通知をスキップ
		}
	}
}
