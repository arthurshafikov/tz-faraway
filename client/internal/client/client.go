package client

import (
	"fmt"
	"io"
	"log"
	"net/http"
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
