package dbo

import (
	"database/sql"
	"encoding/json"
)

type NullInt16 struct {
	sql.NullInt16
}

func (n *NullInt16) UnmarshalJSON(b []byte) error {
	var value int16

	if err := json.Unmarshal(b, &value); err != nil {
		return err
	}

	return n.Scan(value)
}

func (n NullInt16) MarshalJSON() ([]byte, error) {
	if n.Valid {
		return json.Marshal(n.Int16)
	}

	return json.Marshal(nil)
}
