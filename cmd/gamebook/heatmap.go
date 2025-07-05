package main

import (
	"fmt"
	"sort"
	"strings"
	"sync"

	"github.com/pterm/pterm"
)

// HeatmapManager manages visit frequency tracking and visualization
type HeatmapManager struct {
	visitCounts map[int]int
	maxVisits   int
	data        *VisualizationData
	mutex       sync.RWMutex
}

// NewHeatmapManager creates a new HeatmapManager instance
func NewHeatmapManager() *HeatmapManager {
	return &HeatmapManager{
		visitCounts: make(map[int]int),
		maxVisits:   0,
	}
}

// Initialize initializes the heatmap manager with visualization data
func (hm *HeatmapManager) Initialize(data *VisualizationData) error {
	if data == nil {
		return fmt.Errorf("visualization data cannot be nil")
	}
	
	hm.mutex.Lock()
	defer hm.mutex.Unlock()
	
	hm.data = data
	
	// Initialize visit counts from gamebook data
	if data.Gamebook != nil {
		for num, para := range data.Gamebook.Paragraphs {
			if para.Visited {
				hm.visitCounts[num] = 1
			} else {
				hm.visitCounts[num] = 0
			}
		}
	}
	
	hm.updateMaxVisits()
	return nil
}

// RecordVisit records a visit to a paragraph
func (hm *HeatmapManager) RecordVisit(paragraphNumber int) {
	hm.mutex.Lock()
	defer hm.mutex.Unlock()
	
	hm.visitCounts[paragraphNumber]++
	if hm.visitCounts[paragraphNumber] > hm.maxVisits {
		hm.maxVisits = hm.visitCounts[paragraphNumber]
	}
}

// GetVisitCount returns the visit count for a paragraph
func (hm *HeatmapManager) GetVisitCount(paragraphNumber int) int {
	hm.mutex.RLock()
	defer hm.mutex.RUnlock()
	
	return hm.visitCounts[paragraphNumber]
}

// GetHeatmapColor returns the appropriate color for a paragraph based on visit frequency
func (hm *HeatmapManager) GetHeatmapColor(paragraphNumber int) *pterm.Style {
	hm.mutex.RLock()
	defer hm.mutex.RUnlock()
	
	visits := hm.visitCounts[paragraphNumber]
	
	// Calculate intensity (0.0 to 1.0)
	var intensity float64
	if hm.maxVisits > 0 {
		intensity = float64(visits) / float64(hm.maxVisits)
	}
	
	return hm.getColorGradient(intensity)
}

// getColorGradient returns a color based on intensity (0.0 to 1.0)
func (hm *HeatmapManager) getColorGradient(intensity float64) *pterm.Style {
	// Color gradient: Gray (unvisited) -> Blue (low) -> Green (medium) -> Yellow (high) -> Red (max)
	if intensity == 0 {
		// Unvisited - gray
		return pterm.NewStyle(pterm.FgGray)
	} else if intensity <= 0.25 {
		// Low visits - blue
		return pterm.NewStyle(pterm.FgBlue)
	} else if intensity <= 0.5 {
		// Medium visits - green
		return pterm.NewStyle(pterm.FgGreen)
	} else if intensity <= 0.75 {
		// High visits - yellow
		return pterm.NewStyle(pterm.FgYellow)
	} else {
		// Maximum visits - red
		return pterm.NewStyle(pterm.FgRed, pterm.Bold)
	}
}

// ApplyHeatmapToFlow applies heatmap colors to flow data
func (hm *HeatmapManager) ApplyHeatmapToFlow(flowData *FlowData) error {
	if flowData == nil {
		return fmt.Errorf("flow data cannot be nil")
	}
	
	hm.mutex.RLock()
	defer hm.mutex.RUnlock()
	
	hm.applyHeatmapToFlowNode(flowData.Root)
	return nil
}

// applyHeatmapToFlowNode recursively applies heatmap colors to flow nodes
func (hm *HeatmapManager) applyHeatmapToFlowNode(node *FlowNode) {
	if node == nil {
		return
	}
	
	// Apply heatmap color to node
	node.Style = hm.GetHeatmapColor(node.ParagraphNumber)
	
	// Update visit count in node
	node.VisitCount = hm.visitCounts[node.ParagraphNumber]
	
	// Recursively apply to children
	for _, child := range node.Children {
		hm.applyHeatmapToFlowNode(child)
	}
}

