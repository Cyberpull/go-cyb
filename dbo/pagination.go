package dbo

import (
	"math"
	"reflect"

	"cyberpull.com/go-cyb/errors"

	"gorm.io/gorm"
)

type ModelType interface {
	*gorm.Model
}

type Pagination[T any] struct {
	Current_page uint `json:"current_page"`
	From         uint `json:"from"`
	Last_page    uint `json:"last_page"`
	Per_page     int  `json:"per_page"`
	To           uint `json:"to"`
	Total        uint `json:"total"`
	Data         []T  `json:"data"`
}

func Paginate[T any](tx *gorm.DB, page uint, limit ...uint) (value *Pagination[T], err error) {
	if len(limit) == 0 {
		limit = append(limit, 20)
	}

	var model T

	vType := reflect.TypeOf(model)

	if vType.Kind() == reflect.Pointer {
		vType = vType.Elem()
		model = reflect.New(vType).Interface().(T)
		tx = tx.Model(model)
	} else {
		tx = tx.Model(&model)
	}

	tx = tx.Offset(0).Limit(-1).Session(&gorm.Session{})

	if vType.Kind() != reflect.Struct {
		err = errors.New("Model should be a struct")
		return
	}

	tmpValue := &Pagination[T]{}
	tmpValue.Data = make([]T, 0)
	tmpValue.Current_page = uint(math.Max(float64(page), 1))
	tmpValue.Per_page = int(math.Max(float64(limit[0]), 1))

	if tmpValue.Current_page == 1 {
		tmpValue.From = tmpValue.Current_page
	} else {
		tmpValue.From = tmpValue.Current_page * uint(tmpValue.Per_page)
	}

	offset := int(tmpValue.From) - 1

	tx = tx.Offset(offset).Limit(tmpValue.Per_page)

	if err = tx.Find(&tmpValue.Data).Error; err != nil {
		return
	}

	tmpValue.To = tmpValue.From + uint(len(tmpValue.Data))

	var total int64

	tx = tx.Offset(0).Limit(-1)

	if err = tx.Count(&total).Error; err != nil {
		return
	}

	lastPage := float64(total) / float64(len(tmpValue.Data))
	tmpValue.Last_page = uint(math.Ceil(lastPage))
	tmpValue.Total = uint(total)

	value = tmpValue

	return
}
