package web

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/mashmb/gorest/internal/app/handlers"
	"github.com/mashmb/gorest/internal/app/settings"
)

type server struct {
	settings settings.Settings
}

func NewServer(stg settings.Settings) *server {
	return &server{
		settings: stg,
	}
}

func (s *server) middlewares() middleware {
	return middlewaresChain(
		logRequestAndResponseMiddleware,
		authorizationMiddleware(s.settings),
	)
}

func (s *server) routes(rtr *http.ServeMux) {
	rtr.HandleFunc("GET /hello", handlers.Hello)
}

func (s *server) Run() {
	addr := fmt.Sprintf("%s:%s", s.settings.Server.Host, s.settings.Server.Port)
	router := http.NewServeMux()
	s.routes(router)
	server := &http.Server{
		Addr:    addr,
		Handler: s.middlewares()(router),
	}
	slog.Info("HTTP server started", "addr", addr)
	server.ListenAndServe()
}
