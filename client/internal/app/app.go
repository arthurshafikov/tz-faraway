package app

import (
	"context"
	"log"

	"github.com/arthurshafikov/tz-faraway/client/internal/client"
)

func Run() {
	if err := client.MakeQuery(context.Background()); err != nil {
		log.Fatalln(err)
	}
}
