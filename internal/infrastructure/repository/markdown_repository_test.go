package repository

import (
	"os"
	"path/filepath"
	"testing"
	"time"

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
	defer os.RemoveAll(tmpDir)

	repo := NewMarkdownRepository(tmpDir)

	t.Run("複数のゲームブック一覧", func(t *testing.T) {
		// Given
		gb1 := domain.NewGamebook("ブック1")
		p1 := domain.NewParagraph(1, "開始")
		_ = gb1.AddParagraph(p1)
		
		gb2 := domain.NewGamebook("ブック2")
		p2 := domain.NewParagraph(1, "開始")
		_ = gb2.AddParagraph(p2)
		
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

func TestMarkdownRepository_ValidateGamebook(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "gamebook_test")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	repo := NewMarkdownRepository(tmpDir)

	t.Run("正常系：有効なゲームブック", func(t *testing.T) {
		// Given
		gb := domain.NewGamebook("有効なブック")
		p1 := domain.NewParagraph(1, "開始")
		p1.AddChoice("進む", 2)
		_ = gb.AddParagraph(p1)

		// When & Then
		assert.NoError(t, repo.Save(gb))
	})

	t.Run("異常系：nilゲームブック", func(t *testing.T) {
		// When
		err := repo.Save(nil)

		// Then
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "データ検証エラー")
		assert.Contains(t, err.Error(), "nilです")
	})

	t.Run("異常系：空タイトル", func(t *testing.T) {
		// Given
		gb := domain.NewGamebook("")

		// When
		err := repo.Save(gb)

		// Then
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "タイトルが空です")
	})

	t.Run("異常系：不正タイトル文字", func(t *testing.T) {
		// Given
		gb := domain.NewGamebook("不正<文字>")

		// When
		err := repo.Save(gb)

		// Then
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "不正な文字が含まれています")
	})

	t.Run("異常系：パラグラフなし", func(t *testing.T) {
		// Given
		gb := domain.NewGamebook("空のブック")

		// When
		err := repo.Save(gb)

		// Then
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "パラグラフが1つも存在しません")
	})
}

func TestMarkdownRepository_LoadValidation(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "gamebook_test")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	repo := NewMarkdownRepository(tmpDir)

	t.Run("異常系：空タイトルで読み込み", func(t *testing.T) {
		// When
		_, err := repo.Load("")

		// Then
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "タイトルが空です")
	})

	t.Run("異常系：存在しないファイル", func(t *testing.T) {
		// When
		_, err := repo.Load("存在しない")

		// Then
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "が見つかりません")
	})

	t.Run("異常系：空ファイル", func(t *testing.T) {
		// Given
		filename := filepath.Join(tmpDir, "空ファイル.md")
		err := os.WriteFile(filename, []byte(""), 0644)
		require.NoError(t, err)

		// When
		_, err = repo.Load("空ファイル")

		// Then
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "ファイルが空です")
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
		p1 := domain.NewParagraph(1, "開始")
		_ = gb.AddParagraph(p1)
		err := repo.Save(gb)
		require.NoError(t, err)
		
		// ファイルが作成されたことを確認
		filePath := filepath.Join(tmpDir, "削除テスト.md")
		assert.FileExists(t, filePath)

		// When
		err = repo.Delete("削除テスト")

		// Then
		assert.NoError(t, err)
		
		// ファイルが削除されたか確認
		assert.NoFileExists(t, filePath)
	})
}

func TestMarkdownRepository_CorruptedDataRecovery(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "gamebook_test")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	repo := NewMarkdownRepository(tmpDir)

	t.Run("破損したMarkdownファイルの読み込み", func(t *testing.T) {
		// Given - 不正なMarkdownファイルを作成
		corruptedContent := `# 破損テスト

## パラグラフ 
- 概要：
- 選択肢：
  - [x] 不正な選択肢
  - [] 空の選択肢
`
		filename := filepath.Join(tmpDir, "破損テスト.md")
		err := os.WriteFile(filename, []byte(corruptedContent), 0644)
		require.NoError(t, err)

		// When
		_, err = repo.Load("破損テスト")

		// Then
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "読み込み後データ検証エラー")
	})

	t.Run("部分的に破損したファイルの読み込み", func(t *testing.T) {
		// Given - 一部データが欠損したMarkdownファイル
		partialContent := `# 部分破損テスト

## パラグラフ 1
- 概要：開始地点
- 選択肢：
  - [x] 進む → 999

## パラグラフ 2
- 概要：
`
		filename := filepath.Join(tmpDir, "部分破損テスト.md")
		err := os.WriteFile(filename, []byte(partialContent), 0644)
		require.NoError(t, err)

		// When
		_, err = repo.Load("部分破損テスト")

		// Then
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "パラグラフ 2 の説明が空です")
	})
}

func TestMarkdownRepository_AtomicSaveFailure(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "gamebook_test")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	repo := NewMarkdownRepository(tmpDir)

	t.Run("原子的保存の故障テスト", func(t *testing.T) {
		// Given - 既存のファイルを作成
		original := domain.NewGamebook("原子的保存テスト")
		p1 := domain.NewParagraph(1, "元の内容")
		_ = original.AddParagraph(p1)
		_ = repo.Save(original)

		// 元のファイルが存在することを確認
		filePath := filepath.Join(tmpDir, "原子的保存テスト.md")
		assert.FileExists(t, filePath)
		
		// 元のファイルの内容を読み込み
		originalContent, _ := os.ReadFile(filePath)

		// When - 無効なゲームブックで保存を試行（失敗するはず）
		invalidGamebook := domain.NewGamebook("原子的保存テスト")
		// パラグラフなしで保存を試行（検証エラーが発生するはず）
		err := repo.Save(invalidGamebook)

		// Then - 保存は失敗するが、元のファイルは保持される
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "パラグラフが1つも存在しません")
		
		// 元のファイルが変更されていないことを確認
		currentContent, _ := os.ReadFile(filePath)
		assert.Equal(t, originalContent, currentContent)
		
		// 一時ファイルが削除されていることを確認
		tempFilePath := filePath + ".tmp"
		assert.NoFileExists(t, tempFilePath)
	})
}

func TestMarkdownRepository_ConcurrentAccess(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "gamebook_test")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	repo := NewMarkdownRepository(tmpDir)

	t.Run("並行アクセス時の安全性", func(t *testing.T) {
		// Given
		gb := domain.NewGamebook("並行アクセステスト")
		p1 := domain.NewParagraph(1, "基本パラグラフ")
		_ = gb.AddParagraph(p1)
		_ = repo.Save(gb)

		// When - 複数のゴルーチンで同時に読み込み
		const numGoroutines = 10
		results := make(chan error, numGoroutines)
		
		for i := 0; i < numGoroutines; i++ {
			go func() {
				_, err := repo.Load("並行アクセステスト")
				results <- err
			}()
		}

		// Then - すべての読み込みが成功する
		for i := 0; i < numGoroutines; i++ {
			select {
			case err := <-results:
				assert.NoError(t, err)
			case <-time.After(time.Second):
				t.Fatal("タイムアウト")
			}
		}
	})
}