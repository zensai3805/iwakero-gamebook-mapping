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
	defer func() { _ = os.RemoveAll(tmpDir) }()

	repo := NewMarkdownRepository(tmpDir)

	t.Run("現在地の永続化テスト", func(t *testing.T) {
		// Given
		gb := domain.NewGamebook("現在地テスト")
		p1 := domain.NewParagraph(1, "開始地点")
		p2 := domain.NewParagraph(2, "目的地")
		_ = gb.AddParagraph(p1)
		_ = gb.AddParagraph(p2)

		// 現在地を設定
		_ = gb.MoveTo(2)

		// When
		err := repo.Save(gb)

		// Then
		assert.NoError(t, err)

		// Load後に現在地が復元されるか確認
		loadedGb, loadErr := repo.Load("現在地テスト")
		assert.NoError(t, loadErr)
		assert.NotNil(t, loadedGb.Current)
		assert.Equal(t, 2, loadedGb.Current.Number)
		assert.Equal(t, "目的地", loadedGb.Current.Description)
	})

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
	defer func() { _ = os.RemoveAll(tmpDir) }()

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

	t.Run("複数選択肢での選択状態読み込み", func(t *testing.T) {
		// Given - より複雑な選択肢構成
		original := domain.NewGamebook("複数選択肢テスト")
		p1 := domain.NewParagraph(1, "村の入り口")
		p1.AddChoice("北へ進む", 2)
		p1.AddChoice("南へ進む", 3)
		p1.AddChoice("東へ進む", 4)
		p1.AddChoice("西へ進む", 5)
		_ = p1.SelectChoice(1) // 南へ進むを選択
		_ = original.AddParagraph(p1)
		_ = repo.Save(original)

		// When
		loaded, err := repo.Load("複数選択肢テスト")

		// Then
		assert.NoError(t, err)
		p, exists := loaded.Paragraphs[1]
		assert.True(t, exists)
		assert.Equal(t, "村の入り口", p.Description)
		assert.Len(t, p.Choices, 4)

		// 選択状態の詳細検証
		assert.False(t, p.Choices[0].Selected, "北へ進むは選択されていない")
		assert.True(t, p.Choices[1].Selected, "南へ進むが選択されている")
		assert.False(t, p.Choices[2].Selected, "東へ進むは選択されていない")
		assert.False(t, p.Choices[3].Selected, "西へ進むは選択されていない")

		// 選択肢の内容検証
		assert.Equal(t, "北へ進む", p.Choices[0].Description)
		assert.Equal(t, "南へ進む", p.Choices[1].Description)
		assert.Equal(t, "東へ進む", p.Choices[2].Description)
		assert.Equal(t, "西へ進む", p.Choices[3].Description)

		// 遷移先の検証
		assert.Equal(t, 2, p.Choices[0].TargetNumber)
		assert.Equal(t, 3, p.Choices[1].TargetNumber)
		assert.Equal(t, 4, p.Choices[2].TargetNumber)
		assert.Equal(t, 5, p.Choices[3].TargetNumber)
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
	defer func() { _ = os.RemoveAll(tmpDir) }()

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
	defer func() { _ = os.RemoveAll(tmpDir) }()

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
