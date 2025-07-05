package main

// DisplayScopeManager は表示スコープの管理を行う
type DisplayScopeManager struct {
	currentScope DisplayScope
}

// NewDisplayScopeManager は新しい表示スコープ管理器を作成する
func NewDisplayScopeManager() *DisplayScopeManager {
	return &DisplayScopeManager{
		currentScope: ScopeConnected, // デフォルトは接続された未定義パラグラフのみ表示
	}
}

// SetScope は表示スコープを設定する
func (dsm *DisplayScopeManager) SetScope(scope DisplayScope) {
	dsm.currentScope = scope
}

// GetScope は現在の表示スコープを取得する
func (dsm *DisplayScopeManager) GetScope() DisplayScope {
	return dsm.currentScope
}

// CycleScope は表示スコープを順次切り替える（Connected → All → None → Connected）
func (dsm *DisplayScopeManager) CycleScope() DisplayScope {
	switch dsm.currentScope {
	case ScopeConnected:
		dsm.currentScope = ScopeAll
	case ScopeAll:
		dsm.currentScope = ScopeNone
	case ScopeNone:
		dsm.currentScope = ScopeConnected
	default:
		// 不正な値の場合はデフォルトにリセット
		dsm.currentScope = ScopeConnected
	}

	return dsm.currentScope
}

// GetScopeDescription は現在のスコープの日本語説明を取得する
func (dsm *DisplayScopeManager) GetScopeDescription() string {
	switch dsm.currentScope {
	case ScopeConnected:
		return "接続のみ"
	case ScopeAll:
		return "全て表示"
	case ScopeNone:
		return "非表示"
	default:
		return "不明"
	}
}

// IsDefaultScope は現在のスコープがデフォルトかどうかを判定する
func (dsm *DisplayScopeManager) IsDefaultScope() bool {
	return dsm.currentScope == ScopeConnected
}
