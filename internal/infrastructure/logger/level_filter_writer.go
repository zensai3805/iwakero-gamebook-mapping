package logger

import (
	"github.com/zensai3805/iwakero-gamebook-mapping/internal/domain"
	"github.com/zensai3805/iwakero-gamebook-mapping/internal/usecase/interfaces"
)

// LevelFilterWriter は任意のライターにレベルフィルタリングを適用するラッパー
type LevelFilterWriter struct {
	writer      interfaces.LogWriter
	levelFilter domain.LogLevel
}

// NewLevelFilterWriter は新しいLevelFilterWriterを生成する
func NewLevelFilterWriter(writer interfaces.LogWriter, levelFilter domain.LogLevel) *LevelFilterWriter {
	return &LevelFilterWriter{
		writer:      writer,
		levelFilter: levelFilter,
	}
}

// Write は複数のログエントリをレベルフィルタリングして書き込む
func (w *LevelFilterWriter) Write(entries []domain.LogEntry) error {
	// レベルフィルタリングを適用
	filteredEntries := make([]domain.LogEntry, 0, len(entries))
	for _, entry := range entries {
		// レベルフィルタリング
		if !entry.Level.IsHigherThan(w.levelFilter) && entry.Level != w.levelFilter {
			continue
		}
		filteredEntries = append(filteredEntries, entry)
	}

	// フィルタリング後のエントリがない場合は何もしない
	if len(filteredEntries) == 0 {
		return nil
	}

	// 元のライターに書き込み
	return w.writer.Write(filteredEntries)
}

// Close は元のライターをクローズする
func (w *LevelFilterWriter) Close() error {
	return w.writer.Close()
}
