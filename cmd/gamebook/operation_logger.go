package main

import (
	"fmt"
	"time"
)

// LogUserOperation はユーザー操作を記録する（構造化版）
func LogUserOperation(operationType string, details map[string]interface{}) {
	LogStructuredOperation("operation", operationType, 
		fmt.Sprintf("ユーザー操作: %s", operationType), details, 0)
}

// LogErrorWithContext はエラーを詳細なコンテキスト情報とともにログ記録する（構造化版）
func LogErrorWithContext(err error, operation string, context map[string]interface{}) {
	LogStructuredError("error", operation, 
		fmt.Sprintf("操作エラー: %s", operation), err, context)
}

// LogValidationError は入力値検証エラーを記録する（構造化版）
func LogValidationError(field string, value interface{}, reason string, context map[string]interface{}) {
	validationContext := map[string]interface{}{
		"field":  field,
		"value":  value,
		"reason": reason,
	}
	// 追加コンテキストをマージ
	for k, v := range context {
		validationContext[k] = v
	}
	
	LogStructuredValidation(field, value, reason, validationContext)
}

// LogCommandResult はコマンド実行結果を記録する（構造化版）
func LogCommandResult(command string, success bool, details map[string]interface{}) {
	if success {
		LogStructuredOperation("command", command, 
			fmt.Sprintf("コマンド実行成功: %s", command), details, 0)
	} else {
		var err error
		if errMsg, ok := details["error"]; ok {
			err = fmt.Errorf("%v", errMsg)
		}
		LogStructuredError("command", command, 
			fmt.Sprintf("コマンド実行失敗: %s", command), err, details)
	}
}

// LogConfigurationChange は設定変更を記録する（構造化版）
func LogConfigurationChange(configType string, oldValue, newValue interface{}, reason string) {
	context := map[string]interface{}{
		"config_type": configType,
		"old_value":   oldValue,
		"new_value":   newValue,
		"reason":      reason,
	}
	
	LogStructuredOperation("config", "change", 
		fmt.Sprintf("設定変更: %s", configType), context, 0)
}

// LogStateTransition はアプリケーション状態の変化を記録する（構造化版）
func LogStateTransition(fromState, toState string, trigger string, context map[string]interface{}) {
	transitionContext := map[string]interface{}{
		"from_state": fromState,
		"to_state":   toState,
		"trigger":    trigger,
	}
	// 追加コンテキストをマージ
	for k, v := range context {
		transitionContext[k] = v
	}
	
	LogStructuredOperation("state", "transition", 
		fmt.Sprintf("状態遷移: %s → %s", fromState, toState), transitionContext, 0)
}

// LogUIInteraction はUI操作を記録する（構造化版）
func LogUIInteraction(interactionType string, details map[string]interface{}) {
	// 操作時間がある場合は取得
	var duration time.Duration
	if timeMs, ok := details["selection_time_ms"]; ok {
		if timeFloat, ok := timeMs.(float64); ok {
			duration = time.Duration(timeFloat * float64(time.Millisecond))
		}
	}
	
	LogStructuredUI(interactionType, 
		fmt.Sprintf("UI操作: %s", interactionType), details, duration)
}