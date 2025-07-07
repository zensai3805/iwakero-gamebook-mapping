package main

import (
	"fmt"
	"sort"

	"github.com/pterm/pterm"
	"github.com/zensai3805/iwakero-gamebook-mapping/internal/domain"
)

// DataConverter はドメインモデルを可視化データに変換する
type DataConverter struct {
	// 設定可能なパラメータ
	mapWidth  int
	mapHeight int
}

// NewDataConverter は新しいDataConverterを作成
func NewDataConverter() *DataConverter {
	return &DataConverter{
		mapWidth:  20, // デフォルトマップ幅
		mapHeight: 15, // デフォルトマップ高さ
	}
}

// ConvertToVisualizationData GamebookをVisualizationDataに変換
func (dc *DataConverter) ConvertToVisualizationData(gamebook *domain.Gamebook) (*VisualizationData, error) {
	if gamebook == nil {
		return nil, fmt.Errorf("gamebook cannot be nil")
	}

	mapData := dc.convertToMapData(gamebook)
	flowData := dc.convertToFlowData(gamebook)

	return &VisualizationData{
		Gamebook:   gamebook,
		CurrentPos: gamebook.Current,
		MapData:    mapData,
		FlowData:   flowData,
		LastEvent:  EventGameLoaded,
	}, nil
}

// convertToMapData マップデータに変換
func (dc *DataConverter) convertToMapData(gamebook *domain.Gamebook) *MapData {
	// グリッドを初期化
	grid := make([][]MapCell, dc.mapHeight)
	for i := range grid {
		grid[i] = make([]MapCell, dc.mapWidth)
		for j := range grid[i] {
			grid[i][j] = MapCell{
				ParagraphNumbers: []int{},
				Symbol:           " ",
				Style:            pterm.NewStyle(),
				Visited:          false,
				IsCurrent:        false,
			}
		}
	}

	positions := make(map[int]Position)

	// パラグラフを自動配置
	paragraphNumbers := make([]int, 0, len(gamebook.Paragraphs))
	for num := range gamebook.Paragraphs {
		paragraphNumbers = append(paragraphNumbers, num)
	}
	sort.Ints(paragraphNumbers)

	// 簡単な配置アルゴリズム（螺旋状）
	centerX := dc.mapWidth / 2
	centerY := dc.mapHeight / 2
	x, y := centerX, centerY
	dx, dy := 0, -1

	for i, num := range paragraphNumbers {
		if 0 <= x && x < dc.mapWidth && 0 <= y && y < dc.mapHeight {
			paragraph := gamebook.Paragraphs[num]

			// セルを設定
			cell := &grid[y][x]
			cell.ParagraphNumbers = append(cell.ParagraphNumbers, num)
			cell.Symbol = fmt.Sprintf("%d", num)
			cell.Visited = paragraph.Visited
			cell.IsCurrent = (gamebook.Current != nil && gamebook.Current.Number == num)

			// スタイルを設定
			switch {
			case cell.IsCurrent:
				cell.Style = pterm.NewStyle(pterm.BgGreen, pterm.FgBlack)
			case cell.Visited:
				cell.Style = pterm.NewStyle(pterm.FgGreen)
			default:
				cell.Style = pterm.NewStyle(pterm.FgGray)
			}

			positions[num] = Position{X: x, Y: y}
		}

		// 螺旋移動の計算
		if x == centerX-dx && y == centerY-dy || x == centerX+dx && y == centerY+dy {
			dx, dy = -dy, dx
		}
		x, y = x+dx, y+dy

		// 境界チェック
		if i > 0 && (x < 0 || x >= dc.mapWidth || y < 0 || y >= dc.mapHeight) {
			break
		}
	}

	return &MapData{
		Grid:      grid,
		Width:     dc.mapWidth,
		Height:    dc.mapHeight,
		Positions: positions,
	}
}

