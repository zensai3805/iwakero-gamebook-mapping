package domain

// Paragraph はゲームブックの1つのパラグラフを表す
type Paragraph struct {
	Number      int
	Description string
	Choices     []Choice
	Visited     bool
}

// Choice はパラグラフで選択可能な選択肢を表す
type Choice struct {
	Description  string
	TargetNumber int
	Selected     bool
	TotalChoices int // このパラグラフの選択肢総数（省略可能）
}

// NewParagraph は新しいパラグラフを作成する
func NewParagraph(number int, description string) *Paragraph {
	return &Paragraph{
		Number:      number,
		Description: description,
		Choices:     []Choice{},
		Visited:     false,
	}
}

// AddChoice はパラグラフに選択肢を追加する
func (p *Paragraph) AddChoice(description string, targetNumber int) {
	choice := Choice{
		Description:  description,
		TargetNumber: targetNumber,
		Selected:     false,
	}
	p.Choices = append(p.Choices, choice)
}

// SelectChoice は指定された選択肢を選択済みにする
func (p *Paragraph) SelectChoice(index int) error {
	if index < 0 || index >= len(p.Choices) {
		return ErrInvalidChoiceIndex
	}
	p.Choices[index].Selected = true
	p.Visited = true
	return nil
}

// GetUnselectedChoices は未選択の選択肢を返す
func (p *Paragraph) GetUnselectedChoices() []Choice {
	var unselected []Choice
	for _, choice := range p.Choices {
		if !choice.Selected {
			unselected = append(unselected, choice)
		}
	}
	return unselected
}
