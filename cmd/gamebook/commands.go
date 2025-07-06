package main

import (
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/zensai3805/iwakero-gamebook-mapping/internal/domain"
)

const (
	// 選択肢ステータス定数
	StatusUnselected = "未選択"
	StatusSelected   = "選択済み"

	// 状態表示シンボル定数
	StatusSymbolUnvisited = "⚪"
	StatusSymbolVisited   = "✅"
)

// CommandExecutor は既存のCLIコマンドを実行するためのインターフェース
type CommandExecutor interface {
	ExecuteNewCommand(title string) error
	ExecuteLoadCommand(title string) error
	ExecuteAddCommand(number int, description string) error
	ExecuteChoiceCommand(paragraphNum int, description string, targetNum int) error
	ExecuteSelectCommand(paragraphNum int, choiceIndex int) error
	ExecuteMoveCommand(targetNum int) error
	ExecuteShowCommand() error
	ExecuteShowCommandWithScope(scope DisplayScope) error
}

// CLIExecutor は実際のCLIコマンドを実行する
type CLIExecutor struct{}

// NewCLIExecutor は新しいCLIExecutorを作成
func NewCLIExecutor() *CLIExecutor {
	return &CLIExecutor{}
}

// ExecuteNewCommand newコマンドの実装を実行
func (e *CLIExecutor) ExecuteNewCommand(title string) error {
	// ユーザー操作記録
	LogUserOperation("new_gamebook", map[string]interface{}{
		"title": title,
		"title_length": len(title),
	})
	
	// 入力値検証
	if title == "" {
		LogValidationError("title", title, "空のタイトル", map[string]interface{}{"command": "new"})
		return fmt.Errorf("エラー: タイトルが空です")
	}
	
	if len(title) > 100 {
		LogValidationError("title", title, "タイトルが長すぎる", map[string]interface{}{"command": "new", "max_length": 100})
		return fmt.Errorf("エラー: タイトルが長すぎます（100文字以内）")
	}

	// 既存のnewコマンドロジックを呼び出し
	currentGame = domain.NewGamebook(title)
	if err := repo.Save(currentGame); err != nil {
		LogErrorWithContext(err, "new_command_save", map[string]interface{}{"title": title})
		LogCommandResult("new_gamebook", false, map[string]interface{}{"title": title, "error": err.Error()})
		return fmt.Errorf("エラー: %v", err)
	}

	// 現在のゲームとして保存
	if err := sessionRepo.SaveCurrentGame(title); err != nil {
		fmt.Fprintf(os.Stderr, "セッション保存エラー: %v\n", err)
		LogErrorWithContext(err, "new_command_session_save", map[string]interface{}{"title": title})
		if logger := GetGlobalLogger(); logger != nil {
			logger.Warn("セッション保存エラー", domain.Field{Key: "error", Value: err.Error()})
		}
	}

	fmt.Printf("新しいゲームブック「%s」を作成しました。\n", title)
	LogCommandResult("new_gamebook", true, map[string]interface{}{"title": title})
	
	// ログにも記録
	if logger := GetGlobalLogger(); logger != nil {
		logger.Info("新しいゲームブック作成", domain.Field{Key: "title", Value: title})
	}
	return nil
}

// ExecuteLoadCommand loadコマンドの実装を実行
func (e *CLIExecutor) ExecuteLoadCommand(title string) error {
	gamebook, err := repo.Load(title)
	if err != nil {
		return fmt.Errorf("ゲームブックの読み込みに失敗しました: %v", err)
	}
	currentGame = gamebook

	if err := sessionRepo.SaveCurrentGame(title); err != nil {
		fmt.Fprintf(os.Stderr, "セッション保存エラー: %v\n", err)
		// ログにも記録
		if logger := GetGlobalLogger(); logger != nil {
			logger.Warn("セッション保存エラー", domain.Field{Key: "error", Value: err.Error()})
		}
	}

	fmt.Printf("ゲームブック '%s' を読み込みました\n", title)
	// ログにも記録
	if logger := GetGlobalLogger(); logger != nil {
		logger.Info("ゲームブック読み込み", domain.Field{Key: "title", Value: title})
	}
	return nil
}

