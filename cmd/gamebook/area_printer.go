package main

import (
	"fmt"
	"math"
	"strings"

	"github.com/iwapc/iwakero-gamebook-mapping/internal/domain"
	"github.com/pterm/pterm"
)

// AreaPrinter implements IVisualizer for 2D grid-based map visualization
type AreaPrinter struct {
	data   *VisualizationData
	width  int
	height int
}

// NewAreaPrinter creates a new AreaPrinter instance
func NewAreaPrinter() *AreaPrinter {
	return &AreaPrinter{}
}

// Initialize initializes the AreaPrinter with visualization data
func (ap *AreaPrinter) Initialize(data *VisualizationData) error {
	if data == nil {
		return fmt.Errorf("visualization data cannot be nil")
	}
	
	ap.data = data
	
	// Auto-generate map if MapData is not provided
	if data.MapData == nil && data.Gamebook != nil {
		data.MapData = ap.autoLayoutParagraphs(data.Gamebook)
	}
	
	return nil
}

// Update processes data updates and refreshes the map visualization
func (ap *AreaPrinter) Update(event VisualizationEvent, data *VisualizationData) error {
	if data == nil {
		return fmt.Errorf("visualization data cannot be nil")
	}
	
	ap.data = data
	
	// Update current position highlighting if current position changed
	if event == EventCurrentChanged && data.MapData != nil {
		ap.updateCurrentPositionHighlight()
	}
	
	return nil
}

// Render generates the 2D map visualization as a string
func (ap *AreaPrinter) Render() (string, error) {
	if ap.data == nil || ap.data.MapData == nil {
		return "", fmt.Errorf("no map data available")
	}
	
	grid := ap.convertMapDataToGrid(ap.data.MapData)
	if grid == nil {
		return "", fmt.Errorf("failed to convert map data to grid")
	}
	
	return ap.renderGrid(grid), nil
}

// GetArea returns the current drawing area dimensions
func (ap *AreaPrinter) GetArea() (width, height int) {
	return ap.width, ap.height
}

// SetArea sets the drawing area dimensions
func (ap *AreaPrinter) SetArea(width, height int) error {
	if width < 0 || height < 0 {
		return fmt.Errorf("width and height must be non-negative")
	}
	
	ap.width = width
	ap.height = height
	return nil
}

// convertMapDataToGrid converts MapData to a displayable grid
func (ap *AreaPrinter) convertMapDataToGrid(mapData *MapData) [][]string {
	if mapData == nil || mapData.Grid == nil {
		return nil
	}
	
	grid := make([][]string, mapData.Height)
	for i := range grid {
		grid[i] = make([]string, mapData.Width)
		for j := range grid[i] {
			if i < len(mapData.Grid) && j < len(mapData.Grid[i]) {
				grid[i][j] = ap.formatCell(mapData.Grid[i][j])
			} else {
				grid[i][j] = " "
			}
		}
	}
	
	return grid
}

// formatCell formats a MapCell for display with appropriate styling
func (ap *AreaPrinter) formatCell(cell MapCell) string {
	if len(cell.ParagraphNumbers) == 0 {
		return " "
	}
	
	// Use the symbol if provided, otherwise use the first paragraph number
	text := cell.Symbol
	if text == "" {
		text = fmt.Sprintf("%d", cell.ParagraphNumbers[0])
	}
	
	// Apply styling based on cell state
	style := ap.getCellStyle(cell)
	if style != nil {
		return style.Sprint(text)
	}
	
	return text
}

// getCellStyle returns the appropriate style for a cell based on its state
func (ap *AreaPrinter) getCellStyle(cell MapCell) *pterm.Style {
	// Current position gets highest priority styling
	if cell.IsCurrent {
		style := pterm.NewStyle(pterm.FgYellow, pterm.Bold, pterm.BgBlue)
		return style
	}
	
	// Visited cells get secondary styling
	if cell.Visited {
		style := pterm.NewStyle(pterm.FgGreen)
		return style
	}
	
	// Unvisited cells get default styling
	style := pterm.NewStyle(pterm.FgLightWhite)
	return style
}

