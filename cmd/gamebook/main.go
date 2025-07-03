package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/iwapc/iwakero-gamebook-mapping/internal/domain"
	"github.com/iwapc/iwakero-gamebook-mapping/internal/infrastructure/repository"
	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Use:   "gamebook",
		Short: "ゲームブック記録支援ツール",
		Long:  `ゲームブックのパラグラフと選択を記録し、全体構造を可視化するツールです。`,
	}

	dataDir     = getDataDir()
	repo        = repository.NewMarkdownRepository(dataDir)
	sessionRepo = repository.NewFileSessionRepository(dataDir)
	currentGame *domain.Gamebook
)

// getDataDir は環境変数またはデフォルトのデータディレクトリを返す
func getDataDir() string {
	if dir := os.Getenv("GAMEBOOK_DATA_DIR"); dir != "" {
		return dir
	}
	return "./data"
}

func main() {
	// 引数なし → 対話モード
	if len(os.Args) == 1 {
		runInteractiveMode()
		return
	}

	// 引数あり → 既存のワンショットコマンド
	runSingleCommand()
}

// runInteractiveMode 対話モードを実行
func runInteractiveMode() {
	shell := NewPTermInteractiveShell()
	shell.Run()
}

// runSingleCommand 既存のワンショットコマンドを実行
func runSingleCommand() {
	// コマンド実行前に現在のゲームを自動ロード
	loadCurrentGameIfExists()

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

// loadCurrentGameIfExists は保存されている現在のゲームが存在すれば自動ロードする
func loadCurrentGameIfExists() {
	if currentTitle, err := sessionRepo.GetCurrentGame(); err == nil && currentTitle != "" {
		if gamebook, err := repo.Load(currentTitle); err == nil {
			currentGame = gamebook
		}
	}
}

func init() {
	rootCmd.AddCommand(newCmd)
	rootCmd.AddCommand(loadCmd)
	rootCmd.AddCommand(addCmd)
	rootCmd.AddCommand(addChoiceCmd)
	rootCmd.AddCommand(selectChoiceCmd)
	rootCmd.AddCommand(moveCmd)
	rootCmd.AddCommand(showCmd)
	rootCmd.AddCommand(listCmd)
}

var newCmd = &cobra.Command{
	Use:   "new [タイトル]",
	Short: "新しいゲームブックを作成",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		title := args[0]
		currentGame = domain.NewGamebook(title)
		if err := repo.Save(currentGame); err != nil {
			fmt.Fprintf(os.Stderr, "エラー: %v\n", err)
			return
		}

		// 現在のゲームとして保存
		if err := sessionRepo.SaveCurrentGame(title); err != nil {
			fmt.Fprintf(os.Stderr, "セッション保存エラー: %v\n", err)
		}

		fmt.Printf("新しいゲームブック「%s」を作成しました。\n", title)
	},
}

var addCmd = &cobra.Command{
	Use:   "add [パラグラフ番号] [説明]",
	Short: "パラグラフを追加",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		if currentGame == nil {
			fmt.Fprintln(os.Stderr, "エラー: ゲームブックが選択されていません。'gamebook new'または'gamebook load'を実行してください。")
			return
		}

		number, err := strconv.Atoi(args[0])
		if err != nil {
			fmt.Fprintf(os.Stderr, "エラー: パラグラフ番号は数値で指定してください: %v\n", err)
			return
		}
		description := args[1]

		p := domain.NewParagraph(number, description)
		if err := currentGame.AddParagraph(p); err != nil {
			fmt.Fprintf(os.Stderr, "エラー: %v\n", err)
			return
		}

		if err := repo.Save(currentGame); err != nil {
			fmt.Fprintf(os.Stderr, "保存エラー: %v\n", err)
			return
		}

		fmt.Printf("パラグラフ %d を追加しました: %s\n", number, description)
	},
}

var showCmd = &cobra.Command{
	Use:   "show",
	Short: "現在の状態を表示",
	Run: func(cmd *cobra.Command, args []string) {
		if currentGame == nil {
			fmt.Fprintln(os.Stderr, "エラー: ゲームブックが選択されていません。")
			return
		}

		fmt.Printf("=== %s ===\n", currentGame.Title)

		// 統計情報
		stats := currentGame.GetExplorationStats()
		fmt.Printf("\n統計情報:\n")
		fmt.Printf("- 総パラグラフ数: %d\n", stats.TotalParagraphs)
		fmt.Printf("- 訪問済み: %d\n", stats.VisitedParagraphs)
		fmt.Printf("- 総選択肢数: %d\n", stats.TotalChoices)
		fmt.Printf("- 選択済み: %d\n", stats.SelectedChoices)

		// 現在位置
		if currentGame.Current != nil {
			fmt.Printf("\n現在位置: パラグラフ %d\n", currentGame.Current.Number)
			fmt.Printf("説明: %s\n", currentGame.Current.Description)

			if len(currentGame.Current.Choices) > 0 {
				fmt.Println("\n選択肢:")
				for i, choice := range currentGame.Current.Choices {
					status := StatusUnselected
					if choice.Selected {
						status = StatusSelected
					}
					fmt.Printf("  %d. %s → %d [%s]\n", i+1, choice.Description, choice.TargetNumber, status)
				}
			}
		}

		// 最近追加されたパラグラフ
		fmt.Println("\n最近のパラグラフ:")
		count := 0
		for i := 1; i <= 1000 && count < 5; i++ {
			if p, exists := currentGame.Paragraphs[i]; exists {
				status := "未訪問"
				if p.Visited {
					status = "訪問済み"
				}
				fmt.Printf("  %d: %s [%s]\n", p.Number, p.Description, status)

				// 選択肢があれば表示
				if len(p.Choices) > 0 {
					fmt.Printf("    選択肢:\n")
					for j, choice := range p.Choices {
						choiceStatus := StatusUnselected
						if choice.Selected {
							choiceStatus = StatusSelected
						}
						fmt.Printf("      %d. %s → %d [%s]\n", j+1, choice.Description, choice.TargetNumber, choiceStatus)
					}
				}
				count++
			}
		}
	},
}
