package domain

// UndefinedAnalysis は未定義パラグラフの分析結果
type UndefinedAnalysis struct {
	Connected []int // 現在位置から接続されている未定義パラグラフ
	Orphaned  []int // 孤立した未定義パラグラフ
}

// UndefinedAnalyzer は未定義パラグラフの分析を行う
type UndefinedAnalyzer struct {
	gamebook *Gamebook
}

// NewUndefinedAnalyzer は新しい未定義パラグラフ分析器を作成する
func NewUndefinedAnalyzer(gamebook *Gamebook) *UndefinedAnalyzer {
	return &UndefinedAnalyzer{
		gamebook: gamebook,
	}
}

// AnalyzeUndefinedParagraphs は未定義パラグラフを分析する
func (ua *UndefinedAnalyzer) AnalyzeUndefinedParagraphs(currentPos int) *UndefinedAnalysis {
	// 全ての未定義パラグラフを取得
	allUndefined := ua.getAllUndefinedParagraphs()

	// 現在位置から接続されている未定義パラグラフを検出
	connected := ua.getConnectedUndefined(currentPos, allUndefined)

	// 孤立した未定義パラグラフを検出
	orphaned := ua.getOrphanedUndefined(connected, allUndefined)

	return &UndefinedAnalysis{
		Connected: connected,
		Orphaned:  orphaned,
	}
}

// IsConnectedToUndefined は指定されたパラグラフが現在位置から接続された未定義パラグラフかを判定する
func (ua *UndefinedAnalyzer) IsConnectedToUndefined(currentPos int, targetNumber int) bool {
	// 対象パラグラフが定義済みの場合は false
	if _, exists := ua.gamebook.Paragraphs[targetNumber]; exists {
		return false
	}

	// 現在位置が不正な場合は false
	if currentPos <= 0 {
		return false
	}

	// 現在位置から直接接続されているかチェック
	currentParagraph, exists := ua.gamebook.Paragraphs[currentPos]
	if !exists {
		return false
	}

	// 選択肢をチェック
	for _, choice := range currentParagraph.Choices {
		if choice.TargetNumber == targetNumber {
			return true
		}
	}

	return false
}

// getAllUndefinedParagraphs は全ての未定義パラグラフ番号を取得する
func (ua *UndefinedAnalyzer) getAllUndefinedParagraphs() []int {
	undefinedSet := make(map[int]bool)

	// 全パラグラフの選択肢をチェック
	for _, paragraph := range ua.gamebook.Paragraphs {
		for _, choice := range paragraph.Choices {
			// 遷移先が未定義の場合
			if _, exists := ua.gamebook.Paragraphs[choice.TargetNumber]; !exists {
				undefinedSet[choice.TargetNumber] = true
			}
		}
	}

	// スライスに変換（空の場合も適切に処理）
	undefined := make([]int, 0)
	for number := range undefinedSet {
		undefined = append(undefined, number)
	}

	return undefined
}

// getConnectedUndefined は現在位置から接続されている未定義パラグラフを取得する
func (ua *UndefinedAnalyzer) getConnectedUndefined(currentPos int, allUndefined []int) []int {
	connected := make([]int, 0)

	if currentPos <= 0 {
		return connected
	}

	currentParagraph, exists := ua.gamebook.Paragraphs[currentPos]
	if !exists {
		return connected
	}

	// 未定義パラグラフセットを作成（高速化）
	undefinedSet := make(map[int]bool)
	for _, num := range allUndefined {
		undefinedSet[num] = true
	}

	// 現在位置の選択肢をチェック
	for _, choice := range currentParagraph.Choices {
		if undefinedSet[choice.TargetNumber] {
			connected = append(connected, choice.TargetNumber)
		}
	}

	return connected
}

// getOrphanedUndefined は孤立した未定義パラグラフを取得する
func (ua *UndefinedAnalyzer) getOrphanedUndefined(connected []int, allUndefined []int) []int {
	connectedSet := make(map[int]bool)
	for _, num := range connected {
		connectedSet[num] = true
	}

	orphaned := make([]int, 0)
	for _, num := range allUndefined {
		if !connectedSet[num] {
			orphaned = append(orphaned, num)
		}
	}

	return orphaned
}