// ExecuteAddCommand addコマンドの実装を実行
func (e *CLIExecutor) ExecuteAddCommand(number int, description string) error {
	// ユーザー操作記録
	context := map[string]interface{}{
		"paragraph_number": number,
		"description_length": len(description),
		"has_current_game": currentGame != nil,
	}
	
	if currentGame != nil {
		context["total_paragraphs"] = len(currentGame.Paragraphs)
	}
	
	LogUserOperation("add_paragraph", context)
	
	if currentGame == nil {
		LogValidationError("current_game", nil, "ゲームブック未選択", context)
		return fmt.Errorf("エラー: ゲームブックが選択されていません。'gamebook new'または'gamebook load'を実行してください。")
	}

	// 入力値検証
	if number <= 0 {
		LogValidationError("paragraph_number", number, "無効なパラグラフ番号", context)
		return fmt.Errorf("エラー: パラグラフ番号は1以上である必要があります")
	}
	
	if description == "" {
		LogValidationError("description", description, "空の説明", context)
		return fmt.Errorf("エラー: パラグラフの説明が空です")
	}

	// 既存のパラグラフを確認
	existing, exists := currentGame.Paragraphs[number]
	isPlaceholder := exists && existing.Description == "(未定義)"
	
	context["paragraph_exists"] = exists
	context["is_placeholder"] = isPlaceholder

	p := domain.NewParagraph(number, description)
	if err := currentGame.AddParagraph(p); err != nil {
		LogErrorWithContext(err, "add_command_paragraph", context)
		LogCommandResult("add_paragraph", false, map[string]interface{}{"number": number, "error": err.Error()})
		return fmt.Errorf("エラー: %v", err)
	}

	if err := repo.Save(currentGame); err != nil {
		LogErrorWithContext(err, "add_command_save", context)
		LogCommandResult("add_paragraph", false, map[string]interface{}{"number": number, "error": err.Error()})
		return fmt.Errorf("保存エラー: %v", err)
	}

	if isPlaceholder {
		fmt.Printf("パラグラフ %d を更新しました: %s（プレースホルダーから更新）\n", number, description)
		LogCommandResult("add_paragraph", true, map[string]interface{}{"number": number, "type": "placeholder_update"})
		// ログにも記録
		if logger := GetGlobalLogger(); logger != nil {
			logger.Info("パラグラフ更新",
				domain.Field{Key: "number", Value: number},
				domain.Field{Key: "description", Value: description},
				domain.Field{Key: "type", Value: "placeholder_update"})
		}
	} else {
		fmt.Printf("パラグラフ %d を追加しました: %s\n", number, description)
		LogCommandResult("add_paragraph", true, map[string]interface{}{"number": number, "type": "new"})
		// ログにも記録
		if logger := GetGlobalLogger(); logger != nil {
			logger.Info("パラグラフ追加",
				domain.Field{Key: "number", Value: number},
				domain.Field{Key: "description", Value: description})
		}
	}
	return nil
}

// ExecuteChoiceCommand choiceコマンドの実装を実行
func (e *CLIExecutor) ExecuteChoiceCommand(paragraphNum int, description string, targetNum int) error {
	if currentGame == nil {
		return fmt.Errorf("エラー: ゲームブックが選択されていません。")
	}

	// 選択肢を追加
	if err := currentGame.AddChoiceToParagraph(paragraphNum, description, targetNum); err != nil {
		return fmt.Errorf("エラー: %v", err)
	}

	// 保存
	if err := repo.Save(currentGame); err != nil {
		return fmt.Errorf("保存エラー: %v", err)
	}

	fmt.Printf("パラグラフ %d に選択肢「%s → %d」を追加しました。\n", paragraphNum, description, targetNum)
	// ログにも記録
	if logger := GetGlobalLogger(); logger != nil {
		logger.Info("選択肢追加",
			domain.Field{Key: "paragraph", Value: paragraphNum},
			domain.Field{Key: "description", Value: description},
			domain.Field{Key: "target", Value: targetNum})
	}
	return nil
}

