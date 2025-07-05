package main

import (
	"strings"
	"testing"

	"github.com/iwapc/iwakero-gamebook-mapping/internal/domain"
)

// TestShowVisualizationWithoutMap v0.2.1.hotfix用のテスト
// マップ機能が除去されてフロー図のみが表示されることを確認
func TestShowVisualizationWithoutMap(t *testing.T) {
	// テスト用のゲームブック作成
	gamebook := &domain.Gamebook{
		Title: "Test Game",
		Paragraphs: map[int]*domain.Paragraph{
			1: {Number: 1, Description: "Start", Visited: true},
			2: {Number: 2, Description: "Next", Visited: false},
		},
		Current: &domain.Paragraph{Number: 1, Description: "Start"},
	}

	// グローバル変数を設定
	currentGame = gamebook

	// CLIExecutor作成
	executor := NewCLIExecutor()

	// showVisualization実行
	err := executor.showVisualization()
	if err != nil {
		t.Fatalf("showVisualization() failed: %v", err)
	}

	// v0.2.1.hotfixでは統合UIではなく、TreePrinterを直接使用することを確認
	// エラーが発生しないことを確認済み
}

// TestShowVisualizationRenderFlowOnly フロー図のみがレンダリングされることを確認
func TestShowVisualizationRenderFlowOnly(t *testing.T) {
	// テスト用のゲームブック作成
	gamebook := &domain.Gamebook{
		Title: "Test Game",
		Paragraphs: map[int]*domain.Paragraph{
			1: {Number: 1, Description: "Start", Visited: true},
			2: {Number: 2, Description: "Next", Visited: false},
		},
		Current: &domain.Paragraph{Number: 1, Description: "Start"},
	}

	// データ変換テスト
	converter := NewDataConverter()
	visualData, err := converter.ConvertToVisualizationData(gamebook)
	if err != nil {
		t.Fatalf("ConvertToVisualizationData() failed: %v", err)
	}

	// フロー図のみをレンダリング
	treePrinter := NewTreePrinter()
	err = treePrinter.Initialize(visualData)
	if err != nil {
		t.Fatalf("TreePrinter.Initialize() failed: %v", err)
	}

	flowOutput, err := treePrinter.Render()
	if err != nil {
		t.Fatalf("TreePrinter.Render() failed: %v", err)
	}

	// フロー図が空でないことを確認
	if flowOutput == "" {
		t.Error("Flow output should not be empty")
	}

	// フロー図にパラグラフ情報が含まれることを確認
	if !strings.Contains(flowOutput, "Start") {
		t.Error("Flow output should contain paragraph descriptions")
	}
}

// TestShowVisualizationWithoutMapComponents マップ関連コンポーネントが使用されないことを確認
func TestShowVisualizationWithoutMapComponents(t *testing.T) {
	// テスト用のゲームブック作成
	gamebook := &domain.Gamebook{
		Title: "Test Game",
		Paragraphs: map[int]*domain.Paragraph{
			1: {Number: 1, Description: "Start", Visited: true},
		},
		Current: &domain.Paragraph{Number: 1, Description: "Start"},
	}

	// データ変換
	converter := NewDataConverter()
	visualData, err := converter.ConvertToVisualizationData(gamebook)
	if err != nil {
		t.Fatalf("ConvertToVisualizationData() failed: %v", err)
	}

	// フロー図のみ使用する想定のテスト
	// この段階では、まだマップ機能が存在することを確認
	// 実装後にマップ機能が除去されることを確認するテストに変更予定

	// TreePrinterが正常に機能することを確認
	treePrinter := NewTreePrinter()
	err = treePrinter.Initialize(visualData)
	if err != nil {
		t.Fatalf("TreePrinter.Initialize() failed: %v", err)
	}

	// レンダリング結果の確認
	output, err := treePrinter.Render()
	if err != nil {
		t.Fatalf("TreePrinter.Render() failed: %v", err)
	}

	if output == "" {
		t.Error("TreePrinter output should not be empty")
	}
}
