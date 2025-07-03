package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/iwapc/iwakero-gamebook-mapping/internal/domain"
)

const (
	// 選択肢ステータス定数
	StatusUnselected = "未選択"
	StatusSelected   = "選択済み"
)

// CommandExecutor は既存のCLIコマンドを実行するためのインターフェース
type CommandExecutor interface {
	ExecuteNewCommand(title string) error
	ExecuteLoadCommand(title string) error
	ExecuteAddCommand(number int, description string) error
	ExecuteChoiceCommand(paragraphNum int, description string, targetNum int) error
	ExecuteSelectCommand(paragraphNum int, choiceIndex int) error
	ExecuteShowCommand() error
}

// CLIExecutor は実際のCLIコマンドを実行する
type CLIExecutor struct{}

// NewCLIExecutor は新しいCLIExecutorを作成
func NewCLIExecutor() *CLIExecutor {
	return &CLIExecutor{}
}

// ExecuteNewCommand newコマンドの実装を実行
func (e *CLIExecutor) ExecuteNewCommand(title string) error {
	// 既存のnewコマンドロジックを呼び出し
	currentGame = domain.NewGamebook(title)
	if err := repo.Save(currentGame); err != nil {
		return fmt.Errorf("エラー: %v", err)
	}

	// 現在のゲームとして保存
	if err := sessionRepo.SaveCurrentGame(title); err != nil {
		fmt.Fprintf(os.Stderr, "セッション保存エラー: %v\n", err)
	}

	fmt.Printf("新しいゲームブック「%s」を作成しました。\n", title)
	return nil
}

// ExecuteLoadCommand loadコマンドの実装を実行
func (e *CLIExecutor) ExecuteLoadCommand(title string) error {
	gamebook, err := repo.Load(title)
	if err != nil {
		return fmt.Errorf("ゲームブックの読み込みに失敗しました: %v", err)
	}
	currentGame = gamebook

	if err := sessionRepo.SaveCurrentGame(title); err != nil {
		fmt.Fprintf(os.Stderr, "セッション保存エラー: %v\n", err)
	}

	fmt.Printf("ゲームブック '%s' を読み込みました\n", title)
	return nil
}

// ExecuteAddCommand addコマンドの実装を実行
func (e *CLIExecutor) ExecuteAddCommand(number int, description string) error {
	if currentGame == nil {
		return fmt.Errorf("エラー: ゲームブックが選択されていません。'gamebook new'または'gamebook load'を実行してください。")
	}

	p := domain.NewParagraph(number, description)
	if err := currentGame.AddParagraph(p); err != nil {
		return fmt.Errorf("エラー: %v", err)
	}

	if err := repo.Save(currentGame); err != nil {
		return fmt.Errorf("保存エラー: %v", err)
	}

	fmt.Printf("パラグラフ %d を追加しました: %s\n", number, description)
	return nil
}

// ExecuteChoiceCommand choiceコマンドの実装を実行
func (e *CLIExecutor) ExecuteChoiceCommand(paragraphNum int, description string, targetNum int) error {
	if currentGame == nil {
		return fmt.Errorf("エラー: ゲームブックが選択されていません。")
	}

	// 選択肢を追加
	if err := currentGame.AddChoiceToParagraph(paragraphNum, description, targetNum); err != nil {
		return fmt.Errorf("エラー: %v", err)
	}

	// 保存
	if err := repo.Save(currentGame); err != nil {
		return fmt.Errorf("保存エラー: %v", err)
	}

	fmt.Printf("パラグラフ %d に選択肢「%s → %d」を追加しました。\n", paragraphNum, description, targetNum)
	return nil
}

// ExecuteSelectCommand selectコマンドの実装を実行
func (e *CLIExecutor) ExecuteSelectCommand(paragraphNum int, choiceIndex int) error {
	if currentGame == nil {
		return fmt.Errorf("エラー: ゲームブックが選択されていません。")
	}

	// 選択肢を選択して移動
	if err := currentGame.SelectChoiceAndMove(paragraphNum, choiceIndex); err != nil {
		return fmt.Errorf("エラー: %v", err)
	}

	// 保存
	if err := repo.Save(currentGame); err != nil {
		return fmt.Errorf("保存エラー: %v", err)
	}

	fmt.Printf("パラグラフ %d の選択肢 %d を選択し、パラグラフ %d に移動しました。\n",
		paragraphNum, choiceIndex+1, currentGame.Current.Number)
	return nil
}

// ExecuteShowCommand showコマンドの実装を実行
func (e *CLIExecutor) ExecuteShowCommand() error {
	if currentGame == nil {
		return fmt.Errorf("エラー: ゲームブックが選択されていません。")
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
	return nil
}

// ParseArguments は引数をパースするヘルパー関数
func ParseArguments(args []string, minArgs int, usage string) ([]string, error) {
	if len(args) < minArgs {
		return nil, fmt.Errorf("使用法: %s", usage)
	}
	return args, nil
}

// ParseNumber は文字列を数値にパースするヘルパー関数
func ParseNumber(str string, name string) (int, error) {
	num, err := strconv.Atoi(str)
	if err != nil {
		return 0, fmt.Errorf("%sは数値で指定してください: %v", name, err)
	}
	return num, nil
}

// JoinDescription は複数の引数を説明文として結合するヘルパー関数
func JoinDescription(args []string, startIndex int) string {
	if startIndex >= len(args) {
		return ""
	}
	return strings.Join(args[startIndex:], " ")
}