// ExecuteSelectCommand selectコマンドの実装を実行
func (e *CLIExecutor) ExecuteSelectCommand(paragraphNum int, choiceIndex int) error {
	if currentGame == nil {
		return fmt.Errorf("エラー: ゲームブックが選択されていません。")
	}

	// 優雅なエラーハンドリングを使用して選択肢を選択・移動
	moveResult := currentGame.SelectChoiceAndMoveWithGracefulHandling(paragraphNum, choiceIndex)

	if !moveResult.Success {
		return fmt.Errorf("エラー: %s", moveResult.WarningMessage)
	}

	// 保存
	if saveErr := repo.Save(currentGame); saveErr != nil {
		return fmt.Errorf("保存エラー: %w", saveErr)
	}

	// 結果の表示
	if moveResult.HasWarning {
		fmt.Printf("⚠️  警告: %s\n", moveResult.WarningMessage)
		fmt.Printf("パラグラフ %d の選択肢 %d を選択し、パラグラフ %d に移動しました。\n",
			paragraphNum, choiceIndex+1, currentGame.Current.Number)
	} else {
		fmt.Printf("パラグラフ %d の選択肢 %d を選択し、パラグラフ %d に移動しました。\n",
			paragraphNum, choiceIndex+1, currentGame.Current.Number)
	}

	return nil
}

// ExecuteMoveCommand moveコマンドの実装を実行
func (e *CLIExecutor) ExecuteMoveCommand(targetNum int) error {
	if currentGame == nil {
		return fmt.Errorf("エラー: ゲームブックが選択されていません。")
	}

	// 直接移動を実行
	moveErr := currentGame.MoveToWithPathSelection(targetNum)
	if moveErr != nil {
		return fmt.Errorf("エラー: %v", moveErr)
	}

	// 保存
	if saveErr := repo.Save(currentGame); saveErr != nil {
		return fmt.Errorf("保存エラー: %w", saveErr)
	}

	// 結果の表示
	fmt.Printf("パラグラフ %d に移動しました。\n", targetNum)
	return nil
}

// ExecuteShowCommand showコマンドの実装を実行
func (e *CLIExecutor) ExecuteShowCommand() error {
	if currentGame == nil {
		return fmt.Errorf("エラー: ゲームブックが選択されていません。")
	}

	fmt.Printf("=== %s ===\n", currentGame.Title)

	// 統計情報（簡潔な1行表示）
	stats := currentGame.GetExplorationStats()
	explorationRate := 0
	if stats.TotalParagraphs > 0 {
		explorationRate = (stats.VisitedParagraphs * 100) / stats.TotalParagraphs
	}
	fmt.Printf("📊 パラグラフ %d/%d | 選択肢 %d/%d | 探索率: %d%%\n",
		stats.VisitedParagraphs, stats.TotalParagraphs,
		stats.SelectedChoices, stats.TotalChoices,
		explorationRate)

	// 現在位置（簡潔表示）
	if currentGame.Current != nil {
		fmt.Printf("\n📍 現在: #%d %s\n", currentGame.Current.Number, currentGame.Current.Description)

		if len(currentGame.Current.Choices) > 0 {
			choicesDisplay := formatChoicesDisplay(currentGame.Current.Choices)
			fmt.Println(choicesDisplay)
		}
	}

	// v0.2.0可視化機能を統合
	if len(currentGame.Paragraphs) > 0 {
		fmt.Println("\n🗺️  可視化マップ & フロー図:")
		if err := e.showVisualization(); err != nil {
			fmt.Printf("可視化エラー: %v\n", err)
		}
	}

	// 保留参照情報の表示
	pendingTargets := currentGame.GetAllPendingTargets()
	if len(pendingTargets) > 0 {
		fmt.Printf("\n⏳ 未定義段落への参照: ")
		for i, target := range pendingTargets {
			if i > 0 {
				fmt.Printf(", ")
			}
			fmt.Printf("#%d", target)
		}
		fmt.Println()
		fmt.Println("   上記の段落を追加すると、参照が自動的に解決されます。")
	}

	// 他のパラグラフ（現在位置以外の最近追加分）
	fmt.Println("\n📚 その他のパラグラフ:")

	// パラグラフキーを収集してソート
	keys := make([]int, 0, len(currentGame.Paragraphs))
	currentNum := 0
	if currentGame.Current != nil {
		currentNum = currentGame.Current.Number
	}

	for key := range currentGame.Paragraphs {
		if key != currentNum {
			keys = append(keys, key)
		}
	}
	sort.Ints(keys)

	// ソートされたキーで最大3件表示
	displayed := 0
	for _, num := range keys {
		if displayed >= 3 {
			break
		}

		p := currentGame.Paragraphs[num]
		status := StatusSymbolUnvisited
		if p.Visited {
			status = StatusSymbolVisited
		}
		choiceCount := len(p.Choices)
		selectedCount := 0
		for _, choice := range p.Choices {
			if choice.Selected {
				selectedCount++
			}
		}

		if choiceCount > 0 {
			fmt.Printf("  %s #%d %s (選択肢 %d/%d)\n", status, p.Number, p.Description, selectedCount, choiceCount)
		} else {
			fmt.Printf("  %s #%d %s\n", status, p.Number, p.Description)
		}
		displayed++
	}
	return nil
}

