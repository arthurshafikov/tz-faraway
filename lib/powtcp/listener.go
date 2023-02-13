package powtcp

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"io"
	"log"
	"math/big"
	"net"
	"strconv"
)

const (
	randomStringLength = 20
	defaultDifficulty  = 15
)

type ProowOfWorkProtectionListener struct {
	TCPListener net.Listener

	POWDifficulty int
}

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
		closeConnection(conn)

		return conn, nil
	}

	randomString := randomString(randomStringLength)
	if _, err := writeToConnection(conn, []byte(randomString+":"+fmt.Sprintf("%v", l.POWDifficulty))); err != nil {
		log.Println(err)
		closeConnection(conn)

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
		if _, err := writeToConnection(conn, []byte("nonce is not numeric value")); err != nil {
			log.Println(err)
		}
		closeConnection(conn)

		return conn, nil
	}

	if !l.checkNonceIsValid(l.POWDifficulty, []byte(randomString), nonce) {
		log.Println("nonce is not valid")
		if _, err := writeToConnection(conn, []byte("nonce is not valid")); err != nil {
			log.Println(err)
		}
		closeConnection(conn)

		return conn, nil
	}

	if _, err := writeToConnection(conn, []byte(OKResult)); err != nil {
		log.Println(err)
		closeConnection(conn)
	}

	return conn, nil
}

func (l *ProowOfWorkProtectionListener) Close() error { return l.TCPListener.Close() }

func (l *ProowOfWorkProtectionListener) Addr() net.Addr { return l.TCPListener.Addr() }

func (l *ProowOfWorkProtectionListener) checkNonceIsValid(difficulty int, data []byte, nonce int) bool {
	hash := sha256.Sum256(bytes.Join([][]byte{data, []byte(fmt.Sprintf("%v", nonce))}, []byte{}))

	target := big.NewInt(1)
	target.Lsh(target, uint(256-difficulty))

	var intHash big.Int
	return intHash.SetBytes(hash[:]).Cmp(target) == -1
}
