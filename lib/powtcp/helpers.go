package powtcp

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"io"
	"log"
	"math/big"
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

func checkNonceIsValid(difficulty uint8, data []byte, nonce int) bool {
	hash := sha256.Sum256(bytes.Join([][]byte{data, []byte(fmt.Sprintf("%v", nonce))}, []byte{}))

	target := big.NewInt(1)
	target.Lsh(target, uint(256-int(difficulty)))

	var intHash big.Int
	return intHash.SetBytes(hash[:]).Cmp(target) == -1
}
