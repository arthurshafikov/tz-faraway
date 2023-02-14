package app

import (
	"context"
	"fmt"
	"log"

	"github.com/arthurshafikov/tz-faraway/client/internal/client"
	"github.com/arthurshafikov/tz-faraway/client/internal/config"
	"github.com/arthurshafikov/tz-faraway/lib/powtcp"
)

func Run() {
	config, err := config.NewConfig(".env")
	if err != nil {
		log.Fatalln(fmt.Errorf("config returned error: %w", err))
	}

	connDialer, err := powtcp.NewConnDialer(powtcp.ConnDialerOptions{
		Address: config.ServerAddress,
	})
	if err != nil {
		log.Fatalln(err)
	}

	client := client.NewClient(connDialer)
	if err := client.MakeQuery(context.Background(), config.ServerAddress); err != nil {
		log.Fatalln(err)
	}
}
