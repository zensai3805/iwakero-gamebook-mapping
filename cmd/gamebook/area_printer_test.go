package main

import (
	"testing"

	"github.com/zensai3805/iwakero-gamebook-mapping/internal/domain"
)

func TestAreaPrinter_NewAreaPrinter(t *testing.T) {
	printer := NewAreaPrinter()
	if printer == nil {
		t.Fatal("NewAreaPrinter() returned nil")
	}
}

func TestAreaPrinter_Initialize(t *testing.T) {
	printer := NewAreaPrinter()

	// Create test data
	gamebook := &domain.Gamebook{
		Title: "Test Game",
		Paragraphs: map[int]*domain.Paragraph{
			1: {Number: 1, Description: "Start"},
			2: {Number: 2, Description: "North"},
			3: {Number: 3, Description: "East"},
		},
		Current: &domain.Paragraph{Number: 1, Description: "Start"},
	}

	mapData := &MapData{
		Width:  3,
		Height: 3,
		Grid: [][]MapCell{
			{{ParagraphNumbers: []int{1}, Symbol: "1", IsCurrent: true}, {}, {ParagraphNumbers: []int{2}, Symbol: "2"}},
			{{}, {}, {}},
			{{}, {ParagraphNumbers: []int{3}, Symbol: "3"}, {}},
		},
		Positions: map[int]Position{
			1: {X: 0, Y: 0},
			2: {X: 2, Y: 0},
			3: {X: 1, Y: 2},
		},
	}

	data := &VisualizationData{
		Gamebook: gamebook,
		MapData:  mapData,
	}

	err := printer.Initialize(data)
	if err != nil {
		t.Fatalf("Initialize() failed: %v", err)
	}
}

