package client

import (
	"context"
	"fmt"
	"net"
	"testing"

	mock_client "github.com/arthurshafikov/tz-faraway/client/internal/client/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

var address = "localhost:3333"

func TestMakeQuery(t *testing.T) {
	ctrl := gomock.NewController(t)
	connDialerMock := mock_client.NewMockConnDialer(ctrl)
	client := NewClient(connDialerMock)

	l, err := net.Listen("tcp", ":0")
	require.NoError(t, err)
	goroutineEndChan := make(chan int)
	go func() {
		conn, err := l.Accept()
		require.NoError(t, err)

		buffer := make([]byte, 300)
		n, err := conn.Read(buffer)
		require.NoError(t, err)

		require.Equal(
			t,
			"GET / HTTP/1.1\r\nHost: localhost:3333\r\nUser-Agent: Go-http-client/1.1\r\nAccept-Encoding: gzip\r\n\r\n",
			string(buffer[:n]),
		)

		n, err = conn.Write(
			[]byte("HTTP/1.1 200 OK\n\nsomeBody"),
		)
		require.NoError(t, err)
		require.Equal(t, 25, n)

		require.NoError(t, conn.Close())

		goroutineEndChan <- 1
	}()

	conn, err := net.Dial("tcp", l.Addr().String())
	require.NoError(t, err)

	gomock.InOrder(
		connDialerMock.EXPECT().Dial("tcp", address).Return(conn, nil),
		connDialerMock.EXPECT().CloseConnection().Return(nil),
	)

	err = client.MakeQuery(context.Background(), address)
	require.NoError(t, err)

	<-goroutineEndChan
}

func TestMakeQueryWrongAddress(t *testing.T) {
	client := NewClient(nil)

	err := client.MakeQuery(context.Background(), "2som2e!@:WrongAddress/.,,")
	require.ErrorContains(t, err, "invalid port")
}

func TestMakeQueryDialReturnsError(t *testing.T) {
	ctrl := gomock.NewController(t)
	connDialerMock := mock_client.NewMockConnDialer(ctrl)
	client := NewClient(connDialerMock)

	expectedError := fmt.Errorf("some error")
	gomock.InOrder(
		connDialerMock.EXPECT().Dial("tcp", address).Return(nil, expectedError),
		connDialerMock.EXPECT().CloseConnection().Return(nil),
	)

	err := client.MakeQuery(context.Background(), address)
	require.ErrorIs(t, err, expectedError)
}
