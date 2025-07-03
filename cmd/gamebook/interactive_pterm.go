package main

import (
	"fmt"
	"strings"

	"github.com/pterm/pterm"
)

// PTermInteractiveShell PTerm対話モードのシェル
type PTermInteractiveShell struct {
	executor CommandExecutor
}

// NewPTermInteractiveShell PTerm対話シェルを作成
func NewPTermInteractiveShell() *PTermInteractiveShell {
	return &PTermInteractiveShell{
		executor: NewCLIExecutor(),
	}
}

// Run PTerm対話シェルを実行
func (s *PTermInteractiveShell) Run() {
	// ウェルカムメッセージを表示
	pterm.DefaultHeader.WithFullWidth().Println("🎮 Gamebook Interactive Mode")
	pterm.Info.Println("リッチなターミナルUI対話モードです")
	pterm.Info.Println("↑↓キーでメニュー選択、Enterで決定、Ctrl+Cで終了")
	pterm.Println()

	// 最後のゲームを自動ロード
	s.autoLoadLastGame()

	for {
		// 現在の状態を表示（ゲームが読み込まれている場合）
		if currentGame != nil {
			s.handleShow()
			pterm.Println() // 間隔をあける
		}

		// プロンプト表示
		prompt := s.getCurrentPrompt()

		// メニュー選択またはテキスト入力
		choice := s.showMainMenu()

		shouldExit := false
		switch choice {
		case "直接コマンド入力":
			// テキスト入力モード
			input, err := pterm.DefaultInteractiveTextInput.
				WithDefaultText("").
				WithTextStyle(pterm.NewStyle(pterm.FgLightBlue)).
				Show(prompt)

			if err != nil {
				if strings.Contains(err.Error(), "interrupt") {
					shouldExit = true
				}
				break
			}

			input = strings.TrimSpace(input)
			if input == "" {
				break
			}

			shouldExit = s.executeCommand(input)
		case "終了":
			shouldExit = true
		default:
			// メニューからの選択
			shouldExit = s.handleMenuChoice(choice)
		}

		if shouldExit {
			break
		}
	}

	pterm.Success.Println("👋 ゲームブック対話モードを終了します")
}

// showMainMenu メインメニューを表示
func (s *PTermInteractiveShell) showMainMenu() string {
	var options []string

	if currentGame != nil {
		// ゲーム読み込み済み：頻繁に使用する操作を上位に
		options = []string{
			"パラグラフ追加",
			"選択肢追加",
			"選択肢選択",
			"直接コマンド入力",
			"ヘルプ表示",
			"新しいゲームブック作成",
			"既存ゲームブック読み込み",
			"終了",
		}
	} else {
		// ゲーム未読み込み：ゲーム開始操作を上位に
		options = []string{
			"新しいゲームブック作成",
			"既存ゲームブック読み込み",
			"直接コマンド入力",
			"ヘルプ表示",
			"終了",
		}
	}

	var defaultOption string
	if currentGame != nil {
		defaultOption = "パラグラフ追加" // ゲーム中は新パラグラフ追加がデフォルト
	} else {
		defaultOption = "新しいゲームブック作成" // 未読み込み時は新規作成がデフォルト
	}

	selectedOption, _ := pterm.DefaultInteractiveSelect.
		WithOptions(options).
		WithDefaultOption(defaultOption).
		WithMaxHeight(10).
		Show("操作を選択してください:")

	return selectedOption
}

// handleMenuChoice メニュー選択を処理
func (s *PTermInteractiveShell) handleMenuChoice(choice string) bool {
	switch choice {
	case "新しいゲームブック作成":
		return s.handleNewFromMenu()
	case "既存ゲームブック読み込み":
		return s.handleLoadFromMenu()
	case "パラグラフ追加":
		return s.handleAddFromMenu()
	case "選択肢追加":
		return s.handleChoiceFromMenu()
	case "選択肢選択":
		return s.handleSelectFromMenu()
	case "ヘルプ表示":
		s.showHelp()
		return false
	default:
		return false
	}
}

