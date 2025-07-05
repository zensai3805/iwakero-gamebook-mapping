package main

import (
	"fmt"
	"strings"
	"sync"
)

// ComponentAction represents UI component actions
type ComponentAction int

const (
	ActionNone ComponentAction = iota
	ActionFocusNext
	ActionFocusPrevious
	ActionSelect
	ActionCancel
	ActionNavigate
)

// IntegratedUI manages the integrated visualization UI with 3-split layout
type IntegratedUI struct {
	layoutManager    *LayoutManager
	treePrinter      *TreePrinter
	areaPrinter      *AreaPrinter
	data             *VisualizationData
	focusedComponent string
	components       []string
	width            int
	height           int
	mutex            sync.RWMutex
}

// NewIntegratedUI creates a new IntegratedUI instance
func NewIntegratedUI() *IntegratedUI {
	return &IntegratedUI{
		components: []string{"map", "flow", "menu"},
		width:      100,
		height:     30,
	}
}

// Initialize initializes the integrated UI with visualization data
func (ui *IntegratedUI) Initialize(data *VisualizationData) error {
	if data == nil {
		return fmt.Errorf("visualization data cannot be nil")
	}
	
	ui.mutex.Lock()
	defer ui.mutex.Unlock()
	
	ui.data = data
	
	// Create layout manager
	ui.layoutManager = NewLayoutManager()
	
	// Create visualization components
	ui.treePrinter = NewTreePrinter()
	ui.areaPrinter = NewAreaPrinter()
	
	// Initialize components with data
	if err := ui.treePrinter.Initialize(data); err != nil {
		return fmt.Errorf("failed to initialize tree printer: %w", err)
	}
	
	if err := ui.areaPrinter.Initialize(data); err != nil {
		return fmt.Errorf("failed to initialize area printer: %w", err)
	}
	
	// Set initial focus
	ui.focusedComponent = "map"
	
	return nil
}

// SetupLayout sets up the 3-split layout
func (ui *IntegratedUI) SetupLayout() error {
	ui.mutex.Lock()
	defer ui.mutex.Unlock()
	
	if ui.layoutManager == nil {
		return fmt.Errorf("layout manager not initialized")
	}
	
	// Get default layout areas
	areas := GetDefaultAreas()
	
	// Add area printer (map) to left side
	if err := ui.layoutManager.AddComponent("map", ui.areaPrinter, areas["map"]); err != nil {
		return fmt.Errorf("failed to add area printer: %w", err)
	}
	
	// Add tree printer (flow) to right side
	if err := ui.layoutManager.AddComponent("flow", ui.treePrinter, areas["flow"]); err != nil {
		return fmt.Errorf("failed to add tree printer: %w", err)
	}
	
	// Menu area is reserved for future menu component
	// For now, we'll just mark it as allocated
	
	return nil
}

// UpdateComponent updates a specific component
func (ui *IntegratedUI) UpdateComponent(component string, event VisualizationEvent, data *VisualizationData) error {
	ui.mutex.RLock()
	defer ui.mutex.RUnlock()
	
	if ui.layoutManager == nil {
		return fmt.Errorf("layout manager not initialized")
	}
	
	// Update specific component through layout manager
	switch component {
	case "map":
		return ui.areaPrinter.Update(event, data)
	case "flow":
		return ui.treePrinter.Update(event, data)
	case "all":
		return ui.layoutManager.Update(event, data)
	default:
		return fmt.Errorf("unknown component: %s", component)
	}
}

// RenderAll renders all components
func (ui *IntegratedUI) RenderAll() (string, error) {
	ui.mutex.RLock()
	defer ui.mutex.RUnlock()
	
	if ui.layoutManager == nil {
		return "", fmt.Errorf("layout manager not initialized")
	}
	
	// For testing purposes, we'll render each component and combine them
	var result strings.Builder
	
	// Render map
	mapContent, err := ui.areaPrinter.Render()
	if err != nil {
		return "", fmt.Errorf("failed to render map: %w", err)
	}
	
	// Render flow
	flowContent, err := ui.treePrinter.Render()
	if err != nil {
		return "", fmt.Errorf("failed to render flow: %w", err)
	}
	
	// Combine with layout information
	result.WriteString("=== Map ===\n")
	result.WriteString(mapContent)
	result.WriteString("\n\n=== Flow ===\n")
	result.WriteString(flowContent)
	result.WriteString("\n\n=== Menu ===\n")
	result.WriteString(fmt.Sprintf("Focus: %s | Tab: Next | Shift+Tab: Previous | Enter: Select", ui.focusedComponent))
	
	return result.String(), nil
}

