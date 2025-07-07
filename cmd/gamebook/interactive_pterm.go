package main

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"github.com/pterm/pterm"
	"github.com/zensai3805/iwakero-gamebook-mapping/internal/domain"
)

const (
	// エラーメッセージ定数
	ErrNoGameLoaded = "ゲームブックが読み込まれていません"
)

// PTermInteractiveShell PTerm対話モードのシェル
type PTermInteractiveShell struct {
	executor  CommandExecutor
	lastError string // 最後のエラーメッセージを保持
	lastInfo  string // 最後の情報メッセージを保持
}

// NewPTermInteractiveShell PTerm対話シェルを作成
func NewPTermInteractiveShell() *PTermInteractiveShell {
	return &PTermInteractiveShell{
		executor: NewCLIExecutor(),
	}
}

// Run PTerm対話シェルを実行
func (s *PTermInteractiveShell) Run() {
	// インタラクティブモード用のログ出力制御
	originalLevel, logOutputChanged := s.setupInteractiveLogging()
	defer s.restoreLogging(originalLevel, logOutputChanged)

	// 初回表示フラグ
	isFirstDisplay := true

	// 最後のゲームを自動ロード（画面クリア前に実行）
	s.autoLoadLastGame()

	for {
		// 画面をクリアして固定位置から表示開始
		s.clearAndShowHeader()

		// 初回のみ説明を表示
		if isFirstDisplay {
			pterm.Info.Println("リッチなターミナルUI対話モードです")
			pterm.Info.Println("↑↓キーでメニュー選択、Enterで決定、Ctrl+Cで終了")
			pterm.Println()

			// ログにも記録
			if logger := GetGlobalLogger(); logger != nil {
				logger.Info("対話モード開始",
					domain.Field{Key: "mode", Value: "interactive"},
					domain.Field{Key: "ui", Value: "pterm"})
			}

			isFirstDisplay = false
		}

		// エラーメッセージがあれば表示
		if s.lastError != "" {
			pterm.Error.Println(s.lastError)

			// ログにも記録
			if logger := GetGlobalLogger(); logger != nil {
				logger.Error("UI表示エラー", domain.Field{Key: "message", Value: s.lastError})
			}

			s.lastError = "" // 表示後はクリア
		}

		// 情報メッセージがあれば表示
		if s.lastInfo != "" {
			pterm.Success.Println(s.lastInfo)

			// ログにも記録
			if logger := GetGlobalLogger(); logger != nil {
				logger.Info("UI表示情報", domain.Field{Key: "message", Value: s.lastInfo})
			}

			s.lastInfo = "" // 表示後はクリア
		}

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

	// ログにも記録
	if logger := GetGlobalLogger(); logger != nil {
		logger.Info("対話モード終了")
	}
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
			"📍 直接移動",
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

	// UI操作の開始時刻を記録
	startTime := time.Now()

	selectedOption, err := pterm.DefaultInteractiveSelect.
		WithOptions(options).
		WithDefaultOption(defaultOption).
		WithMaxHeight(10).
		Show("操作を選択してください:")

	// UI操作の記録（軽量版）
	LogUIInteraction("menu_selection", map[string]interface{}{
		"selected_option":   selectedOption,
		"default_option":    defaultOption,
		"options_count":     len(options),
		"has_current_game":  currentGame != nil,
		"selection_time_ms": float64(time.Since(startTime).Nanoseconds()) / 1000000,
	})

	if err != nil {
		LogErrorWithContext(err, "menu_selection_error", map[string]interface{}{
			"options_count":    len(options),
			"has_current_game": currentGame != nil,
		})
	}

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
	case "📍 直接移動":
		return s.handleMoveFromMenu()
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
			s.lastInfo = fmt.Sprintf("前回のゲームブック '%s' を自動読み込みしました", currentTitle)
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
		{"move <番号>", "指定パラグラフに直接移動（経路上の選択肢を自動選択）"},
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
		s.lastError = fmt.Sprintf("不明なコマンド: %s ('help' でコマンド一覧を確認できます)", command)
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
		s.lastError = ErrNoGameLoaded
		return false
	}

	// 入力支援機能を使用
	helper := NewInputHelper(currentGame)

	// パラグラフ番号入力（未定義パラグラフの場合のみデフォルト値提供）
	prompt := "パラグラフ番号を入力してください"
	defaultText := helper.GetDefaultTextForAdd()
	if defaultText != "" {
		prompt = fmt.Sprintf("パラグラフ番号を入力してください (Enter: %s)", defaultText)
	}

	number, err := pterm.DefaultInteractiveTextInput.
		WithDefaultText(prompt).
		Show()
	if err != nil {
		return false
	}

	// 空入力時は現在地を自動入力（未定義パラグラフの場合のみ）
	if number == "" && defaultText != "" {
		number = defaultText
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
		s.lastError = ErrNoGameLoaded
		return false
	}

	// 入力支援機能を使用
	helper := NewInputHelper(currentGame)

	// パラグラフ番号入力（空Enter時は現在地）
	prompt := "パラグラフ番号"
	if helper.GetDefaultText() != "" {
		prompt = fmt.Sprintf("パラグラフ番号 (Enter: %s)", helper.GetDefaultText())
	}
	number, err := pterm.DefaultInteractiveTextInput.
		WithDefaultText(prompt).
		Show()
	if err != nil {
		return false
	}

	// 空入力時は現在地を自動入力
	number = helper.ProcessEmptyInput(number)

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
		s.lastError = ErrNoGameLoaded
		return false
	}

	// 入力支援機能を使用
	helper := NewInputHelper(currentGame)

	// パラグラフ番号入力（空Enter時は現在地）
	prompt := "パラグラフ番号"
	if helper.GetDefaultText() != "" {
		prompt = fmt.Sprintf("パラグラフ番号 (Enter: %s)", helper.GetDefaultText())
	}
	number, err := pterm.DefaultInteractiveTextInput.
		WithDefaultText(prompt).
		Show()
	if err != nil {
		return false
	}

	// 空入力時は現在地を自動入力
	number = helper.ProcessEmptyInput(number)

	// パラグラフ番号をintに変換
	paragraphNum, parseErr := ParseNumber(number, "パラグラフ番号")
	if parseErr != nil {
		s.lastError = parseErr.Error()
		return false
	}

	// 指定されたパラグラフを取得
	paragraph, getErr := currentGame.GetParagraph(paragraphNum)
	if getErr != nil {
		s.lastError = getErr.Error()
		return false
	}

	// 選択肢がない場合の処理
	if len(paragraph.Choices) == 0 {
		s.lastError = "選択肢がありません"
		return false
	}

	// 選択肢を表示形式に変換
	choiceOptions := make([]string, len(paragraph.Choices))
	for i, choice := range paragraph.Choices {
		status := "未選択"
		if choice.Selected {
			status = "選択済み"
		}
		choiceOptions[i] = fmt.Sprintf("%s → #%d [%s]", choice.Description, choice.TargetNumber, status)
	}

	// PTerm interactive selectを使用して選択肢を選択
	selectedOption, selectErr := pterm.DefaultInteractiveSelect.
		WithOptions(choiceOptions).
		Show("選択肢を選択してください:")
	if selectErr != nil {
		s.lastError = selectErr.Error()
		return false
	}

	// 選択された選択肢のインデックスを取得
	selectedIndex := -1
	for i, option := range choiceOptions {
		if option == selectedOption {
			selectedIndex = i
			break
		}
	}

	if selectedIndex == -1 {
		s.lastError = "選択肢の特定に失敗しました"
		return false
	}

	// 選択処理を実行
	if execErr := s.executor.ExecuteSelectCommand(paragraphNum, selectedIndex); execErr != nil {
		s.lastError = execErr.Error()
	} else {
		s.lastInfo = fmt.Sprintf("選択肢 %d を選択しました", selectedIndex+1)
	}

	return false
}

