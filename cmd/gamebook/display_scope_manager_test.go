package main

import (
	"testing"
)

func TestDisplayScopeManager_CycleScope(t *testing.T) {
	tests := []struct {
		name         string
		initialScope DisplayScope
		expectedNext DisplayScope
		description  string
	}{
		{
			name:         "Connected から All へ",
			initialScope: ScopeConnected,
			expectedNext: ScopeAll,
			description:  "接続のみから全て表示へ切り替え",
		},
		{
			name:         "All から None へ",
			initialScope: ScopeAll,
			expectedNext: ScopeNone,
			description:  "全て表示から非表示へ切り替え",
		},
		{
			name:         "None から Connected へ",
			initialScope: ScopeNone,
			expectedNext: ScopeConnected,
			description:  "非表示から接続のみへ切り替え（ループ）",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			manager := NewDisplayScopeManager()
			manager.SetScope(tt.initialScope)

			result := manager.CycleScope()

			if result != tt.expectedNext {
				t.Errorf("Expected %s, got %s for %s", tt.expectedNext, result, tt.description)
			}

			// 内部状態も正しく更新されているか確認
			currentScope := manager.GetScope()
			if currentScope != tt.expectedNext {
				t.Errorf("Internal scope not updated. Expected %s, got %s", tt.expectedNext, currentScope)
			}
		})
	}
}

func TestDisplayScopeManager_GetScopeDescription(t *testing.T) {
	tests := []struct {
		name     string
		scope    DisplayScope
		expected string
	}{
		{
			name:     "Connected scope description",
			scope:    ScopeConnected,
			expected: "接続のみ",
		},
		{
			name:     "All scope description",
			scope:    ScopeAll,
			expected: "全て表示",
		},
		{
			name:     "None scope description",
			scope:    ScopeNone,
			expected: "非表示",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			manager := NewDisplayScopeManager()
			manager.SetScope(tt.scope)

			result := manager.GetScopeDescription()

			if result != tt.expected {
				t.Errorf("Expected '%s', got '%s'", tt.expected, result)
			}
		})
	}
}

func TestDisplayScopeManager_GetScope(t *testing.T) {
	manager := NewDisplayScopeManager()

	// デフォルトスコープの確認
	defaultScope := manager.GetScope()
	if defaultScope != ScopeConnected {
		t.Errorf("Expected default scope %s, got %s", ScopeConnected, defaultScope)
	}

	// 設定と取得の確認
	manager.SetScope(ScopeAll)
	currentScope := manager.GetScope()
	if currentScope != ScopeAll {
		t.Errorf("Expected %s, got %s", ScopeAll, currentScope)
	}
}

func TestDisplayScopeManager_IsDefaultScope(t *testing.T) {
	tests := []struct {
		name     string
		scope    DisplayScope
		expected bool
	}{
		{
			name:     "Connected is default",
			scope:    ScopeConnected,
			expected: true,
		},
		{
			name:     "All is not default",
			scope:    ScopeAll,
			expected: false,
		},
		{
			name:     "None is not default",
			scope:    ScopeNone,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			manager := NewDisplayScopeManager()
			manager.SetScope(tt.scope)

			result := manager.IsDefaultScope()

			if result != tt.expected {
				t.Errorf("Expected %v, got %v for scope %s", tt.expected, result, tt.scope)
			}
		})
	}
}