// showVisualization v0.2.3 フィルタリング機能付きフロー図を表示
func (e *CLIExecutor) showVisualization() error {
	// データ変換
	converter := NewDataConverter()
	visualData, err := converter.ConvertToVisualizationData(currentGame)
	if err != nil {
		return fmt.Errorf("可視化データ変換エラー: %w", err)
	}

	// フィルタリング機能付きフロー図を作成・表示
	treePrinter := NewTreePrinterWithFilter()
	treePrinter.SetGamebook(currentGame)
	if initErr := treePrinter.Initialize(visualData); initErr != nil {
		return fmt.Errorf("フロー図初期化エラー: %w", initErr)
	}

	// フロー図をレンダリング
	flowContent, renderErr := treePrinter.Render()
	if renderErr != nil {
		return fmt.Errorf("フロー図レンダリングエラー: %w", renderErr)
	}

	// スコープ情報付きでフロー図を表示
	scopeManager := NewDisplayScopeManager()
	fmt.Printf("=== フロー図 (未定義表示: %s) ===\n", scopeManager.GetScopeDescription())
	fmt.Println(flowContent)
	return nil
}

// ExecuteShowCommandWithScope スコープ指定付きでshowコマンドの実装を実行
func (e *CLIExecutor) ExecuteShowCommandWithScope(scope DisplayScope) error {
	if currentGame == nil {
		return fmt.Errorf("エラー: ゲームブックが選択されていません。")
	}

	fmt.Printf("=== %s ===\n", currentGame.Title)

	// 統計情報（簡潔な1行表示）
	stats := currentGame.GetExplorationStats()
	explorationRate := 0
	if stats.TotalParagraphs > 0 {
		explorationRate = (stats.VisitedParagraphs * 100) / stats.TotalParagraphs
	}
	fmt.Printf("📊 パラグラフ %d/%d | 選択肢 %d/%d | 探索率: %d%%\n",
		stats.VisitedParagraphs, stats.TotalParagraphs,
		stats.SelectedChoices, stats.TotalChoices,
		explorationRate)

	// 現在位置（簡潔表示）
	if currentGame.Current != nil {
		fmt.Printf("\n📍 現在: #%d %s\n", currentGame.Current.Number, currentGame.Current.Description)

		if len(currentGame.Current.Choices) > 0 {
			choicesDisplay := formatChoicesDisplay(currentGame.Current.Choices)
			fmt.Println(choicesDisplay)
		}
	}

	// v0.2.3可視化機能（フィルタリング対応）を統合
	if len(currentGame.Paragraphs) > 0 {
		fmt.Println("\n🗺️  可視化マップ & フロー図:")
		if err := e.showVisualizationWithScope(scope); err != nil {
			fmt.Printf("可視化エラー: %v\n", err)
		}
	}

	// 保留参照情報の表示
	pendingTargets := currentGame.GetAllPendingTargets()
	if len(pendingTargets) > 0 {
		fmt.Printf("\n⏳ 未定義段落への参照: ")
		for i, target := range pendingTargets {
			if i > 0 {
				fmt.Printf(", ")
			}
			fmt.Printf("#%d", target)
		}
		fmt.Println()
		fmt.Println("   上記の段落を追加すると、参照が自動的に解決されます。")
	}

	// 他のパラグラフ（現在位置以外の最近追加分）
	fmt.Println("\n📚 その他のパラグラフ:")

	// パラグラフキーを収集してソート
	keys := make([]int, 0, len(currentGame.Paragraphs))
	currentNum := 0
	if currentGame.Current != nil {
		currentNum = currentGame.Current.Number
	}

	for key := range currentGame.Paragraphs {
		if key != currentNum {
			keys = append(keys, key)
		}
	}
	sort.Ints(keys)

	// ソートされたキーで最大3件表示
	displayed := 0
	for _, num := range keys {
		if displayed >= 3 {
			break
		}

		p := currentGame.Paragraphs[num]
		status := StatusSymbolUnvisited
		if p.Visited {
			status = StatusSymbolVisited
		}
		choiceCount := len(p.Choices)
		selectedCount := 0
		for _, choice := range p.Choices {
			if choice.Selected {
				selectedCount++
			}
		}

		if choiceCount > 0 {
			fmt.Printf("  %s #%d %s (選択肢 %d/%d)\n", status, p.Number, p.Description, selectedCount, choiceCount)
		} else {
			fmt.Printf("  %s #%d %s\n", status, p.Number, p.Description)
		}
		displayed++
	}
	return nil
}

