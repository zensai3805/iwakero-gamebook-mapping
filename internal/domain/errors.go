package domain

import "errors"

var (
	// ErrInvalidChoiceIndex は無効な選択肢インデックスが指定された場合のエラー
	ErrInvalidChoiceIndex = errors.New("invalid choice index")

	// ErrParagraphNotFound は指定されたパラグラフが見つからない場合のエラー
	ErrParagraphNotFound = errors.New("paragraph not found")

	// ErrDuplicateParagraph は同じ番号のパラグラフが既に存在する場合のエラー
	ErrDuplicateParagraph = errors.New("paragraph with this number already exists")

	// ErrInvalidParagraphNumber は無効なパラグラフ番号が指定された場合のエラー
	ErrInvalidParagraphNumber = errors.New("invalid paragraph number")

	// ErrSameFromToNavigation は移動元と移動先が同じ場合のエラー
	ErrSameFromToNavigation = errors.New("from and to must be different")
)
