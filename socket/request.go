package socket

import (
	"reflect"
	"time"

	"cyberpull.com/go-cyb/errors"
	"cyberpull.com/go-cyb/uuid"
)

type Request struct {
	BaseData

	UUID    string `json:"uuid" validator:"required"`
	Method  string `json:"method" validator:"required"`
	Channel string `json:"channel" validator:"required"`
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

/*****************************************/

func MakeRequest[T any](c *Client, method, channel string, data any, timeout ...time.Duration) (value T, err error) {
	if c == nil {
		err = errors.New("Invalid Client instance")
		return
	}

	var out *Output

	if out, err = c.request(method, channel, data, timeout...); err != nil {
		return
	}

	var tmpValue T

	vType := reflect.TypeOf(value)

	if vType.Kind() == reflect.Pointer {
		tmpValue = reflect.New(vType.Elem()).Interface().(T)
		err = out.ParseData(tmpValue)
	} else {
		err = out.ParseData(&tmpValue)
	}

	if err != nil {
		return
	}

	value = tmpValue

	return
}
