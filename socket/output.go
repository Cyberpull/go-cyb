package socket

import "cyberpull.com/go-cyb/errors"

type Output struct {
	BaseData

	uuid    string
	Method  string `json:"method"`
	Channel string `json:"channel"`
	Code    int    `json:"code"`
}

func (o Output) GetError() (err error) {
	if o.Code >= 200 && o.Code < 300 {
		return
	}

	var message string

	if err = o.ParseData(&message); err != nil {
		return
	}

	err = errors.New(message, o.Code)

	return
}
