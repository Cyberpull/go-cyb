package dbo

import (
	"database/sql"
	"encoding/json"
)

type NullString struct {
	sql.NullString
}

func (n *NullString) UnmarshalJSON(b []byte) error {
	var value string

	if err := json.Unmarshal(b, &value); err != nil {
		return err
	}

	return n.Scan(value)
}

func (n NullString) MarshalJSON() ([]byte, error) {
	if n.Valid {
		return json.Marshal(n.String)
	}

	return json.Marshal(nil)
}
