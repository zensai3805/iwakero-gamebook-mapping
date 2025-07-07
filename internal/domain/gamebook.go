package domain

import "fmt"

// Gamebook はゲームブック全体を管理する
type Gamebook struct {
	Title             string
	Paragraphs        map[int]*Paragraph
	Current           *Paragraph
	pendingReferences *PendingReferenceManager // 保留参照管理
	NavigationHistory []NavigationStep         // 移動履歴
}

// NewGamebook は新しいゲームブックを作成する
func NewGamebook(title string) *Gamebook {
	g := &Gamebook{
		Title:             title,
		Paragraphs:        make(map[int]*Paragraph),
		Current:           nil,
		pendingReferences: NewPendingReferenceManager(),
		NavigationHistory: make([]NavigationStep, 0),
	}
	// パラグラフ1を未定義として自動作成し、現在地に設定
	p1 := NewParagraph(1, "(未定義)")
	p1.Visited = true
	g.Paragraphs[1] = p1
	g.Current = p1

	return g
}

// AddParagraph はゲームブックにパラグラフを追加する
func (g *Gamebook) AddParagraph(p *Paragraph) error {
	if existing, exists := g.Paragraphs[p.Number]; exists {
		// プレースホルダーの場合は説明を更新
		if existing.Description == "(未定義)" {
			existing.Description = p.Description
			// 新しいパラグラフの選択肢も追加
			existing.Choices = append(existing.Choices, p.Choices...)
			// 保留参照を解決
			_ = g.pendingReferences.ResolveReference(p.Number)
			return nil
		}
		return ErrDuplicateParagraph
	}
	g.Paragraphs[p.Number] = p

	// 保留参照を解決
	_ = g.pendingReferences.ResolveReference(p.Number)

	return nil
}

// GetParagraph は指定された番号のパラグラフを取得する
func (g *Gamebook) GetParagraph(number int) (*Paragraph, error) {
	p, exists := g.Paragraphs[number]
	if !exists {
		return nil, ErrParagraphNotFound
	}
	return p, nil
}

// MoveTo は指定された番号のパラグラフに移動する
func (g *Gamebook) MoveTo(number int) error {
	p, err := g.GetParagraph(number)
	if err != nil {
		return err
	}
	g.Current = p
	p.Visited = true
	return nil
}

// MoveToOrCreatePlaceholder は指定されたパラグラフに移動。存在しない場合はプレースホルダーを作成
func (g *Gamebook) MoveToOrCreatePlaceholder(number int) error {
	// 既存のパラグラフがあれば通常の移動
	if p, exists := g.Paragraphs[number]; exists {
		g.Current = p
		p.Visited = true
		return nil
	}

	// 未定義パラグラフの場合、プレースホルダーを作成
	placeholder := NewParagraph(number, "(未定義)")
	placeholder.Visited = true
	g.Paragraphs[number] = placeholder
	g.Current = placeholder

	return nil
}

// GetExplorationStats は探索統計を返す
func (g *Gamebook) GetExplorationStats() ExplorationStats {
	totalParagraphs := len(g.Paragraphs)
	visitedParagraphs := 0
	totalChoices := 0
	selectedChoices := 0

	for _, p := range g.Paragraphs {
		if p.Visited {
			visitedParagraphs++
		}
		for _, c := range p.Choices {
			totalChoices++
			if c.Selected {
				selectedChoices++
			}
		}
	}

	return ExplorationStats{
		TotalParagraphs:   totalParagraphs,
		VisitedParagraphs: visitedParagraphs,
		TotalChoices:      totalChoices,
		SelectedChoices:   selectedChoices,
	}
}

// AddChoiceToParagraph は指定されたパラグラフに選択肢を追加する
func (g *Gamebook) AddChoiceToParagraph(paragraphNumber int, description string, targetNumber int) error {
	p, err := g.GetParagraph(paragraphNumber)
	if err != nil {
		return err
	}

	// 選択肢を追加
	p.AddChoice(description, targetNumber)

	// 遷移先が未定義の場合は保留参照として記録
	if _, targetExists := g.Paragraphs[targetNumber]; !targetExists {
		_ = g.pendingReferences.AddReference(paragraphNumber, description, targetNumber)
	}

	return nil
}

// SelectChoiceAndMove は選択肢を選択し、対象パラグラフに移動する
func (g *Gamebook) SelectChoiceAndMove(paragraphNumber int, choiceIndex int) error {
	p, err := g.GetParagraph(paragraphNumber)
	if err != nil {
		return err
	}

	if err := p.SelectChoice(choiceIndex); err != nil {
		return err
	}

	// 選択された選択肢の遷移先に移動
	targetNumber := p.Choices[choiceIndex].TargetNumber
	return g.MoveTo(targetNumber)
}

// ExplorationStats は探索の統計情報
type ExplorationStats struct {
	TotalParagraphs   int
	VisitedParagraphs int
	TotalChoices      int
	SelectedChoices   int
}

// MoveResult は移動結果を表す
type MoveResult struct {
	Success        bool   // 移動が成功したか
	HasWarning     bool   // 警告があるか
	WarningMessage string // 警告メッセージ
}

