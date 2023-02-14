package app

import (
	"log"
	"time"

	"github.com/arthurshafikov/tz-faraway/lib/powtcp"
	server "github.com/arthurshafikov/tz-faraway/server/internal/transport/http"
	"github.com/arthurshafikov/tz-faraway/server/internal/transport/http/handler"
)

func Run() {
	proowOfWorkProtectionListener, err := powtcp.NewProowOfWorkProtectionListener(powtcp.ListenerOptions{
		Address:              ":3333",
		ReadTimeoutDuration:  time.Second * 5,
		WriteTimeoutDuration: time.Second * 5,
	})
	if err != nil {
		log.Fatalln(err)
	}

	server := server.NewServer(handler.NewHandler(), proowOfWorkProtectionListener)
	if err := server.Serve(); err != nil {
		log.Fatalln(err)
	}
}
