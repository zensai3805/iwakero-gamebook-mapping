package main

import (
	"fmt"

	"github.com/pterm/pterm"
	"github.com/pterm/pterm/putils"
	"github.com/zensai3805/iwakero-gamebook-mapping/internal/domain"
)

// TreePrinterWithFilter はフィルタリング機能付きのTreePrinter
type TreePrinterWithFilter struct {
	data     *VisualizationData
	width    int
	height   int
	tree     *pterm.TreePrinter
	scope    DisplayScope
	gamebook *domain.Gamebook
	filter   *DisplayFilter
	flowData *FlowData // フローデータの参照を保持
}

// NewTreePrinterWithFilter は新しいフィルタリング機能付きTreePrinterを作成する
func NewTreePrinterWithFilter() *TreePrinterWithFilter {
	return &TreePrinterWithFilter{
		tree:  &pterm.TreePrinter{},
		scope: ScopeConnected, // デフォルトは接続された未定義パラグラフのみ表示
	}
}

// SetDisplayScope は表示スコープを設定する
func (tp *TreePrinterWithFilter) SetDisplayScope(scope DisplayScope) {
	tp.scope = scope
	tp.updateFilter()
}

// SetGamebook はゲームブックを設定する
func (tp *TreePrinterWithFilter) SetGamebook(gamebook *domain.Gamebook) {
	tp.gamebook = gamebook
	tp.updateFilter()
}

// Initialize はTreePrinterを可視化データで初期化する
func (tp *TreePrinterWithFilter) Initialize(data *VisualizationData) error {
	if data == nil {
		return fmt.Errorf("visualization data cannot be nil")
	}

	tp.data = data
	tp.updateFilter()
	return nil
}

// Update はデータ更新を処理し、ツリー可視化を更新する
func (tp *TreePrinterWithFilter) Update(_ VisualizationEvent, data *VisualizationData) error {
	if data == nil {
		return fmt.Errorf("visualization data cannot be nil")
	}

	tp.data = data
	tp.updateFilter()
	return nil
}

// Render はツリー可視化を文字列として生成する
func (tp *TreePrinterWithFilter) Render() (string, error) {
	if tp.data == nil || tp.data.FlowData == nil {
		return "", fmt.Errorf("no flow data available")
	}

	// デバッグ: FlowDataの内容を確認
	// fmt.Printf("DEBUG: FlowData nodes count: %d\n", len(tp.data.FlowData.Nodes))
	// for _, node := range tp.data.FlowData.Nodes {
	// 	fmt.Printf("DEBUG: Node %d: %s (choices: %d)\n", node.ParagraphNumber, node.Description, len(node.Choices))
	// }

	tp.flowData = tp.data.FlowData // flowDataを設定
	tree := tp.convertFlowDataToTreeWithFilter(tp.data.FlowData)
	if tree == nil {
		return "", fmt.Errorf("failed to convert flow data to tree")
	}

	// LeveledListを使用してツリーをレンダリング
	result, err := pterm.DefaultTree.WithRoot(putils.TreeFromLeveledList(*tree)).Srender()
	if err != nil {
		return "", fmt.Errorf("failed to render tree: %w", err)
	}

	return result, nil
}

// GetArea は現在の描画エリアサイズを返す
func (tp *TreePrinterWithFilter) GetArea() (width, height int) {
	return tp.width, tp.height
}

// SetArea は描画エリアサイズを設定する
func (tp *TreePrinterWithFilter) SetArea(width, height int) error {
	if width < 0 || height < 0 {
		return fmt.Errorf("width and height must be non-negative")
	}

	tp.width = width
	tp.height = height
	return nil
}

// updateFilter はフィルターを更新する
func (tp *TreePrinterWithFilter) updateFilter() {
	if tp.gamebook == nil {
		tp.filter = nil
		return
	}

	// 現在位置を取得
	currentPos := 0
	if tp.gamebook.Current != nil {
		currentPos = tp.gamebook.Current.Number
	} else if len(tp.gamebook.Paragraphs) > 0 {
		// 現在位置が未設定の場合、最小番号のパラグラフを開始位置とする
		minNum := -1
		for num := range tp.gamebook.Paragraphs {
			if minNum == -1 || num < minNum {
				minNum = num
			}
		}
		if minNum != -1 {
			currentPos = minNum
		}
	}

	// 未定義パラグラフ分析を実行
	analyzer := domain.NewUndefinedAnalyzer(tp.gamebook)
	analysis := analyzer.AnalyzeUndefinedParagraphs(currentPos)

	// デバッグ: 分析結果を確認
	// fmt.Printf("DEBUG: Current position: %d\n", currentPos)
	// fmt.Printf("DEBUG: Connected undefined: %v\n", analysis.Connected)
	// fmt.Printf("DEBUG: Orphaned undefined: %v\n", analysis.Orphaned)

	// フィルターを作成
	tp.filter = NewDisplayFilter(tp.scope, analysis)
}

// convertFlowDataToTreeWithFilter はFlowDataをフィルタリング付きでPTerm LeveledList形式に変換する
func (tp *TreePrinterWithFilter) convertFlowDataToTreeWithFilter(flowData *FlowData) *pterm.LeveledList {
	if flowData == nil || flowData.Root == nil {
		return nil
	}

	tp.flowData = flowData // フローデータを保存
	// fmt.Printf("DEBUG: flowData has %d nodes\n", len(flowData.Nodes))
	// for _, node := range flowData.Nodes {
	// 	fmt.Printf("DEBUG: - Node %d: %s\n", node.ParagraphNumber, node.Description)
	// }
	tree := pterm.LeveledList{}
	visited := make(map[int]bool)
	tp.addNodeToTreeWithFilterRecursive(&tree, flowData.Root, 0, visited)
	return &tree
}