// HandleFocusChange handles focus changes between components
func (ui *IntegratedUI) HandleFocusChange(component string) error {
	ui.mutex.Lock()
	defer ui.mutex.Unlock()
	
	// Validate component
	found := false
	for _, c := range ui.components {
		if c == component {
			found = true
			break
		}
	}
	
	if !found {
		return fmt.Errorf("invalid component: %s", component)
	}
	
	ui.focusedComponent = component
	return nil
}

// SyncSelection synchronizes selection across components
func (ui *IntegratedUI) SyncSelection(paragraphNumber int) error {
	ui.mutex.Lock()
	defer ui.mutex.Unlock()
	
	if ui.data == nil {
		return fmt.Errorf("no visualization data available")
	}
	
	// Update current position in flow data
	if ui.data.FlowData != nil {
		ui.updateFlowDataCurrent(ui.data.FlowData.Root, paragraphNumber)
	}
	
	// Update current position in map data
	if ui.data.MapData != nil {
		ui.updateMapDataCurrent(paragraphNumber)
	}
	
	// Update both components
	if err := ui.treePrinter.Update(EventCurrentChanged, ui.data); err != nil {
		return fmt.Errorf("failed to update tree printer: %w", err)
	}
	
	if err := ui.areaPrinter.Update(EventCurrentChanged, ui.data); err != nil {
		return fmt.Errorf("failed to update area printer: %w", err)
	}
	
	return nil
}

// updateFlowDataCurrent recursively updates current position in flow data
func (ui *IntegratedUI) updateFlowDataCurrent(node *FlowNode, targetNumber int) {
	if node == nil {
		return
	}
	
	// Update current flag
	node.IsCurrent = (node.ParagraphNumber == targetNumber)
	
	// Recursively update children
	for _, child := range node.Children {
		ui.updateFlowDataCurrent(child, targetNumber)
	}
}

// updateMapDataCurrent updates current position in map data
func (ui *IntegratedUI) updateMapDataCurrent(targetNumber int) {
	if ui.data.MapData == nil {
		return
	}
	
	// Clear all current flags
	for i := range ui.data.MapData.Grid {
		for j := range ui.data.MapData.Grid[i] {
			ui.data.MapData.Grid[i][j].IsCurrent = false
		}
	}
	
	// Set new current position
	if pos, exists := ui.data.MapData.Positions[targetNumber]; exists {
		if pos.Y < len(ui.data.MapData.Grid) && pos.X < len(ui.data.MapData.Grid[pos.Y]) {
			ui.data.MapData.Grid[pos.Y][pos.X].IsCurrent = true
		}
	}
}

// HandleKeyInput handles keyboard input
func (ui *IntegratedUI) HandleKeyInput(key string) ComponentAction {
	switch key {
	case "tab":
		return ActionFocusNext
	case "shift+tab":
		return ActionFocusPrevious
	case "enter":
		return ActionSelect
	case "escape":
		return ActionCancel
	case "up", "down", "left", "right":
		return ActionNavigate
	default:
		return ActionNone
	}
}

// ResizeLayout handles layout resize
func (ui *IntegratedUI) ResizeLayout(width, height int) error {
	ui.mutex.Lock()
	defer ui.mutex.Unlock()
	
	ui.width = width
	ui.height = height
	
	if ui.layoutManager != nil {
		return ui.layoutManager.Resize(width, height)
	}
	
	return nil
}

// GetFocusedComponent returns the currently focused component
func (ui *IntegratedUI) GetFocusedComponent() string {
	ui.mutex.RLock()
	defer ui.mutex.RUnlock()
	
	return ui.focusedComponent
}

// CycleFocus cycles focus to the next component
func (ui *IntegratedUI) CycleFocus(forward bool) {
	ui.mutex.Lock()
	defer ui.mutex.Unlock()
	
	// Find current focus index
	currentIndex := -1
	for i, c := range ui.components {
		if c == ui.focusedComponent {
			currentIndex = i
			break
		}
	}
	
	if currentIndex == -1 {
		// Default to first component
		ui.focusedComponent = ui.components[0]
		return
	}
	
	// Calculate next index
	var nextIndex int
	if forward {
		nextIndex = (currentIndex + 1) % len(ui.components)
	} else {
		nextIndex = (currentIndex - 1 + len(ui.components)) % len(ui.components)
	}
	
	ui.focusedComponent = ui.components[nextIndex]
}