// getCurrentPrompt 現在の状態に応じたプロンプトを取得
func (s *PTermInteractiveShell) getCurrentPrompt() string {
	if currentGame != nil {
		return fmt.Sprintf("[%s] > ", currentGame.Title)
	}
	return "> "
}

// autoLoadLastGame 最後のゲームを自動ロード
func (s *PTermInteractiveShell) autoLoadLastGame() {
	if currentTitle, err := sessionRepo.GetCurrentGame(); err == nil && currentTitle != "" {
		if err := s.executor.ExecuteLoadCommand(currentTitle); err == nil {
			pterm.Success.Printf("前回のゲームブック '%s' を自動読み込みしました\n", currentTitle)
		}
	}
}

// showHelp ヘルプを表示
func (s *PTermInteractiveShell) showHelp() {
	pterm.DefaultSection.Println("📚 利用可能なコマンド")

	commands := [][]string{
		{"new <ゲーム名>", "新しいゲームブックを作成"},
		{"load <ゲーム名>", "既存のゲームブックを読み込み"},
		{"add <番号> <説明>", "パラグラフを追加"},
		{"choice <番号> <説明> <遷移先>", "選択肢を追加"},
		{"select <番号> <選択肢番号>", "選択肢を選択して移動"},
		{"show", "現在の状態を表示"},
		{"help", "このヘルプを表示"},
		{"exit", "対話モードを終了"},
	}

	_ = pterm.DefaultTable.WithHasHeader().WithData(
		append([][]string{{"コマンド", "説明"}}, commands...),
	).Render()
	pterm.Println()
}

// executeCommand コマンドを実行
func (s *PTermInteractiveShell) executeCommand(input string) bool {
	args := strings.Fields(input)
	if len(args) == 0 {
		return false
	}

	command := args[0]
	commandArgs := args[1:]

	switch command {
	case "help", "h":
		s.showHelp()
	case "exit", "quit", "q":
		return true
	case "new":
		s.handleNew(commandArgs)
	case "load":
		s.handleLoad(commandArgs)
	case "add":
		s.handleAdd(commandArgs)
	case "choice":
		s.handleChoice(commandArgs)
	case "select":
		s.handleSelect(commandArgs)
	case "show":
		s.handleShow()
	default:
		pterm.Error.Printf("不明なコマンド: %s\n", command)
		pterm.Info.Println("'help' でコマンド一覧を確認できます")
	}
	return false
}

// Menu-based handlers
func (s *PTermInteractiveShell) handleNewFromMenu() bool {
	title, err := pterm.DefaultInteractiveTextInput.
		WithDefaultText("").
		Show("新しいゲームブック名を入力してください:")

	if err != nil {
		return false
	}

	s.handleNew([]string{title})
	return false
}

func (s *PTermInteractiveShell) handleLoadFromMenu() bool {
	// 既存ゲーム一覧を取得（実装は簡略化）
	title, err := pterm.DefaultInteractiveTextInput.
		WithDefaultText("").
		Show("読み込むゲームブック名を入力してください:")

	if err != nil {
		return false
	}

	s.handleLoad([]string{title})
	return false
}

func (s *PTermInteractiveShell) handleAddFromMenu() bool {
	if currentGame == nil {
		pterm.Error.Println("ゲームブックが読み込まれていません")
		return false
	}

	number, err := pterm.DefaultInteractiveTextInput.
		WithDefaultText("").
		Show("パラグラフ番号を入力してください:")
	if err != nil {
		return false
	}

	description, err := pterm.DefaultInteractiveTextInput.
		WithDefaultText("").
		Show("パラグラフの説明を入力してください:")
	if err != nil {
		return false
	}

	s.handleAdd([]string{number, description})
	return false
}

