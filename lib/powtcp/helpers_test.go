package powtcp

import (
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"
	"regexp"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestReadFromConnection(t *testing.T) {
	l, err := net.Listen("tcp", ":0")
	require.NoError(t, err)
	defer func() {
		require.NoError(t, l.Close())
	}()

	acceptFinished := make(chan int)
	go func() {
		conn, err := l.Accept()
		require.NoError(t, err)

		_, err = conn.Write([]byte("someString"))
		require.NoError(t, err)
		require.NoError(t, conn.Close())

		acceptFinished <- 1
	}()

	conn, err := net.Dial("tcp", l.Addr().String())
	require.NoError(t, err)

	res, err := readFromConnection(conn, 5)

	require.NoError(t, err)
	require.Equal(t, "someS", string(res))

	<-acceptFinished
}

func TestReadFromConnectionTimeout(t *testing.T) {
	l, err := net.Listen("tcp", ":0")
	require.NoError(t, err)
	defer func() {
		require.NoError(t, l.Close())
	}()

	acceptFinished := make(chan int)
	go func() {
		conn, err := l.Accept()
		require.NoError(t, err)

		_, err = conn.Write([]byte(randomString(1234)))
		require.NoError(t, err)
		require.NoError(t, conn.Close())

		acceptFinished <- 1
	}()

	conn, err := net.Dial("tcp", l.Addr().String())
	require.NoError(t, err)

	require.NoError(t, conn.SetDeadline(time.Now().Add(-1*time.Hour)))
	_, err = readFromConnection(conn, 5)
	require.ErrorContains(t, err, "i/o timeout")

	<-acceptFinished
}

func TestWriteToConnection(t *testing.T) {
	l, err := net.Listen("tcp", ":0")
	require.NoError(t, err)
	defer func() {
		require.NoError(t, l.Close())
	}()

	acceptFinished := make(chan int)
	go func() {
		conn, err := l.Accept()
		require.NoError(t, err)

		buffer := make([]byte, 20)
		n, err := conn.Read(buffer)
		require.NoError(t, err)
		require.Equal(t, n, 10)
		require.Equal(t, "someString", string(buffer[:n]))
		require.NoError(t, conn.Close())

		acceptFinished <- 1
	}()

	conn, err := net.Dial("tcp", l.Addr().String())
	require.NoError(t, err)

	n, err := writeToConnection(conn, []byte("someString"))
	require.NoError(t, err)
	require.Equal(t, n, 10)

	<-acceptFinished
}

func TestWriteToConnectionWriteTimeout(t *testing.T) {
	l, err := net.Listen("tcp", ":0")
	require.NoError(t, err)
	defer func() {
		require.NoError(t, l.Close())
	}()

	acceptFinished := make(chan int)
	go func() {
		conn, err := l.Accept()
		require.NoError(t, err)

		require.NoError(t, conn.Close())
		acceptFinished <- 1
	}()

	conn, err := net.Dial("tcp", l.Addr().String())
	require.NoError(t, err)

	require.NoError(t, conn.SetWriteDeadline(time.Now()))
	n, err := writeToConnection(conn, []byte("someString"))
	require.ErrorContains(t, err, "i/o timeout")
	require.Zero(t, n)

	<-acceptFinished
}

func TestCloseConnection(t *testing.T) {
	l, err := net.Listen("tcp", ":0")
	require.NoError(t, err)

	go func() {
		_, err := net.Dial("tcp", l.Addr().String())
		require.NoError(t, err)
	}()

	conn, err := l.Accept()
	require.NoError(t, err)

	logOutput := wrapLogOutput(t, func() {
		closeConnection(conn)
	})

	require.Equal(t, "", logOutput)
}

type fakeConn struct{}

func (f fakeConn) Close() error {
	return fmt.Errorf("some close error!")
}

func (f fakeConn) Read(b []byte) (n int, err error)   { panic("not implemented") }
func (f fakeConn) Write(b []byte) (n int, err error)  { panic("not implemented") }
func (f fakeConn) LocalAddr() net.Addr                { panic("not implemented") }
func (f fakeConn) RemoteAddr() net.Addr               { panic("not implemented") }
func (f fakeConn) SetDeadline(t time.Time) error      { panic("not implemented") }
func (f fakeConn) SetReadDeadline(t time.Time) error  { panic("not implemented") }
func (f fakeConn) SetWriteDeadline(t time.Time) error { panic("not implemented") }

func TestCloseConnectionReturnsError(t *testing.T) {
	logOutput := wrapLogOutput(t, func() {
		closeConnection(fakeConn{})
	})

	require.Contains(t, logOutput, "closeConnection error: some close error!")
}

func TestRandomString(t *testing.T) {
	randomString := randomString(200)

	matched, err := regexp.Match("^[a-zA-Z0-9]{200}$", []byte(randomString))
	require.NoError(t, err)
	require.True(t, matched)
}

func wrapLogOutput(t *testing.T, callback func()) string {
	t.Helper()

	reader, writer, err := os.Pipe()
	require.NoError(t, err)
	log.SetOutput(writer)

	callback()

	writer.Close()
	out, err := ioutil.ReadAll(reader)
	require.NoError(t, err)

	return string(out)
}
