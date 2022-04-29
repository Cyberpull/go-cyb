package socket

import (
	"bytes"
	"fmt"
	"os"
	"os/signal"

	"cyberpull.com/go-cyb/errors"
	"cyberpull.com/go-cyb/log"
	"cyberpull.com/go-cyb/objects"
)

type ServerClientInstance struct {
	srv       *Server
	ref       *ServerClientRef
	updater   *ServerClientUpdater
	sig       chan os.Signal
	isRunning bool
	isExiting bool
	isStopped bool
}

func (s *ServerClientInstance) Start() {
	s.isRunning = true

	s.sig = make(chan os.Signal, 1)
	signal.Notify(s.sig, os.Interrupt)

	defer func() {
		signal.Stop(s.sig)
		close(s.sig)
	}()

	go s.beginInstance()

	<-s.sig

	s.ref.conn.Close()
}

func (s *ServerClientInstance) beginInstance() {
	var err error

	defer func() {
		s.isRunning = false
		s.isExiting = false
		s.isStopped = false

		if r := recover(); r != nil {
			err = errors.From(r)
		}

		if err != nil {
			log.Errorfln("ServerClientInstance: %s", err)
		}
	}()

	if err = s.ref.sendIdentifier(); err != nil {
		return
	}

	s.srv.addClientInstance(s)

	defer s.srv.removeClientInstance(s)

	if err = s.srv.execClientInit(s.updater); err != nil {
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

	s.Stop()
}

func (s *ServerClientInstance) processInput(input []byte) {
	var err error

	defer func() {
		if r := recover(); r != nil {
			err = errors.From(r)
		}

		if err != nil {
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
	s.sig <- os.Interrupt
}

/********************************************/

func newServerClientInstance(srv *Server, ref *ServerClientRef) *ServerClientInstance {
	return &ServerClientInstance{
		srv:     srv,
		ref:     ref,
		updater: newServerClientUpdater(ref),
	}
}
