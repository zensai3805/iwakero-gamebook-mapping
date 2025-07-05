package main

import (
	"fmt"
	"sync"

	"github.com/pterm/pterm"
)

// LayoutManager はレイアウト管理の実装
type LayoutManager struct {
	components map[string]ComponentInfo
	areas      map[string]LayoutArea
	width      int
	height     int
	mutex      sync.RWMutex
}

// ComponentInfo はコンポーネント情報
type ComponentInfo struct {
	Component IVisualizer
	Area      LayoutArea
}

// NewLayoutManager は新しいLayoutManagerを作成
func NewLayoutManager() *LayoutManager {
	return &LayoutManager{
		components: make(map[string]ComponentInfo),
		areas:      make(map[string]LayoutArea),
		width:      100, // デフォルト幅
		height:     30,  // デフォルト高さ
	}
}

// AddComponent コンポーネントを追加する
func (lm *LayoutManager) AddComponent(name string, component IVisualizer, area LayoutArea) error {
	lm.mutex.Lock()
	defer lm.mutex.Unlock()

	if _, exists := lm.components[name]; exists {
		return fmt.Errorf("component '%s' already exists", name)
	}

	// エリアのサイズを計算して設定
	actualWidth := (area.Width * lm.width) / 100
	actualHeight := (area.Height * lm.height) / 100

	if err := component.SetArea(actualWidth, actualHeight); err != nil {
		return fmt.Errorf("failed to set component area: %w", err)
	}

	lm.components[name] = ComponentInfo{
		Component: component,
		Area:      area,
	}
	lm.areas[name] = area

	return nil
}

// RemoveComponent コンポーネントを削除する
func (lm *LayoutManager) RemoveComponent(name string) error {
	lm.mutex.Lock()
	defer lm.mutex.Unlock()

	if _, exists := lm.components[name]; !exists {
		return fmt.Errorf("component '%s' not found", name)
	}

	delete(lm.components, name)
	delete(lm.areas, name)
	return nil
}

// Update 全体を更新する
func (lm *LayoutManager) Update(event VisualizationEvent, data *VisualizationData) error {
	lm.mutex.RLock()
	defer lm.mutex.RUnlock()

	for name, info := range lm.components {
		if err := info.Component.Update(event, data); err != nil {
			return fmt.Errorf("failed to update component '%s': %w", name, err)
		}
	}
	return nil
}

// Render 全体を描画する
func (lm *LayoutManager) Render() error {
	lm.mutex.RLock()
	defer lm.mutex.RUnlock()

	// 画面をクリア
	pterm.Printf("\033[2J\033[H")

	// 各コンポーネントを描画
	for name, info := range lm.components {
		content, err := info.Component.Render()
		if err != nil {
			return fmt.Errorf("failed to render component '%s': %w", name, err)
		}

		// 位置を設定して描画
		x := (info.Area.X * lm.width) / 100
		y := (info.Area.Y * lm.height) / 100

		pterm.Printf("\033[%d;%dH%s", y+1, x+1, content)
	}

	return nil
}

// Resize 画面サイズ変更に対応する
func (lm *LayoutManager) Resize(width, height int) error {
	lm.mutex.Lock()
	defer lm.mutex.Unlock()

	lm.width = width
	lm.height = height

	// 全コンポーネントのサイズを再計算
	for name, info := range lm.components {
		actualWidth := (info.Area.Width * width) / 100
		actualHeight := (info.Area.Height * height) / 100

		if err := info.Component.SetArea(actualWidth, actualHeight); err != nil {
			return fmt.Errorf("failed to resize component '%s': %w", name, err)
		}
	}

	return nil
}

// GetDefaultAreas デフォルトの3分割レイアウト設定を取得
func GetDefaultAreas() map[string]LayoutArea {
	return map[string]LayoutArea{
		"map": {
			X:      0,
			Y:      0,
			Width:  50,
			Height: 70,
			Name:   "map",
		},
		"flow": {
			X:      50,
			Y:      0,
			Width:  50,
			Height: 70,
			Name:   "flow",
		},
		"menu": {
			X:      0,
			Y:      70,
			Width:  100,
			Height: 30,
			Name:   "menu",
		},
	}
}
