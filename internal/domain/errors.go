package domain

import "errors"

var (
	// ErrInvalidChoiceIndex は無効な選択肢インデックスが指定された場合のエラー
	ErrInvalidChoiceIndex = errors.New("選択肢の番号が無効です")
	
	// ErrParagraphNotFound は指定されたパラグラフが見つからない場合のエラー
	ErrParagraphNotFound = errors.New("指定されたパラグラフが見つかりません")
	
	// ErrDuplicateParagraph は同じ番号のパラグラフが既に存在する場合のエラー
	ErrDuplicateParagraph = errors.New("パラグラフ番号が重複しています。既に同じ番号のパラグラフが存在します")
	
	// ErrUndefinedTargetParagraph は選択肢の遷移先パラグラフが未定義の場合の警告
	ErrUndefinedTargetParagraph = errors.New("警告: 遷移先パラグラフが未定義です。後でパラグラフを追加してください")
)