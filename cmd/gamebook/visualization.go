package main

import (
	"github.com/iwapc/iwakero-gamebook-mapping/internal/domain"
	"github.com/pterm/pterm"
)

// VisualizationEvent は可視化システムのイベント型
type VisualizationEvent int

const (
	EventGameLoaded VisualizationEvent = iota
	EventParagraphAdded
	EventChoiceAdded
	EventParagraphVisited
	EventChoiceSelected
	EventCurrentChanged
)

// VisualizationData は可視化に必要なデータを保持する構造体
type VisualizationData struct {
	Gamebook      *domain.Gamebook
	CurrentPos    *domain.Paragraph
	MapData       *MapData
	FlowData      *FlowData
	LastEvent     VisualizationEvent
	LastEventData interface{}
}

// MapData は2Dマップ表示のためのデータ構造
type MapData struct {
	Grid      [][]MapCell      // 格子状マップ
	Width     int              // マップ幅
	Height    int              // マップ高さ
	Positions map[int]Position // パラグラフ番号→座標のマッピング
}

// MapCell はマップの各セルを表す
type MapCell struct {
	ParagraphNumbers []int        // このセルに配置されたパラグラフ番号
	Symbol           string       // 表示文字
	Style            *pterm.Style // 表示スタイル
	Visited          bool         // 訪問済みフラグ
	IsCurrent        bool         // 現在位置フラグ
}

// Position は2D座標を表す
type Position struct {
	X int
	Y int
}

// FlowData はフロー図表示のためのデータ構造
type FlowData struct {
	Nodes []FlowNode // ノード一覧
	Edges []FlowEdge // エッジ一覧
	Root  *FlowNode  // ルートノード
}

// FlowNode はフロー図のノードを表す
type FlowNode struct {
	ParagraphNumber int             // パラグラフ番号
	Description     string          // 説明文
	Children        []*FlowNode     // 子ノード
	Choices         []domain.Choice // このノードの選択肢情報
	Style           *pterm.Style    // 表示スタイル
	Visited         bool            // 訪問済みフラグ
	IsCurrent       bool            // 現在位置フラグ
	VisitCount      int             // 訪問回数
}

// FlowEdge はフロー図のエッジ（選択肢）を表す
type FlowEdge struct {
	From        *FlowNode    // 開始ノード
	To          *FlowNode    // 終了ノード
	Description string       // 選択肢説明
	Selected    bool         // 選択済みフラグ
	Style       *pterm.Style // 表示スタイル
}

// IVisualizer は可視化コンポーネントのインターフェース
type IVisualizer interface {
	// Initialize コンポーネントを初期化する
	Initialize(data *VisualizationData) error

	// Update データ更新を処理する
	Update(event VisualizationEvent, data *VisualizationData) error

	// Render コンポーネントを描画する
	Render() (string, error)

	// GetArea 描画領域のサイズを取得する
	GetArea() (width, height int)

	// SetArea 描画領域のサイズを設定する
	SetArea(width, height int) error
}

// ILayoutManager はレイアウト管理のインターフェース
type ILayoutManager interface {
	// AddComponent コンポーネントを追加する
	AddComponent(name string, component IVisualizer, area LayoutArea) error

	// RemoveComponent コンポーネントを削除する
	RemoveComponent(name string) error

	// Update 全体を更新する
	Update(event VisualizationEvent, data *VisualizationData) error

	// Render 全体を描画する
	Render() error

	// Resize 画面サイズ変更に対応する
	Resize(width, height int) error
}

// LayoutArea はレイアウト領域の定義
type LayoutArea struct {
	X      int    // X座標（比率）
	Y      int    // Y座標（比率）
	Width  int    // 幅（比率）
	Height int    // 高さ（比率）
	Name   string // 領域名
}

// IEventSystem はイベント通知システムのインターフェース
type IEventSystem interface {
	// Subscribe イベント購読を登録する
	Subscribe(eventType VisualizationEvent, handler func(VisualizationEvent, interface{})) error

	// Unsubscribe イベント購読を解除する
	Unsubscribe(eventType VisualizationEvent, handler func(VisualizationEvent, interface{})) error

	// Publish イベントを発行する
	Publish(eventType VisualizationEvent, data interface{}) error
}
