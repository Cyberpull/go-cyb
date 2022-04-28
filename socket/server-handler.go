package socket

import (
	"strings"

	"cyberpull.com/go-cyb/errors"
)

type ServerHandler func(ctx *Context) *Output

type ServerHandlerCollection struct {
	handlers map[string]ServerHandler
}

func (s *ServerHandlerCollection) key(method, channel string) string {
	method = strings.ToUpper(method)
	return method + "::" + channel
}

func (s *ServerHandlerCollection) Has(method, channel string) bool {
	key := s.key(method, channel)
	_, ok := s.handlers[key]
	return ok
}

func (s *ServerHandlerCollection) Get(method, channel string) (handler ServerHandler, err error) {
	key := s.key(method, channel)

	handler, ok := s.handlers[key]

	if !ok {
		err = errors.Newf(`No action found for "%s" -> "%s"`, 400, method, channel)
		return
	}

	return
}

func (s *ServerHandlerCollection) On(method, channel string, handler ServerHandler) (err error) {
	if s.Has(method, channel) {
		err = errors.Newf(`Action already exists for "%s" -> "%s"`, 500, method, channel)
		return
	}

	key := s.key(method, channel)
	s.handlers[key] = handler

	return
}

func (s *ServerHandlerCollection) Off(method, channel string) {
	if s.Has(method, channel) {
		key := s.key(method, channel)
		delete(s.handlers, key)
	}
}

/**********************************************/

func newServerHandlerCollection() *ServerHandlerCollection {
	return &ServerHandlerCollection{
		handlers: make(map[string]ServerHandler),
	}
}
