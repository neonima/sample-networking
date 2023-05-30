package tcp

import (
	"net"

	"github.com/neonima/sample-networking/pkg/server"

	"github.com/rs/xid"
)

// Router holds data for the router to work
type Router struct {
	Clients    map[string]net.Conn
	Register   chan *Client
	Unregister chan string
	ErrorChan  chan error
	Broadcast  chan Broadcast
	Logger     server.Logger
	closeChan  chan interface{}
}

// Broadcast holds the data and the connection IDs to be sent to
type Broadcast struct {
	Data        []byte
	Connections []string
}

// NewRouter returns a new router
func NewRouter(logger server.Logger) *Router {
	return &Router{
		Clients:    make(map[string]net.Conn),
		Register:   make(chan *Client),
		Unregister: make(chan string),
		ErrorChan:  make(chan error),
		Broadcast:  make(chan Broadcast),
		Logger:     logger,
		closeChan:  make(chan interface{}, 1),
	}
}

// Run starts the router
func (h *Router) Run() {
	for {
		select {
		case <-h.closeChan:
			return

		case client := <-h.Register:
			connID := xid.New().String()
			h.Clients[connID] = client.Conn
			go client.HandleConnection(h.Unregister, h.ErrorChan, h.Broadcast, connID)

		case connID := <-h.Unregister:
			delete(h.Clients, connID)

		case caster := <-h.Broadcast:
			for _, connID := range caster.Connections {
				_, err := h.Clients[connID].Write(caster.Data)
				if err != nil {
					h.Logger.Error(err.Error())
				}
			}

		case err := <-h.ErrorChan:
			h.Logger.Error(err.Error())
		}
	}
}

// Close implements io.Closer
func (h *Router) Close() error {
	h.closeChan <- struct{}{}
	return nil
}