// ApplyHeatmapToMap applies heatmap colors to map data
func (hm *HeatmapManager) ApplyHeatmapToMap(mapData *MapData) error {
	if mapData == nil {
		return fmt.Errorf("map data cannot be nil")
	}
	
	hm.mutex.RLock()
	defer hm.mutex.RUnlock()
	
	// Apply heatmap colors to all cells
	for i := range mapData.Grid {
		for j := range mapData.Grid[i] {
			cell := &mapData.Grid[i][j]
			if len(cell.ParagraphNumbers) > 0 {
				// Use the color of the first paragraph in the cell
				cell.Style = hm.GetHeatmapColor(cell.ParagraphNumbers[0])
			}
		}
	}
	
	return nil
}

// GenerateLegend generates a legend for the heatmap
func (hm *HeatmapManager) GenerateLegend() string {
	var legend strings.Builder
	
	legend.WriteString("=== Heatmap Legend ===\n")
	
	// Create legend entries with colors
	entries := []struct {
		label string
		style *pterm.Style
	}{
		{"Unvisited", pterm.NewStyle(pterm.FgGray)},
		{"Low", pterm.NewStyle(pterm.FgBlue)},
		{"Medium", pterm.NewStyle(pterm.FgGreen)},
		{"High", pterm.NewStyle(pterm.FgYellow)},
		{"Maximum", pterm.NewStyle(pterm.FgRed, pterm.Bold)},
	}
	
	for _, entry := range entries {
		legend.WriteString(entry.style.Sprint("■"))
		legend.WriteString(" ")
		legend.WriteString(entry.label)
		legend.WriteString("\n")
	}
	
	return legend.String()
}

// GenerateStatisticsPanel generates a statistics panel
func (hm *HeatmapManager) GenerateStatisticsPanel() string {
	hm.mutex.RLock()
	defer hm.mutex.RUnlock()
	
	var stats strings.Builder
	
	stats.WriteString("=== Exploration Statistics ===\n")
	
	// Calculate statistics
	totalParagraphs := 0
	visitedParagraphs := 0
	totalVisits := 0
	mostVisitedParagraph := 0
	mostVisits := 0
	
	if hm.data != nil && hm.data.Gamebook != nil {
		totalParagraphs = len(hm.data.Gamebook.Paragraphs)
	}
	
	for para, visits := range hm.visitCounts {
		if visits > 0 {
			visitedParagraphs++
			totalVisits += visits
		}
		if visits > mostVisits {
			mostVisits = visits
			mostVisitedParagraph = para
		}
	}
	
	// Calculate exploration percentage
	explorationPercentage := 0.0
	if totalParagraphs > 0 {
		explorationPercentage = float64(visitedParagraphs) / float64(totalParagraphs) * 100
	}
	
	// Build statistics display
	stats.WriteString(fmt.Sprintf("Total Paragraphs: %d\n", totalParagraphs))
	stats.WriteString(fmt.Sprintf("Visited: %d (%.1f%%)\n", visitedParagraphs, explorationPercentage))
	stats.WriteString(fmt.Sprintf("Total Visits: %d\n", totalVisits))
	
	if mostVisits > 0 {
		stats.WriteString(fmt.Sprintf("Most Visited: #%d (%d visits)\n", mostVisitedParagraph, mostVisits))
	}
	
	// Add visit distribution
	stats.WriteString("\nVisit Distribution:\n")
	distribution := hm.getVisitDistribution()
	for _, entry := range distribution {
		barLength := int(float64(entry.count) / float64(len(hm.visitCounts)) * 20)
		bar := strings.Repeat("█", barLength)
		stats.WriteString(fmt.Sprintf("%d visits: %s %d paragraphs\n", entry.visits, bar, entry.count))
	}
	
	return stats.String()
}

// visitDistributionEntry represents a visit count distribution entry
type visitDistributionEntry struct {
	visits int
	count  int
}

// getVisitDistribution returns the distribution of visit counts
func (hm *HeatmapManager) getVisitDistribution() []visitDistributionEntry {
	// Count paragraphs by visit count
	distribution := make(map[int]int)
	for _, visits := range hm.visitCounts {
		distribution[visits]++
	}
	
	// Convert to sorted slice
	var entries []visitDistributionEntry
	for visits, count := range distribution {
		entries = append(entries, visitDistributionEntry{visits: visits, count: count})
	}
	
	// Sort by visit count
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].visits < entries[j].visits
	})
	
	return entries
}

// updateMaxVisits updates the maximum visit count
func (hm *HeatmapManager) updateMaxVisits() {
	hm.maxVisits = 0
	for _, visits := range hm.visitCounts {
		if visits > hm.maxVisits {
			hm.maxVisits = visits
		}
	}
}