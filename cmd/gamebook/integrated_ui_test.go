package main

import (
	"testing"

	"github.com/zensai3805/iwakero-gamebook-mapping/internal/domain"
)

func TestIntegratedUI_NewIntegratedUI(t *testing.T) {
	ui := NewIntegratedUI()
	if ui == nil {
		t.Fatal("NewIntegratedUI() returned nil")
	}
}

func TestIntegratedUI_Initialize(t *testing.T) {
	ui := NewIntegratedUI()

	// Create test data
	gamebook := &domain.Gamebook{
		Title: "Test Game",
		Paragraphs: map[int]*domain.Paragraph{
			1: {Number: 1, Description: "Start", Visited: true},
			2: {Number: 2, Description: "Next", Visited: false},
		},
		Current: &domain.Paragraph{Number: 1, Description: "Start"},
	}

	data := &VisualizationData{
		Gamebook: gamebook,
	}

	err := ui.Initialize(data)
	if err != nil {
		t.Fatalf("Initialize() failed: %v", err)
	}

	// Check that components are initialized
	if ui.treePrinter == nil {
		t.Error("TreePrinter component not initialized")
	}

	if ui.areaPrinter == nil {
		t.Error("AreaPrinter component not initialized")
	}

	if ui.layoutManager == nil {
		t.Error("LayoutManager not initialized")
	}
}

func TestIntegratedUI_SetupLayout(t *testing.T) {
	ui := NewIntegratedUI()

	gamebook := &domain.Gamebook{
		Title: "Test Game",
		Paragraphs: map[int]*domain.Paragraph{
			1: {Number: 1, Description: "Start"},
		},
		Current: &domain.Paragraph{Number: 1, Description: "Start"},
	}

	data := &VisualizationData{
		Gamebook: gamebook,
	}

	err := ui.Initialize(data)
	if err != nil {
		t.Fatalf("Initialize() failed: %v", err)
	}

	// Test layout setup
	err = ui.SetupLayout()
	if err != nil {
		t.Fatalf("SetupLayout() failed: %v", err)
	}
}

func TestIntegratedUI_UpdateComponent(t *testing.T) {
	ui := NewIntegratedUI()

	gamebook := &domain.Gamebook{
		Title: "Test Game",
		Paragraphs: map[int]*domain.Paragraph{
			1: {Number: 1, Description: "Start"},
			2: {Number: 2, Description: "Next"},
		},
		Current: &domain.Paragraph{Number: 1, Description: "Start"},
	}

	data := &VisualizationData{
		Gamebook: gamebook,
	}

	err := ui.Initialize(data)
	if err != nil {
		t.Fatalf("Initialize() failed: %v", err)
	}
	err = ui.SetupLayout()
	if err != nil {
		t.Fatalf("SetupLayout() failed: %v", err)
	}

	// Test component update
	err = ui.UpdateComponent("map", EventCurrentChanged, data)
	if err != nil {
		t.Fatalf("UpdateComponent() failed: %v", err)
	}
}

func TestIntegratedUI_RenderAll(t *testing.T) {
	ui := NewIntegratedUI()

	gamebook := &domain.Gamebook{
		Title: "Test Game",
		Paragraphs: map[int]*domain.Paragraph{
			1: {Number: 1, Description: "Start", Visited: true},
		},
		Current: &domain.Paragraph{Number: 1, Description: "Start"},
	}

	// Create flow data for tree printer
	flowData := &FlowData{
		Root: &FlowNode{
			ParagraphNumber: 1,
			Description:     "Start",
			IsCurrent:       true,
			Visited:         true,
		},
	}

	// Create map data for area printer
	mapData := &MapData{
		Width:  3,
		Height: 3,
		Grid: [][]MapCell{
			{{ParagraphNumbers: []int{1}, Symbol: "1", IsCurrent: true, Visited: true}, {}, {}},
			{{}, {}, {}},
			{{}, {}, {}},
		},
		Positions: map[int]Position{
			1: {X: 0, Y: 0},
		},
	}

	data := &VisualizationData{
		Gamebook: gamebook,
		FlowData: flowData,
		MapData:  mapData,
	}

	err := ui.Initialize(data)
	if err != nil {
		t.Fatalf("Initialize() failed: %v", err)
	}
	err = ui.SetupLayout()
	if err != nil {
		t.Fatalf("SetupLayout() failed: %v", err)
	}

	// Test render
	output, err := ui.RenderAll()
	if err != nil {
		t.Fatalf("RenderAll() failed: %v", err)
	}

	if output == "" {
		t.Error("RenderAll() returned empty string")
	}
}

func TestIntegratedUI_HandleFocusChange(t *testing.T) {
	ui := NewIntegratedUI()

	gamebook := &domain.Gamebook{
		Title: "Test Game",
		Paragraphs: map[int]*domain.Paragraph{
			1: {Number: 1, Description: "Start"},
		},
		Current: &domain.Paragraph{Number: 1, Description: "Start"},
	}

	data := &VisualizationData{
		Gamebook: gamebook,
	}

	err := ui.Initialize(data)
	if err != nil {
		t.Fatalf("Initialize() failed: %v", err)
	}
	err = ui.SetupLayout()
	if err != nil {
		t.Fatalf("SetupLayout() failed: %v", err)
	}

	// Test focus changes
	testCases := []struct {
		name      string
		component string
		wantError bool
	}{
		{
			name:      "Focus on map",
			component: "map",
			wantError: false,
		},
		{
			name:      "Focus on flow",
			component: "flow",
			wantError: false,
		},
		{
			name:      "Focus on menu",
			component: "menu",
			wantError: false,
		},
		{
			name:      "Focus on invalid component",
			component: "invalid",
			wantError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := ui.HandleFocusChange(tc.component)
			if tc.wantError && err == nil {
				t.Error("HandleFocusChange() expected error but got nil")
			}
			if !tc.wantError && err != nil {
				t.Errorf("HandleFocusChange() unexpected error: %v", err)
			}

			if !tc.wantError && ui.focusedComponent != tc.component {
				t.Errorf("focusedComponent = %s, want %s", ui.focusedComponent, tc.component)
			}
		})
	}
}

