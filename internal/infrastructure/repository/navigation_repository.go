package repository

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/zensai3805/iwakero-gamebook-mapping/internal/domain"
	"github.com/zensai3805/iwakero-gamebook-mapping/internal/usecase/interfaces"
)

// NavigationRepository はNavigationRepositoryインターフェースの実装
type NavigationRepository struct {
	dataDir string
}

// NewNavigationRepository は新しいNavigationRepositoryを作成する
func NewNavigationRepository(dataDir string) interfaces.NavigationRepository {
	return &NavigationRepository{
		dataDir: dataDir,
	}
}

// SaveNavigationHistory は移動履歴を保存する
func (r *NavigationRepository) SaveNavigationHistory(gamebookTitle string, history []domain.NavigationStep) error {
	// ディレクトリが存在しない場合は作成
	if err := os.MkdirAll(r.dataDir, 0755); err != nil {
		return err
	}

	filename := filepath.Join(r.dataDir, gamebookTitle+"_history.md")
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer func() {
		_ = file.Close()
	}()

	// ヘッダーを書き込み
	_, _ = fmt.Fprintf(file, "# %s 移動履歴\n\n", gamebookTitle)

	// 移動履歴を書き込み
	for _, step := range history {
		_, _ = fmt.Fprintf(file, "- %d -> %d\n", step.From, step.To)
	}

	return nil
}

// LoadNavigationHistory は移動履歴を読み込む
func (r *NavigationRepository) LoadNavigationHistory(gamebookTitle string) ([]domain.NavigationStep, error) {
	filename := filepath.Join(r.dataDir, gamebookTitle+"_history.md")
	content, err := os.ReadFile(filename)
	if err != nil {
		// ファイルが存在しない場合は空の履歴を返す
		if os.IsNotExist(err) {
			return []domain.NavigationStep{}, nil
		}
		return nil, err
	}

	var history []domain.NavigationStep
	lines := strings.Split(string(content), "\n")

	for _, line := range lines {
		trimmedLine := strings.TrimSpace(line)
		
		// "- 1 -> 2" 形式の行をパースする
		if strings.HasPrefix(trimmedLine, "- ") && strings.Contains(trimmedLine, " -> ") {
			parts := strings.Split(strings.TrimPrefix(trimmedLine, "- "), " -> ")
			if len(parts) == 2 {
				from, fromErr := strconv.Atoi(strings.TrimSpace(parts[0]))
				to, toErr := strconv.Atoi(strings.TrimSpace(parts[1]))
				
				if fromErr == nil && toErr == nil {
					step := domain.NewNavigationStep(from, to, []int{})
					history = append(history, *step)
				}
			}
		}
	}

	return history, nil
}