package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/iwapc/iwakero-gamebook-mapping/internal/domain"
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

		// 選択肢を追加（遷移先検証付き）
		if err := currentGame.AddChoiceToParagraphWithValidation(paragraphNum, description, targetNum); err != nil {
			if err == domain.ErrUndefinedTargetParagraph {
				// 警告だが処理は継続（選択肢は追加済み）
				fmt.Fprintf(os.Stderr, "%v\n", err)
			} else {
				// その他のエラーは処理を停止
				fmt.Fprintf(os.Stderr, "エラー: %v\n", err)
				return
			}
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

		// 選択肢を選択して移動
		if err := currentGame.SelectChoiceAndMove(paragraphNum, choiceIndex); err != nil {
			fmt.Fprintf(os.Stderr, "エラー: %v\n", err)
			return
		}

		// 保存
		if err := repo.Save(currentGame); err != nil {
			fmt.Fprintf(os.Stderr, "保存エラー: %v\n", err)
			return
		}

		fmt.Printf("パラグラフ %d の選択肢 %d を選択し、パラグラフ %d に移動しました。\n", 
			paragraphNum, choiceNum, currentGame.Current.Number)
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