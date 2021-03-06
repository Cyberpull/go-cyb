package socket

import (
	"bytes"
	"fmt"

	"cyberpull.com/go-cyb/errors"
	"cyberpull.com/go-cyb/log"
	"cyberpull.com/go-cyb/objects"
)

type ServerClientInstance struct {
	srv       *Server
	ref       *ServerClientRef
	updater   *ServerClientUpdater
	isRunning bool
	isExiting bool
	isStopped bool
}

func (s *ServerClientInstance) Start() {
	s.isRunning = true
	s.beginInstance()
}

func (s *ServerClientInstance) beginInstance() {
	var err error

	defer func() {
		s.isRunning = false
		s.isExiting = false
		s.isStopped = false

		if r := recover(); r != nil {
			err = errors.From(r)
			log.Errorln(err)
		}
	}()

	if err = s.ref.sendIdentifier(); err != nil {
		log.Errorln(err)
		return
	}

	s.srv.addClientInstance(s)

	defer s.srv.removeClientInstance(s)

	if err = s.srv.execClientInit(s.updater); err != nil {
		log.Errorln(err)
		return
	}

	var input []byte

	for {
		if input, err = s.ref.ReadBytes('\n'); err != nil {
			if s.isStopped {
				err = nil
			}

			break
		}

		go s.processInput(input)
	}
}

func (s *ServerClientInstance) processInput(input []byte) {
	var err error

	defer func() {
		if r := recover(); r != nil {
			err = errors.From(r)
			log.Errorfln("ServerClientInstance: %s", err)
		}
	}()

	input = bytes.TrimSpace(input)

	if len(input) == 0 {
		return
	}

	req := &Request{}

	if err = objects.ParseJSON(input, req); err != nil {
		msg := fmt.Sprintf("Invalid Request: %s", input)
		s.ref.sendError(msg)
	}

	var output *Output

	defer func() {
		if output != nil {
			s.ref.writeResponse(output)
		}
	}()

	ctx := newContext(s, req)

	handler, err := s.srv.handlerCollection.Get(req.Method, req.Channel)

	if err != nil {
		output = ctx.Error(err)
		return
	}

	output = handler(ctx)
}

func (s *ServerClientInstance) Update(args ...any) (err error) {
	return s.updater.Update(args...)
}

func (s *ServerClientInstance) Stop() {
	s.isStopped = true
	s.ref.close()
}

/********************************************/

func newServerClientInstance(srv *Server, ref *ServerClientRef) *ServerClientInstance {
	return &ServerClientInstance{
		srv:     srv,
		ref:     ref,
		updater: newServerClientUpdater(ref),
	}
}
