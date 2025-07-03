package repository

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/iwapc/iwakero-gamebook-mapping/internal/domain"
)

// MarkdownRepository はMarkdownファイルでゲームブックを永続化する
type MarkdownRepository struct {
	baseDir string
}

// NewMarkdownRepository は新しいMarkdownRepositoryを作成する
func NewMarkdownRepository(baseDir string) *MarkdownRepository {
	return &MarkdownRepository{
		baseDir: baseDir,
	}
}

// Save はゲームブックをMarkdownファイルとして保存する
func (r *MarkdownRepository) Save(gamebook *domain.Gamebook) error {
	// データ検証
	if err := r.validateGamebook(gamebook); err != nil {
		return fmt.Errorf("データ検証エラー: %w", err)
	}

	// ディレクトリが存在しない場合は作成
	if err := os.MkdirAll(r.baseDir, 0755); err != nil {
		return fmt.Errorf("ディレクトリ作成エラー: %w", err)
	}

	filename := filepath.Join(r.baseDir, gamebook.Title+".md")
	tempFilename := filename + ".tmp"

	// 一時ファイルに書き込み
	file, err := os.Create(tempFilename)
	if err != nil {
		return fmt.Errorf("一時ファイル作成エラー: %w", err)
	}
	defer func() {
		file.Close()
		// エラー時は一時ファイルを削除
		if err != nil {
			os.Remove(tempFilename)
		}
	}()

	// ヘッダーを書き込み
	fmt.Fprintf(file, "# %s\n\n", gamebook.Title)

	// 各パラグラフを書き込み
	r.writeParagraphs(file, gamebook)

	// Mermaid形式のフロー図を追加
	r.writeMermaidDiagram(file, gamebook)

	// ファイル同期
	if err := file.Sync(); err != nil {
		return fmt.Errorf("ファイル同期エラー: %w", err)
	}

	// ファイルクローズ
	if err := file.Close(); err != nil {
		return fmt.Errorf("ファイルクローズエラー: %w", err)
	}

	// 原子的リネーム
	if err := os.Rename(tempFilename, filename); err != nil {
		return fmt.Errorf("ファイルリネームエラー: %w", err)
	}

	return nil
}

// Load は指定されたタイトルのゲームブックを読み込む
func (r *MarkdownRepository) Load(title string) (*domain.Gamebook, error) {
	if title == "" {
		return nil, fmt.Errorf("タイトルが空です")
	}

	filename := filepath.Join(r.baseDir, title+".md")
	content, err := os.ReadFile(filename)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("ゲームブック '%s' が見つかりません", title)
		}
		return nil, fmt.Errorf("ファイル読み込みエラー: %w", err)
	}

	if len(content) == 0 {
		return nil, fmt.Errorf("ファイルが空です: %s", filename)
	}

	gamebook := domain.NewGamebook(title)
	r.parseMarkdownContent(gamebook, string(content))

	// 読み込み後のデータ検証
	if err := r.validateGamebook(gamebook); err != nil {
		return nil, fmt.Errorf("読み込み後データ検証エラー: %w", err)
	}

	return gamebook, nil
}

// List は保存されているゲームブックのタイトル一覧を返す
func (r *MarkdownRepository) List() ([]string, error) {
	files, err := os.ReadDir(r.baseDir)
	if err != nil {
		if os.IsNotExist(err) {
			return []string{}, nil
		}
		return nil, err
	}

	var titles []string
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".md") {
			title := strings.TrimSuffix(file.Name(), ".md")
			titles = append(titles, title)
		}
	}

	return titles, nil
}

// Delete は指定されたタイトルのゲームブックを削除する
func (r *MarkdownRepository) Delete(title string) error {
	filename := filepath.Join(r.baseDir, title+".md")
	return os.Remove(filename)
}

// truncate は文字列を指定された長さに切り詰める
func truncate(s string, maxLen int) string {
	runes := []rune(s)
	if len(runes) <= maxLen {
		return s
	}
	return string(runes[:maxLen-3]) + "..."
}

