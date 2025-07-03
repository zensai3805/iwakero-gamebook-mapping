package domain

import "errors"

var (
	// ErrInvalidChoiceIndex は無効な選択肢インデックスが指定された場合のエラー
	ErrInvalidChoiceIndex = errors.New("invalid choice index")
	
	// ErrParagraphNotFound は指定されたパラグラフが見つからない場合のエラー
	ErrParagraphNotFound = errors.New("paragraph not found")
	
	// ErrDuplicateParagraph は同じ番号のパラグラフが既に存在する場合のエラー
	ErrDuplicateParagraph = errors.New("パラグラフ番号が重複しています。既に同じ番号のパラグラフが存在します")
)