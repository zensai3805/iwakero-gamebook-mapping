package main

import (
	"strings"
	"testing"

	"github.com/zensai3805/iwakero-gamebook-mapping/internal/domain"
)

func TestTreePrinter_NewTreePrinter(t *testing.T) {
	printer := NewTreePrinter()
	if printer == nil {
		t.Fatal("NewTreePrinter() returned nil")
	}
}

func TestTreePrinter_Initialize(t *testing.T) {
	printer := NewTreePrinter()

	// Create test data
	gamebook := &domain.Gamebook{
		Title: "Test Game",
		Paragraphs: map[int]*domain.Paragraph{
			1: {Number: 1, Description: "Start"},
			2: {Number: 2, Description: "Next"},
		},
		Current: &domain.Paragraph{Number: 1, Description: "Start"},
	}

	flowData := &FlowData{
		Root: &FlowNode{
			ParagraphNumber: 1,
			Description:     "Start",
			Children: []*FlowNode{
				{
					ParagraphNumber: 2,
					Description:     "Next",
					Children:        []*FlowNode{},
				},
			},
		},
	}

	data := &VisualizationData{
		Gamebook: gamebook,
		FlowData: flowData,
	}

	err := printer.Initialize(data)
	if err != nil {
		t.Fatalf("Initialize() failed: %v", err)
	}
}

func TestTreePrinter_Update(t *testing.T) {
	printer := NewTreePrinter()

	// Initialize first
	gamebook := &domain.Gamebook{
		Title: "Test Game",
		Paragraphs: map[int]*domain.Paragraph{
			1: {Number: 1, Description: "Start"},
		},
		Current: &domain.Paragraph{Number: 1, Description: "Start"},
	}

	flowData := &FlowData{
		Root: &FlowNode{
			ParagraphNumber: 1,
			Description:     "Start",
			IsCurrent:       true,
		},
	}

	data := &VisualizationData{
		Gamebook: gamebook,
		FlowData: flowData,
	}

	err := printer.Initialize(data)
	if err != nil {
		t.Fatalf("Initialize() failed: %v", err)
	}

	// Test update
	err = printer.Update(EventCurrentChanged, data)
	if err != nil {
		t.Fatalf("Update() failed: %v", err)
	}
}

func TestTreePrinter_Render(t *testing.T) {
	printer := NewTreePrinter()

	// Create test data
	flowData := &FlowData{
		Root: &FlowNode{
			ParagraphNumber: 1,
			Description:     "Start",
			IsCurrent:       true,
			Children: []*FlowNode{
				{
					ParagraphNumber: 2,
					Description:     "Next",
					Visited:         false,
				},
			},
		},
	}

	data := &VisualizationData{
		FlowData: flowData,
	}

	err := printer.Initialize(data)
	if err != nil {
		t.Fatalf("Initialize() failed: %v", err)
	}

	result, err := printer.Render()
	if err != nil {
		t.Fatalf("Render() failed: %v", err)
	}

	if result == "" {
		t.Error("Render() returned empty string")
	}
}

func TestTreePrinter_SetGetArea(t *testing.T) {
	printer := NewTreePrinter()

	// Set area
	err := printer.SetArea(80, 24)
	if err != nil {
		t.Fatalf("SetArea() failed: %v", err)
	}

	// Get area
	width, height := printer.GetArea()
	if width != 80 || height != 24 {
		t.Errorf("GetArea() = (%d, %d), want (80, 24)", width, height)
	}
}

func TestTreePrinter_ConvertFlowDataToTree(t *testing.T) {
	printer := NewTreePrinter()

	// Create test flow data
	flowData := &FlowData{
		Root: &FlowNode{
			ParagraphNumber: 1,
			Description:     "Start",
			IsCurrent:       true,
			Children: []*FlowNode{
				{
					ParagraphNumber: 2,
					Description:     "Choice A",
					Visited:         true,
					Children: []*FlowNode{
						{
							ParagraphNumber: 3,
							Description:     "End A",
							Visited:         false,
						},
					},
				},
				{
					ParagraphNumber: 4,
					Description:     "Choice B",
					Visited:         false,
				},
			},
		},
	}

	data := &VisualizationData{
		FlowData: flowData,
	}

	err := printer.Initialize(data)
	if err != nil {
		t.Fatalf("Initialize() failed: %v", err)
	}

	// This method should be accessible for testing
	tree := printer.convertFlowDataToTree(flowData)
	if tree == nil {
		t.Error("convertFlowDataToTree() returned nil")
		return
	}

	if len(*tree) == 0 {
		t.Error("convertFlowDataToTree() returned empty tree")
	}
}

