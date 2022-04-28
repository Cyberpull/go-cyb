package socket

import (
	"cyberpull.com/go-cyb/errors"
)

type Context struct {
	ref      *ServerClientRef
	instance *ServerClientInstance
	request  *Request
}

func (c *Context) createOutput() *Output {
	return &Output{uuid: c.request.UUID}
}

func (c *Context) Error(v any, code ...int) *Output {
	value := c.createOutput()

	if x, ok := v.(error); ok {
		xErr := errors.From(x, code...)

		value.Code = xErr.Code()
		value.SetData(xErr.Error())

		return value
	}

	if len(code) == 0 {
		code = append(code, 500)
	}

	value.Code = code[0]
	value.SetData(v)

	return value
}

func (c *Context) Success(v any) *Output {
	value := c.createOutput()
	value.Code = 200
	value.SetData(v)

	return value
}

/**********************************************/

func newContext(i *ServerClientInstance, r *Request) *Context {
	return &Context{
		instance: i,
		ref:      i.ref,
		request:  r,
	}
}
