package http

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"io"
	"log"
	"math/big"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func RunServer(handler http.Handler) {
	log.Println("Starting the server on port 3333")

	tcpListener, err := net.Listen("tcp", ":3333")
	if err != nil {
		panic(err)
	}
	defer tcpListener.Close()

	if err := http.Serve(&CustomListener{tcpListener}, handler); err != nil {
		if !errors.Is(err, http.ErrServerClosed) {
			log.Fatalln(err)
		}
	}
}

type CustomListener struct {
	TCPListener net.Listener
}

func (l *CustomListener) Accept() (net.Conn, error) {
	fmt.Printf("%#v\n", "accept")
	c, err := l.TCPListener.Accept()
	if err != nil {
		panic(err)
	}

	buffer := make([]byte, 300)
	n, err := c.Read(buffer)
	if err != nil {
		if err != io.EOF {
			log.Fatalln("read error:", err)
		}
	}

	protectionCodeBytes := append(make([]byte, 0, n), buffer[:n]...)

	fmt.Printf("protectionCodeString: %s\n", string(protectionCodeBytes))
	if !checkHashIsValid(string(protectionCodeBytes)) {
		c.Close()
		fmt.Printf("%#v\n", "declined")

		return c, nil
	}
	fmt.Printf("%#v\n", "accepted")

	return c, err
}

func (l *CustomListener) Close() error { return l.TCPListener.Close() }

func (l *CustomListener) Addr() net.Addr { return l.TCPListener.Addr() }

func checkHashIsValid(protectionCode string) bool {
	if protectionCode == "" {
		log.Println("protectionCode == \"\"")

		return false
	}

	protectionCodeParts := strings.Split(protectionCode, ":")
	if len(protectionCodeParts) != 5 {
		log.Println("!= 5")

		return false
	}

	version := protectionCodeParts[0]
	if version != "1" { // const LIB::LIB_VERSION
		log.Println("const LIB::LIB_VERSION")

		return false
	}

	difficulty, err := strconv.Atoi(protectionCodeParts[1])
	if err != nil {
		log.Println(err)

		return false
	}
	if difficulty < 10 { // const LIB::DIFFICULTY??? libClient->newClient(difficulty????)
		log.Println("low difficulty")

		return false
	}

	dateTime, err := time.Parse("0102061504", protectionCodeParts[2])
	if err != nil {
		log.Println(err)

		return false
	}

	if dateTime.Add(time.Hour).Before(time.Now().UTC()) { // newClient() -> validity time ???? depends on difficulty
		log.Println(dateTime, time.Now().UTC())
		log.Println("OUTDATED!!!")

		return false
	}

	if len(protectionCodeParts[3]) < 5 { // const LIB::Minimal_Word_Len
		log.Println("minimal len is 5")

		return false
	}

	if _, err := strconv.Atoi(protectionCodeParts[4]); err != nil { // check that counter is a number
		log.Println("strconv.Atoi(string(counterDecodedBytes))")
		log.Println(err)

		return false
	}

	protectionCodeSum256 := sha256.Sum256([]byte(protectionCode))

	fmt.Printf("%x\n", protectionCodeSum256)
	target := big.NewInt(1)
	target.Lsh(target, uint(256-10)) // const LIB::DIFFICULTY???

	fmt.Printf("%x\n", target)
	var intHash big.Int
	if intHash.SetBytes(protectionCodeSum256[:]).Cmp(target) == -1 {
		return true
	}

	log.Println(intHash.SetBytes(protectionCodeSum256[:]).Cmp(target) != -1)

	return false
}
