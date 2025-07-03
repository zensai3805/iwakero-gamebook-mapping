package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/chzyer/readline"
	"github.com/iwapc/iwakero-gamebook-mapping/internal/domain"
	"github.com/iwapc/iwakero-gamebook-mapping/internal/infrastructure/repository"
)

// InteractiveShell 対話モードのシェル
type InteractiveShell struct {
	rl          *readline.Instance
	currentGame *domain.Gamebook
	repository  domain.GamebookRepository
	sessionRepo domain.SessionRepository
}

// NewInteractiveShell 対話シェルを作成
func NewInteractiveShell() (*InteractiveShell, error) {
	// ホームディレクトリの履歴ファイルパスを取得
	homeDir, err := os.UserHomeDir()
	if err != nil {
		homeDir = "."
	}
	historyFile := filepath.Join(homeDir, ".gamebook_history")

	// readline設定
	rl, err := readline.NewEx(&readline.Config{
		Prompt:       "> ",
		HistoryFile:  historyFile,
		HistoryLimit: 500,
		AutoComplete: setupCompleter(),
	})
	if err != nil {
		return nil, err
	}

	// リポジトリ初期化
	dataDir := getDataDir()
	markdownRepo := repository.NewMarkdownRepository(dataDir)
	sessionRepo := repository.NewFileSessionRepository(dataDir)

	return &InteractiveShell{
		rl:          rl,
		repository:  markdownRepo,
		sessionRepo: sessionRepo,
	}, nil
}

// Run 対話シェルを実行
func (s *InteractiveShell) Run() error {
	defer s.rl.Close()

	fmt.Println("Gamebook Interactive Mode")
	fmt.Println("Type 'help' for available commands, 'exit' to quit.")
	fmt.Println("")

	for {
		line, err := s.rl.Readline()
		if err != nil {
			if err == readline.ErrInterrupt {
				if len(line) == 0 {
					break
				} else {
					continue
				}
			} else if err == io.EOF {
				break
			}
			return err
		}

		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		if err := s.executeCommand(line); err != nil {
			if err.Error() == "exit" {
				break
			}
			fmt.Fprintf(os.Stderr, "エラー: %v\n", err)
		}
	}

	return nil
}

// executeCommand コマンドを実行
func (s *InteractiveShell) executeCommand(line string) error {
	args := strings.Fields(line)
	if len(args) == 0 {
		return nil
	}

	command := args[0]
	commandArgs := args[1:]

	switch command {
	case "help", "h":
		return s.showHelp()
	case "exit", "quit", "q":
		return fmt.Errorf("exit")
	case "new":
		return s.handleNew(commandArgs)
	case "load":
		return s.handleLoad(commandArgs)
	case "add":
		return s.handleAdd(commandArgs)
	case "choice":
		return s.handleChoice(commandArgs)
	case "select":
		return s.handleSelect(commandArgs)
	case "show":
		return s.handleShow(commandArgs)
	default:
		return fmt.Errorf("不明なコマンド: %s (ヘルプは 'help' を入力)", command)
	}
}

// showHelp ヘルプを表示
func (s *InteractiveShell) showHelp() error {
	fmt.Println("利用可能なコマンド:")
	fmt.Println("  new <ゲーム名>              - 新しいゲームブックを作成")
	fmt.Println("  load <ゲーム名>             - 既存のゲームブックを読み込み")
	fmt.Println("  add <番号> <説明>           - パラグラフを追加")
	fmt.Println("  choice <番号> <説明> <遷移先> - 選択肢を追加")
	fmt.Println("  select <番号> <選択肢番号>    - 選択肢を選択して移動")
	fmt.Println("  show                       - 現在の状態を表示")
	fmt.Println("  help, h                    - このヘルプを表示")
	fmt.Println("  exit, quit, q              - 対話モードを終了")
	fmt.Println("")
	fmt.Println("Ctrl+C または Ctrl+D でも終了できます")
	return nil
}

// updatePrompt プロンプトを更新
func (s *InteractiveShell) updatePrompt() {
	if s.currentGame != nil {
		s.rl.SetPrompt(fmt.Sprintf("[%s] > ", s.currentGame.Title))
	} else {
		s.rl.SetPrompt("> ")
	}
}

// setupCompleter タブ補完を設定
func setupCompleter() readline.AutoCompleter {
	return readline.NewPrefixCompleter(
		readline.PcItem("new"),
		readline.PcItem("load"),
		readline.PcItem("add"),
		readline.PcItem("choice"),
		readline.PcItem("select"),
		readline.PcItem("show"),
		readline.PcItem("help"),
		readline.PcItem("exit"),
		readline.PcItem("quit"),
	)
}

// 既存のコマンドハンドラーを簡単な形で実装
func (s *InteractiveShell) handleNew(args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("使用法: new <ゲーム名>")
	}
	title := args[0]
	s.currentGame = domain.NewGamebook(title)
	s.updatePrompt()
	fmt.Printf("新しいゲームブック '%s' を作成しました\n", title)
	return nil
}

func (s *InteractiveShell) handleLoad(args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("使用法: load <ゲーム名>")
	}
	title := args[0]
	gamebook, err := s.repository.Load(title)
	if err != nil {
		return fmt.Errorf("ゲームブックの読み込みに失敗しました: %v", err)
	}
	s.currentGame = gamebook
	s.updatePrompt()
	fmt.Printf("ゲームブック '%s' を読み込みました\n", title)
	return nil
}

func (s *InteractiveShell) handleAdd(args []string) error {
	if s.currentGame == nil {
		return fmt.Errorf("ゲームブックが読み込まれていません。'new' または 'load' を実行してください")
	}
	if len(args) < 2 {
		return fmt.Errorf("使用法: add <番号> <説明>")
	}

	// 既存のaddCommand実装を参考に実装
	return fmt.Errorf("add コマンドは実装中です")
}

func (s *InteractiveShell) handleChoice(args []string) error {
	if s.currentGame == nil {
		return fmt.Errorf("ゲームブックが読み込まれていません")
	}
	if len(args) < 3 {
		return fmt.Errorf("使用法: choice <番号> <説明> <遷移先>")
	}

	return fmt.Errorf("choice コマンドは実装中です")
}

func (s *InteractiveShell) handleSelect(args []string) error {
	if s.currentGame == nil {
		return fmt.Errorf("ゲームブックが読み込まれていません")
	}
	if len(args) < 2 {
		return fmt.Errorf("使用法: select <番号> <選択肢番号>")
	}

	return fmt.Errorf("select コマンドは実装中です")
}

func (s *InteractiveShell) handleShow(_ []string) error {
	if s.currentGame == nil {
		return fmt.Errorf("ゲームブックが読み込まれていません")
	}

	return fmt.Errorf("show コマンドは実装中です")
}
