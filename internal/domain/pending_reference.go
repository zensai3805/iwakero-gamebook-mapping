package domain

// PendingReference は未定義段落への参照を表す
type PendingReference struct {
	FromParagraph     int    // 参照元段落番号
	ChoiceDescription string // 選択肢の説明
	TargetParagraph   int    // 参照先段落番号
}

// PendingReferenceManager は未定義段落への参照を管理する
type PendingReferenceManager struct {
	references map[int][]PendingReference // 対象段落番号 -> 参照のリスト
}

// NewPendingReferenceManager は新しい保留参照管理を作成する
func NewPendingReferenceManager() *PendingReferenceManager {
	return &PendingReferenceManager{
		references: make(map[int][]PendingReference),
	}
}

// AddReference は保留参照を追加する
func (m *PendingReferenceManager) AddReference(fromParagraph int, choiceDescription string, targetParagraph int) error {
	ref := PendingReference{
		FromParagraph:     fromParagraph,
		ChoiceDescription: choiceDescription,
		TargetParagraph:   targetParagraph,
	}

	m.references[targetParagraph] = append(m.references[targetParagraph], ref)
	return nil
}

// GetReferences は指定された対象段落への保留参照を取得する
func (m *PendingReferenceManager) GetReferences(targetParagraph int) []PendingReference {
	refs, exists := m.references[targetParagraph]
	if !exists {
		return []PendingReference{}
	}
	// スライスのコピーを返して外部からの変更を防ぐ
	result := make([]PendingReference, len(refs))
	copy(result, refs)
	return result
}

// ResolveReference は保留参照を解決して削除する
func (m *PendingReferenceManager) ResolveReference(targetParagraph int) error {
	delete(m.references, targetParagraph)
	return nil
}

// GetAllPendingTargets は全ての保留対象段落番号を取得する
func (m *PendingReferenceManager) GetAllPendingTargets() []int {
	var targets []int
	for target := range m.references {
		targets = append(targets, target)
	}
	return targets
}

// HasPendingReferences は保留参照が存在するかを確認する
func (m *PendingReferenceManager) HasPendingReferences() bool {
	return len(m.references) > 0
}
