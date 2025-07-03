package repository

import (
	"os"
	"path/filepath"
	"strings"
)

// FileSessionRepository はファイルベースでセッション状態を管理する
type FileSessionRepository struct {
	baseDir string
}

// NewFileSessionRepository は新しいFileSessionRepositoryを作成する
func NewFileSessionRepository(baseDir string) *FileSessionRepository {
	return &FileSessionRepository{
		baseDir: baseDir,
	}
}

// SaveCurrentGame は現在のゲームブックタイトルを保存する
func (r *FileSessionRepository) SaveCurrentGame(title string) error {
	if err := os.MkdirAll(r.baseDir, 0755); err != nil {
		return err
	}
	
	sessionFile := filepath.Join(r.baseDir, ".current_game")
	return os.WriteFile(sessionFile, []byte(title), 0644)
}

// GetCurrentGame は現在のゲームブックタイトルを取得する
func (r *FileSessionRepository) GetCurrentGame() (string, error) {
	sessionFile := filepath.Join(r.baseDir, ".current_game")
	content, err := os.ReadFile(sessionFile)
	if err != nil {
		if os.IsNotExist(err) {
			return "", nil // ファイルが存在しない場合は空文字を返す
		}
		return "", err
	}
	
	return strings.TrimSpace(string(content)), nil
}

// ClearCurrentGame は現在のゲーム設定をクリアする
func (r *FileSessionRepository) ClearCurrentGame() error {
	sessionFile := filepath.Join(r.baseDir, ".current_game")
	if err := os.Remove(sessionFile); err != nil && !os.IsNotExist(err) {
		return err
	}
	return nil
}