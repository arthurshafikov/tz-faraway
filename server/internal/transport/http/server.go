package http

import (
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/arthurshafikov/tz-faraway/lib/powtcp"
)

func RunServer(handler http.Handler, port string) {
	log.Println("Starting the server on port " + port)

	proowOfWorkProtectionListener, err := powtcp.NewProowOfWorkProtectionListener(powtcp.ListenerOptions{
		Address:              ":" + port,
		ReadTimeoutDuration:  time.Second * 5,
		WriteTimeoutDuration: time.Second * 5,
	})
	if err != nil {
		log.Fatalln(err)
	}

	if err := http.Serve(proowOfWorkProtectionListener, handler); err != nil {
		if !errors.Is(err, http.ErrServerClosed) {
			log.Fatalln(err)
		}
	}
}
