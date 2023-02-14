package client

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
)

type ConnDialer interface {
	Dial(network, addr string) (net.Conn, error)
	CloseConnection() error
}

type Client struct {
	connDialer ConnDialer
}

func NewClient(connDialer ConnDialer) *Client {
	return &Client{
		connDialer: connDialer,
	}
}

func (c *Client) MakeQuery(ctx context.Context, address string) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://"+address, nil)
	if err != nil {
		return err
	}

	defer func() {
		if err := c.connDialer.CloseConnection(); err != nil {
			log.Println(err)
		}
	}()
	client := http.Client{Transport: &http.Transport{Dial: c.connDialer.Dial}}
	res, err := client.Do(req)
	if err != nil {
		return err
	}

	bodyRes, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}
	defer func() {
		if err := res.Body.Close(); err != nil {
			log.Println(err)
		}
	}()

	fmt.Printf("%s %s", res.Status, string(bodyRes))

	return nil
}
