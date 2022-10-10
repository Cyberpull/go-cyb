package dbo

import (
	"database/sql"
	"encoding/json"
)

type NullBool struct {
	sql.NullBool
}

func (n *NullBool) UnmarshalJSON(b []byte) error {
	var value bool

	if err := json.Unmarshal(b, &value); err != nil {
		return err
	}

	return n.Scan(value)
}

func (n NullBool) MarshalJSON() ([]byte, error) {
	if n.Valid {
		return json.Marshal(n.Bool)
	}

	return json.Marshal(nil)
}
