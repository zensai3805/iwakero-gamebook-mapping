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
		if currentGame == nil {
			fmt.Fprintln(os.Stderr, "エラー: ゲームブックが選択されていません。")
			return
		}

		// パラグラフ番号をパース
		paragraphNum, err := strconv.Atoi(args[0])
		if err != nil {
			fmt.Fprintf(os.Stderr, "エラー: パラグラフ番号は数値で指定してください: %v\n", err)
			return
		}

		description := args[1]

		// 遷移先パラグラフ番号をパース
		targetNum, err := strconv.Atoi(args[2])
		if err != nil {
			fmt.Fprintf(os.Stderr, "エラー: 遷移先パラグラフ番号は数値で指定してください: %v\n", err)
			return
		}

		// 選択肢を追加
		if err := currentGame.AddChoiceToParagraph(paragraphNum, description, targetNum); err != nil {
			fmt.Fprintf(os.Stderr, "エラー: %v\n", err)
			return
		}

		// 保存
		if err := repo.Save(currentGame); err != nil {
			fmt.Fprintf(os.Stderr, "保存エラー: %v\n", err)
			return
		}

		fmt.Printf("パラグラフ %d に選択肢「%s → %d」を追加しました。\n", paragraphNum, description, targetNum)
	},
}

var selectChoiceCmd = &cobra.Command{
	Use:   "select [パラグラフ番号] [選択肢番号]",
	Short: "パラグラフの選択肢を選択し、遷移先に移動",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		if currentGame == nil {
			fmt.Fprintln(os.Stderr, "エラー: ゲームブックが選択されていません。")
			return
		}

		// パラグラフ番号をパース
		paragraphNum, err := strconv.Atoi(args[0])
		if err != nil {
			fmt.Fprintf(os.Stderr, "エラー: パラグラフ番号は数値で指定してください: %v\n", err)
			return
		}

		// 選択肢番号をパース（1ベースから0ベースに変換）
		choiceNum, err := strconv.Atoi(args[1])
		if err != nil {
			fmt.Fprintf(os.Stderr, "エラー: 選択肢番号は数値で指定してください: %v\n", err)
			return
		}
		choiceIndex := choiceNum - 1 // 1ベースから0ベースに変換

		// 優雅なエラーハンドリングを使用して選択肢を選択・移動
		moveResult := currentGame.SelectChoiceAndMoveWithGracefulHandling(paragraphNum, choiceIndex)

		if !moveResult.Success {
			fmt.Fprintf(os.Stderr, "エラー: %s\n", moveResult.WarningMessage)
			return
		}

		// 保存
		if err := repo.Save(currentGame); err != nil {
			fmt.Fprintf(os.Stderr, "保存エラー: %v\n", err)
			return
		}

		// 結果の表示
		if moveResult.HasWarning {
			fmt.Printf("⚠️  警告: %s\n", moveResult.WarningMessage)
			fmt.Printf("パラグラフ %d の選択肢 %d を選択しました（移動先が未定義のため現在位置を維持）。\n",
				paragraphNum, choiceNum)
		} else {
			fmt.Printf("パラグラフ %d の選択肢 %d を選択し、パラグラフ %d に移動しました。\n",
				paragraphNum, choiceNum, currentGame.Current.Number)
		}
	},
}

var moveCmd = &cobra.Command{
	Use:   "move [パラグラフ番号]",
	Short: "指定されたパラグラフに移動（選択なし）",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if currentGame == nil {
			fmt.Fprintln(os.Stderr, "エラー: ゲームブックが選択されていません。")
			return
		}

		// パラグラフ番号をパース
		paragraphNum, err := strconv.Atoi(args[0])
		if err != nil {
			fmt.Fprintf(os.Stderr, "エラー: パラグラフ番号は数値で指定してください: %v\n", err)
			return
		}

		// 移動
		if err := currentGame.MoveTo(paragraphNum); err != nil {
			fmt.Fprintf(os.Stderr, "エラー: %v\n", err)
			return
		}

		// 保存
		if err := repo.Save(currentGame); err != nil {
			fmt.Fprintf(os.Stderr, "保存エラー: %v\n", err)
			return
		}

		fmt.Printf("パラグラフ %d に移動しました。\n", paragraphNum)
	},
}
