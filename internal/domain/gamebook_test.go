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

func TestGamebook_MoveToWithPathSelection(t *testing.T) {
	t.Run("正常系：定義済みパラグラフへの直接移動", func(t *testing.T) {
		// Given
		gb := NewGamebook("テストブック")
		p1 := NewParagraph(1, "開始")
		p2 := NewParagraph(2, "中間")
		p3 := NewParagraph(3, "目的地")
		
		// 選択肢を設定して経路を作成
		p1.AddChoice("中間へ", 2)
		p2.AddChoice("目的地へ", 3)
		
		_ = gb.AddParagraph(p1)
		_ = gb.AddParagraph(p2)
		_ = gb.AddParagraph(p3)
		
		// パラグラフ1を現在地に設定
		_ = gb.MoveTo(1)
		
		// When: パラグラフ3に直接移動
		err := gb.MoveToWithPathSelection(3)
		
		// Then
		assert.NoError(t, err)
		assert.Equal(t, 3, gb.Current.Number)
		assert.True(t, gb.Current.Visited) // 移動先は訪問済みに
		
		// 経路上の選択肢が選択済みになることを確認
		assert.True(t, p1.Choices[0].Selected) // 1->2の選択肢
		assert.True(t, p2.Choices[0].Selected) // 2->3の選択肢
		assert.True(t, p2.Visited) // 中間パラグラフも訪問済みに
	})
	
	t.Run("異常系：未定義パラグラフへの移動", func(t *testing.T) {
		// Given
		gb := NewGamebook("テストブック")
		p1 := NewParagraph(1, "開始")
		_ = gb.AddParagraph(p1)
		_ = gb.MoveTo(1)
		
		// When: 未定義パラグラフに移動を試行
		err := gb.MoveToWithPathSelection(999)
		
		// Then: 未定義でも移動可能（仕様）、しかし現在地は変更されない
		assert.Error(t, err)
		assert.Equal(t, 1, gb.Current.Number) // 現在地は変更されない
	})
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

// TestGamebook_AddChoiceToParagraph_WithPendingReference 保留参照を使った選択肢追加テスト
func TestGamebook_AddChoiceToParagraph_WithPendingReference(t *testing.T) {
	t.Run("正常系：未定義段落への選択肢追加を許可", func(t *testing.T) {
		// Given
		gb := NewGamebook("テストブック")
		p1 := NewParagraph(1, "開始")
		_ = gb.AddParagraph(p1)

		// When: 未定義段落23への選択肢を追加
		err := gb.AddChoiceToParagraph(1, "洞窟に入る", 23)

		// Then: エラーなく追加され、保留参照として記録される
		assert.NoError(t, err)
		assert.Len(t, p1.Choices, 1)
		assert.Equal(t, "洞窟に入る", p1.Choices[0].Description)
		assert.Equal(t, 23, p1.Choices[0].TargetNumber)

		// 保留参照が記録されていることを確認
		pendingTargets := gb.GetAllPendingTargets()
		assert.Contains(t, pendingTargets, 23)
	})

	t.Run("正常系：保留参照の自動解決", func(t *testing.T) {
		// Given
		gb := NewGamebook("テストブック")
		p1 := NewParagraph(1, "開始")
		_ = gb.AddParagraph(p1)

		// 未定義段落への選択肢を追加
		_ = gb.AddChoiceToParagraph(1, "洞窟に入る", 23)

		// When: 対象段落を追加
		p23 := NewParagraph(23, "暗い洞窟")
		err := gb.AddParagraph(p23)

		// Then: 保留参照が自動解決される
		assert.NoError(t, err)
		pendingTargets := gb.GetAllPendingTargets()
		assert.NotContains(t, pendingTargets, 23)
	})
}

// TestGamebook_MoveToWithGracefulHandling 未定義段落への移動の優雅な処理テスト
func TestGamebook_MoveToWithGracefulHandling(t *testing.T) {
	t.Run("警告付きで未定義段落への移動を記録", func(t *testing.T) {
		// Given
		gb := NewGamebook("テストブック")
		p1 := NewParagraph(1, "開始")
		p1.AddChoice("洞窟に入る", 23)
		_ = gb.AddParagraph(p1)
		_ = gb.MoveTo(1)

		// When: 未定義段落への移動を試行
		moveResult := gb.SelectChoiceAndMoveWithGracefulHandling(1, 0)

		// Then: 警告情報と共に処理が続行される
		assert.True(t, moveResult.Success)
		assert.True(t, moveResult.HasWarning)
		assert.Contains(t, moveResult.WarningMessage, "段落23")
		assert.Contains(t, moveResult.WarningMessage, "未定義")

		// 選択は記録される
		assert.True(t, p1.Choices[0].Selected)

		// 現在位置は変わらない（移動先が未定義のため）
		assert.Equal(t, p1, gb.Current)
	})
}
