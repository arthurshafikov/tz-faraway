package powtcp

import (
	"bytes"
	"fmt"
	"math"
	"net"
	"strconv"
	"time"
)

type ConnDialer struct {
	conn                 net.Conn
	readTimeoutDuration  time.Duration
	writeTimeoutDuration time.Duration
}

type ConnDialerOptions struct {
	Address string

	ReadTimeoutDuration  time.Duration
	WriteTimeoutDuration time.Duration
}

func NewConnDialer(opts ConnDialerOptions) (*ConnDialer, error) {
	conn, err := net.Dial("tcp", opts.Address)
	if err != nil {
		return nil, err
	}

	return &ConnDialer{
		conn:                 conn,
		readTimeoutDuration:  opts.ReadTimeoutDuration,
		writeTimeoutDuration: opts.WriteTimeoutDuration,
	}, nil
}

func (cd ConnDialer) CloseConnection() error {
	return cd.conn.Close()
}

func (cd ConnDialer) Dial(network, addr string) (net.Conn, error) {
	if err := cd.SetReadDeadline(); err != nil {
		return nil, err
	}
	result, err := readFromConnection(cd.conn, randomStringLength+4) // randomStringLength:256
	if err != nil {
		return nil, err
	}

	dataWithDifficulty := bytes.Split(result, []byte(":"))
	if len(dataWithDifficulty) != 2 {
		return nil, fmt.Errorf("wrong data with difficulty came from host: %s", string(result))
	}

	difficulty, err := strconv.Atoi(string(dataWithDifficulty[1]))
	if err != nil || difficulty < 1 || difficulty > 255 {
		return nil, fmt.Errorf("wrong difficulty came from host: %s", string(dataWithDifficulty[1]))
	}

	nonce := cd.findNonce(dataWithDifficulty[0], difficulty)
	if err := cd.SetWriteDeadline(); err != nil {
		return nil, err
	}
	if _, err := writeToConnection(cd.conn, []byte(fmt.Sprintf("%v", nonce))); err != nil {
		return nil, err
	}

	if err := cd.SetReadDeadline(); err != nil {
		return nil, err
	}
	result, err = readFromConnection(cd.conn, 300)
	if err != nil {
		return nil, err
	}

	if string(result) != OKResult {
		return nil, fmt.Errorf("host responded: %s", string(result))
	}

	return cd.conn, nil
}

func (cd *ConnDialer) SetReadDeadline() error {
	return setReadDeadline(cd.conn, cd.readTimeoutDuration)
}

func (cd *ConnDialer) SetWriteDeadline() error {
	return setWriteDeadline(cd.conn, cd.writeTimeoutDuration)
}

func (cd *ConnDialer) findNonce(data []byte, difficulty int) int {
	for nonce := 0; nonce < math.MaxInt64; nonce++ {
		if checkNonceIsValid(uint8(difficulty), data, nonce) {
			return nonce
		}
	}

	return 0
}
