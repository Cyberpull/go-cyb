package dbo

import (
	"database/sql/driver"
	"encoding/json"
	"reflect"

	"cyberpull.com/go-cyb/errors"
)

type Null[T comparable] struct {
	Data  T
	Valid bool
}

func (n *Null[T]) Scan(value any) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errors.From(r)
		}
	}()

	xType := reflect.TypeOf(n.Data)

	xValue := reflect.ValueOf(value)
	xValueType := xValue.Type()

	if !xValueType.AssignableTo(xType) && !xValueType.ConvertibleTo(xType) {
		return errors.New("Invalid value")
	}

	n.Data = xValue.Convert(xType).Interface().(T)
	n.Valid = !xValue.IsZero()

	return
}

func (n Null[T]) Value() (driver.Value, error) {
	xValue := reflect.ValueOf(n.Data)

	if n.Valid && !xValue.IsZero() {
		return n.Data, nil
	}

	return nil, nil
}

func (n *Null[T]) UnmarshalJSON(b []byte) error {
	var value uint

	if err := json.Unmarshal(b, &value); err != nil {
		return err
	}

	return n.Scan(value)
}

func (n Null[T]) MarshalJSON() ([]byte, error) {
	if n.Valid {
		return json.Marshal(n.Data)
	}

	return json.Marshal(nil)
}
