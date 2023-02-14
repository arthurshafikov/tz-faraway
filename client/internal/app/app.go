package app

import (
	"context"
	"log"

	"github.com/arthurshafikov/tz-faraway/client/internal/client"
	"github.com/arthurshafikov/tz-faraway/lib/powtcp"
)

func Run() {
	address := "localhost:3333"
	connDialer, err := powtcp.NewConnDialer(powtcp.ConnDialerOptions{
		Address: address,
	})
	if err != nil {
		log.Fatalln(err)
	}

	client := client.NewClient(connDialer)
	if err := client.MakeQuery(context.Background(), address); err != nil {
		log.Fatalln(err)
	}
}