// validateGamebook はゲームブックのデータ整合性を検証する
func (r *MarkdownRepository) validateGamebook(gamebook *domain.Gamebook) error {
	if gamebook == nil {
		return fmt.Errorf("ゲームブックがnilです")
	}

	if gamebook.Title == "" {
		return fmt.Errorf("ゲームブックのタイトルが空です")
	}

	// タイトルに不正な文字が含まれていないかチェック
	if strings.ContainsAny(gamebook.Title, "/\\:*?\"<>|") {
		return fmt.Errorf("タイトルに不正な文字が含まれています: %s", gamebook.Title)
	}

	if gamebook.Paragraphs == nil {
		return fmt.Errorf("パラグラフマップがnilです")
	}

	if len(gamebook.Paragraphs) == 0 {
		return fmt.Errorf("パラグラフが1つも存在しません")
	}

	// 各パラグラフの整合性をチェック
	for number, paragraph := range gamebook.Paragraphs {
		if paragraph == nil {
			return fmt.Errorf("パラグラフ %d がnilです", number)
		}

		if paragraph.Number != number {
			return fmt.Errorf("パラグラフ番号の不整合: マップキー=%d, パラグラフ番号=%d", number, paragraph.Number)
		}

		if paragraph.Description == "" {
			return fmt.Errorf("パラグラフ %d の説明が空です", number)
		}

		// 選択肢の整合性をチェック
		for i, choice := range paragraph.Choices {
			if choice.Description == "" {
				return fmt.Errorf("パラグラフ %d の選択肢 %d の説明が空です", number, i)
			}

			if choice.TargetNumber <= 0 {
				return fmt.Errorf("パラグラフ %d の選択肢 %d の遷移先番号が無効です: %d", number, i, choice.TargetNumber)
			}
		}
	}

	return nil
}

// writeParagraphs はパラグラフ情報をファイルに書き込む
func (r *MarkdownRepository) writeParagraphs(file *os.File, gamebook *domain.Gamebook) {
	for number := 1; number <= 1000; number++ { // 適当な上限
		p, exists := gamebook.Paragraphs[number]
		if !exists {
			continue
		}

		fmt.Fprintf(file, "## パラグラフ %d\n", p.Number)
		fmt.Fprintf(file, "- 概要：%s\n", p.Description)

		if len(p.Choices) > 0 {
			fmt.Fprintf(file, "- 選択肢：\n")
			for _, choice := range p.Choices {
				selected := " "
				if choice.Selected {
					selected = "x"
				}
				fmt.Fprintf(file, "  - [%s] %s → %d\n", selected, choice.Description, choice.TargetNumber)
			}
		}

		if p.Visited {
			fmt.Fprintf(file, "- 訪問済み：はい\n")
		}

		fmt.Fprintln(file)
	}
}

// writeMermaidDiagram はMermaid形式のフロー図をファイルに書き込む
func (r *MarkdownRepository) writeMermaidDiagram(file *os.File, gamebook *domain.Gamebook) {
	fmt.Fprintln(file, "## フロー図")
	fmt.Fprintln(file, "```mermaid")
	fmt.Fprintln(file, "graph TD")

	// ノードの定義
	for number := 1; number <= 1000; number++ {
		p, exists := gamebook.Paragraphs[number]
		if !exists {
			continue
		}

		nodeStyle := ""
		if p.Visited {
			nodeStyle = ":::"
		}
		fmt.Fprintf(file, "    %d[%d: %s]%s\n", p.Number, p.Number, truncate(p.Description, 20), nodeStyle)
	}

	fmt.Fprintln(file)

	// エッジの定義
	for number := 1; number <= 1000; number++ {
		p, exists := gamebook.Paragraphs[number]
		if !exists {
			continue
		}

		for _, choice := range p.Choices {
			arrow := "-.->|%s|"
			if choice.Selected {
				arrow = "-->|%s|"
			}
			fmt.Fprintf(file, "    %d %s %d\n", p.Number, fmt.Sprintf(arrow, truncate(choice.Description, 15)), choice.TargetNumber)
		}
	}

	fmt.Fprintln(file, "```")
}

