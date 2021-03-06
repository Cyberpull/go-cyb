package socket

import (
	"crypto/tls"
	"net"
	"sync"

	"cyberpull.com/go-cyb/cert"
	"cyberpull.com/go-cyb/errors"
	"cyberpull.com/go-cyb/log"
	"cyberpull.com/go-cyb/objects"
)

type ServerBootSubscriber func() (err error)
type ServerAuthSubscriber func(ref *ServerClientRef) (err error)
type ServerHandlerSubscriber func(subscriber *ServerHandlerCollection) (err error)
type ServerSetupSubscriber func(ref *ServerClientRef) (err error)
type ServerCleanupSubscriber func(ref *ServerClientRef) (err error)

type ServerClientInitHandler func(updater *ServerClientUpdater) (err error)

type Server struct {
	err                error
	opts               ServerOptions
	mutex              sync.Mutex
	listener           net.Listener
	clients            *objects.Array[*ServerClientInstance]
	handlerCollection  *ServerHandlerCollection
	clientInitHandlers []ServerClientInitHandler
	bootSubscribers    []ServerBootSubscriber
	authSubscribers    []ServerAuthSubscriber
	handlerSubscribers []ServerHandlerSubscriber
	setupSubscribers   []ServerSetupSubscriber
	cleanupSubscribers []ServerCleanupSubscriber
	isStarting         bool
	isListening        bool
}

func (s *Server) Boot(subs ...ServerBootSubscriber) {
	s.bootSubscribers = append(s.bootSubscribers, subs...)
}

func (s *Server) execBoot() (err error) {
	for _, subscriber := range s.bootSubscribers {
		err = subscriber()

		if err != nil {
			break
		}
	}

	return
}

func (s *Server) Handlers(subs ...ServerHandlerSubscriber) {
	s.handlerSubscribers = append(s.handlerSubscribers, subs...)
}

func (s *Server) execHandlers() (err error) {
	for _, subscriber := range s.handlerSubscribers {
		err = subscriber(s.handlerCollection)

		if err != nil {
			break
		}
	}

	return
}

func (s *Server) Auth(subs ...ServerAuthSubscriber) {
	s.authSubscribers = append(s.authSubscribers, subs...)
}

func (s *Server) execAuth(ref *ServerClientRef) (err error) {
	for _, subscriber := range s.authSubscribers {
		err = subscriber(ref)

		if err != nil {
			break
		}
	}

	return
}

func (s *Server) Setup(subs ...ServerSetupSubscriber) {
	s.setupSubscribers = append(s.setupSubscribers, subs...)
}

func (s *Server) execSetup(ref *ServerClientRef) (err error) {
	for _, subscriber := range s.setupSubscribers {
		err = subscriber(ref)

		if err != nil {
			break
		}
	}

	return
}

func (s *Server) Cleanup(subs ...ServerCleanupSubscriber) {
	s.cleanupSubscribers = append(s.cleanupSubscribers, subs...)
}

func (s *Server) execCleanup(ref *ServerClientRef) (err error) {
	for _, subscriber := range s.cleanupSubscribers {
		err = subscriber(ref)

		if err != nil {
			break
		}
	}

	return
}

func (s *Server) ClientInit(handlers ...ServerClientInitHandler) {
	s.clientInitHandlers = append(s.clientInitHandlers, handlers...)
}

func (s *Server) execClientInit(updater *ServerClientUpdater) (err error) {
	for _, handler := range s.clientInitHandlers {
		err = handler(updater)

		if err != nil {
			break
		}
	}

	return
}

func (s *Server) addClientInstance(i *ServerClientInstance) {
	s.mutex.Lock()

	defer s.mutex.Unlock()

	s.clients.Push(i)
}

func (s *Server) removeClientInstance(i *ServerClientInstance) {
	s.mutex.Lock()

	defer s.mutex.Unlock()

	index := s.clients.IndexOf(i)

	if index >= 0 {
		s.clients.Splice(index, 1)
	}
}

func (s *Server) Listen(errChan ...chan error) {
	var err error

	defer func() {
		s.isStarting = false
		s.isListening = false

		if r := recover(); r != nil {
			err = errors.From(r)
			writeOne(errChan, err)
		}
	}()

	s.isStarting = true

	if err = sanitizeNameAndAlias(&s.opts); err != nil {
		writeOne(errChan, err)
		return
	}

	if err = s.execBoot(); err != nil {
		writeOne(errChan, err)
		return
	}

	if err = s.execHandlers(); err != nil {
		writeOne(errChan, err)
		return
	}

	if err = sanitizeTlsConfig(&s.opts); err != nil {
		writeOne(errChan, err)
		return
	}

	address := address(&s.opts)

	if cert.IsEnabled() {
		s.listener, err = tls.Listen("tcp", address, s.opts.TlsConfig)
	} else {
		s.listener, err = net.Listen("tcp", address)
	}

	if err != nil {
		writeOne(errChan, err)
		return
	}

	s.isStarting = false
	s.isListening = true

	defer s.listener.Close()

	log.Successfln("%s listening on %s", s.opts.Name, address)

	writeOne(errChan, nil)

	var conn net.Conn

	for {
		var err error

		if conn, err = s.listener.Accept(); err != nil {
			break
		}

		go s.handleIncomingConnection(conn)
	}
}

func (s *Server) handleIncomingConnection(conn net.Conn) {
	var (
		err error
		ref *ServerClientRef
	)

	defer func() {
		if r := recover(); r != nil {
			err = errors.From(r)
		}

		if err != nil {
			if ref != nil {
				ref.sendError(err.Error())
			} else {
				log.Errorln(err)
			}
		}
	}()

	if ref, err = newServerClientRef(conn, s.opts); err != nil {
		return
	}

	if err = s.execAuth(ref); err != nil {
		return
	}

	defer s.execCleanup(ref)

	if err = s.execSetup(ref); err != nil {
		return
	}

	instance := newServerClientInstance(s, ref)

	instance.Start()
}

func (s *Server) Stop() error {
	for {
		if s.clients.Length() == 0 {
			break
		}

		s.clients.Pop().Stop()
	}

	return s.listener.Close()
}

func (s *Server) UpdateAll(args ...any) {
	s.clients.ForEach(func(client *ServerClientInstance, _ int) {
		go client.Update(args...)
	})
}

func (s *Server) IsListening() bool {
	return s.isListening
}

/****************************************************/

func NewServer(opts ServerOptions) *Server {
	srv := &Server{
		opts:               opts,
		clients:            objects.NewArray[*ServerClientInstance](),
		handlerCollection:  newServerHandlerCollection(),
		bootSubscribers:    make([]ServerBootSubscriber, 0),
		handlerSubscribers: make([]ServerHandlerSubscriber, 0),
		authSubscribers:    make([]ServerAuthSubscriber, 0),
		setupSubscribers:   make([]ServerSetupSubscriber, 0),
		cleanupSubscribers: make([]ServerCleanupSubscriber, 0),
	}

	return srv
}
