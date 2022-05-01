package dbo

import "gorm.io/gorm"

type ModelType interface {
	*gorm.Model
}

type Pagination[T ModelType] struct {
	Current_page uint `json:"current_page"`
	From         uint `json:"from"`
	Last_page    uint `json:"last_page"`
	Per_page     int  `json:"per_page"`
	To           uint `json:"to"`
	Total        uint `json:"total"`
	Data         []T  `json:"data"`
}

func Paginate[T any](tx *TxDB) (value T, err error) {
	return
}
