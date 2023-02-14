package http

import (
	"errors"
	"log"
	"net"
	"net/http"
)

type Server struct {
	server   *http.Server
	listener net.Listener
}

func NewServer(handler http.Handler, listener net.Listener) *Server {
	return &Server{
		server: &http.Server{
			Handler: handler,
		},
		listener: listener,
	}
}

func (s *Server) Serve() error {
	log.Println("Starting the server on address: " + s.listener.Addr().String())

	if err := s.server.Serve(s.listener); err != nil {
		if !errors.Is(err, http.ErrServerClosed) {
			return err
		}
	}

	return nil
}
