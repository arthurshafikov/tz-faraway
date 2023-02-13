package powtcp

import (
	"fmt"
	"io"
	"log"
	"math/rand"
	"net"
	"time"
)

const (
	OKResult                         = "OK"
	charsetForRandomStringGeneration = "abcdefghijklmnopqrstuvwxyz" +
		"ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
)

func readFromConnection(conn net.Conn, bufferSize int) ([]byte, error) {
	buffer := make([]byte, bufferSize)
	n, err := conn.Read(buffer)
	if err != nil {
		if err != io.EOF {
			return nil, fmt.Errorf("readFromConnection error: %w", err)
		}
	}

	return buffer[:n], nil
}

func writeToConnection(conn net.Conn, data []byte) (int, error) {
	n, err := conn.Write(data)
	if err != nil {
		return 0, fmt.Errorf("writeToConnection error: %w", err)
	}

	return n, nil
}

func closeConnection(conn net.Conn) {
	if err := conn.Close(); err != nil {
		log.Println(fmt.Errorf("closeConnection error: %w", err))
	}
}

func setWriteDeadline(conn net.Conn, duration time.Duration) error {
	if duration.Nanoseconds() != 0 {
		return conn.SetWriteDeadline(time.Now().Add(duration))
	}

	return nil
}

func setReadDeadline(conn net.Conn, duration time.Duration) error {
	if duration.Nanoseconds() != 0 {
		return conn.SetReadDeadline(time.Now().Add(duration))
	}

	return nil
}

func randomString(length int) string {
	rand.Seed(time.Now().UnixNano())

	b := make([]byte, length)
	for i := range b {
		b[i] = charsetForRandomStringGeneration[rand.Intn(len(charsetForRandomStringGeneration))]
	}

	return string(b)
}