// Command handlers - CLIラッパー
func (s *PTermInteractiveShell) handleNew(args []string) {
	if len(args) != 1 {
		s.lastError = "使用法: new <ゲーム名>"
		return
	}

	if err := s.executor.ExecuteNewCommand(args[0]); err != nil {
		s.lastError = err.Error()
	} else {
		s.lastInfo = fmt.Sprintf("新しいゲームブック '%s' を作成しました", args[0])
	}
}

func (s *PTermInteractiveShell) handleLoad(args []string) {
	if len(args) != 1 {
		s.lastError = "使用法: load <ゲーム名>"
		return
	}

	if err := s.executor.ExecuteLoadCommand(args[0]); err != nil {
		s.lastError = err.Error()
	} else {
		s.lastInfo = fmt.Sprintf("ゲームブック '%s' を読み込みました", args[0])
	}
}

func (s *PTermInteractiveShell) handleAdd(args []string) {
	if len(args) < 2 {
		s.lastError = "使用法: add <番号> <説明>"
		return
	}

	number, err := ParseNumber(args[0], "パラグラフ番号")
	if err != nil {
		s.lastError = err.Error()
		return
	}

	description := JoinDescription(args, 1)
	if err := s.executor.ExecuteAddCommand(number, description); err != nil {
		s.lastError = err.Error()
	} else {
		s.lastInfo = fmt.Sprintf("パラグラフ %d を追加しました", number)
	}
}

