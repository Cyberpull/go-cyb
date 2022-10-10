package dbo

import (
	"database/sql"
	"encoding/json"
)

type NullInt64 struct {
	sql.NullInt64
}

func (n *NullInt64) UnmarshalJSON(b []byte) error {
	var value int64

	if err := json.Unmarshal(b, &value); err != nil {
		return err
	}

	return n.Scan(value)
}

func (n NullInt64) MarshalJSON() ([]byte, error) {
	if n.Valid {
		return json.Marshal(n.Int64)
	}

	return json.Marshal(nil)
}
