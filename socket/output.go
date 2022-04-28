package socket

import "cyberpull.com/go-cyb/objects"

type Output struct {
	uuid string
	Code int    `json:"code"`
	Data []byte `json:"data"`
}

func (o *Output) SetData(v any) (err error) {
	data, err := objects.ToJSON(o.Data)

	if err != nil {
		return
	}

	o.Data = data

	return
}

func (o *Output) Parse(v any) error {
	return objects.ParseJSON(o.Data, v)
}
