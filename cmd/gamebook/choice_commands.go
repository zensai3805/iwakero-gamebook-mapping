package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/spf13/cobra"
)

var addChoiceCmd = &cobra.Command{
	Use:   "choice [パラグラフ番号] [選択肢の説明] [遷移先パラグラフ番号]",
	Short: "パラグラフに選択肢を追加",
	Args:  cobra.ExactArgs(3),
	Run: func(cmd *cobra.Command, args []string) {
		// パラグラフ番号をパース
		paragraphNum, paragraphErr := strconv.Atoi(args[0])
		if paragraphErr != nil {
			fmt.Fprintf(os.Stderr, "エラー: パラグラフ番号は数値で指定してください: %v\n", paragraphErr)
			return
		}

		description := args[1]

		// 遷移先パラグラフ番号をパース
		targetNum, targetErr := strconv.Atoi(args[2])
		if targetErr != nil {
			fmt.Fprintf(os.Stderr, "エラー: 遷移先パラグラフ番号は数値で指定してください: %v\n", targetErr)
			return
		}

		executor := NewCLIExecutor()
		if err := executor.ExecuteChoiceCommand(paragraphNum, description, targetNum); err != nil {
			fmt.Fprintf(os.Stderr, "エラー: %v\n", err)
			return
		}
	},
}

var selectChoiceCmd = &cobra.Command{
	Use:   "select [パラグラフ番号] [選択肢番号]",
	Short: "パラグラフの選択肢を選択し、遷移先に移動",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		// パラグラフ番号をパース
		paragraphNum, paragraphErr := strconv.Atoi(args[0])
		if paragraphErr != nil {
			fmt.Fprintf(os.Stderr, "エラー: パラグラフ番号は数値で指定してください: %v\n", paragraphErr)
			return
		}

		// 選択肢番号をパース（1ベースから0ベースに変換）
		choiceNum, choiceErr := strconv.Atoi(args[1])
		if choiceErr != nil {
			fmt.Fprintf(os.Stderr, "エラー: 選択肢番号は数値で指定してください: %v\n", choiceErr)
			return
		}
		choiceIndex := choiceNum - 1 // 1ベースから0ベースに変換

		executor := NewCLIExecutor()
		if err := executor.ExecuteSelectCommand(paragraphNum, choiceIndex); err != nil {
			fmt.Fprintf(os.Stderr, "エラー: %v\n", err)
			return
		}
	},
}

var moveCmd = &cobra.Command{
	Use:   "move [パラグラフ番号]",
	Short: "指定されたパラグラフに直接移動（経路上の選択肢を自動選択）",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		targetNum, err := strconv.Atoi(args[0])
		if err != nil {
			fmt.Fprintf(os.Stderr, "パラグラフ番号は数値で指定してください: %v\n", err)
			return
		}

		executor := NewCLIExecutor()
		if err := executor.ExecuteMoveCommand(targetNum); err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
	},
}
