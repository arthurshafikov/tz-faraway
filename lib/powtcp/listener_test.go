package powtcp

import (
	"net"
	"testing"
	"time"

	"github.com/arthurshafikov/tz-faraway/lib/powtcp/mocks"
	"github.com/stretchr/testify/require"
	"github.com/tkuchiki/faketime"
)

func TestNewProowOfWorkProtectionListener(t *testing.T) {
	l, err := NewProowOfWorkProtectionListener(ListenerOptions{
		Address: ":0",
	})
	require.NoError(t, err)
	require.Equal(t, l.powDifficulty, defaultDifficulty)

	_, err = net.Dial("tcp", l.Addr().String())
	require.NoError(t, err)
}

func TestNewProowOfWorkProtectionListenerEmptyAddress(t *testing.T) {
	_, err := NewProowOfWorkProtectionListener(ListenerOptions{})
	require.ErrorContains(t, err, "empty address for custom listener")
}

func TestNewProowOfWorkProtectionListenerPortOccupied(t *testing.T) {
	anotherListener, err := net.Listen("tcp", ":0")
	require.NoError(t, err)
	_, err = NewProowOfWorkProtectionListener(ListenerOptions{
		Address: anotherListener.Addr().String(),
	})
	require.ErrorContains(t, err, "bind: address already in use")
}

func TestAccept(t *testing.T) {
	// this is how we get an equal random string every time
	timeNow := time.Date(2022, time.March, 10, 23, 0, 0, 0, time.UTC)
	f := faketime.NewFaketimeWithTime(timeNow)
	defer f.Undo()
	f.Do()

	l, err := NewProowOfWorkProtectionListener(ListenerOptions{
		Address:    ":0",
		Difficulty: 1,
	})
	require.NoError(t, err)

	goroutineEndChan := make(chan int)
	go func() {
		conn, err := net.Dial("tcp", l.Addr().String())
		require.NoError(t, err)

		res, err := readFromConnection(conn, 30)
		require.NoError(t, err)
		require.Len(t, res, randomStringLength+2)

		_, err = writeToConnection(conn, []byte("20"))
		require.NoError(t, err)

		res, err = readFromConnection(conn, 300)
		require.NoError(t, err)
		require.Equal(t, string(res), OKResult)

		goroutineEndChan <- 1
	}()

	conn, err := l.Accept()
	require.NoError(t, err)
	require.NotNil(t, conn)

	<-goroutineEndChan
}

func TestAcceptError(t *testing.T) {
	logOutput := wrapLogOutput(t, func() {
		powListener := ProowOfWorkProtectionListener{
			tcpListener: &mocks.FakeListener{},
		}

		_, err := powListener.Accept()
		require.NoError(t, err)
	})

	require.Contains(t, logOutput, "ProowOfWorkProtectionListener.Accept() error: accept error")
}

func TestAcceptWriteTimeout(t *testing.T) {
	logOutput := wrapLogOutput(t, func() {
		l, err := NewProowOfWorkProtectionListener(ListenerOptions{
			Address:              ":0",
			Difficulty:           1,
			WriteTimeoutDuration: time.Nanosecond,
		})
		require.NoError(t, err)

		goroutineEndChan := make(chan int)
		go func() {
			_, err := net.Dial("tcp", l.Addr().String())
			require.NoError(t, err)

			goroutineEndChan <- 1
		}()

		conn, err := l.Accept()
		require.NoError(t, err)
		require.NotNil(t, conn)

		<-goroutineEndChan
	})

	require.Contains(t, logOutput, "i/o timeout")
}

func TestAcceptReadTimeout(t *testing.T) {
	logOutput := wrapLogOutput(t, func() {
		l, err := NewProowOfWorkProtectionListener(ListenerOptions{
			Address:             ":0",
			Difficulty:          1,
			ReadTimeoutDuration: time.Nanosecond,
		})
		require.NoError(t, err)

		goroutineEndChan := make(chan int)
		go func() {
			conn, err := net.Dial("tcp", l.Addr().String())
			require.NoError(t, err)

			res, err := readFromConnection(conn, 30)
			require.NoError(t, err)
			require.Len(t, res, randomStringLength+2)

			goroutineEndChan <- 1
		}()

		conn, err := l.Accept()
		require.NoError(t, err)
		require.NotNil(t, conn)

		<-goroutineEndChan
	})

	require.Contains(t, logOutput, "i/o timeout")
}

func TestAcceptNonceNotNumeric(t *testing.T) {
	l, err := NewProowOfWorkProtectionListener(ListenerOptions{
		Address:    ":0",
		Difficulty: 1,
	})
	require.NoError(t, err)

	goroutineEndChan := make(chan int)
	go func() {
		conn, err := net.Dial("tcp", l.Addr().String())
		require.NoError(t, err)

		res, err := readFromConnection(conn, 30)
		require.NoError(t, err)
		require.Len(t, res, randomStringLength+2)

		_, err = writeToConnection(conn, []byte("fasafs"))
		require.NoError(t, err)

		res, err = readFromConnection(conn, 300)
		require.NoError(t, err)
		require.Equal(t, string(res), "nonce is not numeric value")

		goroutineEndChan <- 1
	}()

	conn, err := l.Accept()
	require.NoError(t, err)
	require.NotNil(t, conn)

	<-goroutineEndChan
}

func TestAcceptNonceInvalid(t *testing.T) {
	// this is how we get an equal random string every time
	timeNow := time.Date(2022, time.March, 10, 23, 0, 0, 0, time.UTC)
	f := faketime.NewFaketimeWithTime(timeNow)
	defer f.Undo()
	f.Do()

	l, err := NewProowOfWorkProtectionListener(ListenerOptions{
		Address:    ":0",
		Difficulty: 255,
	})
	require.NoError(t, err)

	goroutineEndChan := make(chan int)
	go func() {
		conn, err := net.Dial("tcp", l.Addr().String())
		require.NoError(t, err)

		res, err := readFromConnection(conn, 30)
		require.NoError(t, err)
		require.Len(t, res, randomStringLength+4)

		_, err = writeToConnection(conn, []byte("1"))
		require.NoError(t, err)

		res, err = readFromConnection(conn, 300)
		require.NoError(t, err)
		require.Equal(t, string(res), "nonce is not valid")

		goroutineEndChan <- 1
	}()

	conn, err := l.Accept()
	require.NoError(t, err)
	require.NotNil(t, conn)

	<-goroutineEndChan
}

func TestClose(t *testing.T) {
	powListener := ProowOfWorkProtectionListener{
		tcpListener: &mocks.FakeListener{},
	}

	err := powListener.Close()
	require.ErrorContains(t, err, "close error")
}
