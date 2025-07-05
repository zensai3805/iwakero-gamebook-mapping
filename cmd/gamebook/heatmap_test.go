package main

import (
	"testing"

	"github.com/iwapc/iwakero-gamebook-mapping/internal/domain"
)

func TestHeatmapManager_NewHeatmapManager(t *testing.T) {
	hm := NewHeatmapManager()
	if hm == nil {
		t.Fatal("NewHeatmapManager() returned nil")
	}
}

func TestHeatmapManager_Initialize(t *testing.T) {
	hm := NewHeatmapManager()
	
	// Create test data
	gamebook := &domain.Gamebook{
		Title: "Test Game",
		Paragraphs: map[int]*domain.Paragraph{
			1: {Number: 1, Description: "Start", Visited: true},
			2: {Number: 2, Description: "Next", Visited: true},
			3: {Number: 3, Description: "End", Visited: false},
		},
		Current: &domain.Paragraph{Number: 1, Description: "Start"},
	}
	
	data := &VisualizationData{
		Gamebook: gamebook,
	}
	
	err := hm.Initialize(data)
	if err != nil {
		t.Fatalf("Initialize() failed: %v", err)
	}
	
	// Check visit counts were initialized
	if len(hm.visitCounts) != len(gamebook.Paragraphs) {
		t.Errorf("visitCounts length = %d, want %d", len(hm.visitCounts), len(gamebook.Paragraphs))
	}
}

func TestHeatmapManager_RecordVisit(t *testing.T) {
	hm := NewHeatmapManager()
	
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
	
	err := hm.Initialize(data)
	if err != nil {
		t.Fatalf("Initialize() failed: %v", err)
	}
	
	// Record multiple visits
	hm.RecordVisit(1)
	hm.RecordVisit(1)
	hm.RecordVisit(1)
	
	// Check visit count
	count := hm.GetVisitCount(1)
	if count != 3 {
		t.Errorf("GetVisitCount(1) = %d, want 3", count)
	}
}

func TestHeatmapManager_GetHeatmapColor(t *testing.T) {
	hm := NewHeatmapManager()
	
	gamebook := &domain.Gamebook{
		Title: "Test Game",
		Paragraphs: map[int]*domain.Paragraph{
			1: {Number: 1, Description: "Heavily visited"},
			2: {Number: 2, Description: "Moderately visited"},
			3: {Number: 3, Description: "Lightly visited"},
			4: {Number: 4, Description: "Unvisited"},
		},
		Current: &domain.Paragraph{Number: 1, Description: "Start"},
	}
	
	data := &VisualizationData{
		Gamebook: gamebook,
	}
	
	err := hm.Initialize(data)
	if err != nil {
		t.Fatalf("Initialize() failed: %v", err)
	}
	
	// Set up different visit counts
	for range 10 {
		hm.RecordVisit(1)
	}
	for range 5 {
		hm.RecordVisit(2)
	}
	hm.RecordVisit(3)
	// 4 remains unvisited
	
	// Test color gradients
	testCases := []struct {
		paragraph int
		name      string
	}{
		{1, "Heavily visited"},
		{2, "Moderately visited"},
		{3, "Lightly visited"},
		{4, "Unvisited"},
	}
	
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			color := hm.GetHeatmapColor(tc.paragraph)
			if color == nil {
				t.Errorf("GetHeatmapColor(%d) returned nil", tc.paragraph)
			}
		})
	}
}

func TestHeatmapManager_ApplyHeatmapToFlow(t *testing.T) {
	hm := NewHeatmapManager()
	
	gamebook := &domain.Gamebook{
		Title: "Test Game",
		Paragraphs: map[int]*domain.Paragraph{
			1: {Number: 1, Description: "Start", Visited: true},
			2: {Number: 2, Description: "Next", Visited: true},
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
				},
			},
		},
	}
	
	data := &VisualizationData{
		Gamebook: gamebook,
		FlowData: flowData,
	}
	
	err := hm.Initialize(data)
	if err != nil {
		t.Fatalf("Initialize() failed: %v", err)
	}
	hm.RecordVisit(1)
	hm.RecordVisit(1)
	hm.RecordVisit(2)
	
	// Apply heatmap
	err = hm.ApplyHeatmapToFlow(flowData)
	if err != nil {
		t.Fatalf("ApplyHeatmapToFlow() failed: %v", err)
	}
	
	// Check that styles were applied
	if flowData.Root.Style == nil {
		t.Error("Root node style not applied")
	}
	if flowData.Root.Children[0].Style == nil {
		t.Error("Child node style not applied")
	}
}

