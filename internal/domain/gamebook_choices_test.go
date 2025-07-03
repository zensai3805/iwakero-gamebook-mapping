package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGamebook_AddChoiceToParagraph(t *testing.T) {
	t.Run("パラグラフに選択肢を追加", func(t *testing.T) {
		// Given
		gb := NewGamebook("テストブック")
		p1 := NewParagraph(1, "開始地点")
		p2 := NewParagraph(2, "北の部屋")
		_ = gb.AddParagraph(p1)
		_ = gb.AddParagraph(p2)

		// When
		err := gb.AddChoiceToParagraph(1, "北へ進む", 2)

		// Then
		assert.NoError(t, err)

		paragraph, _ := gb.GetParagraph(1)
		assert.Len(t, paragraph.Choices, 1)
		assert.Equal(t, "北へ進む", paragraph.Choices[0].Description)
		assert.Equal(t, 2, paragraph.Choices[0].TargetNumber)
	})

	t.Run("存在しないパラグラフに選択肢を追加", func(t *testing.T) {
		// Given
		gb := NewGamebook("テストブック")

		// When
		err := gb.AddChoiceToParagraph(999, "どこかへ", 1)

		// Then
		assert.Equal(t, ErrParagraphNotFound, err)
	})
}

func TestGamebook_SelectChoice(t *testing.T) {
	t.Run("選択肢を選択", func(t *testing.T) {
		// Given
		gb := NewGamebook("テストブック")
		p1 := NewParagraph(1, "開始地点")
		p2 := NewParagraph(2, "北の部屋")
		_ = gb.AddParagraph(p1)
		_ = gb.AddParagraph(p2)
		_ = gb.AddChoiceToParagraph(1, "北へ進む", 2)
		_ = gb.AddChoiceToParagraph(1, "南へ進む", 3)

		// When
		err := gb.SelectChoiceAndMove(1, 0) // 最初の選択肢（北へ進む）を選択

		// Then
		assert.NoError(t, err)

		// 選択肢が選択済みになっている
		p1Updated, _ := gb.GetParagraph(1)
		assert.True(t, p1Updated.Choices[0].Selected)
		assert.False(t, p1Updated.Choices[1].Selected)

		// 現在位置が移動している
		assert.Equal(t, 2, gb.Current.Number)
		assert.True(t, gb.Current.Visited)
	})
}
