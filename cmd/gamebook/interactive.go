package main

import (
	"fmt"
	"strings"

	"github.com/pterm/pterm"
)

// InteractiveShell PTerm対話モードのシェル
type InteractiveShell struct {
	executor CommandExecutor
	commands []string
}

// NewInteractiveShell PTerm対話シェルを作成
func NewInteractiveShell() (*InteractiveShell, error) {
	// 利用可能コマンド定義
	commands := []string{
		"new <ゲーム名>      - 新しいゲームブックを作成",
		"load <ゲーム名>     - 既存のゲームブックを読み込み",
		"add <番号> <説明>   - パラグラフを追加",
		"choice <番号> <説明> <遷移先> - 選択肢を追加",
		"select <番号> <選択肢番号> - 選択肢を選択して移動",
		"show               - 現在の状態を表示",
		"help               - このヘルプを表示",
		"exit               - 対話モードを終了",
	}

	return &InteractiveShell{
		executor: NewCLIExecutor(),
		commands: commands,
	}, nil
}

// Run PTerm対話シェルを実行
func (s *InteractiveShell) Run() error {
	// ウェルカムメッセージを表示
	pterm.DefaultHeader.WithFullWidth().Println("🎮 Gamebook Interactive Mode")
	pterm.Info.Println("リッチなターミナルUI対話モードです")
	pterm.Info.Println("'help'でコマンド一覧、'exit'で終了")
	pterm.Println()

	// 最後のゲームを自動ロード
	s.autoLoadLastGame()

	for {
		// 現在の状態に応じたプロンプトを表示
		prompt := s.getCurrentPrompt()

		// コマンド入力（リアルタイム補完付き）
		input, err := pterm.DefaultInteractiveTextInput.
			WithDefaultText("").
			WithTextStyle(pterm.NewStyle(pterm.FgLightBlue)).
			Show(prompt)

		if err != nil {
			if strings.Contains(err.Error(), "interrupt") {
				pterm.Success.Println("👋 ゲームブック対話モードを終了します")
				break
			}
			return err
		}

		input = strings.TrimSpace(input)
		if input == "" {
			continue
		}

		// コマンド実行
		if shouldExit := s.executeCommand(input); shouldExit {
			break
		}
	}

	return nil
}

// executeCommand コマンドを実行
func (s *InteractiveShell) executeCommand(line string) bool {
	args := strings.Fields(line)
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
		s.handleShow(commandArgs)
	default:
		pterm.Error.Printf("不明なコマンド: %s\n", command)
		pterm.Info.Println("'help' でコマンド一覧を確認できます")
	}
	return false
}

// showHelp ヘルプを表示
func (s *InteractiveShell) showHelp() {
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

// getCurrentPrompt 現在の状態に応じたプロンプトを取得
func (s *InteractiveShell) getCurrentPrompt() string {
	if currentGame != nil {
		return fmt.Sprintf("[%s] > ", currentGame.Title)
	}
	return "> "
}

// autoLoadLastGame 最後のゲームを自動ロード
func (s *InteractiveShell) autoLoadLastGame() {
	if currentTitle, err := sessionRepo.GetCurrentGame(); err == nil && currentTitle != "" {
		if err := s.executor.ExecuteLoadCommand(currentTitle); err == nil {
			pterm.Success.Printf("前回のゲームブック '%s' を自動読み込みしました\n", currentTitle)
		}
	}
}

// CLIコマンドハンドラー
func (s *InteractiveShell) handleNew(args []string) {
	if len(args) != 1 {
		pterm.Error.Println("使用法: new <ゲーム名>")
		return
	}

	if err := s.executor.ExecuteNewCommand(args[0]); err != nil {
		pterm.Error.Println(err.Error())
	}
}

func (s *InteractiveShell) handleLoad(args []string) {
	if len(args) != 1 {
		pterm.Error.Println("使用法: load <ゲーム名>")
		return
	}

	if err := s.executor.ExecuteLoadCommand(args[0]); err != nil {
		pterm.Error.Println(err.Error())
	}
}

func (s *InteractiveShell) handleAdd(args []string) {
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

func (s *InteractiveShell) handleChoice(args []string) {
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

func (s *InteractiveShell) handleSelect(args []string) {
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

func (s *InteractiveShell) handleShow(_ []string) {
	if err := s.executor.ExecuteShowCommand(); err != nil {
		pterm.Error.Println(err.Error())
	}
}
