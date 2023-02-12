package main

import (
	"errors"
	"fmt"
	"html"
	"io/ioutil"
	"net"
	"net/http"
	"time"
)

type connDialer struct {
	c net.Conn
}

func (cd connDialer) Dial(network, addr string) (net.Conn, error) {
	return cd.c, nil
}

type listener struct{ c net.Conn }

func (l *listener) Accept() (net.Conn, error) {
	if l.c != nil {
		c := l.c
		l.c = nil
		return c, nil
	}
	time.Sleep(time.Minute)
	return nil, errors.New("Blah")
}

func (l *listener) Close() error { return nil }

func (l *listener) Addr() net.Addr { return nil }

func main1() {

	http.HandleFunc("/foo", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
	})

	conn, err := net.Dial("tcp", "127.0.0.1:8533")
	if err != nil {
		panic(err)
	}

	fmt.Println("Listening")
	if err := http.Serve(&listener{conn}, nil); err != nil {
		panic(err)
	}
}

func main2() {
	listener, err := net.Listen("tcp", "127.0.0.1:8533")
	if err != nil {
		panic(err)
	}

	conn, err := listener.Accept()
	if err != nil {
		panic(err)
	}

	client := http.Client{Transport: &http.Transport{Dial: connDialer{conn}.Dial}}

	fmt.Println("Making request")
	res, err := client.Get("http://www.shouldNotMatter.com:8080/foo")
	if err != nil {
		fmt.Println(err)
	}

	////////////////////////  HANGS HERE ////////////////////////////
	fmt.Println("Received response")
	defer res.Body.Close()

	if res.StatusCode == http.StatusOK {
		bodyBytes, err := ioutil.ReadAll(res.Body)
		if err != nil {
			panic(err)
		}
		bodyString := string(bodyBytes)
		fmt.Println(bodyString)
	} else {
		fmt.Println(res)
	}
}
