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
			var err error
			title, err = sessionRepo.GetCurrentGame()
			if err != nil || title == "" {
				fmt.Fprintln(os.Stderr, "エラー: 最後に使用したゲームブックが見つかりません。")
				return
			}
		} else {
			title = args[0]
		}

		gamebook, err := repo.Load(title)
		if err != nil {
			fmt.Fprintf(os.Stderr, "読み込みエラー: %v\n", err)
			return
		}
		currentGame = gamebook

		// 現在のゲームとして保存
		if err := sessionRepo.SaveCurrentGame(title); err != nil {
			fmt.Fprintf(os.Stderr, "セッション保存エラー: %v\n", err)
		}

		fmt.Printf("ゲームブック「%s」を読み込みました。\n", title)
	},
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "保存されているゲームブック一覧を表示",
	Run: func(cmd *cobra.Command, args []string) {
		titles, err := repo.List()
		if err != nil {
			fmt.Fprintf(os.Stderr, "エラー: %v\n", err)
			return
		}

		if len(titles) == 0 {
			fmt.Println("保存されているゲームブックはありません。")
			return
		}

		fmt.Println("=== 保存されているゲームブック ===")
		for i, title := range titles {
			fmt.Printf("%d. %s\n", i+1, title)
		}
	},
}
