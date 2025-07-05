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

	// 統計情報（簡潔な1行表示）
	stats := currentGame.GetExplorationStats()
	fmt.Printf("📊 パラグラフ %d/%d | 選択肢 %d/%d | 訪問済み: %d\n",
		stats.VisitedParagraphs, stats.TotalParagraphs,
		stats.SelectedChoices, stats.TotalChoices,
		stats.VisitedParagraphs)

	// 現在位置（簡潔表示）
	if currentGame.Current != nil {
		fmt.Printf("\n📍 現在: #%d %s\n", currentGame.Current.Number, currentGame.Current.Description)

		if len(currentGame.Current.Choices) > 0 {
			fmt.Printf("🎯 選択肢: ")
			for i, choice := range currentGame.Current.Choices {
				status := "⚪"
				if choice.Selected {
					status = "✅"
				}
				if i > 0 {
					fmt.Printf(" | ")
				}
				fmt.Printf("%s %d.%s→#%d", status, i+1, choice.Description, choice.TargetNumber)
			}
			fmt.Println()
		}
	}

	// 他のパラグラフ（現在位置以外の最近追加分）
	fmt.Println("\n📚 その他のパラグラフ:")
	count := 0
	currentNum := 0
	if currentGame.Current != nil {
		currentNum = currentGame.Current.Number
	}

	for i := 1; i <= 1000 && count < 3; i++ {
		if p, exists := currentGame.Paragraphs[i]; exists && p.Number != currentNum {
			status := "⚪"
			if p.Visited {
				status = "✅"
			}
			choiceCount := len(p.Choices)
			selectedCount := 0
			for _, choice := range p.Choices {
				if choice.Selected {
					selectedCount++
				}
			}

			if choiceCount > 0 {
				fmt.Printf("  %s #%d %s (選択肢 %d/%d)\n", status, p.Number, p.Description, selectedCount, choiceCount)
			} else {
				fmt.Printf("  %s #%d %s\n", status, p.Number, p.Description)
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
