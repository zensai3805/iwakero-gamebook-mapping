package main

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/pterm/pterm"
	"github.com/zensai3805/iwakero-gamebook-mapping/internal/domain"
)

// CandidateType 候補の種類
type CandidateType int

const (
	CandidateTypeExisting CandidateType = iota
	CandidateTypeUndefined
	CandidateTypeRecommended
)

// Candidate 候補情報
type Candidate struct {
	Value       string
	Type        CandidateType
	Description string
	Priority    int
	Color       pterm.Color
}

// ParagraphSuggestions パラグラフ候補情報
type ParagraphSuggestions struct {
	Existing    []int // 既存パラグラフ
	Undefined   []int // 未定義パラグラフ
	Recommended []int // 推奨番号
	Next        int   // 次の論理的番号
}

// SuggestionEngine 候補生成エンジン
type SuggestionEngine struct {
	gamebook *domain.Gamebook
}

// NewSuggestionEngine 候補生成エンジンを作成
func NewSuggestionEngine(gamebook *domain.Gamebook) *SuggestionEngine {
	return &SuggestionEngine{
		gamebook: gamebook,
	}
}

// GenerateParagraphSuggestions パラグラフ候補を生成
func (se *SuggestionEngine) GenerateParagraphSuggestions() *ParagraphSuggestions {
	suggestions := &ParagraphSuggestions{
		Existing:    []int{},
		Undefined:   []int{},
		Recommended: []int{},
		Next:        1,
	}

	if se.gamebook == nil {
		return suggestions
	}

	// 既存パラグラフを取得
	maxNum := 0
	for num, p := range se.gamebook.Paragraphs {
		suggestions.Existing = append(suggestions.Existing, num)
		if p.Number > maxNum {
			maxNum = p.Number
		}
	}

	// 未定義パラグラフを取得
	undefinedNums := make(map[int]bool)
	for _, p := range se.gamebook.Paragraphs {
		for _, choice := range p.Choices {
			if _, err := se.gamebook.GetParagraph(choice.TargetNumber); err != nil {
				undefinedNums[choice.TargetNumber] = true
			}
		}
	}

	for num := range undefinedNums {
		suggestions.Undefined = append(suggestions.Undefined, num)
	}

	// 推奨番号を生成（現在の最大値+5, +10, +15）
	for i := 5; i <= 15; i += 5 {
		suggestions.Recommended = append(suggestions.Recommended, maxNum+i)
	}

	// 次の論理的番号を設定
	suggestions.Next = maxNum + 1

	// ソート
	sort.Ints(suggestions.Existing)
	sort.Ints(suggestions.Undefined)
	sort.Ints(suggestions.Recommended)

	return suggestions
}

// TabCompleter Tab補完機能
type TabCompleter struct {
	engine *SuggestionEngine
}

// NewTabCompleter Tab補完機能を作成
func NewTabCompleter(engine *SuggestionEngine) *TabCompleter {
	return &TabCompleter{
		engine: engine,
	}
}

// Complete Tab補完候補を生成
func (tc *TabCompleter) Complete(input string) []Candidate {
	suggestions := tc.engine.GenerateParagraphSuggestions()
	candidates := []Candidate{}

	// 入力が空の場合は全候補を返す
	if input == "" {
		// 既存パラグラフ
		for _, num := range suggestions.Existing {
			candidates = append(candidates, Candidate{
				Value:       strconv.Itoa(num),
				Type:        CandidateTypeExisting,
				Description: "既存",
				Priority:    1,
				Color:       pterm.FgGreen,
			})
		}

		// 未定義パラグラフ
		for _, num := range suggestions.Undefined {
			candidates = append(candidates, Candidate{
				Value:       strconv.Itoa(num),
				Type:        CandidateTypeUndefined,
				Description: "未定義",
				Priority:    2,
				Color:       pterm.FgYellow,
			})
		}

		// 推奨パラグラフ
		for _, num := range suggestions.Recommended {
			candidates = append(candidates, Candidate{
				Value:       strconv.Itoa(num),
				Type:        CandidateTypeRecommended,
				Description: "推奨",
				Priority:    3,
				Color:       pterm.FgLightBlue,
			})
		}

		return candidates
	}

	// 入力に基づく候補フィルタリング
	allNums := make([]int, 0, len(suggestions.Existing)+len(suggestions.Undefined)+len(suggestions.Recommended))
	allNums = append(allNums, suggestions.Existing...)
	allNums = append(allNums, suggestions.Undefined...)
	allNums = append(allNums, suggestions.Recommended...)

	for _, num := range allNums {
		numStr := strconv.Itoa(num)
		if strings.HasPrefix(numStr, input) {
			candidateType := CandidateTypeExisting
			description := "既存"
			color := pterm.FgGreen
			priority := 1

			// 種類を判定
			if containsInt(suggestions.Undefined, num) {
				candidateType = CandidateTypeUndefined
				description = "未定義"
				color = pterm.FgYellow
				priority = 2
			} else if containsInt(suggestions.Recommended, num) {
				candidateType = CandidateTypeRecommended
				description = "推奨"
				color = pterm.FgLightBlue
				priority = 3
			}

			candidates = append(candidates, Candidate{
				Value:       numStr,
				Type:        candidateType,
				Description: description,
				Priority:    priority,
				Color:       color,
			})
		}
	}

	// 優先度でソート
	sort.Slice(candidates, func(i, j int) bool {
		return candidates[i].Priority < candidates[j].Priority
	})

	return candidates
}

