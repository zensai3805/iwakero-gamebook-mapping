package main

import (
	"testing"

	"github.com/iwapc/iwakero-gamebook-mapping/internal/domain"
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
	
	printer.Initialize(data)
	
	// Test update
	err := printer.Update(EventCurrentChanged, data)
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
	
	printer.Initialize(data)
	
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
	
	printer.Initialize(data)
	
	// This method should be accessible for testing
	tree := printer.convertFlowDataToTree(flowData)
	if tree == nil {
		t.Error("convertFlowDataToTree() returned nil")
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