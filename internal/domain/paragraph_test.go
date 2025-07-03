package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewParagraph(t *testing.T) {
	// Given
	number := 1
	description := "冒険の始まり"

	// When
	p := NewParagraph(number, description)

	// Then
	assert.Equal(t, number, p.Number)
	assert.Equal(t, description, p.Description)
	assert.Empty(t, p.Choices)
	assert.False(t, p.Visited)
}

func TestParagraph_AddChoice(t *testing.T) {
	// Given
	p := NewParagraph(1, "冒険の始まり")

	// When
	p.AddChoice("北へ進む", 5)
	p.AddChoice("南へ進む", 10)

	// Then
	assert.Len(t, p.Choices, 2)
	assert.Equal(t, "北へ進む", p.Choices[0].Description)
	assert.Equal(t, 5, p.Choices[0].TargetNumber)
	assert.False(t, p.Choices[0].Selected)
	assert.Equal(t, "南へ進む", p.Choices[1].Description)
	assert.Equal(t, 10, p.Choices[1].TargetNumber)
	assert.False(t, p.Choices[1].Selected)
}

func TestParagraph_SelectChoice(t *testing.T) {
	t.Run("正常系：選択肢を選択", func(t *testing.T) {
		// Given
		p := NewParagraph(1, "冒険の始まり")
		p.AddChoice("北へ進む", 5)
		p.AddChoice("南へ進む", 10)

		// When
		err := p.SelectChoice(0)

		// Then
		assert.NoError(t, err)
		assert.True(t, p.Choices[0].Selected)
		assert.False(t, p.Choices[1].Selected)
		assert.True(t, p.Visited)
	})

	t.Run("異常系：無効なインデックス（負の数）", func(t *testing.T) {
		// Given
		p := NewParagraph(1, "冒険の始まり")
		p.AddChoice("北へ進む", 5)

		// When
		err := p.SelectChoice(-1)

		// Then
		assert.Equal(t, ErrInvalidChoiceIndex, err)
		assert.False(t, p.Visited)
	})

	t.Run("異常系：無効なインデックス（範囲外）", func(t *testing.T) {
		// Given
		p := NewParagraph(1, "冒険の始まり")
		p.AddChoice("北へ進む", 5)

		// When
		err := p.SelectChoice(1)

		// Then
		assert.Equal(t, ErrInvalidChoiceIndex, err)
		assert.False(t, p.Visited)
	})
}

func TestParagraph_GetUnselectedChoices(t *testing.T) {
	// Given
	p := NewParagraph(1, "冒険の始まり")
	p.AddChoice("北へ進む", 5)
	p.AddChoice("南へ進む", 10)
	p.AddChoice("東へ進む", 15)
	
	// 北へ進むを選択
	_ = p.SelectChoice(0)

	// When
	unselected := p.GetUnselectedChoices()

	// Then
	assert.Len(t, unselected, 2)
	assert.Equal(t, "南へ進む", unselected[0].Description)
	assert.Equal(t, "東へ進む", unselected[1].Description)
}