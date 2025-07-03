package domain

// Gamebook はゲームブック全体を管理する
type Gamebook struct {
	Title      string
	Paragraphs map[int]*Paragraph
	Current    *Paragraph
}

// NewGamebook は新しいゲームブックを作成する
func NewGamebook(title string) *Gamebook {
	return &Gamebook{
		Title:      title,
		Paragraphs: make(map[int]*Paragraph),
		Current:    nil,
	}
}

// AddParagraph はゲームブックにパラグラフを追加する
func (g *Gamebook) AddParagraph(p *Paragraph) error {
	if _, exists := g.Paragraphs[p.Number]; exists {
		return ErrDuplicateParagraph
	}
	g.Paragraphs[p.Number] = p
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
	p.AddChoice(description, targetNumber)
	return nil
}

// AddChoiceToParagraphWithValidation は遷移先検証付きで選択肢を追加する
func (g *Gamebook) AddChoiceToParagraphWithValidation(paragraphNumber int, description string, targetNumber int) error {
	p, err := g.GetParagraph(paragraphNumber)
	if err != nil {
		return err
	}

	// 遷移先パラグラフの存在確認
	if _, exists := g.Paragraphs[targetNumber]; !exists {
		// 警告を返すが、選択肢は追加する（SPECIFICATION.mdの要件に従い即座に通知）
		p.AddChoice(description, targetNumber)
		return ErrUndefinedTargetParagraph
	}

	p.AddChoice(description, targetNumber)
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
