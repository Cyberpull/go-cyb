package dbo

import "database/sql/driver"

type NullUint struct {
	Uint  uint
	Valid bool
}

func (n *NullUint) Scan(value any) (err error) {
	if value == nil {
		n.Uint, n.Valid = 0, false
		return
	}

	if v, ok := value.(*uint); ok {
		n.Uint, n.Valid = *v, *v > 0
		return
	}

	if v, ok := value.(uint); ok {
		n.Uint, n.Valid = v, v > 0
	}

	return
}

func (n NullUint) Value() (driver.Value, error) {
	if n.Valid && n.Uint > 0 {
		return n.Uint, nil
	}

	return nil, nil
}
