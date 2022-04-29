package socket

import "cyberpull.com/go-cyb/objects"

type BaseData struct {
	Data    []byte `json:"data"`
}

func (b *BaseData) SetData(v any) (err error) {
	data, err := objects.ToJSON(v)

	if err != nil {
		return
	}

	b.Data = data

	return
}

func (b *BaseData) ParseData(v any) error {
	return objects.ParseJSON(b.Data, v)
}
