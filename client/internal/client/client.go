package client

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"io"
	"log"
	"math"
	"math/big"
	"net"
	"net/http"
	"strconv"
)

func MakeQuery() {
	address := "localhost:3333"
	req, err := http.NewRequest(http.MethodGet, "http://"+address, nil)
	if err != nil {
		log.Fatalln(err)
	}

	connDialer, err := NewConnDialer(address)
	if err != nil {
		log.Fatalln(err)
	}
	defer func() {
		if err := connDialer.CloseConnection(); err != nil {
			log.Fatalln(err)
		}
	}()
	client := http.Client{Transport: &http.Transport{Dial: connDialer.Dial}}
	res, err := client.Do(req)
	if err != nil {
		log.Fatalln(err)
	}

	bodyRes, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatalln(err)
	}
	defer res.Body.Close()

	fmt.Printf("%s %s\n", res.Status, string(bodyRes))
}

type ConnDialer struct {
	c net.Conn
}

func NewConnDialer(address string) (*ConnDialer, error) {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		return nil, err
	}

	return &ConnDialer{
		c: conn,
	}, nil
}

func (cd ConnDialer) CloseConnection() error {
	return cd.c.Close()
}

func (cd ConnDialer) Dial(network, addr string) (net.Conn, error) {
	buffer := make([]byte, 300)
	n, err := cd.c.Read(buffer)
	if err != nil {
		if err != io.EOF {
			return nil, fmt.Errorf("accept() read error: %w", err)
		}
	}

	dataWithDifficulty := bytes.Split(buffer[:n], []byte(":"))
	if len(dataWithDifficulty) != 2 {
		return nil, fmt.Errorf("wrong data with difficulty came from host: %s", string(buffer[:n]))
	}

	data := dataWithDifficulty[0]
	difficulty, err := strconv.Atoi(string(dataWithDifficulty[1]))
	if err != nil || difficulty < 1 {
		return nil, fmt.Errorf("wrong difficulty came from host: %s", string(dataWithDifficulty[1]))
	}

	nonce := findNonce(data, difficulty)
	if _, err := cd.c.Write([]byte(fmt.Sprintf("%v", nonce))); err != nil {
		return nil, fmt.Errorf("connection write error: %w", err)
	}

	return cd.c, nil
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
			fmt.Printf("Nonce = %v\n", nonce)
			fmt.Printf("Result = %s\n", intHash.String())
			fmt.Printf("Target = %s\n", target.String())
			return nonce
		}
	}

	return 0
}
