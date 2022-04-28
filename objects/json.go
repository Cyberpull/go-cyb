package objects

import (
	"encoding/json"

	"cyberpull.com/go-cyb/validator"
)

func ParseJSON(data []byte, v any) (err error) {
	if err = json.Unmarshal(data, v); err != nil {
		return
	}

	if err = validator.Validate(v); err != nil {
		return
	}

	return
}

func ToJSON(v any) (value []byte, err error) {
	if err = validator.Validate(v); err != nil {
		return
	}

	value, err = json.Marshal(v)

	return
}