func (s *PTermInteractiveShell) handleChoiceFromMenu() bool {
	if currentGame == nil {
		pterm.Error.Println("ゲームブックが読み込まれていません")
		return false
	}

	number, err := pterm.DefaultInteractiveTextInput.Show("パラグラフ番号:")
	if err != nil {
		return false
	}

	description, err := pterm.DefaultInteractiveTextInput.Show("選択肢の説明:")
	if err != nil {
		return false
	}

	target, err := pterm.DefaultInteractiveTextInput.Show("遷移先パラグラフ番号:")
	if err != nil {
		return false
	}

	s.handleChoice([]string{number, description, target})
	return false
}

func (s *PTermInteractiveShell) handleSelectFromMenu() bool {
	if currentGame == nil {
		pterm.Error.Println("ゲームブックが読み込まれていません")
		return false
	}

	number, err := pterm.DefaultInteractiveTextInput.Show("パラグラフ番号:")
	if err != nil {
		return false
	}

	choiceIndex, err := pterm.DefaultInteractiveTextInput.Show("選択肢番号:")
	if err != nil {
		return false
	}

	s.handleSelect([]string{number, choiceIndex})
	return false
}

// Command handlers - CLIラッパー
func (s *PTermInteractiveShell) handleNew(args []string) {
	if len(args) != 1 {
		pterm.Error.Println("使用法: new <ゲーム名>")
		return
	}

	if err := s.executor.ExecuteNewCommand(args[0]); err != nil {
		pterm.Error.Println(err.Error())
	}
}

func (s *PTermInteractiveShell) handleLoad(args []string) {
	if len(args) != 1 {
		pterm.Error.Println("使用法: load <ゲーム名>")
		return
	}

	if err := s.executor.ExecuteLoadCommand(args[0]); err != nil {
		pterm.Error.Println(err.Error())
	}
}

func (s *PTermInteractiveShell) handleAdd(args []string) {
	if len(args) < 2 {
		pterm.Error.Println("使用法: add <番号> <説明>")
		return
	}

	number, err := ParseNumber(args[0], "パラグラフ番号")
	if err != nil {
		pterm.Error.Println(err.Error())
		return
	}

	description := JoinDescription(args, 1)
	if err := s.executor.ExecuteAddCommand(number, description); err != nil {
		pterm.Error.Println(err.Error())
	}
}

func (s *PTermInteractiveShell) handleChoice(args []string) {
	if len(args) < 3 {
		pterm.Error.Println("使用法: choice <番号> <説明> <遷移先>")
		return
	}

	paragraphNum, err := ParseNumber(args[0], "パラグラフ番号")
	if err != nil {
		pterm.Error.Println(err.Error())
		return
	}

	targetNum, err := ParseNumber(args[2], "遷移先パラグラフ番号")
	if err != nil {
		pterm.Error.Println(err.Error())
		return
	}

	if err := s.executor.ExecuteChoiceCommand(paragraphNum, args[1], targetNum); err != nil {
		pterm.Error.Println(err.Error())
	}
}

func (s *PTermInteractiveShell) handleSelect(args []string) {
	if len(args) < 2 {
		pterm.Error.Println("使用法: select <番号> <選択肢番号>")
		return
	}

	paragraphNum, err := ParseNumber(args[0], "パラグラフ番号")
	if err != nil {
		pterm.Error.Println(err.Error())
		return
	}

	choiceNum, err := ParseNumber(args[1], "選択肢番号")
	if err != nil {
		pterm.Error.Println(err.Error())
		return
	}
	choiceIndex := choiceNum - 1 // 1ベースから0ベースに変換

	if err := s.executor.ExecuteSelectCommand(paragraphNum, choiceIndex); err != nil {
		pterm.Error.Println(err.Error())
	}
}

func (s *PTermInteractiveShell) handleShow() {
	if err := s.executor.ExecuteShowCommand(); err != nil {
		pterm.Error.Println(err.Error())
	}
}
