package main

import (
	"fmt"
	"sync"
)

// EventHandler はイベントハンドラーの型
type EventHandler func(VisualizationEvent, interface{})

// EventSystem はイベント通知システムの実装
type EventSystem struct {
	subscribers map[VisualizationEvent][]EventHandler
	mutex       sync.RWMutex
}

// NewEventSystem は新しいEventSystemを作成
func NewEventSystem() *EventSystem {
	return &EventSystem{
		subscribers: make(map[VisualizationEvent][]EventHandler),
	}
}

// Subscribe イベント購読を登録する
func (es *EventSystem) Subscribe(eventType VisualizationEvent, handler EventHandler) error {
	es.mutex.Lock()
	defer es.mutex.Unlock()

	if handler == nil {
		return fmt.Errorf("handler cannot be nil")
	}

	es.subscribers[eventType] = append(es.subscribers[eventType], handler)
	return nil
}

// Unsubscribe イベント購読を解除する
func (es *EventSystem) Unsubscribe(eventType VisualizationEvent, handler EventHandler) error {
	es.mutex.Lock()
	defer es.mutex.Unlock()

	handlers, exists := es.subscribers[eventType]
	if !exists {
		return fmt.Errorf("no subscribers for event type %d", eventType)
	}

	// ハンドラーを削除（ポインタ比較）
	for i, h := range handlers {
		if &h == &handler {
			es.subscribers[eventType] = append(handlers[:i], handlers[i+1:]...)
			return nil
		}
	}

	return fmt.Errorf("handler not found for event type %d", eventType)
}

// Publish イベントを発行する
func (es *EventSystem) Publish(eventType VisualizationEvent, data interface{}) error {
	es.mutex.RLock()
	defer es.mutex.RUnlock()

	handlers, exists := es.subscribers[eventType]
	if !exists {
		// 購読者がいない場合はエラーではない
		return nil
	}

	// 全ハンドラーを呼び出し
	for _, handler := range handlers {
		go func(h EventHandler) {
			defer func() {
				if r := recover(); r != nil {
					// パニックをキャッチしてログ出力（実際の実装では適切なロガーを使用）
					fmt.Printf("Event handler panic: %v\n", r)
				}
			}()
			h(eventType, data)
		}(handler)
	}

	return nil
}

// GetSubscriberCount 購読者数を取得（テスト用）
func (es *EventSystem) GetSubscriberCount(eventType VisualizationEvent) int {
	es.mutex.RLock()
	defer es.mutex.RUnlock()

	handlers, exists := es.subscribers[eventType]
	if !exists {
		return 0
	}
	return len(handlers)
}
