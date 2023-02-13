package powtcp

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"log"
	"math/big"
	"net"
	"strconv"
	"time"
)

const (
	randomStringLength = 20
	defaultDifficulty  = 15
)

type ProowOfWorkProtectionListener struct {
	TCPListener net.Listener

	POWDifficulty int

	readTimeoutDuration  time.Duration
	writeTimeoutDuration time.Duration
}

// required - Address, optional - Difficulty (Default: 15), ReadTimeoutDuration, WriteTimeoutDuration
type ListenerOptions struct {
	Address              string
	Difficulty           int
	ReadTimeoutDuration  time.Duration
	WriteTimeoutDuration time.Duration
}

func NewProowOfWorkProtectionListener(opts ListenerOptions) (*ProowOfWorkProtectionListener, error) {
	if opts.Address == "" {
		return nil, fmt.Errorf("empty address for custom listener")
	}
	if opts.Difficulty == 0 {
		opts.Difficulty = defaultDifficulty
	} else if opts.Difficulty > 256 || opts.Difficulty < 1 {
		return nil, fmt.Errorf("difficulty can be from 1 to 256")
	}

	tcpListener, err := net.Listen("tcp", opts.Address)
	if err != nil {
		return nil, err
	}

	return &ProowOfWorkProtectionListener{
		TCPListener:          tcpListener,
		POWDifficulty:        opts.Difficulty,
		readTimeoutDuration:  opts.ReadTimeoutDuration,
		writeTimeoutDuration: opts.WriteTimeoutDuration,
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
	if err := l.writeTextToConn(conn, fmt.Sprintf("%s:%v", randomString, l.POWDifficulty)); err != nil {
		log.Println(err)
		closeConnection(conn)

		return conn, nil
	}

	if err := setReadDeadline(conn, l.readTimeoutDuration); err != nil {
		return nil, err
	}
	res, err := readFromConnection(conn, 300)
	if err != nil {
		log.Println(err)
		closeConnection(conn)

		return conn, nil
	}

	nonce, err := strconv.Atoi(string(res))
	if err != nil {
		if err := l.writeTextToConn(conn, "nonce is not numeric value"); err != nil {
			log.Println(err)
		}
		closeConnection(conn)

		return conn, nil
	}

	if !l.checkNonceIsValid(l.POWDifficulty, []byte(randomString), nonce) {
		if err := l.writeTextToConn(conn, "nonce is not valid"); err != nil {
			log.Println(err)
		}
		closeConnection(conn)

		return conn, nil
	}

	if err := l.writeTextToConn(conn, OKResult); err != nil {
		log.Println(err)

		closeConnection(conn)
	}

	return conn, nil
}

func (l *ProowOfWorkProtectionListener) writeTextToConn(conn net.Conn, text string) error {
	if err := setWriteDeadline(conn, l.writeTimeoutDuration); err != nil {
		return err
	}
	_, err := writeToConnection(conn, []byte(text))

	return err
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