func (s *PTermInteractiveShell) handleChoice(args []string) {
	if len(args) < 3 {
		s.lastError = "使用法: choice <番号> <説明> <遷移先>"
		return
	}

	paragraphNum, err := ParseNumber(args[0], "パラグラフ番号")
	if err != nil {
		s.lastError = err.Error()
		return
	}

	targetNum, err := ParseNumber(args[2], "遷移先パラグラフ番号")
	if err != nil {
		s.lastError = err.Error()
		return
	}

	if err := s.executor.ExecuteChoiceCommand(paragraphNum, args[1], targetNum); err != nil {
		s.lastError = err.Error()
	} else {
		s.lastInfo = fmt.Sprintf("選択肢を追加しました: %s → %d", args[1], targetNum)
	}
}

func (s *PTermInteractiveShell) handleSelect(args []string) {
	if len(args) < 2 {
		s.lastError = "使用法: select <番号> <選択肢番号>"
		return
	}

	paragraphNum, err := ParseNumber(args[0], "パラグラフ番号")
	if err != nil {
		s.lastError = err.Error()
		return
	}

	choiceNum, err := ParseNumber(args[1], "選択肢番号")
	if err != nil {
		s.lastError = err.Error()
		return
	}
	choiceIndex := choiceNum - 1 // 1ベースから0ベースに変換

	if err := s.executor.ExecuteSelectCommand(paragraphNum, choiceIndex); err != nil {
		s.lastError = err.Error()
	} else {
		s.lastInfo = fmt.Sprintf("選択肢 %d を選択しました", choiceNum)
	}
}

func (s *PTermInteractiveShell) handleShow() {
	if err := s.executor.ExecuteShowCommand(); err != nil {
		pterm.Error.Println(err.Error())
	}
}

