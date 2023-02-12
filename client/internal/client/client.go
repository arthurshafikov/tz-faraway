package client

import (
	"crypto/sha256"
	"fmt"
	"io"
	"math"
	"math/big"
	"math/rand"
	"net"
	"net/http"
	"strings"
	"time"
)

func MakeQuery() {
	conn, err := net.Dial("tcp", ":3333")
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	client := http.Client{Transport: &http.Transport{Dial: connDialer{conn}.Dial}}

	req, err := http.NewRequest(http.MethodGet, "http://localhost:3333", nil)
	if err != nil {
		panic(err)
	}

	res, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	bodyRes, err := io.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()

	fmt.Printf("%#v\n", res.Status)
	fmt.Printf("%#v\n", string(bodyRes))
}

type connDialer struct {
	c net.Conn
}

func (cd connDialer) Dial(network, addr string) (net.Conn, error) {
	fmt.Printf("%#v\n", "DIAL")
	hash := getHashProtection()
	fmt.Printf("%#v\n", hash)

	if _, err := cd.c.Write([]byte(hash)); err != nil {
		panic(err)
	}

	return cd.c, nil
}

func getHashProtection() string {
	rand.Seed(time.Now().Unix())

	protectionHeaderParts := []string{
		"1",                                   // const LIB::Version,
		"20",                                  // LIB::GetDifficulty()
		time.Now().UTC().Format("0102061504"), // LIB::GetDateFormat
		RandomString(10),                      // LIB::GetRandomStringLenght
	}

	var counter uint64

	var protectionHeader string

	target := big.NewInt(1)
	target.Lsh(target, uint(256-10)) // LIB::GetDifficulty()

	var hash [32]byte
	var intHash big.Int
	for counter = 0; counter < math.MaxUint64; counter++ {
		protectionHeader = strings.Join(append(protectionHeaderParts, fmt.Sprintf("%v", counter)), ":")
		fmt.Printf("%s\n", protectionHeader)

		hash = sha256.Sum256([]byte(protectionHeader))

		intHash.SetBytes(hash[:])

		if intHash.Cmp(target) == -1 {
			fmt.Printf("Nonce = %v\n", counter)
			fmt.Printf("Result = %s\n", intHash.String())
			fmt.Printf("Target = %s\n", target.String())

			break
		}
		fmt.Println(counter)
		fmt.Printf("Result = %s\n", intHash.String())
		fmt.Printf("Target = %s\n", target.String())
	}

	fmt.Println("WON, RESULT IS")
	fmt.Println(protectionHeader)

	return protectionHeader
}

const charset = "abcdefghijklmnopqrstuvwxyz" +
	"ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func RandomString(length int) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}