// GetAllPendingTargets は全ての保留対象段落番号を取得する
func (g *Gamebook) GetAllPendingTargets() []int {
	return g.pendingReferences.GetAllPendingTargets()
}

// MoveToWithPathSelection は指定されたパラグラフに直接移動し、経路上の選択肢を自動選択する
func (g *Gamebook) MoveToWithPathSelection(targetNumber int) error {
	// 目的地が存在しない場合はプレースホルダーを作成して移動
	_, exists := g.Paragraphs[targetNumber]
	if !exists {
		return g.MoveToOrCreatePlaceholder(targetNumber)
	}

	// 現在地が未設定の場合は直接移動
	if g.Current == nil {
		return g.MoveTo(targetNumber)
	}

	// 同じ位置の場合は何もしない
	if g.Current.Number == targetNumber {
		return nil
	}

	// 最短経路を探索
	path := g.findShortestPath(g.Current.Number, targetNumber)
	if len(path) == 0 {
		// 経路が見つからない場合は直接移動
		return g.MoveTo(targetNumber)
	}

	// 経路上の選択肢を選択して移動
	for i := 0; i < len(path)-1; i++ {
		currentNum := path[i]
		nextNum := path[i+1]

		// 現在パラグラフから次のパラグラフへの選択肢を探して選択
		currentParagraph := g.Paragraphs[currentNum]
		for choiceIndex, choice := range currentParagraph.Choices {
			if choice.TargetNumber == nextNum {
				// 選択肢を選択
				if err := currentParagraph.SelectChoice(choiceIndex); err != nil {
					return fmt.Errorf("選択肢の選択に失敗: %w", err)
				}
				// 次のパラグラフに移動
				if err := g.MoveTo(nextNum); err != nil {
					return fmt.Errorf("パラグラフ %d への移動に失敗: %w", nextNum, err)
				}
				break
			}
		}
	}

	return nil
}

// findShortestPath BFSで最短経路を探索
func (g *Gamebook) findShortestPath(start, target int) []int {
	if start == target {
		return []int{start}
	}

	// BFSのためのキューと訪問済みセット
	queue := [][]int{{start}}
	visited := make(map[int]bool)
	visited[start] = true

	for len(queue) > 0 {
		path := queue[0]
		queue = queue[1:]
		current := path[len(path)-1]

		// 現在パラグラフからの全選択肢を確認
		currentParagraph, exists := g.Paragraphs[current]
		if !exists {
			continue
		}

		for _, choice := range currentParagraph.Choices {
			next := choice.TargetNumber

			// 目的地に到達
			if next == target {
				return append(path, next)
			}

			// 未訪問かつ存在するパラグラフのみをキューに追加
			if !visited[next] {
				if _, nextExists := g.Paragraphs[next]; nextExists {
					visited[next] = true
					newPath := make([]int, len(path)+1)
					copy(newPath, path)
					newPath[len(path)] = next
					queue = append(queue, newPath)
				}
			}
		}
	}

	// 経路が見つからない
	return []int{}
}

// SelectChoiceAndMoveWithGracefulHandling は選択肢を選択し、優雅な移動処理を行う
func (g *Gamebook) SelectChoiceAndMoveWithGracefulHandling(paragraphNumber int, choiceIndex int) MoveResult {
	p, err := g.GetParagraph(paragraphNumber)
	if err != nil {
		return MoveResult{
			Success:        false,
			HasWarning:     true,
			WarningMessage: "段落が見つかりません: " + err.Error(),
		}
	}

	if selectErr := p.SelectChoice(choiceIndex); selectErr != nil {
		return MoveResult{
			Success:        false,
			HasWarning:     true,
			WarningMessage: "選択肢の選択に失敗: " + selectErr.Error(),
		}
	}

	// 選択された選択肢の遷移先に移動を試行（未定義でもプレースホルダー作成）
	targetNumber := p.Choices[choiceIndex].TargetNumber
	moveErr := g.MoveToOrCreatePlaceholder(targetNumber)

	if moveErr != nil {
		// エラーが発生した場合（通常は起こらない）
		return MoveResult{
			Success:        false,
			HasWarning:     true,
			WarningMessage: fmt.Sprintf("移動エラー: %v", moveErr),
		}
	}

	// 移動成功（プレースホルダーが作成された場合も含む）
	if g.Paragraphs[targetNumber].Description == "(未定義)" {
		return MoveResult{
			Success:        true,
			HasWarning:     true,
			WarningMessage: fmt.Sprintf("段落%dは未定義でしたが、プレースホルダーを作成して移動しました。", targetNumber),
		}
	}

	return MoveResult{
		Success:        true,
		HasWarning:     false,
		WarningMessage: "",
	}
}

// GetNavigationHistory は移動履歴を取得する
func (g *Gamebook) GetNavigationHistory() []NavigationStep {
	return g.NavigationHistory
}

// AddNavigationStep は移動履歴を追加する
func (g *Gamebook) AddNavigationStep(step NavigationStep) error {
	g.NavigationHistory = append(g.NavigationHistory, step)
	return nil
}
