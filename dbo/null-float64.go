package dbo

import (
	"database/sql"
	"encoding/json"
)

type NullFloat64 struct {
	sql.NullFloat64
}

func (n *NullFloat64) UnmarshalJSON(b []byte) error {
	var value float64

	if err := json.Unmarshal(b, &value); err != nil {
		return err
	}

	return n.Scan(value)
}

func (n NullFloat64) MarshalJSON() ([]byte, error) {
	if n.Valid {
		return json.Marshal(n.Float64)
	}

	return json.Marshal(nil)
}
