package main

import (
	"fmt"

	"github.com/iwapc/iwakero-gamebook-mapping/internal/domain"
	"github.com/pterm/pterm"
	"github.com/pterm/pterm/putils"
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
	}

	// 未定義パラグラフ分析を実行
	analyzer := domain.NewUndefinedAnalyzer(tp.gamebook)
	analysis := analyzer.AnalyzeUndefinedParagraphs(currentPos)

	// フィルターを作成
	tp.filter = NewDisplayFilter(tp.scope, analysis)
}

// convertFlowDataToTreeWithFilter はFlowDataをフィルタリング付きでPTerm LeveledList形式に変換する
func (tp *TreePrinterWithFilter) convertFlowDataToTreeWithFilter(flowData *FlowData) *pterm.LeveledList {
	if flowData == nil || flowData.Root == nil {
		return nil
	}

	tree := pterm.LeveledList{}
	tp.addNodeToTreeWithFilter(&tree, flowData.Root, 0)
	return &tree
}

// addNodeToTreeWithFilter はフィルタリング機能付きでノードをツリー構造に追加する
func (tp *TreePrinterWithFilter) addNodeToTreeWithFilter(tree *pterm.LeveledList, node *FlowNode, level int) {
	if node == nil {
		return
	}

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
		}
	} else {
		// フィルター未設定の場合は全ての選択肢を表示
		for _, choice := range node.Choices {
			choiceText := tp.formatChoiceText(choice)
			*tree = append(*tree, pterm.LeveledListItem{
				Level: level + 1,
				Text:  choiceText,
			})
		}
	}

	// 子ノードを再帰的に追加
	for _, child := range node.Children {
		tp.addNodeToTreeWithFilter(tree, child, level+1)
	}
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
