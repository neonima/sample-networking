package tcp

import (
	"fmt"
	"net"
	"os"

	"github.com/neonima/sample-networking/pkg/server"

	"github.com/neonima/sample-networking/pkg/store"
	"github.com/neonima/sample-networking/pkg/store/light"

	"github.com/francoispqt/onelog"

	"errors"
)

type TCPServer struct {
	Port      string
	Listener  net.Listener
	Router    *Router
	Logger    server.Logger
	Store     store.UserStore
	closeChan chan interface{}
}

// NewTCPServer return a tcp server
func NewTCPServer(port int) server.Server {
	logger := onelog.New(
		os.Stdout,
		onelog.ERROR|onelog.FATAL|onelog.INFO)

	return &TCPServer{
		Port:      fmt.Sprintf(":%v", port),
		Logger:    logger,
		Store:     light.New(),
		closeChan: make(chan interface{}, 1),
		Router:    NewRouter(logger),
	}
}

// Start starts the server
func (s *TCPServer) Start() error {
	listener, err := net.Listen("tcp", s.Port)
	if err != nil {
		return err
	}
	s.Logger.Info(fmt.Sprintf("ALO started on port %v", s.Port))
	s.Listener = listener

	s.run()

	return nil
}

func (s *TCPServer) run() {
	go s.Router.Run()
	go s.listen()
}

func (s *TCPServer) listen() {
	var addrErr *net.AddrError
	for {
		select {
		case <-s.closeChan:
			return
		default:
			if s.Listener == nil {
				if err := s.Close(); err != nil {
					s.Logger.Error(err.Error())
				}
				continue
			}
			c, err := s.Listener.Accept()
			if err != nil && errors.Is(err, addrErr) {
				s.Logger.Error(err.Error())
			}
			s.Router.Register <- &Client{
				Conn: c,
				UserContainer: &store.UserContainer{
					Metadata: &store.Metadata{},
				},
				Store:  s.Store,
				Logger: s.Logger,
			}
		}

	}
}

func (s *TCPServer) Close() error {
	return func() error {
		s.closeChan <- struct{}{}
		var err error
		if s.Listener != nil {
			err = s.Listener.Close()
			if err != nil {
				err = fmt.Errorf("problem while closing the server, %w", err)
			}
		}

		if s.Router != nil {
			if err = s.Router.Close(); err != nil {
				err = fmt.Errorf("problem while closing the hub, %w", err)
			}
		}

		s.Logger.Info("server ended")
		return err
	}()
}

// WithLogger allows to set a custom logger
func (s *TCPServer) WithLogger(logger server.Logger) server.Server {
	s.Logger = logger
	s.Router.Logger = logger
	return s
}

// WithStore allows to set a custom a store
func (s *TCPServer) WithStore(store store.UserStore) server.Server {
	s.Store = store
	return s
}
