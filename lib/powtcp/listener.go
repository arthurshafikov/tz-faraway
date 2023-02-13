package powtcp

import (
	"fmt"
	"log"
	"net"
	"strconv"
	"time"
)

const (
	randomStringLength       = 20
	defaultDifficulty  uint8 = 15
)

type ProowOfWorkProtectionListener struct {
	tcpListener net.Listener

	powDifficulty uint8

	readTimeoutDuration  time.Duration
	writeTimeoutDuration time.Duration
}

// required - Address, optional - Difficulty (Default: 15), ReadTimeoutDuration, WriteTimeoutDuration.
type ListenerOptions struct {
	Address              string
	Difficulty           uint8
	ReadTimeoutDuration  time.Duration
	WriteTimeoutDuration time.Duration
}

func NewProowOfWorkProtectionListener(opts ListenerOptions) (*ProowOfWorkProtectionListener, error) {
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
		tcpListener:          tcpListener,
		powDifficulty:        opts.Difficulty,
		readTimeoutDuration:  opts.ReadTimeoutDuration,
		writeTimeoutDuration: opts.WriteTimeoutDuration,
	}, nil
}

func (l *ProowOfWorkProtectionListener) Accept() (net.Conn, error) {
	conn, err := l.tcpListener.Accept()
	if err != nil {
		log.Println(fmt.Errorf("ProowOfWorkProtectionListener.Accept() error: %w", err))
		closeConnection(conn)

		return conn, nil
	}

	randomString := randomString(randomStringLength)
	if err := l.writeTextToConn(conn, fmt.Sprintf("%s:%v", randomString, l.powDifficulty)); err != nil {
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

	if !checkNonceIsValid(l.powDifficulty, []byte(randomString), nonce) {
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

func (l *ProowOfWorkProtectionListener) Close() error { return l.tcpListener.Close() }

func (l *ProowOfWorkProtectionListener) Addr() net.Addr { return l.tcpListener.Addr() }
