package powtcp

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"math"
	"math/big"
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
	if err != nil || difficulty < 1 || difficulty > 256 {
		return nil, fmt.Errorf("wrong difficulty came from host: %s", string(dataWithDifficulty[1]))
	}

	nonce := findNonce(dataWithDifficulty[0], difficulty)
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
	if cd.readTimeoutDuration.Nanoseconds() != 0 {
		return cd.conn.SetReadDeadline(time.Now().Add(cd.readTimeoutDuration))
	}

	return nil
}

func (cd *ConnDialer) SetWriteDeadline() error {
	if cd.writeTimeoutDuration.Nanoseconds() != 0 {
		return cd.conn.SetWriteDeadline(time.Now().Add(cd.writeTimeoutDuration))
	}

	return nil
}

func findNonce(data []byte, difficulty int) int {
	target := big.NewInt(1)
	target.Lsh(target, uint(256-difficulty))

	var hash [32]byte
	var intHash big.Int
	for nonce := 0; nonce < math.MaxInt64; nonce++ {
		hash = sha256.Sum256(bytes.Join([][]byte{data, []byte(fmt.Sprintf("%v", nonce))}, []byte{}))

		intHash.SetBytes(hash[:])
		if intHash.Cmp(target) == -1 {
			return nonce
		}
	}

	return 0
}
