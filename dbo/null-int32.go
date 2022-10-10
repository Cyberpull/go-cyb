package dbo

import (
	"database/sql"
	"encoding/json"
)

type NullInt32 struct {
	sql.NullInt32
}

func (n *NullInt32) UnmarshalJSON(b []byte) error {
	var value int32

	if err := json.Unmarshal(b, &value); err != nil {
		return err
	}

	return n.Scan(value)
}

func (n NullInt32) MarshalJSON() ([]byte, error) {
	if n.Valid {
		return json.Marshal(n.Int32)
	}

	return json.Marshal(nil)
}
