package tcp

import (
	"fmt"
	"net"

	"github.com/neonima/sample-networking/pkg/server"

	"github.com/neonima/sample-networking/pkg/store"

	"github.com/francoispqt/gojay"

	"github.com/neonima/sample-networking/pkg/model"
)

// Client hold data for a client ot work
type Client struct {
	*store.UserContainer
	Conn   net.Conn
	Store  store.UserStore
	Logger server.Logger
	ConnID string
}

// HandleConnection handles a client connection and its business logic
func (c *Client) HandleConnection(unregisterChan chan<- string, errChan chan<- error, broadChan chan<- Broadcast, connID string) {
	defer func() {
		if c.Conn != nil {
			c.Conn.Close()
		}
	}()

	c.Metadata.ConnID = connID
	userChan := model.UserStream(make(chan *model.User))
	decoder := gojay.Stream.BorrowDecoder(c.Conn)

	go func() {
		if err := decoder.DecodeStream(userChan); err != nil {
			errChan <- err
		}
	}()

	for {
		select {
		case statusUpdate := <-userChan:

			if c.User == nil {
				c.User = statusUpdate
				c.Online = true
				if err := c.Store.AddUser(c.UserContainer); err != nil {
					errChan <- err
					continue
				}
				data, connIDS, err := c.GetSendFriendsStatusUpdateData()
				if err != nil {
					errChan <- err
					continue
				}
				BroadcastSignal(data, connIDS, broadChan)

				c.Logger.Info(fmt.Sprintf("user id %v connected", c.ID))
				continue
			}

			if !statusUpdate.Online {
				c.Online = false

				data, connIDS, err := c.GetSendFriendsStatusUpdateData()
				if err != nil {
					errChan <- err
					continue
				}
				BroadcastSignal(data, connIDS, broadChan)

				unregisterChan <- c.ConnID
				c.Logger.Info(fmt.Sprintf("user id %v offline", c.ID))
				return
			}
		case <-decoder.Done():
			if c.User != nil {
				c.Online = false
			}
			unregisterChan <- c.ConnID
			return
		}
	}
}

func (c *Client) GetSendFriendsStatusUpdateData() ([]byte, []string, error) {
	users, err := c.Store.FriendWith(c.UserContainer)
	if err != nil {
		return nil, nil, err
	}
	var connIDs []string
	data, err := gojay.Marshal(c.User)
	if err != nil {
		return nil, nil, err
	}
	for _, u := range users {
		if u.Online {
			connIDs = append(connIDs, u.Metadata.ConnID)
		}
	}
	return data, connIDs, nil
}

func BroadcastSignal(data []byte, connectionsToSendMessage []string, broadChan chan<- Broadcast) {
	broadChan <- Broadcast{
		Data:        data,
		Connections: connectionsToSendMessage,
	}
}
