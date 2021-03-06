package socket

import "strings"

type ClientUpdateHandler func(out *Output)

type ClientUpdateHandlerCollection struct {
	handlers map[string][]ClientUpdateHandler
}

func (c *ClientUpdateHandlerCollection) key(method, channel string) string {
	method = strings.ToUpper(method)
	return method + "::" + channel
}

func (c *ClientUpdateHandlerCollection) On(method, channel string, handler ClientUpdateHandler) {
	key := c.key(method, channel)

	if _, ok := c.handlers[key]; !ok {
		c.handlers[key] = make([]ClientUpdateHandler, 0)
	}

	c.handlers[key] = append(c.handlers[key], handler)
}

func (c *ClientUpdateHandlerCollection) updateAll(out *Output) {
	key := c.key(out.Method, out.Channel)

	handlers, ok := c.handlers[key]

	if !ok {
		return
	}

	for _, handler := range handlers {
		go handler(out)
	}
}

/******************************************************/

func newClientUpdateHandlerCollection() *ClientUpdateHandlerCollection {
	return &ClientUpdateHandlerCollection{
		handlers: make(map[string][]ClientUpdateHandler),
	}
}
