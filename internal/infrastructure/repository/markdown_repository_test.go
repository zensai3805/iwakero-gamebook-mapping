package repository

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/iwapc/iwakero-gamebook-mapping/internal/domain"
)

func TestMarkdownRepository_Save(t *testing.T) {
	// テスト用の一時ディレクトリを作成
	tmpDir, err := os.MkdirTemp("", "gamebook_test")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	repo := NewMarkdownRepository(tmpDir)

	t.Run("新規ゲームブックの保存", func(t *testing.T) {
		// Given
		gb := domain.NewGamebook("テストブック")
		p1 := domain.NewParagraph(1, "開始地点")
		p1.AddChoice("北へ進む", 2)
		p1.AddChoice("南へ進む", 3)
		_ = gb.AddParagraph(p1)

		// When
		err := repo.Save(gb)

		// Then
		assert.NoError(t, err)
		
		// ファイルが作成されたか確認
		filePath := filepath.Join(tmpDir, "テストブック.md")
		assert.FileExists(t, filePath)
	})

	t.Run("既存ゲームブックの更新", func(t *testing.T) {
		// Given
		gb := domain.NewGamebook("更新テスト")
		p1 := domain.NewParagraph(1, "最初の内容")
		_ = gb.AddParagraph(p1)
		_ = repo.Save(gb)

		// 内容を更新
		p2 := domain.NewParagraph(2, "追加内容")
		_ = gb.AddParagraph(p2)

		// When
		err := repo.Save(gb)

		// Then
		assert.NoError(t, err)
	})
}

func TestMarkdownRepository_Load(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "gamebook_test")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	repo := NewMarkdownRepository(tmpDir)

	t.Run("保存したゲームブックの読み込み", func(t *testing.T) {
		// Given - 先にゲームブックを保存
		original := domain.NewGamebook("読み込みテスト")
		p1 := domain.NewParagraph(1, "開始地点")
		p1.AddChoice("北へ進む", 2)
		_ = p1.SelectChoice(0) // 北へ進むを選択
		_ = original.AddParagraph(p1)
		_ = repo.Save(original)

		// When
		loaded, err := repo.Load("読み込みテスト")

		// Then
		assert.NoError(t, err)
		assert.Equal(t, "読み込みテスト", loaded.Title)
		assert.Len(t, loaded.Paragraphs, 1)
		
		p, exists := loaded.Paragraphs[1]
		assert.True(t, exists)
		assert.Equal(t, "開始地点", p.Description)
		assert.Len(t, p.Choices, 1)
		if len(p.Choices) > 0 {
			assert.True(t, p.Choices[0].Selected)
		}
	})

	t.Run("存在しないゲームブックの読み込み", func(t *testing.T) {
		// When
		_, err := repo.Load("存在しない")

		// Then
		assert.Error(t, err)
	})
}

func TestMarkdownRepository_List(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "gamebook_test")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	repo := NewMarkdownRepository(tmpDir)

	t.Run("複数のゲームブック一覧", func(t *testing.T) {
		// Given
		gb1 := domain.NewGamebook("ブック1")
		gb2 := domain.NewGamebook("ブック2")
		_ = repo.Save(gb1)
		_ = repo.Save(gb2)

		// When
		titles, err := repo.List()

		// Then
		assert.NoError(t, err)
		assert.Len(t, titles, 2)
		assert.Contains(t, titles, "ブック1")
		assert.Contains(t, titles, "ブック2")
	})
}

func TestMarkdownRepository_Delete(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "gamebook_test")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	repo := NewMarkdownRepository(tmpDir)

	t.Run("ゲームブックの削除", func(t *testing.T) {
		// Given
		gb := domain.NewGamebook("削除テスト")
		_ = repo.Save(gb)

		// When
		err := repo.Delete("削除テスト")

		// Then
		assert.NoError(t, err)
		
		// ファイルが削除されたか確認
		filePath := filepath.Join(tmpDir, "削除テスト.md")
		assert.NoFileExists(t, filePath)
	})
}