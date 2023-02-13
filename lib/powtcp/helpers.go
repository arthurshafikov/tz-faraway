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
	OKResult = "OK"
	charset  = "abcdefghijklmnopqrstuvwxyz" +
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

func randomString(length int) string {
	rand.Seed(time.Now().Unix() + 123123)

	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}

	return string(b)
}
