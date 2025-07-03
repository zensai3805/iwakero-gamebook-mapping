.PHONY: test test-v test-cover build run clean lint

# テスト実行
test:
	go test ./...

# 詳細なテスト実行
test-v:
	go test -v ./...

# カバレッジ付きテスト実行
test-cover:
	go test -cover ./...

# カバレッジレポートをHTMLで表示
test-cover-html:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out

# ビルド
build:
	go build -o bin/gamebook cmd/gamebook/main.go

# 実行
run:
	go run cmd/gamebook/main.go

# クリーンアップ
clean:
	rm -rf bin/
	rm -f coverage.out

# 依存関係の取得
deps:
	go mod download
	go mod tidy

# リント（golangci-lintがインストールされている場合）
lint:
	golangci-lint run

# フォーマット
fmt:
	go fmt ./...