// ValidationType 検証結果の種類
type ValidationType int

const (
	ValidationTypeValid ValidationType = iota
	ValidationTypeDuplicate
	ValidationTypeUndefined
	ValidationTypeInvalid
)

// ValidationResult 検証結果
type ValidationResult struct {
	IsValid     bool
	Type        ValidationType
	Message     string
	Suggestions []string
}

// InputValidator 入力検証機能
type InputValidator struct {
	gamebook *domain.Gamebook
}

// NewInputValidator 入力検証機能を作成
func NewInputValidator(gamebook *domain.Gamebook) *InputValidator {
	return &InputValidator{
		gamebook: gamebook,
	}
}

// ValidateParagraphNumber パラグラフ番号を検証
func (iv *InputValidator) ValidateParagraphNumber(num int) *ValidationResult {
	result := &ValidationResult{
		IsValid:     true,
		Type:        ValidationTypeValid,
		Message:     "",
		Suggestions: []string{},
	}

	// 無効な番号のチェック
	if num <= 0 {
		result.IsValid = false
		result.Type = ValidationTypeInvalid
		result.Message = "パラグラフ番号は1以上である必要があります"
		return result
	}

	// 既存パラグラフの重複チェック
	if iv.gamebook != nil {
		if _, err := iv.gamebook.GetParagraph(num); err == nil {
			result.IsValid = false
			result.Type = ValidationTypeDuplicate
			result.Message = fmt.Sprintf("パラグラフ %d は既に存在します", num)

			// 代替案を提案
			engine := NewSuggestionEngine(iv.gamebook)
			suggestions := engine.GenerateParagraphSuggestions()
			result.Suggestions = append(result.Suggestions, fmt.Sprintf("次の番号: %d", suggestions.Next))

			return result
		}
	}

	return result
}

// HintRenderer ヒント表示機能
type HintRenderer struct {
	suggestions *ParagraphSuggestions
}

// NewHintRenderer ヒント表示機能を作成
func NewHintRenderer(suggestions *ParagraphSuggestions) *HintRenderer {
	return &HintRenderer{
		suggestions: suggestions,
	}
}

// RenderHints ヒント文字列を生成
func (hr *HintRenderer) RenderHints(_ string) string {
	if hr.suggestions == nil {
		return ""
	}

	var hints []string

	// 既存パラグラフのヒント
	if len(hr.suggestions.Existing) > 0 {
		hints = append(hints, pterm.FgGreen.Sprintf("既存: %v", hr.suggestions.Existing[:min(3, len(hr.suggestions.Existing))]))
	}

	// 未定義パラグラフのヒント
	if len(hr.suggestions.Undefined) > 0 {
		hints = append(hints, pterm.FgYellow.Sprintf("未定義: %v", hr.suggestions.Undefined[:min(3, len(hr.suggestions.Undefined))]))
	}

	// 推奨パラグラフのヒント
	if len(hr.suggestions.Recommended) > 0 {
		hints = append(hints, pterm.FgLightBlue.Sprintf("推奨: %v", hr.suggestions.Recommended[:min(3, len(hr.suggestions.Recommended))]))
	}

	// 次の論理的番号のヒント
	hints = append(hints, pterm.FgCyan.Sprintf("次: %d", hr.suggestions.Next))

	return strings.Join(hints, " | ")
}

// EnhancedInput 拡張入力コンポーネント
type EnhancedInput struct {
	suggestionEngine *SuggestionEngine
	tabCompleter     *TabCompleter
	hintRenderer     *HintRenderer
	validator        *InputValidator
}

// NewEnhancedInput 拡張入力コンポーネントを作成
func NewEnhancedInput(gamebook *domain.Gamebook) *EnhancedInput {
	engine := NewSuggestionEngine(gamebook)
	suggestions := engine.GenerateParagraphSuggestions()

	return &EnhancedInput{
		suggestionEngine: engine,
		tabCompleter:     NewTabCompleter(engine),
		hintRenderer:     NewHintRenderer(suggestions),
		validator:        NewInputValidator(gamebook),
	}
}

// ShowWithSuggestions 候補表示付きの入力を表示
func (ei *EnhancedInput) ShowWithSuggestions(prompt string) (string, error) {
	// 候補とヒントを表示
	hints := ei.hintRenderer.RenderHints("")

	// ヒントを表示
	if hints != "" {
		pterm.Info.Println("候補: " + hints)
	}

	// 通常の入力を表示
	input, err := pterm.DefaultInteractiveTextInput.
		WithDefaultText("").
		WithTextStyle(pterm.NewStyle(pterm.FgLightBlue)).
		Show(prompt)

	if err != nil {
		return "", err
	}

	// 入力値を検証
	if num, parseErr := strconv.Atoi(input); parseErr == nil {
		if validationResult := ei.validator.ValidateParagraphNumber(num); !validationResult.IsValid {
			pterm.Warning.Println(validationResult.Message)
			if len(validationResult.Suggestions) > 0 {
				pterm.Info.Println("提案: " + strings.Join(validationResult.Suggestions, ", "))
			}
		}
	}

	return input, nil
}

// ヘルパー関数
func containsInt(slice []int, item int) bool {
	for _, v := range slice {
		if v == item {
			return true
		}
	}
	return false
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
