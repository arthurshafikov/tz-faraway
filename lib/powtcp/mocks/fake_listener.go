package mocks

import (
	"fmt"
	"net"
)

type FakeListener struct{}

func (f *FakeListener) Accept() (net.Conn, error) {
	return &net.TCPConn{}, fmt.Errorf("accept error")
}

func (f *FakeListener) Close() error {
	return fmt.Errorf("close error")
}

func (f *FakeListener) Addr() net.Addr {
	return nil
}
