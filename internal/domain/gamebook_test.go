package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewGamebook(t *testing.T) {
	// Given
	title := "ファイティングファンタジー"

	// When
	gb := NewGamebook(title)

	// Then
	assert.Equal(t, title, gb.Title)
	assert.NotNil(t, gb.Paragraphs)
	assert.Empty(t, gb.Paragraphs)
	assert.Nil(t, gb.Current)
}

func TestGamebook_AddParagraph(t *testing.T) {
	t.Run("正常系：パラグラフを追加", func(t *testing.T) {
		// Given
		gb := NewGamebook("テストブック")
		p := NewParagraph(1, "開始")

		// When
		err := gb.AddParagraph(p)

		// Then
		assert.NoError(t, err)
		assert.Len(t, gb.Paragraphs, 1)
		assert.Equal(t, p, gb.Paragraphs[1])
	})

	t.Run("異常系：重複するパラグラフ番号", func(t *testing.T) {
		// Given
		gb := NewGamebook("テストブック")
		p1 := NewParagraph(1, "開始")
		p2 := NewParagraph(1, "別の開始")
		_ = gb.AddParagraph(p1)

		// When
		err := gb.AddParagraph(p2)

		// Then
		assert.Equal(t, ErrDuplicateParagraph, err)
		assert.Len(t, gb.Paragraphs, 1)
		// 元のパラグラフは変更されていない
		assert.Equal(t, "開始", gb.Paragraphs[1].Description)
	})

	t.Run("異常系：複数の重複検出", func(t *testing.T) {
		// Given
		gb := NewGamebook("テストブック")
		p1 := NewParagraph(1, "開始")
		p2 := NewParagraph(2, "中間")
		p3 := NewParagraph(1, "重複1")
		p4 := NewParagraph(2, "重複2")
		
		_ = gb.AddParagraph(p1)
		_ = gb.AddParagraph(p2)

		// When & Then
		err1 := gb.AddParagraph(p3)
		assert.Equal(t, ErrDuplicateParagraph, err1)
		
		err2 := gb.AddParagraph(p4)
		assert.Equal(t, ErrDuplicateParagraph, err2)
		
		// 元のパラグラフのみ存在
		assert.Len(t, gb.Paragraphs, 2)
		assert.Equal(t, "開始", gb.Paragraphs[1].Description)
		assert.Equal(t, "中間", gb.Paragraphs[2].Description)
	})
}

func TestGamebook_GetParagraph(t *testing.T) {
	t.Run("正常系：存在するパラグラフを取得", func(t *testing.T) {
		// Given
		gb := NewGamebook("テストブック")
		p := NewParagraph(1, "開始")
		_ = gb.AddParagraph(p)

		// When
		result, err := gb.GetParagraph(1)

		// Then
		assert.NoError(t, err)
		assert.Equal(t, p, result)
	})

	t.Run("異常系：存在しないパラグラフ", func(t *testing.T) {
		// Given
		gb := NewGamebook("テストブック")

		// When
		result, err := gb.GetParagraph(999)

		// Then
		assert.Equal(t, ErrParagraphNotFound, err)
		assert.Nil(t, result)
	})
}

func TestGamebook_MoveTo(t *testing.T) {
	t.Run("正常系：パラグラフに移動", func(t *testing.T) {
		// Given
		gb := NewGamebook("テストブック")
		p1 := NewParagraph(1, "開始")
		p2 := NewParagraph(2, "次の場所")
		_ = gb.AddParagraph(p1)
		_ = gb.AddParagraph(p2)

		// When
		err := gb.MoveTo(2)

		// Then
		assert.NoError(t, err)
		assert.Equal(t, p2, gb.Current)
		assert.True(t, p2.Visited)
	})

	t.Run("異常系：存在しないパラグラフへの移動", func(t *testing.T) {
		// Given
		gb := NewGamebook("テストブック")

		// When
		err := gb.MoveTo(999)

		// Then
		assert.Equal(t, ErrParagraphNotFound, err)
		assert.Nil(t, gb.Current)
	})
}

func TestGamebook_GetExplorationStats(t *testing.T) {
	// Given
	gb := NewGamebook("テストブック")
	
	// パラグラフを追加
	p1 := NewParagraph(1, "開始")
	p1.AddChoice("北へ", 2)
	p1.AddChoice("南へ", 3)
	
	p2 := NewParagraph(2, "北の部屋")
	p2.AddChoice("東へ", 4)
	
	p3 := NewParagraph(3, "南の部屋")
	
	_ = gb.AddParagraph(p1)
	_ = gb.AddParagraph(p2)
	_ = gb.AddParagraph(p3)
	
	// いくつか訪問と選択を行う
	_ = gb.MoveTo(1)
	_ = p1.SelectChoice(0) // 北へを選択
	_ = gb.MoveTo(2)

	// When
	stats := gb.GetExplorationStats()

	// Then
	assert.Equal(t, 3, stats.TotalParagraphs)
	assert.Equal(t, 2, stats.VisitedParagraphs) // 1と2を訪問
	assert.Equal(t, 3, stats.TotalChoices)      // 2+1+0
	assert.Equal(t, 1, stats.SelectedChoices)   // 北へのみ選択
}