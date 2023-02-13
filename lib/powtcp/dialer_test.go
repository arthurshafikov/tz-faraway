package powtcp

import (
	"net"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/arthurshafikov/tz-faraway/lib/powtcp/mocks"
)

func TestNewConnDialer(t *testing.T) {
	l, err := net.Listen("tcp", ":0")
	require.NoError(t, err)

	goroutineEndChan := make(chan int)
	go func() {
		connDialer, err := NewConnDialer(ConnDialerOptions{
			Address: l.Addr().String(),
		})
		require.NoError(t, err)

		buffer := make([]byte, 10)
		n, err := connDialer.conn.Read(buffer)
		require.NoError(t, err)
		require.Equal(t, "someString", string(buffer[:n]))

		goroutineEndChan <- 1
	}()

	conn, err := l.Accept()
	require.NoError(t, err)
	defer func() {
		require.NoError(t, conn.Close())
	}()

	_, err = conn.Write([]byte("someString"))
	require.NoError(t, err)

	<-goroutineEndChan
}

func TestNewConnDialerReturnsError(t *testing.T) {
	l, err := net.Listen("tcp", ":0")
	require.NoError(t, err)
	require.NoError(t, l.Close())

	_, err = NewConnDialer(ConnDialerOptions{
		Address: l.Addr().String(),
	})
	require.ErrorContains(t, err, "connect: connection refused")
}

func TestConnDialerCloseConnection(t *testing.T) {
	cd := ConnDialer{
		conn: mocks.FakeConn{},
	}

	err := cd.CloseConnection()
	require.ErrorContains(t, err, "some close error")
}

func TestDial(t *testing.T) {
	l, err := net.Listen("tcp", ":0")
	require.NoError(t, err)

	goroutineEndChan := make(chan int)
	go func() {
		connDialer, err := NewConnDialer(ConnDialerOptions{
			Address: l.Addr().String(),
		})
		require.NoError(t, err)

		conn, err := connDialer.Dial("tcp", connDialer.conn.RemoteAddr().String())
		require.NoError(t, err)
		require.NotNil(t, conn)

		goroutineEndChan <- 1
	}()

	conn, err := l.Accept()
	require.NoError(t, err)
	defer func() {
		require.NoError(t, conn.Close())
	}()

	n, err := conn.Write([]byte("randomString:8"))
	require.NoError(t, err)
	require.Equal(t, len("randomString:8"), n)

	buffer := make([]byte, 20)
	n, err = conn.Read(buffer)
	require.NoError(t, err)

	nonce, err := strconv.Atoi(string(buffer[:n]))
	require.NoError(t, err)
	require.NotZero(t, nonce)

	_, err = conn.Write([]byte(OKResult))
	require.NoError(t, err)

	<-goroutineEndChan
}

func TestDialReadTimeout(t *testing.T) {
	l, err := net.Listen("tcp", ":0")
	require.NoError(t, err)

	goroutineEndChan := make(chan int)
	go func() {
		connDialer, err := NewConnDialer(ConnDialerOptions{
			Address:             l.Addr().String(),
			ReadTimeoutDuration: time.Microsecond,
		})
		require.NoError(t, err)

		conn, err := connDialer.Dial("tcp", connDialer.conn.RemoteAddr().String())
		require.ErrorContains(t, err, "i/o timeout")
		require.Nil(t, conn)

		goroutineEndChan <- 1
	}()

	conn, err := l.Accept()
	require.NoError(t, err)
	defer func() {
		require.NoError(t, conn.Close())
	}()

	<-goroutineEndChan
}

func TestDialWrongDataWithDifficulty(t *testing.T) {
	l, err := net.Listen("tcp", ":0")
	require.NoError(t, err)

	goroutineEndChan := make(chan int)
	go func() {
		connDialer, err := NewConnDialer(ConnDialerOptions{
			Address: l.Addr().String(),
		})
		require.NoError(t, err)

		conn, err := connDialer.Dial("tcp", connDialer.conn.RemoteAddr().String())
		require.ErrorContains(t, err, "wrong data with difficulty came from host: randomStringdsajdsa")
		require.Nil(t, conn)

		goroutineEndChan <- 1
	}()

	conn, err := l.Accept()
	require.NoError(t, err)
	defer func() {
		require.NoError(t, conn.Close())
	}()

	n, err := conn.Write([]byte("randomStringdsajdsa"))
	require.NoError(t, err)
	require.Equal(t, len("randomStringdsajdsa"), n)

	<-goroutineEndChan
}

func TestDialDifficultyTooHigh(t *testing.T) {
	l, err := net.Listen("tcp", ":0")
	require.NoError(t, err)

	goroutineEndChan := make(chan int)
	go func() {
		connDialer, err := NewConnDialer(ConnDialerOptions{
			Address: l.Addr().String(),
		})
		require.NoError(t, err)

		conn, err := connDialer.Dial("tcp", connDialer.conn.RemoteAddr().String())
		require.ErrorContains(t, err, "wrong difficulty came from host: 232321")
		require.Nil(t, conn)

		goroutineEndChan <- 1
	}()

	conn, err := l.Accept()
	require.NoError(t, err)
	defer func() {
		require.NoError(t, conn.Close())
	}()

	n, err := conn.Write([]byte("randomString:232321"))
	require.NoError(t, err)
	require.Equal(t, len("randomString:232321"), n)

	<-goroutineEndChan
}

func TestDialWriteTimeout(t *testing.T) {
	l, err := net.Listen("tcp", ":0")
	require.NoError(t, err)

	goroutineEndChan := make(chan int)
	go func() {
		connDialer, err := NewConnDialer(ConnDialerOptions{
			Address:              l.Addr().String(),
			WriteTimeoutDuration: time.Nanosecond,
		})
		require.NoError(t, err)

		conn, err := connDialer.Dial("tcp", connDialer.conn.RemoteAddr().String())
		require.ErrorContains(t, err, "i/o timeout")
		require.Nil(t, conn)

		goroutineEndChan <- 1
	}()

	conn, err := l.Accept()
	require.NoError(t, err)
	defer func() {
		require.NoError(t, conn.Close())
	}()

	n, err := conn.Write([]byte("randomString:8"))
	require.NoError(t, err)
	require.Equal(t, len("randomString:8"), n)

	<-goroutineEndChan
}

func TestDialNonOkResult(t *testing.T) {
	l, err := net.Listen("tcp", ":0")
	require.NoError(t, err)

	goroutineEndChan := make(chan int)
	go func() {
		connDialer, err := NewConnDialer(ConnDialerOptions{
			Address: l.Addr().String(),
		})
		require.NoError(t, err)

		conn, err := connDialer.Dial("tcp", connDialer.conn.RemoteAddr().String())
		require.ErrorContains(t, err, "host responded: some wrong result")
		require.Nil(t, conn)

		goroutineEndChan <- 1
	}()

	conn, err := l.Accept()
	require.NoError(t, err)
	defer func() {
		require.NoError(t, conn.Close())
	}()

	n, err := conn.Write([]byte("randomString:8"))
	require.NoError(t, err)
	require.Equal(t, len("randomString:8"), n)

	buffer := make([]byte, 20)
	n, err = conn.Read(buffer)
	require.NoError(t, err)

	nonce, err := strconv.Atoi(string(buffer[:n]))
	require.NoError(t, err)
	require.NotZero(t, nonce)

	_, err = conn.Write([]byte("some wrong result"))
	require.NoError(t, err)

	<-goroutineEndChan
}
