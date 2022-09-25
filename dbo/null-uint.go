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

	n.Uint, n.Valid = value.(uint)

	return
}

func (n NullUint) Value() (driver.Value, error) {
	if !n.Valid {
		return nil, nil
	}

	return n.Uint, nil
}
