package domain

// NavigationStep は移動履歴の1ステップを表す
type NavigationStep struct {
	From     int   // 移動元パラグラフ番号
	To       int   // 移動先パラグラフ番号
	ViaPaths []int // 経路上の中間パラグラフ（選択肢移動時は空、ジャンプ時は経路）
}

// NewNavigationStep は新しいNavigationStepを作成する
func NewNavigationStep(from, to int, viaPaths []int) *NavigationStep {
	return &NavigationStep{
		From:     from,
		To:       to,
		ViaPaths: viaPaths,
	}
}

// NewNavigationStepWithValidation はバリデーション付きでNavigationStepを作成する
func NewNavigationStepWithValidation(from, to int, viaPaths []int) (*NavigationStep, error) {
	if from < 1 {
		return nil, ErrInvalidParagraphNumber
	}
	if to < 1 {
		return nil, ErrInvalidParagraphNumber
	}
	if from == to {
		return nil, ErrSameFromToNavigation
	}
	return &NavigationStep{
		From:     from,
		To:       to,
		ViaPaths: viaPaths,
	}, nil
}
