package app

import (
	server "github.com/arthurshafikov/tz-faraway/server/internal/transport/http"
	"github.com/arthurshafikov/tz-faraway/server/internal/transport/http/handler"
)

func Run() {
	server.RunServer(handler.NewHandler(), "3333")
}
