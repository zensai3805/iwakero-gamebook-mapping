package repository

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFileSessionRepository(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "session_test")
	require.NoError(t, err)
	defer func() { _ = os.RemoveAll(tmpDir) }()

	repo := NewFileSessionRepository(tmpDir)

	t.Run("現在のゲームの保存と取得", func(t *testing.T) {
		// Given
		title := "テストゲーム"

		// When - 保存
		err := repo.SaveCurrentGame(title)

		// Then
		assert.NoError(t, err)

		// When - 取得
		retrieved, err := repo.GetCurrentGame()

		// Then
		assert.NoError(t, err)
		assert.Equal(t, title, retrieved)
	})

	t.Run("セッションファイルが存在しない場合", func(t *testing.T) {
		// Given - 新しいディレクトリ
		tmpDir2, err := os.MkdirTemp("", "session_test2")
		require.NoError(t, err)
		defer func() { _ = os.RemoveAll(tmpDir2) }()

		repo2 := NewFileSessionRepository(tmpDir2)

		// When
		retrieved, err := repo2.GetCurrentGame()

		// Then
		assert.NoError(t, err)
		assert.Empty(t, retrieved)
	})

	t.Run("現在のゲーム設定のクリア", func(t *testing.T) {
		// Given
		title := "削除テスト"
		_ = repo.SaveCurrentGame(title)

		// When
		err := repo.ClearCurrentGame()

		// Then
		assert.NoError(t, err)

		// 取得しても空になっていることを確認
		retrieved, err := repo.GetCurrentGame()
		assert.NoError(t, err)
		assert.Empty(t, retrieved)
	})
}
