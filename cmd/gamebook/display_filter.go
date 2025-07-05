package main

import (
	"github.com/iwapc/iwakero-gamebook-mapping/internal/domain"
)

// DisplayScope は未定義パラグラフの表示範囲を定義する
type DisplayScope string

const (
	// ScopeConnected は接続された未定義パラグラフのみ表示
	ScopeConnected DisplayScope = "connected"
	// ScopeAll は全ての未定義パラグラフを表示
	ScopeAll DisplayScope = "all"
	// ScopeNone は未定義パラグラフを非表示
	ScopeNone DisplayScope = "none"
)

// DisplayFilter は表示フィルタリングを行う
type DisplayFilter struct {
	scope    DisplayScope
	analysis *domain.UndefinedAnalysis
}

// NewDisplayFilter は新しい表示フィルターを作成する
func NewDisplayFilter(scope DisplayScope, analysis *domain.UndefinedAnalysis) *DisplayFilter {
	return &DisplayFilter{
		scope:    scope,
		analysis: analysis,
	}
}

// ShouldDisplayUndefined は指定されたパラグラフを表示すべきかを判定する
func (df *DisplayFilter) ShouldDisplayUndefined(gamebook *domain.Gamebook, paragraphNumber int) bool {
	// 定義済みパラグラフは常に表示
	if _, exists := gamebook.Paragraphs[paragraphNumber]; exists {
		return true
	}

	// 未定義パラグラフの場合、スコープに応じて判定
	switch df.scope {
	case ScopeNone:
		return false
	case ScopeAll:
		return true
	case ScopeConnected:
		// 接続された未定義パラグラフのみ表示
		return df.isInConnectedList(paragraphNumber)
	default:
		// デフォルトは接続された未定義パラグラフのみ表示
		return df.isInConnectedList(paragraphNumber)
	}
}

// FilterChoices は選択肢をフィルタリングする
func (df *DisplayFilter) FilterChoices(gamebook *domain.Gamebook, choices []domain.Choice) []domain.Choice {
	if df.scope == ScopeAll {
		// 全て表示の場合はそのまま返す
		return choices
	}

	filtered := make([]domain.Choice, 0)
	for _, choice := range choices {
		if df.ShouldDisplayUndefined(gamebook, choice.TargetNumber) {
			filtered = append(filtered, choice)
		}
	}

	return filtered
}

// isInConnectedList は指定されたパラグラフが接続された未定義パラグラフリストに含まれているかを判定する
func (df *DisplayFilter) isInConnectedList(paragraphNumber int) bool {
	if df.analysis == nil {
		return false
	}

	for _, connectedNum := range df.analysis.Connected {
		if connectedNum == paragraphNumber {
			return true
		}
	}

	return false
}
