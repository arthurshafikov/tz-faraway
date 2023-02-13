package http

import (
	"errors"
	"log"
	"net/http"
)

func RunServer(handler http.Handler, port string) {
	log.Println("Starting the server on port " + port)

	proowOfWorkProtectionListener, err := NewProowOfWorkProtectionListener(Options{
		Address: ":" + port,
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