// convertToFlowData フローデータに変換
func (dc *DataConverter) convertToFlowData(gamebook *domain.Gamebook) *FlowData {
	// 移動履歴から選択された経路を特定
	selectedPath := dc.calculateSelectedPath(gamebook.GetNavigationHistory())
	nodes := make([]FlowNode, 0, len(gamebook.Paragraphs))
	nodeMap := make(map[int]*FlowNode)
	edges := make([]FlowEdge, 0)

	// まず、全ての選択肢の遷移先を収集（未定義パラグラフ特定のため）
	undefinedTargets := make(map[int]bool)
	for _, paragraph := range gamebook.Paragraphs {
		for _, choice := range paragraph.Choices {
			if _, exists := gamebook.Paragraphs[choice.TargetNumber]; !exists {
				undefinedTargets[choice.TargetNumber] = true
			}
		}
	}

	// ノードを作成（定義済みパラグラフ）
	for num, paragraph := range gamebook.Paragraphs {
		// 選択肢情報を移動履歴に基づいて更新
		choices := make([]domain.Choice, len(paragraph.Choices))
		for i, choice := range paragraph.Choices {
			edgeKey := fmt.Sprintf("%d->%d", num, choice.TargetNumber)
			choices[i] = domain.Choice{
				Description:  choice.Description,
				TargetNumber: choice.TargetNumber,
				Selected:     selectedPath[edgeKey], // 移動履歴ベースの選択状態
			}
		}
		
		node := FlowNode{
			ParagraphNumber: num,
			Description:     paragraph.Description,
			Children:        []*FlowNode{},
			Choices:         choices, // 更新された選択肢情報
			Visited:         paragraph.Visited,
			IsCurrent:       (gamebook.Current != nil && gamebook.Current.Number == num),
			VisitCount:      1, // 基本的には1、後で拡張可能
		}

		// スタイルを設定
		switch {
		case node.IsCurrent:
			node.Style = pterm.NewStyle(pterm.BgGreen, pterm.FgBlack)
		case node.Visited:
			node.Style = pterm.NewStyle(pterm.FgGreen)
		default:
			node.Style = pterm.NewStyle(pterm.FgGray)
		}

		nodes = append(nodes, node)
		nodeMap[num] = &nodes[len(nodes)-1]
	}

	// 未定義パラグラフのノードを作成
	for targetNum := range undefinedTargets {
		node := FlowNode{
			ParagraphNumber: targetNum,
			Description:     "(未定義)",
			Children:        []*FlowNode{},
			Choices:         nil,
			Visited:         false,
			IsCurrent:       false,
			VisitCount:      0,
			Style:           pterm.NewStyle(pterm.FgRed), // 未定義は赤色で表示
		}
		nodes = append(nodes, node)
		nodeMap[targetNum] = &nodes[len(nodes)-1]
	}

	// エッジを作成
	for _, paragraph := range gamebook.Paragraphs {
		fromNode, exists := nodeMap[paragraph.Number]
		if !exists {
			continue
		}

		for _, choice := range paragraph.Choices {
			toNode, exists := nodeMap[choice.TargetNumber]
			if !exists {
				continue
			}

			// 親子関係を設定
			fromNode.Children = append(fromNode.Children, toNode)

			// 移動履歴に基づいて選択状態を判定
			edgeKey := fmt.Sprintf("%d->%d", paragraph.Number, choice.TargetNumber)
			isSelectedInHistory := selectedPath[edgeKey]

			// エッジを作成
			edge := FlowEdge{
				From:        fromNode,
				To:          toNode,
				Description: choice.Description,
				Selected:    isSelectedInHistory, // 移動履歴ベースの選択状態
			}

			// スタイルを設定
			if edge.Selected {
				edge.Style = pterm.NewStyle(pterm.FgGreen)
			} else {
				edge.Style = pterm.NewStyle(pterm.FgGray)
			}

			edges = append(edges, edge)
		}
	}

	// ルートノードを特定（最小番号）
	var root *FlowNode
	if len(nodes) > 0 {
		minNum := nodes[0].ParagraphNumber
		for i := range nodes {
			if nodes[i].ParagraphNumber < minNum {
				minNum = nodes[i].ParagraphNumber
			}
		}
		root = nodeMap[minNum]
	}

	return &FlowData{
		Nodes: nodes,
		Edges: edges,
		Root:  root,
	}
}

// calculateSelectedPath 移動履歴から選択された経路を計算
func (dc *DataConverter) calculateSelectedPath(history []domain.NavigationStep) map[string]bool {
	selectedPath := make(map[string]bool)
	
	for _, step := range history {
		// From->To の移動を記録
		key := fmt.Sprintf("%d->%d", step.From, step.To)
		selectedPath[key] = true
	}
	
	return selectedPath
}

// UpdateVisualizationData イベントに基づいてデータを更新
func (dc *DataConverter) UpdateVisualizationData(data *VisualizationData, event VisualizationEvent, eventData interface{}) error {
	if data == nil {
		return fmt.Errorf("visualization data cannot be nil")
	}

	data.LastEvent = event
	data.LastEventData = eventData

	// イベント種別に応じて更新
	switch event {
	case EventCurrentChanged:
		data.CurrentPos = data.Gamebook.Current
		return dc.updateCurrentPosition(data)
	case EventParagraphAdded, EventChoiceAdded, EventParagraphVisited, EventChoiceSelected:
		// データ全体を再変換
		newData, err := dc.ConvertToVisualizationData(data.Gamebook)
		if err != nil {
			return err
		}
		data.MapData = newData.MapData
		data.FlowData = newData.FlowData
		data.CurrentPos = newData.CurrentPos
	}

	return nil
}

// updateCurrentPosition 現在位置の更新
func (dc *DataConverter) updateCurrentPosition(data *VisualizationData) error {
	// マップデータの現在位置を更新
	for y := range data.MapData.Grid {
		for x := range data.MapData.Grid[y] {
			cell := &data.MapData.Grid[y][x]
			cell.IsCurrent = false

			// 現在位置をチェック
			if data.CurrentPos != nil {
				for _, num := range cell.ParagraphNumbers {
					if num == data.CurrentPos.Number {
						cell.IsCurrent = true
						cell.Style = pterm.NewStyle(pterm.BgGreen, pterm.FgBlack)
						break
					}
				}
			}

			// 現在位置でない場合はスタイルを復元
			if !cell.IsCurrent {
				switch {
				case cell.Visited:
					cell.Style = pterm.NewStyle(pterm.FgGreen)
				default:
					cell.Style = pterm.NewStyle(pterm.FgGray)
				}
			}
		}
	}

	// フローデータの現在位置を更新
	for i := range data.FlowData.Nodes {
		node := &data.FlowData.Nodes[i]
		node.IsCurrent = (data.CurrentPos != nil && node.ParagraphNumber == data.CurrentPos.Number)

		// スタイルを更新
		switch {
		case node.IsCurrent:
			node.Style = pterm.NewStyle(pterm.BgGreen, pterm.FgBlack)
		case node.Visited:
			node.Style = pterm.NewStyle(pterm.FgGreen)
		default:
			node.Style = pterm.NewStyle(pterm.FgGray)
		}
	}

	return nil
}
