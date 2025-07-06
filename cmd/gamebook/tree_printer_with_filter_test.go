package main

import (
	"strings"
	"testing"

	"github.com/zensai3805/iwakero-gamebook-mapping/internal/domain"
)

func TestTreePrinter_WithDisplayFilter(t *testing.T) {
	tests := []struct {
		name             string
		setup            func() (*domain.Gamebook, *VisualizationData)
		scope            DisplayScope
		shouldContain    []string
		shouldNotContain []string
		description      string
	}{
		{
			name: "Connected scope - 接続された未定義パラグラフのみ表示",
			setup: func() (*domain.Gamebook, *VisualizationData) {
				gb := domain.NewGamebook("Test")

				// パラグラフ1（現在位置）
				p1 := domain.NewParagraph(1, "冒険の始まり")
				_ = gb.AddParagraph(p1)
				_ = gb.AddChoiceToParagraph(1, "北へ進む", 10) // 接続された未定義
				_ = gb.AddChoiceToParagraph(1, "南へ進む", 2)  // 定義済み
				gb.Current = p1

				// パラグラフ2（定義済み）
				p2 := domain.NewParagraph(2, "次の場所")
				_ = gb.AddParagraph(p2)

				// パラグラフ5（孤立）
				p5 := domain.NewParagraph(5, "別の場所")
				_ = gb.AddParagraph(p5)
				_ = gb.AddChoiceToParagraph(5, "東へ進む", 20) // 孤立した未定義

				// VisualizationDataを作成
				flowData := createFlowDataFromGamebook(gb)
				vizData := &VisualizationData{
					FlowData: flowData,
				}

				return gb, vizData
			},
			scope: ScopeConnected,
			shouldContain: []string{
				"1: 冒険の始まり",
				"2: 次の場所",
				"[ ] 北へ進む → 10", // 接続された未定義選択肢
				"[ ] 南へ進む → 2",  // 定義済み選択肢
			},
			shouldNotContain: []string{
				"[ ] 東へ進む → 20", // 孤立した未定義選択肢
				"20:",           // 孤立した未定義パラグラフ
			},
			description: "接続された未定義パラグラフと選択肢のみ表示",
		},
		{
			name: "None scope - 未定義パラグラフと選択肢を非表示",
			setup: func() (*domain.Gamebook, *VisualizationData) {
				gb := domain.NewGamebook("Test")

				p1 := domain.NewParagraph(1, "冒険の始まり")
				_ = gb.AddParagraph(p1)
				_ = gb.AddChoiceToParagraph(1, "未定義へ", 10) // 未定義選択肢
				_ = gb.AddChoiceToParagraph(1, "定義済みへ", 2) // 定義済み選択肢
				gb.Current = p1

				p2 := domain.NewParagraph(2, "次の場所")
				_ = gb.AddParagraph(p2)

				flowData := createFlowDataFromGamebook(gb)
				vizData := &VisualizationData{
					FlowData: flowData,
				}

				return gb, vizData
			},
			scope: ScopeNone,
			shouldContain: []string{
				"1: 冒険の始まり",
				"2: 次の場所",
				"[ ] 定義済みへ → 2", // 定義済み選択肢のみ
			},
			shouldNotContain: []string{
				"[ ] 未定義へ → 10", // 未定義選択肢
				"10:",           // 未定義パラグラフ
			},
			description: "未定義パラグラフと選択肢を完全に非表示",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gb, vizData := tt.setup()

			// TreePrinterWithFilterを使用
			printer := NewTreePrinterWithFilter()
			printer.SetDisplayScope(tt.scope)

			initErr := printer.Initialize(vizData)
			if initErr != nil {
				t.Fatalf("Failed to initialize printer: %v", initErr)
			}

			// gamebookを設定
			printer.SetGamebook(gb)

			result, renderErr := printer.Render()
			if renderErr != nil {
				t.Fatalf("Failed to render: %v", renderErr)
			}

			// 含まれるべきテキストの確認
			for _, shouldContain := range tt.shouldContain {
				if !strings.Contains(result, shouldContain) {
					t.Errorf("Result should contain '%s' but doesn't.\nResult:\n%s", shouldContain, result)
				}
			}

			// 含まれてはいけないテキストの確認
			for _, shouldNotContain := range tt.shouldNotContain {
				if strings.Contains(result, shouldNotContain) {
					t.Errorf("Result should not contain '%s' but does.\nResult:\n%s", shouldNotContain, result)
				}
			}
		})
	}
}

// createFlowDataFromGamebook は簡単なFlowData作成のヘルパー関数
func createFlowDataFromGamebook(gb *domain.Gamebook) *FlowData {
	if gb.Current == nil {
		return nil
	}

	// 現在位置をルートノードとして作成
	root := &FlowNode{
		ParagraphNumber: gb.Current.Number,
		Description:     gb.Current.Description,
		IsCurrent:       true,
		Visited:         gb.Current.Visited,
		Choices:         gb.Current.Choices,
		Children:        []*FlowNode{},
	}

	// 定義済みパラグラフを子ノードとして追加
	for _, choice := range gb.Current.Choices {
		if targetParagraph, exists := gb.Paragraphs[choice.TargetNumber]; exists {
			childNode := &FlowNode{
				ParagraphNumber: targetParagraph.Number,
				Description:     targetParagraph.Description,
				IsCurrent:       false,
				Visited:         targetParagraph.Visited,
				Choices:         targetParagraph.Choices,
				Children:        []*FlowNode{},
			}
			root.Children = append(root.Children, childNode)
		}
	}

	return &FlowData{
		Root: root,
	}
}
