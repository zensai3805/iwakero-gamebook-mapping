package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/zensai3805/iwakero-gamebook-mapping/internal/domain"
	"github.com/zensai3805/iwakero-gamebook-mapping/internal/infrastructure/repository"
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
	// ログシステムを初期化
	_, loggerInitErr := initializeApplicationLogger()
	if loggerInitErr != nil {
		fmt.Fprintf(os.Stderr, "ログシステム初期化エラー: %v\n", loggerInitErr)
		os.Exit(1)
	}

	// アプリケーション終了時にクリーンアップ
	defer func() {
		if cleanupErr := cleanupApplicationLogger(); cleanupErr != nil {
			fmt.Fprintf(os.Stderr, "ログシステムクリーンアップエラー: %v\n", cleanupErr)
		}
	}()

	// パニック時のログ出力
	defer func() {
		if r := recover(); r != nil {
			logger := GetGlobalLogger()
			if logger != nil {
				logger.Fatal("アプリケーションパニック",
					domain.Field{Key: "panic", Value: fmt.Sprintf("%v", r)})
			}
			fmt.Fprintf(os.Stderr, "パニック: %v\n", r)
			os.Exit(1)
		}
	}()

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
		executor := NewCLIExecutor()
		if err := executor.ExecuteNewCommand(args[0]); err != nil {
			fmt.Fprintf(os.Stderr, "エラー: %v\n", err)
			return
		}
	},
}

var addCmd = &cobra.Command{
	Use:   "add [パラグラフ番号] [説明]",
	Short: "パラグラフを追加",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		number, numberErr := strconv.Atoi(args[0])
		if numberErr != nil {
			fmt.Fprintf(os.Stderr, "エラー: パラグラフ番号は数値で指定してください: %v\n", numberErr)
			return
		}
		description := args[1]

		executor := NewCLIExecutor()
		if err := executor.ExecuteAddCommand(number, description); err != nil {
			fmt.Fprintf(os.Stderr, "エラー: %v\n", err)
			return
		}
	},
}

var showCmd = &cobra.Command{
	Use:   "show",
	Short: "現在の状態を表示",
	Run: func(cmd *cobra.Command, args []string) {
		executor := NewCLIExecutor()
		if err := executor.ExecuteShowCommand(); err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
	},
}

// initializeApplicationLogger はアプリケーション用のログシステムを初期化する
func initializeApplicationLogger() (domain.Logger, error) {
	// AI開発モード時は詳細ログを有効化
	if IsAIDevelopmentMode() {
		os.Setenv("LOG_LEVEL", "DEBUG")
	}

	return InitializeLogger()
}

// cleanupApplicationLogger はアプリケーション用のログシステムをクリーンアップする
func cleanupApplicationLogger() error {
	return CleanupLogger()
}