func TestTreePrinter_StyleApplication(t *testing.T) {
	printer := NewTreePrinter()

	// Test different node states
	testCases := []struct {
		name      string
		node      *FlowNode
		wantStyle bool
	}{
		{
			name: "Current position",
			node: &FlowNode{
				ParagraphNumber: 1,
				Description:     "Current",
				IsCurrent:       true,
			},
			wantStyle: true,
		},
		{
			name: "Visited node",
			node: &FlowNode{
				ParagraphNumber: 2,
				Description:     "Visited",
				Visited:         true,
			},
			wantStyle: true,
		},
		{
			name: "Unvisited node",
			node: &FlowNode{
				ParagraphNumber: 3,
				Description:     "Unvisited",
				Visited:         false,
			},
			wantStyle: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			style := printer.getNodeStyle(tc.node)
			// Style should always be returned (never nil with new implementation)
			if tc.wantStyle {
				// Test that style contains some formatting
				styled := style.Sprint("test")
				if styled == "" {
					t.Errorf("getNodeStyle() returned empty styled text for %s", tc.name)
				}
			}
		})
	}
}

// TestTreePrinter_ChoiceDisplayInTree RED phase: 選択肢情報表示のテスト
func TestTreePrinter_ChoiceDisplayInTree(t *testing.T) {
	printer := NewTreePrinter()

	// 選択肢を持つテストデータを作成
	gamebook := &domain.Gamebook{
		Title: "Choice Test Game",
		Paragraphs: map[int]*domain.Paragraph{
			1: {
				Number:      1,
				Description: "Village entrance",
				Choices: []domain.Choice{
					{Description: "北へ進む", TargetNumber: 5, Selected: true},
					{Description: "南へ進む", TargetNumber: 10, Selected: false},
					{Description: "東へ進む", TargetNumber: 15, Selected: false},
				},
				Visited: true,
			},
			5: {
				Number:      5,
				Description: "Dark forest",
				Choices: []domain.Choice{
					{Description: "道を進む", TargetNumber: 20, Selected: false},
					{Description: "脇道に入る", TargetNumber: 25, Selected: true},
				},
				Visited: true,
			},
		},
		Current: &domain.Paragraph{Number: 5, Description: "Dark forest"},
	}

	// データ変換
	converter := NewDataConverter()
	visualData, err := converter.ConvertToVisualizationData(gamebook)
	if err != nil {
		t.Fatalf("ConvertToVisualizationData() failed: %v", err)
	}

	// TreePrinter初期化
	err = printer.Initialize(visualData)
	if err != nil {
		t.Fatalf("Initialize() failed: %v", err)
	}

	// レンダリング実行
	result, err := printer.Render()
	if err != nil {
		t.Fatalf("Render() failed: %v", err)
	}

	// 選択肢情報が含まれているかテスト
	expectedContent := []string{
		"1: Village entrance", // パラグラフ1
		"[✓] 北へ進む → 5",        // 選択済み選択肢
		"[ ] 南へ進む → 10",       // 未選択選択肢
		"[ ] 東へ進む → 15",       // 未選択選択肢
		"5: Dark forest",      // パラグラフ5
		"[ ] 道を進む → 20",       // 未選択選択肢
		"[✓] 脇道に入る → 25",      // 選択済み選択肢
	}

	for _, content := range expectedContent {
		if !strings.Contains(result, content) {
			t.Errorf("Expected content '%s' not found in rendered output", content)
		}
	}
}

// TestTreePrinter_ChoiceStatus GREEN phase: 選択肢状態表示のテスト
func TestTreePrinter_ChoiceStatus(t *testing.T) {
	printer := NewTreePrinter()

	testCases := []struct {
		name           string
		selected       bool
		expectedSymbol string
	}{
		{
			name:           "Selected choice",
			selected:       true,
			expectedSymbol: "[✓]",
		},
		{
			name:           "Unselected choice",
			selected:       false,
			expectedSymbol: "[ ]",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			choice := domain.Choice{
				Description:  "Test choice",
				TargetNumber: 99,
				Selected:     tc.selected,
			}

			// formatChoiceText関数をテスト
			result := printer.formatChoiceText(choice)

			// 期待されるシンボルが含まれているかテスト
			if !strings.Contains(result, tc.expectedSymbol) {
				t.Errorf("Expected symbol '%s' not found in result '%s'", tc.expectedSymbol, result)
			}

			// 選択肢の説明と遷移先が含まれているかテスト
			if !strings.Contains(result, "Test choice") {
				t.Errorf("Expected description 'Test choice' not found in result '%s'", result)
			}
			if !strings.Contains(result, "→ 99") {
				t.Errorf("Expected target '→ 99' not found in result '%s'", result)
			}
		})
	}
}
