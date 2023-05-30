package tcp_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/neonima/sample-networking/pkg/server/tcp"

	"github.com/neonima/sample-networking/pkg/store/light"

	"github.com/francoispqt/onelog"

	"github.com/stretchr/testify/require"
)

func TestNewTCPServer(t *testing.T) {
	port := 9090
	srv := tcp.NewTCPServer(port)
	require.Equal(t, fmt.Sprintf(":%v", port), srv.(*tcp.TCPServer).Port)
}

func TestTCPServer_Start(t *testing.T) {
	t.Run("should start without error", func(t *testing.T) {
		port := 9090
		srv := tcp.NewTCPServer(port)
		require.NoError(t, srv.Start())
		require.NoError(t, srv.Close())
	})

	t.Run("should accept new connection", func(t *testing.T) {
		t.Skip("not implemented")
	})

	t.Run("should send new client to the router", func(t *testing.T) {
		t.Skip("not implemented")
	})

	t.Run("should fail if address already in use", func(t *testing.T) {
		t.Skip("not implemented")
	})
}

func TestTCPServer_WithLogger(t *testing.T) {
	port := 9090
	srv := tcp.NewTCPServer(port)
	myLogger := onelog.New(os.Stderr, onelog.FATAL)
	srv.WithLogger(myLogger)
	require.Equal(t, myLogger, srv.(*tcp.TCPServer).Logger)
}

func TestTCPServer_WithStore(t *testing.T) {
	port := 9090
	srv := tcp.NewTCPServer(port)
	myStore := light.New()
	srv.WithStore(myStore)
	require.Equal(t, myStore, srv.(*tcp.TCPServer).Store)
}
