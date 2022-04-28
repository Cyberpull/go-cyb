package socket

import "strings"

type ClientUpdateHander func(out *Output)

type ClientUpdateHandlerCollection struct {
	handlers map[string][]ClientUpdateHander
}

func (c *ClientUpdateHandlerCollection) key(method, channel string) string {
	method = strings.ToUpper(method)
	return method + "::" + channel
}

func (c *ClientUpdateHandlerCollection) On(method, channel string, handler ClientUpdateHander) {
	key := c.key(method, channel)

	if _, ok := c.handlers[key]; !ok {
		c.handlers[key] = make([]ClientUpdateHander, 0)
	}

	c.handlers[key] = append(c.handlers[key], handler)
}

func (c *ClientUpdateHandlerCollection) updateAll(method, channel string, out *Output) {
	key := c.key(method, channel)

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
		handlers: make(map[string][]ClientUpdateHander),
	}
}
