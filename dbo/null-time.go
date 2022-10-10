package dbo

import (
	"database/sql"
	"encoding/json"
	"time"
)

type NullTime struct {
	sql.NullTime
}

func (n *NullTime) UnmarshalJSON(b []byte) error {
	var value time.Time

	if err := json.Unmarshal(b, &value); err != nil {
		return err
	}

	return n.Scan(value)
}

func (n NullTime) MarshalJSON() ([]byte, error) {
	if n.Valid {
		return json.Marshal(n.Time)
	}

	return json.Marshal(nil)
}
