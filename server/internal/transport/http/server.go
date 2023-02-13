package http

import (
	"bytes"
	"crypto/sha256"
	"errors"
	"fmt"
	"io"
	"log"
	"math/big"
	"math/rand"
	"net"
	"net/http"
	"strconv"
	"time"
)

func RunServer(handler http.Handler, port string) {
	log.Println("Starting the server on port " + port)

	proowOfWorkProtectionListener, err := NewProowOfWorkProtectionListener(Options{
		Address: ":" + port,
	})
	if err != nil {
		log.Fatalln(err)
	}

	if err := http.Serve(proowOfWorkProtectionListener, handler); err != nil {
		if !errors.Is(err, http.ErrServerClosed) {
			log.Fatalln(err)
		}
	}
}

type ProowOfWorkProtectionListener struct {
	TCPListener net.Listener

	POWDifficulty int
}

const (
	randomStringLength = 20
	defaultDifficulty  = 15
)

// required - Address, optional - Difficulty (Default: 15)
type Options struct {
	Address    string
	Difficulty int
}

func NewProowOfWorkProtectionListener(opts Options) (*ProowOfWorkProtectionListener, error) {
	if opts.Address == "" {
		return nil, fmt.Errorf("empty address for custom listener")
	}
	if opts.Difficulty == 0 {
		opts.Difficulty = defaultDifficulty
	}

	tcpListener, err := net.Listen("tcp", opts.Address)
	if err != nil {
		return nil, err
	}

	return &ProowOfWorkProtectionListener{
		TCPListener:   tcpListener,
		POWDifficulty: opts.Difficulty,
	}, nil
}

func (l *ProowOfWorkProtectionListener) Accept() (net.Conn, error) {
	conn, err := l.TCPListener.Accept()
	if err != nil {
		log.Println(fmt.Errorf("ProowOfWorkProtectionListener.Accept() error: %w", err))
		l.CloseConnection(conn)

		return conn, nil
	}

	hashWord := randomString(randomStringLength)
	if _, err := l.WriteToConnection(conn, []byte(hashWord+":"+fmt.Sprintf("%v", l.POWDifficulty))); err != nil {
		log.Println(err)
		l.CloseConnection(conn)

		return conn, nil
	}

	buffer := make([]byte, 300)
	n, err := conn.Read(buffer)
	if err != nil {
		if err != io.EOF {
			log.Println(fmt.Errorf("ProowOfWorkProtectionListener.Accept() read error: %w", err))
		}
	}

	nonce, err := strconv.Atoi(string(buffer[:n]))
	if err != nil {
		log.Println("nonce is not numeric value")
		if _, err := l.WriteToConnection(conn, []byte("nonce is not numeric value")); err != nil {
			log.Println(err)
		}
		l.CloseConnection(conn)

		return conn, nil
	}

	if !checkNonceIsValid(l.POWDifficulty, []byte(hashWord), nonce) {
		log.Println("nonce is not valid")
		if _, err := l.WriteToConnection(conn, []byte("nonce is not valid")); err != nil {
			log.Println(err)
		}
		l.CloseConnection(conn)

		return conn, nil
	}

	return conn, err
}

func (l *ProowOfWorkProtectionListener) CloseConnection(conn net.Conn) {
	if err := conn.Close(); err != nil {
		log.Println(fmt.Errorf("CloseConnection() error: %w", err))
	}
}

func (l *ProowOfWorkProtectionListener) WriteToConnection(conn net.Conn, data []byte) (int, error) {
	n, err := conn.Write(data)
	if err != nil {
		return 0, fmt.Errorf("write to client error: %w", err)
	}

	return n, nil
}

func (l *ProowOfWorkProtectionListener) Close() error { return l.TCPListener.Close() }

func (l *ProowOfWorkProtectionListener) Addr() net.Addr { return l.TCPListener.Addr() }

func checkNonceIsValid(difficulty int, data []byte, nonce int) bool {
	hash := sha256.Sum256(bytes.Join([][]byte{data, []byte(fmt.Sprintf("%v", nonce))}, []byte{}))

	target := big.NewInt(1)
	target.Lsh(target, uint(256-difficulty))

	var intHash big.Int
	return intHash.SetBytes(hash[:]).Cmp(target) == -1
}

const charset = "abcdefghijklmnopqrstuvwxyz" +
	"ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func randomString(length int) string {
	rand.Seed(time.Now().Unix())

	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}

	return string(b)
}
