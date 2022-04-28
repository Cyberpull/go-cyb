package socket

import (
	"cyberpull.com/go-cyb/objects"
	"cyberpull.com/go-cyb/uuid"
)

type Request struct {
	UUID    string `json:"uuid" validator:"required"`
	Method  string `json:"method" validator:"required"`
	Channel string `json:"channel" validator:"required"`
	Data    []byte `json:"data" validator:"required"`
}

func (r *Request) SetData(v any) (err error) {
	data, err := objects.ToJSON(r.Data)

	if err != nil {
		return
	}

	r.Data = data

	return
}

/*****************************************/

func newRequest(c *Client) (value *Request, err error) {
	var uniqueId string

	if uniqueId, err = uuid.Generate(); err != nil {
		return
	}

	value = &Request{
		UUID: c.uuid + "||" + uniqueId,
	}

	return
}