// renderGrid renders the grid with borders and formatting
func (ap *AreaPrinter) renderGrid(grid [][]string) string {
	if len(grid) == 0 {
		return ""
	}
	
	var result strings.Builder
	width := len(grid[0])
	
	// Top border
	result.WriteString(ap.generateGridBorder(width, true))
	result.WriteString("\n")
	
	// Grid rows with borders
	for i, row := range grid {
		result.WriteString("|")
		for _, cell := range row {
			result.WriteString(" ")
			result.WriteString(cell)
			result.WriteString(" |")
		}
		result.WriteString("\n")
		
		// Add horizontal separator between rows (except after last row)
		if i < len(grid)-1 {
			result.WriteString(ap.generateGridBorder(width, false))
			result.WriteString("\n")
		}
	}
	
	// Bottom border
	result.WriteString(ap.generateGridBorder(width, true))
	
	return result.String()
}

// generateGridBorder generates horizontal border lines for the grid
func (ap *AreaPrinter) generateGridBorder(width int, _ bool) string {
	var result strings.Builder
	
	for i := 0; i < width; i++ {
		if i == 0 {
			result.WriteString("+")
		}
		result.WriteString("---+")
	}
	
	return result.String()
}

// updateCurrentPositionHighlight updates the current position highlighting in the map
func (ap *AreaPrinter) updateCurrentPositionHighlight() {
	if ap.data.MapData == nil || ap.data.Gamebook == nil || ap.data.Gamebook.Current == nil {
		return
	}
	
	currentNum := ap.data.Gamebook.Current.Number
	
	// Clear all current position flags
	for i := range ap.data.MapData.Grid {
		for j := range ap.data.MapData.Grid[i] {
			ap.data.MapData.Grid[i][j].IsCurrent = false
		}
	}
	
	// Set current position flag
	if pos, exists := ap.data.MapData.Positions[currentNum]; exists {
		if pos.Y < len(ap.data.MapData.Grid) && pos.X < len(ap.data.MapData.Grid[pos.Y]) {
			ap.data.MapData.Grid[pos.Y][pos.X].IsCurrent = true
		}
	}
}

// autoLayoutParagraphs automatically generates a 2D layout for paragraphs
func (ap *AreaPrinter) autoLayoutParagraphs(gamebook *domain.Gamebook) *MapData {
	if gamebook == nil || len(gamebook.Paragraphs) == 0 {
		return &MapData{
			Width:     1,
			Height:    1,
			Grid:      [][]MapCell{{{}}},
			Positions: make(map[int]Position),
		}
	}
	
	// Calculate grid size based on paragraph count
	numParagraphs := len(gamebook.Paragraphs)
	size := int(math.Ceil(math.Sqrt(float64(numParagraphs))))
	
	// Ensure minimum grid size
	if size < 3 {
		size = 3
	}
	
	// Initialize grid
	grid := make([][]MapCell, size)
	for i := range grid {
		grid[i] = make([]MapCell, size)
	}
	
	positions := make(map[int]Position)
	
	// Simple grid placement algorithm (left-to-right, top-to-bottom)
	placedCount := 0
	for num := range gamebook.Paragraphs {
		if placedCount >= size*size {
			break
		}
		
		// Calculate grid position
		x := placedCount % size
		y := placedCount / size
		
		// Place paragraph at current position
		positions[num] = Position{X: x, Y: y}
		paragraph := gamebook.Paragraphs[num]
		
		grid[y][x] = MapCell{
			ParagraphNumbers: []int{num},
			Symbol:           fmt.Sprintf("%d", num),
			Visited:          paragraph.Visited,
			IsCurrent:        gamebook.Current != nil && gamebook.Current.Number == num,
		}
		
		placedCount++
	}
	
	return &MapData{
		Width:     size,
		Height:    size,
		Grid:      grid,
		Positions: positions,
	}
}