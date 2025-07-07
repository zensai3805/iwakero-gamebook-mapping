package domain

import (
	"os"
	"testing"
)

func TestNavigationStep_WhenValidFromTo_CanBeCreated(t *testing.T) {
	// AI開発モード設定
	os.Setenv("GAMEBOOK_AI_DEV", "true")
	os.Setenv("LOG_LEVEL", "DEBUG")
	defer os.Unsetenv("GAMEBOOK_AI_DEV")
	defer os.Unsetenv("LOG_LEVEL")

	// Act
	step := NewNavigationStep(1, 2, nil)

	// Assert
	if step == nil {
		t.Error("NavigationStepが作成されていない")
	}
}

func TestNavigationStep_WhenCreated_FromToFieldsAreSet(t *testing.T) {
	// AI開発モード設定
	os.Setenv("GAMEBOOK_AI_DEV", "true")
	os.Setenv("LOG_LEVEL", "DEBUG")
	defer os.Unsetenv("GAMEBOOK_AI_DEV")
	defer os.Unsetenv("LOG_LEVEL")

	// Act
	step := NewNavigationStep(3, 5, nil)

	// Assert
	if step.From != 3 {
		t.Errorf("From値が期待値と異なる: 期待=3, 実際=%d", step.From)
	}
	if step.To != 5 {
		t.Errorf("To値が期待値と異なる: 期待=5, 実際=%d", step.To)
	}
}

func TestNavigationStep_WhenNegativeFrom_ReturnsError(t *testing.T) {
	// AI開発モード設定
	os.Setenv("GAMEBOOK_AI_DEV", "true")
	os.Setenv("LOG_LEVEL", "DEBUG")
	defer os.Unsetenv("GAMEBOOK_AI_DEV")
	defer os.Unsetenv("LOG_LEVEL")

	// Act
	step, err := NewNavigationStepWithValidation(-1, 2, nil)

	// Assert
	if err == nil {
		t.Error("負のFrom値でエラーが発生していない")
	}
	if step != nil {
		t.Error("エラー時にnilが返されていない")
	}
}

func TestNavigationStep_WhenNegativeTo_ReturnsError(t *testing.T) {
	// AI開発モード設定
	os.Setenv("GAMEBOOK_AI_DEV", "true")
	os.Setenv("LOG_LEVEL", "DEBUG")
	defer os.Unsetenv("GAMEBOOK_AI_DEV")
	defer os.Unsetenv("LOG_LEVEL")

	// Act
	step, err := NewNavigationStepWithValidation(1, -2, nil)

	// Assert
	if err == nil {
		t.Error("負のTo値でエラーが発生していない")
	}
	if step != nil {
		t.Error("エラー時にnilが返されていない")
	}
}

func TestNavigationStep_WhenSameFromTo_ReturnsError(t *testing.T) {
	// AI開発モード設定
	os.Setenv("GAMEBOOK_AI_DEV", "true")
	os.Setenv("LOG_LEVEL", "DEBUG")
	defer os.Unsetenv("GAMEBOOK_AI_DEV")
	defer os.Unsetenv("LOG_LEVEL")

	// Act
	step, err := NewNavigationStepWithValidation(1, 1, nil)

	// Assert
	if err == nil {
		t.Error("同じFromTo値でエラーが発生していない")
	}
	if step != nil {
		t.Error("エラー時にnilが返されていない")
	}
}