func TestIntegratedUI_SyncSelection(t *testing.T) {
	ui := NewIntegratedUI()

	gamebook := &domain.Gamebook{
		Title: "Test Game",
		Paragraphs: map[int]*domain.Paragraph{
			1: {Number: 1, Description: "Start", Visited: true},
			2: {Number: 2, Description: "Next", Visited: false},
		},
		Current: &domain.Paragraph{Number: 1, Description: "Start"},
	}

	// Create comprehensive visualization data
	flowData := &FlowData{
		Root: &FlowNode{
			ParagraphNumber: 1,
			Description:     "Start",
			IsCurrent:       true,
			Visited:         true,
			Children: []*FlowNode{
				{
					ParagraphNumber: 2,
					Description:     "Next",
					Visited:         false,
				},
			},
		},
	}

	mapData := &MapData{
		Width:  2,
		Height: 1,
		Grid: [][]MapCell{
			{
				{ParagraphNumbers: []int{1}, Symbol: "1", IsCurrent: true, Visited: true},
				{ParagraphNumbers: []int{2}, Symbol: "2", Visited: false},
			},
		},
		Positions: map[int]Position{
			1: {X: 0, Y: 0},
			2: {X: 1, Y: 0},
		},
	}

	data := &VisualizationData{
		Gamebook: gamebook,
		FlowData: flowData,
		MapData:  mapData,
	}

	err := ui.Initialize(data)
	if err != nil {
		t.Fatalf("Initialize() failed: %v", err)
	}
	err = ui.SetupLayout()
	if err != nil {
		t.Fatalf("SetupLayout() failed: %v", err)
	}

	// Test selection sync - change current to paragraph 2
	gamebook.Current = gamebook.Paragraphs[2]
	data.CurrentPos = gamebook.Paragraphs[2]

	err = ui.SyncSelection(2)
	if err != nil {
		t.Fatalf("SyncSelection() failed: %v", err)
	}

	// Verify both components were updated
	if !data.FlowData.Root.Children[0].IsCurrent {
		t.Error("Flow data not synced - paragraph 2 should be current")
	}

	if !data.MapData.Grid[0][1].IsCurrent {
		t.Error("Map data not synced - paragraph 2 should be current")
	}
}

func TestIntegratedUI_HandleKeyInput(t *testing.T) {
	ui := NewIntegratedUI()

	gamebook := &domain.Gamebook{
		Title: "Test Game",
		Paragraphs: map[int]*domain.Paragraph{
			1: {Number: 1, Description: "Start"},
		},
		Current: &domain.Paragraph{Number: 1, Description: "Start"},
	}

	data := &VisualizationData{
		Gamebook: gamebook,
	}

	err := ui.Initialize(data)
	if err != nil {
		t.Fatalf("Initialize() failed: %v", err)
	}
	err = ui.SetupLayout()
	if err != nil {
		t.Fatalf("SetupLayout() failed: %v", err)
	}

	// Test key navigation
	testCases := []struct {
		name   string
		key    string
		result ComponentAction
	}{
		{
			name:   "Tab key",
			key:    "tab",
			result: ActionFocusNext,
		},
		{
			name:   "Shift+Tab key",
			key:    "shift+tab",
			result: ActionFocusPrevious,
		},
		{
			name:   "Enter key",
			key:    "enter",
			result: ActionSelect,
		},
		{
			name:   "Escape key",
			key:    "escape",
			result: ActionCancel,
		},
		{
			name:   "Arrow key",
			key:    "up",
			result: ActionNavigate,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			action := ui.HandleKeyInput(tc.key)
			if action != tc.result {
				t.Errorf("HandleKeyInput(%s) = %v, want %v", tc.key, action, tc.result)
			}
		})
	}
}

func TestIntegratedUI_ResizeLayout(t *testing.T) {
	ui := NewIntegratedUI()

	gamebook := &domain.Gamebook{
		Title: "Test Game",
		Paragraphs: map[int]*domain.Paragraph{
			1: {Number: 1, Description: "Start"},
		},
		Current: &domain.Paragraph{Number: 1, Description: "Start"},
	}

	data := &VisualizationData{
		Gamebook: gamebook,
	}

	err := ui.Initialize(data)
	if err != nil {
		t.Fatalf("Initialize() failed: %v", err)
	}
	err = ui.SetupLayout()
	if err != nil {
		t.Fatalf("SetupLayout() failed: %v", err)
	}

	// Test resize
	err = ui.ResizeLayout(120, 40)
	if err != nil {
		t.Fatalf("ResizeLayout() failed: %v", err)
	}

	// Verify dimensions were updated
	if ui.width != 120 || ui.height != 40 {
		t.Errorf("ResizeLayout() dimensions = (%d, %d), want (120, 40)", ui.width, ui.height)
	}
}
