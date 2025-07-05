package main

import (
	"fmt"

	"github.com/pterm/pterm"
	"github.com/pterm/pterm/putils"
)

// TreePrinter implements IVisualizer for tree-based flow visualization
type TreePrinter struct {
	data   *VisualizationData
	width  int
	height int
	tree   *pterm.TreePrinter
}

// NewTreePrinter creates a new TreePrinter instance
func NewTreePrinter() *TreePrinter {
	return &TreePrinter{
		tree: &pterm.TreePrinter{},
	}
}

// Initialize initializes the TreePrinter with visualization data
func (tp *TreePrinter) Initialize(data *VisualizationData) error {
	if data == nil {
		return fmt.Errorf("visualization data cannot be nil")
	}
	
	tp.data = data
	return nil
}

// Update processes data updates and refreshes the tree visualization
func (tp *TreePrinter) Update(event VisualizationEvent, data *VisualizationData) error {
	if data == nil {
		return fmt.Errorf("visualization data cannot be nil")
	}
	
	tp.data = data
	return nil
}

// Render generates the tree visualization as a string
func (tp *TreePrinter) Render() (string, error) {
	if tp.data == nil || tp.data.FlowData == nil {
		return "", fmt.Errorf("no flow data available")
	}
	
	tree := tp.convertFlowDataToTree(tp.data.FlowData)
	if tree == nil {
		return "", fmt.Errorf("failed to convert flow data to tree")
	}
	
	// Use LeveledList to render the tree
	result, err := pterm.DefaultTree.WithRoot(putils.TreeFromLeveledList(*tree)).Srender()
	if err != nil {
		return "", fmt.Errorf("failed to render tree: %w", err)
	}
	
	return result, nil
}

// GetArea returns the current drawing area dimensions
func (tp *TreePrinter) GetArea() (width, height int) {
	return tp.width, tp.height
}

// SetArea sets the drawing area dimensions
func (tp *TreePrinter) SetArea(width, height int) error {
	if width < 0 || height < 0 {
		return fmt.Errorf("width and height must be non-negative")
	}
	
	tp.width = width
	tp.height = height
	return nil
}

// convertFlowDataToTree converts FlowData to PTerm LeveledList format
func (tp *TreePrinter) convertFlowDataToTree(flowData *FlowData) *pterm.LeveledList {
	if flowData == nil || flowData.Root == nil {
		return nil
	}
	
	tree := pterm.LeveledList{}
	tp.addNodeToTree(&tree, flowData.Root, 0)
	return &tree
}

// addNodeToTree recursively adds nodes to the tree structure
func (tp *TreePrinter) addNodeToTree(tree *pterm.LeveledList, node *FlowNode, level int) {
	if node == nil {
		return
	}
	
	// Create node text with styling
	nodeText := tp.formatNodeText(node)
	
	// Add current node to tree
	*tree = append(*tree, pterm.LeveledListItem{
		Level: level,
		Text:  nodeText,
	})
	
	// Recursively add children
	for _, child := range node.Children {
		tp.addNodeToTree(tree, child, level+1)
	}
}

// formatNodeText formats the node text with appropriate styling
func (tp *TreePrinter) formatNodeText(node *FlowNode) string {
	if node == nil {
		return ""
	}
	
	// Base text format
	text := fmt.Sprintf("%d: %s", node.ParagraphNumber, node.Description)
	
	// Apply styling based on node state
	style := tp.getNodeStyle(node)
	return style.Sprint(text)
}

// getNodeStyle returns the appropriate style for a node based on its state
func (tp *TreePrinter) getNodeStyle(node *FlowNode) *pterm.Style {
	if node == nil {
		return pterm.NewStyle()
	}
	
	// Current position gets highest priority styling
	if node.IsCurrent {
		return pterm.NewStyle(pterm.FgYellow, pterm.Bold, pterm.BgBlue)
	}
	
	// Visited nodes get secondary styling
	if node.Visited {
		return pterm.NewStyle(pterm.FgGreen)
	}
	
	// Unvisited nodes get default styling
	return pterm.NewStyle(pterm.FgLightWhite)
}