// parseMarkdownContent はMarkdown内容を解析してゲームブックを構築する
func (r *MarkdownRepository) parseMarkdownContent(gamebook *domain.Gamebook, content string) {
	lines := strings.Split(content, "\n")

	var currentParagraph *domain.Paragraph
	inChoices := false

	for _, line := range lines {
		trimmedLine := strings.TrimSpace(line)

		// フロー図セクションに達したら終了
		if trimmedLine == "## フロー図" {
			break
		}

		currentParagraph, inChoices = r.processLine(gamebook, currentParagraph, line, trimmedLine, inChoices)
	}

	// 最後のパラグラフを追加
	if currentParagraph != nil {
		_ = gamebook.AddParagraph(currentParagraph)
	}
}

// processLine は1行を処理してパラグラフ情報を更新する
func (r *MarkdownRepository) processLine(gamebook *domain.Gamebook, currentParagraph *domain.Paragraph, line, trimmedLine string, inChoices bool) (*domain.Paragraph, bool) {
	// パラグラフの開始
	if strings.HasPrefix(trimmedLine, "## パラグラフ ") {
		if currentParagraph != nil {
			_ = gamebook.AddParagraph(currentParagraph)
		}

		var number int
		if _, err := fmt.Sscanf(trimmedLine, "## パラグラフ %d", &number); err != nil {
			return currentParagraph, inChoices // パラグラフ番号が読み取れない場合はスキップ
		}
		currentParagraph = domain.NewParagraph(number, "")
		return currentParagraph, false
	}

	if currentParagraph == nil {
		return currentParagraph, inChoices
	}

	// 概要
	if strings.HasPrefix(trimmedLine, "- 概要：") {
		currentParagraph.Description = strings.TrimPrefix(trimmedLine, "- 概要：")
		return currentParagraph, inChoices
	}

	// 選択肢セクション
	if trimmedLine == "- 選択肢：" {
		return currentParagraph, true
	}

	// 選択肢の内容（元の行でチェック）
	if inChoices && strings.HasPrefix(line, "  - [") {
		if err := r.parseChoice(currentParagraph, line); err != nil {
			return currentParagraph, inChoices // エラーの場合はスキップ
		}
		return currentParagraph, inChoices
	}

	// 訪問済み
	if strings.HasPrefix(trimmedLine, "- 訪問済み：はい") {
		currentParagraph.Visited = true
		return currentParagraph, false // 選択肢セクション終了
	}

	// 他の "- " で始まる行も選択肢セクション終了の合図
	if inChoices && strings.HasPrefix(trimmedLine, "- ") && !strings.HasPrefix(line, "  - [") {
		inChoices = false
	}

	return currentParagraph, inChoices
}

// parseChoice は選択肢行を解析してパラグラフに追加する
func (r *MarkdownRepository) parseChoice(paragraph *domain.Paragraph, line string) error {
	var description string
	var targetNumber int
	var selected bool

	// 選択状態を確認
	if strings.Contains(line, "[x]") {
		selected = true
		// "  - [x] " を除去して矢印で分割
		parts := strings.Split(strings.TrimPrefix(line, "  - [x] "), " → ")
		if len(parts) == 2 {
			description = parts[0]
			if _, err := fmt.Sscanf(parts[1], "%d", &targetNumber); err != nil {
				return err
			}
		}
	} else {
		selected = false
		// "  - [ ] " を除去して矢印で分割
		parts := strings.Split(strings.TrimPrefix(line, "  - [ ] "), " → ")
		if len(parts) == 2 {
			description = parts[0]
			if _, err := fmt.Sscanf(parts[1], "%d", &targetNumber); err != nil {
				return err
			}
		}
	}

	if description != "" && targetNumber > 0 {
		paragraph.AddChoice(description, targetNumber)
		if selected {
			_ = paragraph.SelectChoice(len(paragraph.Choices) - 1)
		}
	}

	return nil
}
