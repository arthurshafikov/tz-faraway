package http

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"io"
	"log"
	"math/big"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func RunServer(handler http.Handler) {
	log.Println("Starting the server on port 3333")

	tcpListener, err := net.Listen("tcp", ":3333")
	if err != nil {
		panic(err)
	}
	defer tcpListener.Close()

	if err := http.Serve(&CustomListener{tcpListener}, handler); err != nil {
		if !errors.Is(err, http.ErrServerClosed) {
			log.Fatalln(err)
		}
	}
}
