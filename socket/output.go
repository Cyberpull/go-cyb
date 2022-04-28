package socket

import "cyberpull.com/go-cyb/objects"

type Output struct {
	uuid    string
	Method  string `json:"method"`
	Channel string `json:"channel"`
	Code    int    `json:"code"`
	Data    []byte `json:"data"`
}

func (o *Output) SetData(v any) (err error) {
	data, err := objects.ToJSON(v)

	if err != nil {
		return
	}

	o.Data = data

	return
}

func (o *Output) ParseData(v any) error {
	return objects.ParseJSON(o.Data, v)
}
