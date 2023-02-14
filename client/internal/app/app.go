package app

import (
	"context"
	"log"
	"time"

	"github.com/arthurshafikov/tz-faraway/client/internal/client"
	"github.com/arthurshafikov/tz-faraway/lib/powtcp"
)

func Run() {
	address := "localhost:3333"
	connDialer, err := powtcp.NewConnDialer(powtcp.ConnDialerOptions{
		Address:              address,
		ReadTimeoutDuration:  time.Second * 5,
		WriteTimeoutDuration: time.Second * 5,
	})
	if err != nil {
		log.Fatalln(err)
	}

	client := client.NewClient(connDialer)
	if err := client.MakeQuery(context.Background(), address); err != nil {
		log.Fatalln(err)
	}
}
