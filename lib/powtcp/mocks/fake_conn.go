package mocks

import (
	"fmt"
	"net"
	"time"
)

type FakeConn struct{}

func (f FakeConn) Close() error {
	return fmt.Errorf("some close error")
}

func (f FakeConn) Read(b []byte) (n int, err error)   { panic("not implemented") }
func (f FakeConn) Write(b []byte) (n int, err error)  { panic("not implemented") }
func (f FakeConn) LocalAddr() net.Addr                { panic("not implemented") }
func (f FakeConn) RemoteAddr() net.Addr               { panic("not implemented") }
func (f FakeConn) SetDeadline(t time.Time) error      { panic("not implemented") }
func (f FakeConn) SetReadDeadline(t time.Time) error  { panic("not implemented") }
func (f FakeConn) SetWriteDeadline(t time.Time) error { panic("not implemented") }