// showVisualizationWithScope スコープ指定付きでフロー図を表示
func (e *CLIExecutor) showVisualizationWithScope(scope DisplayScope) error {
	// データ変換
	converter := NewDataConverter()
	visualData, err := converter.ConvertToVisualizationData(currentGame)
	if err != nil {
		return fmt.Errorf("可視化データ変換エラー: %w", err)
	}

	// フィルタリング機能付きフロー図を作成・表示
	treePrinter := NewTreePrinterWithFilter()
	treePrinter.SetDisplayScope(scope)
	treePrinter.SetGamebook(currentGame)
	if initErr := treePrinter.Initialize(visualData); initErr != nil {
		return fmt.Errorf("フロー図初期化エラー: %w", initErr)
	}

	// フロー図をレンダリング
	flowContent, renderErr := treePrinter.Render()
	if renderErr != nil {
		return fmt.Errorf("フロー図レンダリングエラー: %w", renderErr)
	}

	// スコープ情報付きでフロー図を表示
	scopeManager := NewDisplayScopeManager()
	scopeManager.SetScope(scope)
	fmt.Printf("=== フロー図 (未定義表示: %s) ===\n", scopeManager.GetScopeDescription())
	fmt.Println(flowContent)
	return nil
}

// ParseArguments は引数をパースするヘルパー関数
func ParseArguments(args []string, minArgs int, usage string) ([]string, error) {
	if len(args) < minArgs {
		return nil, fmt.Errorf("使用法: %s", usage)
	}
	return args, nil
}

// ParseNumber は文字列を数値にパースするヘルパー関数
func ParseNumber(str string, name string) (int, error) {
	num, err := strconv.Atoi(str)
	if err != nil {
		return 0, fmt.Errorf("%sは数値で指定してください: %v", name, err)
	}
	return num, nil
}

// JoinDescription は複数の引数を説明文として結合するヘルパー関数
func JoinDescription(args []string, startIndex int) string {
	if startIndex >= len(args) {
		return ""
	}
	return strings.Join(args[startIndex:], " ")
}

// formatChoicesDisplay は選択肢を改善されたフォーマットで表示する
func formatChoicesDisplay(choices []domain.Choice) string {
	if len(choices) == 0 {
		return ""
	}

	var result strings.Builder
	result.WriteString("🎯 選択肢:")

	for i, choice := range choices {
		// 選択済みかどうかでシンボルを決定
		statusSymbol := StatusSymbolUnvisited
		if choice.Selected {
			statusSymbol = StatusSymbolVisited
		}

		// 選択肢の一行を組み立て
		choiceLine := fmt.Sprintf("\n  [%d] %s %s → #%d",
			i+1, statusSymbol, choice.Description, choice.TargetNumber)
		result.WriteString(choiceLine)
	}

	return result.String()
}