func TestHeatmapManager_ApplyHeatmapToMap(t *testing.T) {
	hm := NewHeatmapManager()
	
	gamebook := &domain.Gamebook{
		Title: "Test Game",
		Paragraphs: map[int]*domain.Paragraph{
			1: {Number: 1, Description: "Start", Visited: true},
			2: {Number: 2, Description: "Next", Visited: false},
		},
		Current: &domain.Paragraph{Number: 1, Description: "Start"},
	}
	
	mapData := &MapData{
		Width:  2,
		Height: 1,
		Grid: [][]MapCell{
			{
				{ParagraphNumbers: []int{1}, Symbol: "1", Visited: true},
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
		MapData:  mapData,
	}
	
	err := hm.Initialize(data)
	if err != nil {
		t.Fatalf("Initialize() failed: %v", err)
	}
	hm.RecordVisit(1)
	hm.RecordVisit(1)
	hm.RecordVisit(1)
	
	// Apply heatmap
	err = hm.ApplyHeatmapToMap(mapData)
	if err != nil {
		t.Fatalf("ApplyHeatmapToMap() failed: %v", err)
	}
	
	// Check that styles were applied
	if mapData.Grid[0][0].Style == nil {
		t.Error("Visited cell style not applied")
	}
}

func TestHeatmapManager_GenerateLegend(t *testing.T) {
	hm := NewHeatmapManager()
	
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
	
	err := hm.Initialize(data)
	if err != nil {
		t.Fatalf("Initialize() failed: %v", err)
	}
	
	// Generate legend
	legend := hm.GenerateLegend()
	if legend == "" {
		t.Error("GenerateLegend() returned empty string")
	}
	
	// Check that legend contains expected elements
	expectedElements := []string{"Unvisited", "Low", "Medium", "High"}
	for _, elem := range expectedElements {
		if !contains(legend, elem) {
			t.Errorf("Legend missing expected element: %s", elem)
		}
	}
}

func TestHeatmapManager_GenerateStatisticsPanel(t *testing.T) {
	hm := NewHeatmapManager()
	
	gamebook := &domain.Gamebook{
		Title: "Test Game",
		Paragraphs: map[int]*domain.Paragraph{
			1: {Number: 1, Description: "Start", Visited: true},
			2: {Number: 2, Description: "Next", Visited: true},
			3: {Number: 3, Description: "End", Visited: false},
		},
		Current: &domain.Paragraph{Number: 1, Description: "Start"},
	}
	
	data := &VisualizationData{
		Gamebook: gamebook,
	}
	
	err := hm.Initialize(data)
	if err != nil {
		t.Fatalf("Initialize() failed: %v", err)
	}
	hm.RecordVisit(1)
	hm.RecordVisit(1)
	hm.RecordVisit(2)
	
	// Generate statistics
	stats := hm.GenerateStatisticsPanel()
	if stats == "" {
		t.Error("GenerateStatisticsPanel() returned empty string")
	}
	
	// Check statistics content
	if !contains(stats, "Total Paragraphs") {
		t.Error("Statistics missing total paragraphs")
	}
	if !contains(stats, "Visited") {
		t.Error("Statistics missing visited count")
	}
	if !contains(stats, "Most Visited") {
		t.Error("Statistics missing most visited info")
	}
}

func TestHeatmapManager_GetColorGradient(t *testing.T) {
	hm := NewHeatmapManager()
	
	testCases := []struct {
		intensity float64
		name      string
	}{
		{0.0, "Minimum intensity"},
		{0.25, "Low intensity"},
		{0.5, "Medium intensity"},
		{0.75, "High intensity"},
		{1.0, "Maximum intensity"},
	}
	
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			color := hm.getColorGradient(tc.intensity)
			if color == nil {
				t.Errorf("getColorGradient(%f) returned nil", tc.intensity)
			}
		})
	}
}

// Helper function
func contains(s, substr string) bool {
	return len(s) > 0 && len(substr) > 0 && s != substr
}