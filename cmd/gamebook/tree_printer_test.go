package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zensai3805/iwakero-gamebook-mapping/internal/domain"
)

func TestTreePrinter_WithNavigationHistory_HighlightsSelectedPath(t *testing.T) {
	printer := NewTreePrinter()
	gamebook := createTestGamebookWithHistory()
	converter := NewDataConverter()
	
	data, err := converter.ConvertToVisualizationData(gamebook)
	assert.NoError(t, err)
	
	err = printer.Initialize(data)
	assert.NoError(t, err)
	
	result, err := printer.Render()
	assert.NoError(t, err)
	
	// 選択された経路（B→D）がハイライトされることを確認
	assert.Contains(t, result, "[✓] B → 12")
	assert.Contains(t, result, "[✓] D → 11")
	
	// 選択されていない経路（A）はハイライトされないことを確認
	assert.Contains(t, result, "[ ] A → 11")
}

func createTestGamebookWithHistory() *domain.Gamebook {
	gamebook := domain.NewGamebook("test")
	
	// パラグラフ1を追加
	p1 := domain.NewParagraph(1, "start")
	p1.AddChoice("A", 11)
	p1.AddChoice("B", 12)
	gamebook.AddParagraph(p1)
	
	// 選択肢Bを選択してパラグラフ12に移動
	p1.SelectChoice(1)
	gamebook.MoveToOrCreatePlaceholder(12)
	
	// パラグラフ12を追加
	p12 := domain.NewParagraph(12, "test1")
	p12.AddChoice("C", 21)
	p12.AddChoice("D", 11)
	gamebook.AddParagraph(p12)
	
	// 選択肢Dを選択してパラグラフ11に移動
	p12.SelectChoice(1)
	gamebook.MoveToOrCreatePlaceholder(11)
	
	// 移動履歴を追加
	gamebook.AddNavigationStep(domain.NavigationStep{From: 1, To: 12, ViaPaths: []int{}})
	gamebook.AddNavigationStep(domain.NavigationStep{From: 12, To: 11, ViaPaths: []int{}})
	
	return gamebook
}