// handleMoveFromMenu メニューからの直接移動を処理
func (s *PTermInteractiveShell) handleMoveFromMenu() bool {
	if currentGame == nil {
		s.lastError = ErrNoGameLoaded
		return false
	}

	// 入力支援機能を使用
	helper := NewInputHelper(currentGame)

	// パラグラフ番号入力（空Enter時は現在地）
	prompt := "📍 移動先パラグラフ番号"
	if helper.GetDefaultText() != "" {
		prompt = fmt.Sprintf("📍 移動先パラグラフ番号 (Enter: %s)", helper.GetDefaultText())
	}
	targetNumStr, err := pterm.DefaultInteractiveTextInput.
		WithDefaultText(prompt).
		WithTextStyle(pterm.NewStyle(pterm.FgLightBlue)).
		Show()

	if err != nil {
		return false
	}

	// 空入力時は現在地を自動入力
	targetNumStr = helper.ProcessEmptyInput(targetNumStr)

	// 数値パース
	targetNum, parseErr := ParseNumber(targetNumStr, "パラグラフ番号")
	if parseErr != nil {
		s.lastError = parseErr.Error()
		return false
	}

	// moveコマンド実行
	if moveErr := s.executor.ExecuteMoveCommand(targetNum); moveErr != nil {
		s.lastError = moveErr.Error()
	} else {
		s.lastInfo = fmt.Sprintf("📍 パラグラフ %d に移動しました", targetNum)
	}

	return false
}

// clearAndShowHeader 画面をクリアしてヘッダーを表示
func (s *PTermInteractiveShell) clearAndShowHeader() {
	s.clearScreen()
	pterm.DefaultHeader.WithFullWidth().Println("🎮 Gamebook Interactive Mode")
}

// clearScreen 画面をクリア（クロスプラットフォーム対応）
func (s *PTermInteractiveShell) clearScreen() {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("cmd", "/c", "cls")
	default:
		cmd = exec.Command("clear")
	}
	cmd.Stdout = os.Stdout
	_ = cmd.Run()
}

// setupInteractiveLogging はインタラクティブモード用のログ設定を行う
func (s *PTermInteractiveShell) setupInteractiveLogging() (domain.LogLevel, bool) {
	// 現在のログレベルを保存
	originalLevel := domain.LogLevelInfo
	if loggingController != nil {
		originalLevel = loggingController.GetLevel()
	}

	logOutputChanged := false

	// コンソール出力の場合は調整が必要
	if IsConsoleOutput() {
		if IsAIDevelopmentMode() {
			// AI開発モード: ファイル出力に切り替え
			interactiveLogFile := "./logs/interactive.log"

			// ログディレクトリを作成
			if dirErr := os.MkdirAll("./logs", 0755); dirErr != nil {
				pterm.Warning.Printf("ログディレクトリ作成に失敗: %v\n", dirErr)
				SetTemporaryLogLevel(domain.LogLevelWarn)
			} else {
				// AI開発モード時は先にDEBUGレベルに設定
				if loggingController != nil {
					loggingController.SetLevel(domain.LogLevelDebug)
				}
				if switchErr := SwitchToFileOutput(loggingController, interactiveLogFile); switchErr != nil {
					// 切り替えに失敗した場合はログレベルを制限
					pterm.Warning.Printf("ログ出力をファイルに切り替えできませんでした。ログレベルをWARNに制限します: %v\n", switchErr)
					SetTemporaryLogLevel(domain.LogLevelWarn)
				} else {
					logOutputChanged = true
					pterm.Info.Printf("AI開発モードのため、ログをファイルに出力します: %s\n", interactiveLogFile)
				}
			}
		} else {
			// 通常モード: ログレベルをWARN以上に制限
			SetTemporaryLogLevel(domain.LogLevelWarn)
			pterm.Info.Println("インタラクティブモード: ログレベルをWARN以上に制限しました")
		}
	}

	return originalLevel, logOutputChanged
}

// restoreLogging はログ設定を元に戻す
func (s *PTermInteractiveShell) restoreLogging(originalLevel domain.LogLevel, logOutputChanged bool) {
	if logOutputChanged {
		// ファイル出力から元の設定に戻すのは複雑なので、
		// ユーザーに再起動を促すメッセージを表示
		pterm.Info.Println("ログ出力設定が変更されました。元の設定に戻すにはアプリケーションを再起動してください。")
	} else {
		// ログレベルのみ元に戻す
		RestoreLogLevel(originalLevel)
	}
}