// addNodeToTreeWithFilterRecursive はフィルタリング機能付きでノードをツリー構造に再帰的に追加する
func (tp *TreePrinterWithFilter) addNodeToTreeWithFilterRecursive(tree *pterm.LeveledList, node *FlowNode, level int, visited map[int]bool) {
	if node == nil {
		return
	}

	// 既に訪問済みの場合はスキップ（無限再帰防止）
	if visited[node.ParagraphNumber] {
		// fmt.Printf("DEBUG: Node %d already visited, skipping\n", node.ParagraphNumber)
		return
	}
	visited[node.ParagraphNumber] = true
	// fmt.Printf("DEBUG: Processing node %d at level %d\n", node.ParagraphNumber, level)

	// フィルタリングチェック（定義済みパラグラフは常に表示）
	if tp.filter != nil && !tp.filter.ShouldDisplayUndefined(tp.gamebook, node.ParagraphNumber) {
		return
	}

	// ノードテキストをスタイリング付きで作成
	nodeText := tp.formatNodeText(node)

	// 現在のノードをツリーに追加
	*tree = append(*tree, pterm.LeveledListItem{
		Level: level,
		Text:  nodeText,
	})

	// フィルタリング済みの選択肢を子アイテムとして追加
	if tp.filter != nil {
		filteredChoices := tp.filter.FilterChoices(tp.gamebook, node.Choices)
		for _, choice := range filteredChoices {
			choiceText := tp.formatChoiceText(choice)
			*tree = append(*tree, pterm.LeveledListItem{
				Level: level + 1,
				Text:  choiceText,
			})

			// 選択肢の遷移先ノードも表示（未定義パラグラフを含む）
			if targetNode := tp.findNodeByNumber(choice.TargetNumber); targetNode != nil {
				// fmt.Printf("DEBUG: Adding target node %d for choice '%s'\n", choice.TargetNumber, choice.Description)
				tp.addNodeToTreeWithFilterRecursive(tree, targetNode, level+2, visited)
			}
		}
	} else {
		// フィルター未設定の場合は全ての選択肢を表示
		for _, choice := range node.Choices {
			choiceText := tp.formatChoiceText(choice)
			*tree = append(*tree, pterm.LeveledListItem{
				Level: level + 1,
				Text:  choiceText,
			})

			// 選択肢の遷移先ノードも表示（未定義パラグラフを含む）
			if targetNode := tp.findNodeByNumber(choice.TargetNumber); targetNode != nil {
				// fmt.Printf("DEBUG: Adding target node %d for choice '%s'\n", choice.TargetNumber, choice.Description)
				tp.addNodeToTreeWithFilterRecursive(tree, targetNode, level+2, visited)
			}
		}
	}

	// 子ノードを再帰的に追加（ただし、選択肢から既に追加されたノードは除く）
	for _, child := range node.Children {
		// 選択肢の遷移先として既に追加されていないか確認
		isChoiceTarget := false
		for _, choice := range node.Choices {
			if choice.TargetNumber == child.ParagraphNumber {
				isChoiceTarget = true
				break
			}
		}
		if !isChoiceTarget {
			tp.addNodeToTreeWithFilterRecursive(tree, child, level+1, visited)
		}
	}
}

// findNodeByNumber はフローデータから指定された番号のノードを検索する
func (tp *TreePrinterWithFilter) findNodeByNumber(number int) *FlowNode {
	if tp.flowData == nil {
		// fmt.Printf("DEBUG: flowData is nil in findNodeByNumber\n")
		return nil
	}
	for i := range tp.flowData.Nodes {
		if tp.flowData.Nodes[i].ParagraphNumber == number {
			return &tp.flowData.Nodes[i]
		}
	}
	return nil
}

// formatNodeText はノードテキストを適切なスタイリング付きでフォーマットする
func (tp *TreePrinterWithFilter) formatNodeText(node *FlowNode) string {
	if node == nil {
		return ""
	}

	// ベーステキストフォーマット
	text := fmt.Sprintf("%d: %s", node.ParagraphNumber, node.Description)

	// ノード状態に基づくスタイリングを適用
	style := tp.getNodeStyle(node)
	return style.Sprint(text)
}

// getNodeStyle はノード状態に基づく適切なスタイルを返す
func (tp *TreePrinterWithFilter) getNodeStyle(node *FlowNode) *pterm.Style {
	if node == nil {
		return pterm.NewStyle()
	}

	// 現在位置が最優先のスタイリング
	if node.IsCurrent {
		return pterm.NewStyle(pterm.FgYellow, pterm.Bold, pterm.BgBlue)
	}

	// 訪問済みノードがセカンダリスタイリング
	if node.Visited {
		return pterm.NewStyle(pterm.FgGreen)
	}

	// 未訪問ノードがデフォルトスタイリング
	return pterm.NewStyle(pterm.FgLightWhite)
}

// formatChoiceText は選択情報を選択状態付きでフォーマットする
func (tp *TreePrinterWithFilter) formatChoiceText(choice domain.Choice) string {
	// 選択状態を示すシンボル
	symbol := "[ ]"
	if choice.Selected {
		symbol = "[✓]"
	}

	// 選択肢テキストをフォーマット
	choiceText := fmt.Sprintf("%s %s → %d", symbol, choice.Description, choice.TargetNumber)

	// スタイルを適用
	style := tp.getChoiceStyle(choice)
	return style.Sprint(choiceText)
}

// getChoiceStyle は選択状態に基づく適切なスタイルを返す
func (tp *TreePrinterWithFilter) getChoiceStyle(choice domain.Choice) *pterm.Style {
	if choice.Selected {
		return pterm.NewStyle(pterm.FgGreen)
	}
	return pterm.NewStyle(pterm.FgGray)
}
