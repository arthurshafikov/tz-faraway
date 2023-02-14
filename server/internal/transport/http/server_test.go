package http

import (
	"context"
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestRunServer(t *testing.T) {
	l, err := net.Listen("tcp", ":0")
	require.NoError(t, err)
	server := NewServer(nil, l)

	go func() {
		time.Sleep(time.Second / 2)
		err := server.server.Shutdown(context.Background())
		require.NoError(t, err)
	}()
	err = server.Serve()
	require.NoError(t, err)
}
