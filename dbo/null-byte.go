package dbo

import (
	"database/sql"
	"encoding/json"
)

type NullByte struct {
	sql.NullByte
}

func (n *NullByte) UnmarshalJSON(b []byte) error {
	var value byte

	if err := json.Unmarshal(b, &value); err != nil {
		return err
	}

	return n.Scan(value)
}

func (n NullByte) MarshalJSON() ([]byte, error) {
	if n.Valid {
		return json.Marshal(n.Byte)
	}

	return json.Marshal(nil)
}
