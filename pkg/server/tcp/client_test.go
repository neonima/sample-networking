package tcp_test

import (
	"net"
	"testing"
	"time"

	"github.com/neonima/sample-networking/pkg/mocks"
	"github.com/neonima/sample-networking/pkg/model"

	"github.com/francoispqt/gojay"

	"github.com/neonima/sample-networking/pkg/server"
	"github.com/neonima/sample-networking/pkg/server/tcp"
	"github.com/neonima/sample-networking/pkg/store"
	"github.com/stretchr/testify/require"
)

func getChans(t *testing.T) (unregisterChan chan string, errChan chan error, broadChan chan tcp.Broadcast) {
	t.Helper()
	return make(chan string), make(chan error), make(chan tcp.Broadcast)
}
func newTestClient(t *testing.T, st store.UserStore, logger server.Logger, client net.Conn) *tcp.Client {
	t.Helper()
	userTest := &tcp.Client{
		UserContainer: &store.UserContainer{
			Metadata: &store.Metadata{
				ConnID: "1",
			},
			User: nil,
		},
		Conn:   client,
		Store:  st,
		Logger: logger,
	}
	return userTest
}

func friends(t *testing.T) []*store.UserContainer {
	t.Helper()
	return []*store.UserContainer{
		{
			User: &model.User{
				ID:      2,
				Friends: nil,
				Online:  true,
			},
			Metadata: &store.Metadata{
				ConnID: "2",
			},
		},
		{
			User: &model.User{
				ID:      3,
				Friends: nil,
				Online:  true,
			},
			Metadata: &store.Metadata{
				ConnID: "3",
			},
		},
		{
			User: &model.User{
				ID:      4,
				Friends: nil,
				Online:  true,
			},
			Metadata: &store.Metadata{
				ConnID: "4",
			},
		},
	}
}

func assertMock(t *testing.T, mockLogger *mocks.Logger, mockStore *mocks.UserStore) {
	mockLogger.AssertExpectations(t)
	mockStore.AssertExpectations(t)
}

func TestClient_HandleConnection(t *testing.T) {
	t.Run("should succesfully be online and let friends know about it", func(t *testing.T) {
		s, client := net.Pipe()
		defer func() {
			require.NoError(t, s.Close())
			require.NoError(t, client.Close())
		}()
		mockStore := &mocks.UserStore{}
		mockLogger := &mocks.Logger{}
		unregister, errChan, broadChan := getChans(t)
		userTest := newTestClient(t, mockStore, mockLogger, client)
		expectedUser := &store.UserContainer{
			User: &model.User{
				ID:      1,
				Friends: []int{2, 3, 4},
				Online:  true,
			},
			Metadata: &store.Metadata{
				ConnID: "1",
			},
		}
		mockStore.On("AddUser", expectedUser).Return(nil)
		mockStore.On("FriendWith", expectedUser).Return(friends(t), nil)
		go userTest.HandleConnection(unregister, errChan, broadChan, userTest.Metadata.ConnID)

		require.NoError(t, gojay.NewEncoder(s).Encode(expectedUser.User))
		time.Sleep(time.Millisecond)
		require.Equal(t, expectedUser.User, userTest.User)
		assertMock(t, mockLogger, mockStore)
	})

	t.Run("should return an error when store AddUser call failed", func(t *testing.T) {
		t.Skip("not implemented")
	})

	t.Run("should return an error when store FriendWith call failed", func(t *testing.T) {
		t.Skip("not implemented")
	})

	t.Run("should success with UserInfo and signal offline", func(t *testing.T) {
		s, client := net.Pipe()
		defer func() {
			require.NoError(t, s.Close())
			require.NoError(t, client.Close())
		}()
		mockStore := &mocks.UserStore{}
		mockLogger := &mocks.Logger{}
		unregister, errChan, broadChan := getChans(t)
		userTest := newTestClient(t, mockStore, mockLogger, client)
		expectedUser := &store.UserContainer{
			User: &model.User{
				ID:      1,
				Friends: []int{2, 3, 4},
				Online:  false,
			},
			Metadata: &store.Metadata{
				ConnID: "1",
			},
		}
		userTest.User = expectedUser.User

		mockStore.On("FriendWith", expectedUser).Return(friends(t), nil)
		go userTest.HandleConnection(unregister, errChan, broadChan, userTest.Metadata.ConnID)

		require.NoError(t, gojay.NewEncoder(s).Encode(expectedUser.User))
		time.Sleep(time.Millisecond)
		require.Equal(t, expectedUser.User, userTest.User)
		assertMock(t, mockLogger, mockStore)
	})

	t.Run("should set online to false when connection is closed", func(t *testing.T) {
		t.Skip("not implemented")
	})
}
