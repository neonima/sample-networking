package server

import (
	"github.com/neonima/sample-networking/pkg/store"
)

// Server is an interface defining the implementation
type Server interface {
	Start() error
	Close() error
	WithLogger(logger Logger) Server
	WithStore(store store.UserStore) Server
}

// Logger is an interface for to print information
type Logger interface {
	Info(string)
	Fatal(string)
	Error(string)
}
