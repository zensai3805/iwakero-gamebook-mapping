package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var loadCmd = &cobra.Command{
	Use:   "load [タイトル]",
	Short: "既存のゲームブックを読み込み（タイトル省略時は最後のゲーム）",
	Args:  cobra.RangeArgs(0, 1),
	Run: func(cmd *cobra.Command, args []string) {
		var title string

		if len(args) == 0 {
			// タイトルが指定されていない場合、最後のゲームを取得
			var loadErr error
			title, loadErr = sessionRepo.GetCurrentGame()
			if loadErr != nil || title == "" {
				fmt.Fprintln(os.Stderr, "エラー: 最後に使用したゲームブックが見つかりません。")
				return
			}
		} else {
			title = args[0]
		}

		executor := NewCLIExecutor()
		if err := executor.ExecuteLoadCommand(title); err != nil {
			fmt.Fprintf(os.Stderr, "エラー: %v\n", err)
			return
		}
	},
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "保存されているゲームブック一覧を表示",
	Run: func(cmd *cobra.Command, args []string) {
		executor := NewCLIExecutor()
		if err := executor.ExecuteListCommand(); err != nil {
			fmt.Fprintf(os.Stderr, "エラー: %v\n", err)
			return
		}
	},
}
