package main

import (
	"testing"
	"time"

	"github.com/zensai3805/iwakero-gamebook-mapping/internal/domain"
)

func TestEventSystem(t *testing.T) {
	es := NewEventSystem()

	// チャネルを使用して同期
	done := make(chan bool, 1)
	var receivedEvent VisualizationEvent
	var receivedData interface{}

	handler := func(event VisualizationEvent, data interface{}) {
		receivedEvent = event
		receivedData = data
		done <- true
	}

	// 購読テスト
	err := es.Subscribe(EventGameLoaded, handler)
	if err != nil {
		t.Errorf("Subscribe failed: %v", err)
	}

	if es.GetSubscriberCount(EventGameLoaded) != 1 {
		t.Errorf("Expected 1 subscriber, got %d", es.GetSubscriberCount(EventGameLoaded))
	}

	// イベント発行テスト
	testData := "test data"
	err = es.Publish(EventGameLoaded, testData)
	if err != nil {
		t.Errorf("Publish failed: %v", err)
	}

	// goroutineの完了を同期的に待機
	select {
	case <-done:
		// 正常完了
	case <-time.After(100 * time.Millisecond):
		t.Error("Handler was not called within timeout")
		return
	}

	if receivedEvent != EventGameLoaded {
		t.Errorf("Expected event %d, got %d", EventGameLoaded, receivedEvent)
	}

	if receivedData != testData {
		t.Errorf("Expected data %v, got %v", testData, receivedData)
	}
}

func TestLayoutManager(t *testing.T) {
	lm := NewLayoutManager()

	// モックコンポーネント
	mockComponent := &MockVisualizer{}
	area := LayoutArea{
		X:      0,
		Y:      0,
		Width:  50,
		Height: 50,
		Name:   "test",
	}

	// コンポーネント追加テスト
	err := lm.AddComponent("test", mockComponent, area)
	if err != nil {
		t.Errorf("AddComponent failed: %v", err)
	}

	// 重複追加テスト
	err = lm.AddComponent("test", mockComponent, area)
	if err == nil {
		t.Error("Expected error for duplicate component")
	}

	// サイズ変更テスト
	err = lm.Resize(200, 60)
	if err != nil {
		t.Errorf("Resize failed: %v", err)
	}

	// コンポーネント削除テスト
	err = lm.RemoveComponent("test")
	if err != nil {
		t.Errorf("RemoveComponent failed: %v", err)
	}

	// 存在しないコンポーネント削除テスト
	err = lm.RemoveComponent("nonexistent")
	if err == nil {
		t.Error("Expected error for nonexistent component")
	}
}

func TestDataConverter(t *testing.T) {
	dc := NewDataConverter()

	// テスト用ゲームブック作成
	gamebook := domain.NewGamebook("Test Game")
	p1 := domain.NewParagraph(1, "Start")
	p2 := domain.NewParagraph(2, "Middle")

	err := gamebook.AddParagraph(p1)
	if err != nil {
		t.Fatalf("Failed to add paragraph: %v", err)
	}

	err = gamebook.AddParagraph(p2)
	if err != nil {
		t.Fatalf("Failed to add paragraph: %v", err)
	}

	err = gamebook.AddChoiceToParagraph(1, "Go to middle", 2)
	if err != nil {
		t.Fatalf("Failed to add choice: %v", err)
	}

	// 変換テスト
	vizData, err := dc.ConvertToVisualizationData(gamebook)
	if err != nil {
		t.Errorf("ConvertToVisualizationData failed: %v", err)
	}

	if vizData == nil {
		t.Fatal("VisualizationData is nil")
	}

	if vizData.Gamebook != gamebook {
		t.Error("Gamebook reference mismatch")
	}

	if vizData.MapData == nil {
		t.Error("MapData is nil")
	}

	if vizData.FlowData == nil {
		t.Error("FlowData is nil")
	}

	// マップデータ検証
	if len(vizData.MapData.Positions) != 2 {
		t.Errorf("Expected 2 positions, got %d", len(vizData.MapData.Positions))
	}

	// フローデータ検証
	if len(vizData.FlowData.Nodes) != 2 {
		t.Errorf("Expected 2 nodes, got %d", len(vizData.FlowData.Nodes))
	}

	if len(vizData.FlowData.Edges) != 1 {
		t.Errorf("Expected 1 edge, got %d", len(vizData.FlowData.Edges))
	}
}

func TestDataConverterWithNilGamebook(t *testing.T) {
	dc := NewDataConverter()

	_, err := dc.ConvertToVisualizationData(nil)
	if err == nil {
		t.Error("Expected error for nil gamebook")
	}
}

func TestUpdateVisualizationData(t *testing.T) {
	dc := NewDataConverter()
	gamebook := domain.NewGamebook("Test Game")
	p1 := domain.NewParagraph(1, "Start")

	err := gamebook.AddParagraph(p1)
	if err != nil {
		t.Fatalf("Failed to add paragraph: %v", err)
	}

	vizData, err := dc.ConvertToVisualizationData(gamebook)
	if err != nil {
		t.Fatalf("ConvertToVisualizationData failed: %v", err)
	}

	// 更新テスト
	err = dc.UpdateVisualizationData(vizData, EventParagraphAdded, nil)
	if err != nil {
		t.Errorf("UpdateVisualizationData failed: %v", err)
	}

	if vizData.LastEvent != EventParagraphAdded {
		t.Errorf("Expected last event %d, got %d", EventParagraphAdded, vizData.LastEvent)
	}
}

// MockVisualizer はテスト用のモックビジュアライザー
type MockVisualizer struct {
	width  int
	height int
}

func (mv *MockVisualizer) Initialize(data *VisualizationData) error {
	return nil
}

func (mv *MockVisualizer) Update(event VisualizationEvent, data *VisualizationData) error {
	return nil
}

func (mv *MockVisualizer) Render() (string, error) {
	return "mock content", nil
}

func (mv *MockVisualizer) GetArea() (width, height int) {
	return mv.width, mv.height
}

func (mv *MockVisualizer) SetArea(width, height int) error {
	mv.width = width
	mv.height = height
	return nil
}