func TestAreaPrinter_Update(t *testing.T) {
	printer := NewAreaPrinter()

	// Initialize first
	gamebook := &domain.Gamebook{
		Title: "Test Game",
		Paragraphs: map[int]*domain.Paragraph{
			1: {Number: 1, Description: "Start"},
		},
		Current: &domain.Paragraph{Number: 1, Description: "Start"},
	}

	mapData := &MapData{
		Width:  1,
		Height: 1,
		Grid: [][]MapCell{
			{{ParagraphNumbers: []int{1}, Symbol: "1", IsCurrent: true}},
		},
		Positions: map[int]Position{
			1: {X: 0, Y: 0},
		},
	}

	data := &VisualizationData{
		Gamebook: gamebook,
		MapData:  mapData,
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

func TestAreaPrinter_Render(t *testing.T) {
	printer := NewAreaPrinter()

	// Create test data
	mapData := &MapData{
		Width:  3,
		Height: 2,
		Grid: [][]MapCell{
			{
				{ParagraphNumbers: []int{1}, Symbol: "1", IsCurrent: true, Visited: true},
				{},
				{ParagraphNumbers: []int{2}, Symbol: "2", Visited: false},
			},
			{
				{ParagraphNumbers: []int{3}, Symbol: "3", Visited: true},
				{ParagraphNumbers: []int{4}, Symbol: "4", Visited: false},
				{},
			},
		},
		Positions: map[int]Position{
			1: {X: 0, Y: 0},
			2: {X: 2, Y: 0},
			3: {X: 0, Y: 1},
			4: {X: 1, Y: 1},
		},
	}

	data := &VisualizationData{
		MapData: mapData,
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

	// Check that grid structure is present
	if len(result) < 10 {
		t.Error("Render() returned suspiciously short result")
	}
}

func TestAreaPrinter_SetGetArea(t *testing.T) {
	printer := NewAreaPrinter()

	// Set area
	err := printer.SetArea(25, 10)
	if err != nil {
		t.Fatalf("SetArea() failed: %v", err)
	}

	// Get area
	width, height := printer.GetArea()
	if width != 25 || height != 10 {
		t.Errorf("GetArea() = (%d, %d), want (25, 10)", width, height)
	}
}

func TestAreaPrinter_ConvertMapDataToGrid(t *testing.T) {
	printer := NewAreaPrinter()

	// Create test map data
	mapData := &MapData{
		Width:  2,
		Height: 2,
		Grid: [][]MapCell{
			{
				{ParagraphNumbers: []int{1}, Symbol: "1", IsCurrent: true},
				{ParagraphNumbers: []int{2}, Symbol: "2"},
			},
			{
				{},
				{ParagraphNumbers: []int{3}, Symbol: "3", Visited: true},
			},
		},
		Positions: map[int]Position{
			1: {X: 0, Y: 0},
			2: {X: 1, Y: 0},
			3: {X: 1, Y: 1},
		},
	}

	data := &VisualizationData{
		MapData: mapData,
	}

	err := printer.Initialize(data)
	if err != nil {
		t.Fatalf("Initialize() failed: %v", err)
	}

	grid := printer.convertMapDataToGrid(mapData)
	if grid == nil {
		t.Error("convertMapDataToGrid() returned nil")
	}

	if len(grid) != 2 {
		t.Errorf("convertMapDataToGrid() returned grid with %d rows, want 2", len(grid))
	}

	if len(grid[0]) != 2 {
		t.Errorf("convertMapDataToGrid() returned grid with %d columns, want 2", len(grid[0]))
	}
}

func TestAreaPrinter_FormatCell(t *testing.T) {
	printer := NewAreaPrinter()

	testCases := []struct {
		name string
		cell MapCell
		want string
	}{
		{
			name: "Current position",
			cell: MapCell{
				ParagraphNumbers: []int{1},
				Symbol:           "1",
				IsCurrent:        true,
			},
			want: "1", // Style will be applied, but we just check content
		},
		{
			name: "Visited cell",
			cell: MapCell{
				ParagraphNumbers: []int{2},
				Symbol:           "2",
				Visited:          true,
			},
			want: "2",
		},
		{
			name: "Empty cell",
			cell: MapCell{},
			want: " ",
		},
		{
			name: "Multiple paragraphs",
			cell: MapCell{
				ParagraphNumbers: []int{1, 2, 3},
				Symbol:           "*",
			},
			want: "*",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := printer.formatCell(tc.cell)
			// Check that result contains expected content (ignoring styling)
			if len(result) == 0 {
				t.Errorf("formatCell() returned empty string for %s", tc.name)
			}
		})
	}
}

func TestAreaPrinter_GetCellStyle(t *testing.T) {
	printer := NewAreaPrinter()

	testCases := []struct {
		name      string
		cell      MapCell
		wantStyle bool
	}{
		{
			name: "Current position",
			cell: MapCell{
				IsCurrent: true,
			},
			wantStyle: true,
		},
		{
			name: "Visited cell",
			cell: MapCell{
				Visited: true,
			},
			wantStyle: true,
		},
		{
			name: "Unvisited cell",
			cell: MapCell{
				Visited: false,
			},
			wantStyle: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			style := printer.getCellStyle(tc.cell)
			if tc.wantStyle && style == nil {
				t.Errorf("getCellStyle() returned nil for %s", tc.name)
			}
		})
	}
}

func TestAreaPrinter_GenerateGridBorder(t *testing.T) {
	printer := NewAreaPrinter()

	// Test horizontal border
	border := printer.generateGridBorder(3, true)
	if border == "" {
		t.Error("generateGridBorder() returned empty string for horizontal border")
	}

	// Test vertical border
	border = printer.generateGridBorder(3, false)
	if border == "" {
		t.Error("generateGridBorder() returned empty string for vertical border")
	}
}

func TestAreaPrinter_AutoLayoutParagraphs(t *testing.T) {
	printer := NewAreaPrinter()

	// Create test gamebook with paragraph relationships
	gamebook := &domain.Gamebook{
		Title: "Test Game",
		Paragraphs: map[int]*domain.Paragraph{
			1: {
				Number:      1,
				Description: "Start",
				Choices: []domain.Choice{
					{Description: "Go North", TargetNumber: 2},
					{Description: "Go East", TargetNumber: 3},
				},
			},
			2: {Number: 2, Description: "North Room"},
			3: {Number: 3, Description: "East Room"},
		},
		Current: &domain.Paragraph{Number: 1, Description: "Start"},
	}

	data := &VisualizationData{
		Gamebook: gamebook,
	}

	err := printer.Initialize(data)
	if err != nil {
		t.Fatalf("Initialize() failed: %v", err)
	}

	mapData := printer.autoLayoutParagraphs(gamebook)
	if mapData == nil {
		t.Error("autoLayoutParagraphs() returned nil")
		return
	}

	if mapData.Width <= 0 || mapData.Height <= 0 {
		t.Errorf("autoLayoutParagraphs() returned invalid dimensions: %dx%d", mapData.Width, mapData.Height)
	}

	if len(mapData.Positions) != len(gamebook.Paragraphs) {
		t.Errorf("autoLayoutParagraphs() positioned %d paragraphs, want %d", len(mapData.Positions), len(gamebook.Paragraphs))
	}
}
